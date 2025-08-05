import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
import { Footer } from './components/footer/footer';
import { Header } from './components/header/header';
import { fadeRouteAnimation } from './route-animation';
import { AuthService } from './services/auth/auth.service';
import { LoaderComponent } from './services/loader/loader.component';
import {
    UserContextService,
    UserProfile,
} from './services/profile/profile-service.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    Header,
    Footer,
    LoaderComponent,
    FormsModule,
    ReactiveFormsModule,
  ],
  animations: [fadeRouteAnimation],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit {
  constructor(
    private http: HttpClient,
    private auth: AuthService,
    private userContext: UserContextService
  ) {}

  ngOnInit() {
    const token = this.auth.getToken();
    if (token) {
      const cachedUser = this.userContext.getUser();
      if (cachedUser) {
      }

      this.http.get<UserProfile>('/api/profile').subscribe({
        next: (profile) => {
          const currentUser = this.userContext.getUser();
          if (
            !currentUser ||
            JSON.stringify(currentUser) !== JSON.stringify(profile)
          ) {
            this.userContext.setUser(profile);
          }
        },
        error: (error) => {
          if (error.status === 401 || error.status === 403) {
            this.auth.logout();
          }
        },
      });
    } else {
      this.userContext.clear();
    }
  }
}
