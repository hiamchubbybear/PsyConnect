import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, Renderer2 } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {
    Event,
    NavigationCancel,
    NavigationEnd,
    NavigationStart,
    Router,
    RouterModule,
    RouterOutlet
} from '@angular/router';
import { filter, Subscription } from 'rxjs';
import { Footer } from './components/footer/footer';
import { Header } from './components/header/header';
import { SidebarComponent } from './components/sidebar/sidebar';
import { fadeRouteAnimation } from './route-animation';
import { AuthStateService } from './services/auth/auth-state.service';
import { AuthService } from './services/auth/auth.service';
import { LoaderService } from './services/loader/loader';
import { LoaderComponent } from './services/loader/loader.component';
import {
    UserContextService,
    UserProfile,
} from './services/profile/profile-service';
import { ThemeService } from './services/theme/theme-service';

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
  ],
  animations: [fadeRouteAnimation],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit {
  showSidebar: boolean = false;
  showFooter = true;

  constructor(
    private renderer: Renderer2,
    private themeService: ThemeService,
    private http: HttpClient,
    private auth: AuthService,
    private userContext: UserContextService,
    private authState: AuthStateService,
    private router: Router,
    private loader: LoaderService
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

    this.router.events.subscribe((event: Event) => {
      if (event instanceof NavigationStart) {
        this.loader.show();
      } else if (
        event instanceof NavigationEnd ||
        event instanceof NavigationCancel ||
        event instanceof NavigationEnd
      ) {
        this.loader.hide();
      }
    });
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
