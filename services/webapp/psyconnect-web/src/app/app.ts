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
import { filter, Subscription } from 'rxjs';
import { Footer } from './components/footer/footer';
import { Header } from './components/header/header';
import { SidebarComponent } from './components/sidebar/sidebar';
import { MobileRequiredComponent } from './pages/mobile/mobile';
import { fadeRouteAnimation } from './route-animation';
import { AuthStateService } from './services/auth/auth-state.service';
import { AuthService } from './services/auth/auth.service';
import { LoaderService } from './services/loader/loader';
import { LoaderComponent } from './services/loader/loader.component';
import { NotificationService } from './services/notification/notification.service';
import {
    UserContextService,
    UserProfile,
} from './services/profile/profile-service';
import { ThemeService } from './services/theme/theme-service';
import { ToastContainerComponent } from './shared/toast/toast-container';
import { ToastType } from './shared/toast/toast-type';
import { ToastService } from './shared/toast/toast.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    Header,
    Footer,
    LoaderComponent,
    FormsModule,
    SidebarComponent,
    ReactiveFormsModule,
    CommonModule,
    RouterModule,
    MobileRequiredComponent,
    ToastContainerComponent,
  ],
  animations: [fadeRouteAnimation],
  templateUrl: './app.html',
})
export class App implements OnInit {
  showSidebar: boolean = false;
  showFooter = true;

  minWidth = 1235;
  minHeight = 277;
  screenOk = true;

  constructor(
    private themeService: ThemeService,
    private http: HttpClient,
    private auth: AuthService,
    private userContext: UserContextService,
    private authState: AuthStateService,
    private router: Router,
    private loader: LoaderService,
    private notification: NotificationService,
    private toastService: ToastService
  ) {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: any) => {
        const url = event.urlAfterRedirects;
        this.showFooter =
          url === '/' || url.startsWith('/about') || url.startsWith('/contact');
      });
  }

  private routerSub?: Subscription;
  isLoading = false;

  ngOnInit() {
    this.setInitialTheme();

    this.checkScreenSize(window.innerWidth, window.innerHeight);

    this.authState.sidebarVisible$.subscribe((visible) => {
      this.showSidebar = visible;
    });

    const token = this.auth.getToken();
    if (token) {
      this.authState.showSidebar();
      this.fetchUserProfile();
    } else {
      this.authState.hideSidebar();
    }

    this.loader.loading$.subscribe((v) => (this.isLoading = v));
    const profileId = this.userContext.getUser()?.profileId;

    this.router.events.subscribe((event: Event) => {
      if (event instanceof NavigationStart) {
        this.loader.show();
      } else if (
        event instanceof NavigationEnd ||
        event instanceof NavigationCancel
      ) {
        this.loader.hide();
      }
    });
    if (profileId) {
      this.notification.init(profileId);
      this.toastService.show(
        'TOAST.error_generic',
        'TOAST.error_notification',
        ToastType.Error
      );
    } else {
      this.toastService.show(
        'TOAST.error_generic',
        'TOAST.error_notification',
        ToastType.Error
      );
    }
  }

  @HostListener('window:resize', ['$event'])
  onResize(event: any) {
    this.checkScreenSize(event.target.innerWidth, event.target.innerHeight);
  }

  private checkScreenSize(width: number, height: number) {
    this.screenOk = width >= this.minWidth && height >= this.minHeight;
  }

  ngOnDestroy() {
    this.routerSub?.unsubscribe();
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
      },
      error: (error) => {
        if (error.status === 401 || error.status === 403) {
          this.auth.logout();
          this.authState.hideSidebar();
        }
      },
    });
  }
}
