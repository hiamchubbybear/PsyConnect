import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    FormsModule,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { environment } from '../../environment';
import { Auth } from '../../services/auth/auth';
import { LoaderService } from '../../services/loader/loader';
import { Profile, ProfileResponse } from '../../services/profile/profile';
import {
    UserContextService,
    UserProfile,
} from '../../services/profile/profile-service.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  loginForm: FormGroup;
  error: string | null = null;
  loading = false;
  avatarUri: string | null = null;
  apiUrl = environment.apiUrl;
  constructor(
    private fb: FormBuilder,
    private auth: Auth,
    private profile: Profile,
    private router: Router,
    private userContext: UserContextService,
    private loaderService: LoaderService
  ) {
    this.loginForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(6)]],
      password: ['', [Validators.required]],
    });
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.loading = true;
    this.error = null;
    this.loaderService.show();
    this.auth.login(this.loginForm.value).subscribe({
      next: () => this.loadProfile(),
      error: () => {
        this.error = 'Invalid email or password.';
        this.loading = false;
      },
    });
  }

  private loadProfile() {
    this.profile.getProfile().subscribe({
      next: (profile: ProfileResponse) => {
        const data = profile.data;
        this.avatarUri = data.avatarUri || null;
        console.log(profile);
        const userProfile: UserProfile = {
          accountId: data.accountId,
          profileId: data.profileId,
          firstName: data.firstName,
          lastName: data.lastName,
          dob: data.dob,
          address: data.address,
          gender: data.gender,
          avatarUri: data.avatarUri,
        };
        console.log(userProfile);
        this.userContext.setUser(userProfile);
        this.loaderService.hide();
        this.router.navigate(['/']);
        this.loading = false;
      },
      error: (err) => {
        console.error('Get profile error', err);
        this.loading = false;
      },
    });
  }

  handleOAuth2(provider: string) {
    console.log(`${this.apiUrl}/oauth2/authorization/${provider}`);
    window.location.href = `${this.apiUrl}/oauth2/authorization/${provider}`;
  }
  ngOnInit() {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const email = params.get('email');
    const provider = params.get('provider');

    if (code && email && provider) {
      this.auth.exchangeOAuth2Code(code, email, provider).subscribe({
        next: (res) => {
          console.log('OAuth2 exchange success:', res);
          localStorage.setItem('token', res.data.token);
          localStorage.setItem('profile', JSON.stringify(res.data));
          this.userContext.setUser(res.data);
          this.router.navigate(['/']);
        },
        error: (err) => {
          console.error('OAuth2 exchange failed', err);
        },
      });
    }
  }
}
