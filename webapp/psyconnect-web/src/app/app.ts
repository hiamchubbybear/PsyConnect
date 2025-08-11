import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, RouterOutlet } from '@angular/router';
import { Footer } from './components/footer/footer';
import { Header } from './components/header/header';
import { SidebarComponent } from './components/sidebar/sidebar';
import { fadeRouteAnimation } from './route-animation';
import { AuthStateService } from './services/auth/auth-state.service';
import { AuthService } from './services/auth/auth.service';
import { LoaderComponent } from './services/loader/loader.component';
import {
  UserContextService,
  UserProfile,
} from './services/profile/profile-service';

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
  constructor(
    private http: HttpClient,
    private auth: AuthService,
    private userContext: UserContextService,
    private authState: AuthStateService
  ) {}

  ngOnInit() {
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
