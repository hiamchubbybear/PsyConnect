import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, HostListener, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {
  Event,
  NavigationCancel,
  NavigationEnd,
  NavigationStart,
  Router,
  RouterModule,
  RouterOutlet,
} from '@angular/router';
import { filter, Subscription, timer } from 'rxjs';
import { AuthHeaderComponent } from './components/auth-header/auth-header';
import { FloatingChatComponent } from './components/chat/floating-chat/floating-chat';
import { Footer } from './components/footer/footer';
import { Header } from './components/header/header';
import { SidebarComponent } from './components/sidebar/sidebar';
import { MobileRequiredComponent } from './pages/mobile/mobile';
import { fadeRouteAnimation } from './route-animation';
import { AuthStateService } from './services/auth/auth-state.service';
import { AuthService } from './services/auth/auth.service';
import { SessionService } from './services/consultation/session.service';
import { LoaderService } from './services/loader/loader';
import { LoaderComponent } from './services/loader/loader.component';
import { NotificationService } from './services/notification/notification.service';
import {
  UserContextService,
  UserProfile,
} from './services/profile/profile-service';
import { ThemeService } from './services/theme/theme-service';
import { ToastContainerComponent } from './shared/toast/toast-container';
import { ToastService } from './shared/toast/toast.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    Header,
    AuthHeaderComponent,
    Footer,
    LoaderComponent,
    FormsModule,
    SidebarComponent,
    ReactiveFormsModule,
    CommonModule,
    RouterModule,
    MobileRequiredComponent,
    ToastContainerComponent,
    FloatingChatComponent,
  ],
  animations: [fadeRouteAnimation],
  templateUrl: './app.html',
})
export class App implements OnInit {
  showSidebar: boolean = false;
  showFooter = true;

  minWidth = 320;
  minHeight = 400;
  screenOk = true;
  isAuthRoute = false;

  private subscriptions = new Subscription();

  constructor(
    private themeService: ThemeService,
    private http: HttpClient,
    private auth: AuthService,
    private userContext: UserContextService,
    private authState: AuthStateService,
    private router: Router,
    private loader: LoaderService,
    private notification: NotificationService,
    private toastService: ToastService,
    private sessionService: SessionService,
  ) {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: any) => {
        const url = event.urlAfterRedirects;
        this.showFooter =
          url === '/' || url.startsWith('/about') || url.startsWith('/contact');

        this.isAuthRoute = url.startsWith('/auth') || url.startsWith('/oauth2');
      });
  }

  private routerSub?: Subscription;
  isLoading = false;

  ngOnInit() {
    this.setInitialTheme();

    this.checkScreenSize(window.innerWidth, window.innerHeight);

    this.subscriptions.add(
      this.authState.sidebarVisible$.subscribe((visible) => {
        this.showSidebar = visible;
      }),
    );

    this.fetchUserProfile();

    this.subscriptions.add(
      this.loader.loading$.subscribe((v) => (this.isLoading = v)),
    );

    this.subscriptions.add(
      this.router.events.subscribe((event: Event) => {
        if (event instanceof NavigationStart) {
          this.loader.show();
        } else if (
          event instanceof NavigationEnd ||
          event instanceof NavigationCancel
        ) {
          this.loader.hide();
        }
      }),
    );
  }

  @HostListener('window:resize', ['$event'])
  onResize(event: any) {
    this.checkScreenSize(event.target.innerWidth, event.target.innerHeight);
  }

  private checkScreenSize(width: number, height: number) {
    this.screenOk = width >= this.minWidth && height >= this.minHeight;
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  setInitialTheme() {
    const savedTheme = localStorage.getItem('app-theme');
    if (savedTheme === 'dark') {
      this.themeService.setTheme('dark');
    } else {
      this.themeService.setTheme('light');
    }
  }

  private fetchUserProfile() {
    this.http.get<UserProfile>('/api/profile').subscribe({
      next: (profile) => {
        this.userContext.setUser(profile);
        this.authState.showSidebar();
        if (profile.profileId) {
          this.notification.init(profile.profileId);
          // Poll every minute
          this.subscriptions.add(
            timer(0, 60000).subscribe(() => this.fetchNextSession()),
          );
        }
      },
      error: (error) => {
        if (error.status === 401 || error.status === 403) {
          this.auth.logout();
          this.authState.hideSidebar();
        }
      },
    });
  }

  private fetchNextSession() {
    this.sessionService.getOverview().subscribe({
      next: (data) => {
        if (!data) return;

        // Find the best candidate for the "next session"
        // 1. Precise next_session from backend
        // 2. Or search in recent_sessions for anything confirmed/upcoming
        let session = data.next_session as any;

        if (
          !session &&
          data.recent_sessions &&
          data.recent_sessions.length > 0
        ) {
          const now = new Date();
          session = data.recent_sessions.find((s: any) => {
            const start = new Date(s.start_time || s.startTime);
            return (
              start > now &&
              (s.status === 'CONFIRMED' || s.status === 'PENDING_PAYMENT')
            );
          });
        }

        if (session) {
          const startTime = session.start_time || session.startTime;
          const start = new Date(startTime);

          let timeText = 'Giờ chưa xác định';
          if (!isNaN(start.getTime())) {
            const now = new Date();
            const diffMs = start.getTime() - now.getTime();
            if (diffMs <= 0) {
              timeText = 'Đang diễn ra';
            } else {
              const minutes = Math.floor(diffMs / 60000);
              const hours = Math.floor(minutes / 60);
              if (hours > 24)
                timeText = `Bắt đầu sau ${Math.floor(hours / 24)} ngày`;
              else if (hours > 0)
                timeText = `Bắt đầu sau ${hours} giờ ${minutes % 60} phút`;
              else timeText = `Bắt đầu sau ${minutes} phút`;
            }
          }

          this.toastService.session({
            id: 'upcoming-session',
            title: `Phiên họp: ${session.title || 'Tư vấn tâm lý'}`,
            message: timeText,
            actions: [
              {
                label: 'Tham gia',
                primary: true,
                action: () => {
                  this.router.navigate([
                    '/feature/chat',
                    session.conversationId ||
                      session.conversation_id ||
                      session.session_id,
                  ]);
                },
              },
            ],
          });
        } else {
          // If no session found, ensure we clear the persistent notification if it exists
          this.toastService.remove('upcoming-session');
        }
      },
      error: (err) => console.error('❌ Error fetching session overview:', err),
    });
  }
}
