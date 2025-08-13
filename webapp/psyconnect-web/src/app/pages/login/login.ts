import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    FormsModule,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { environment } from '../../../environments/environment';
import { Auth } from '../../services/auth/auth';
import { AuthStateService } from '../../services/auth/auth-state.service';
import { LoaderService } from '../../services/loader/loader';
import { Profile, ProfileResponse } from '../../services/profile/profile';
import {
    UserContextService,
    UserProfile,
} from '../../services/profile/profile-service';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, RouterModule],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login implements OnInit {
  loginForm: FormGroup;
  error: string | null = null;
  loading = false;
  apiUrl = environment.apiUrl;

  constructor(
    private fb: FormBuilder,
    private auth: Auth,
    private profile: Profile,
    private router: Router,
    private userContext: UserContextService,
    private loaderService: LoaderService,
    private toastService: ToastService,
    private authState: AuthStateService
  ) {
    this.loginForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(6)]],
      password: ['', [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.checkOAuth2Callback();
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.setLoading(true);

    this.auth.login(this.loginForm.value).subscribe({
      next: () => this.loadUserProfile(),
      error: () => {
        this.authState.showSidebar();
        this.setError('Invalid email or password.');
        this.setLoading(false);
        this.toastService.show(`${this.error}`, 'Failed', ToastType.Error);
      },
    });
  }

  handleOAuth2(provider: string): void {
    window.location.href = `${this.apiUrl}/oauth2/authorization/${provider}`;
  }

  private checkOAuth2Callback(): void {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const email = params.get('email');
    const provider = params.get('provider');
    const platform = params.get('platform');

    if (code && email && provider && platform) {
      this.setLoading(true);
      this.auth.exchangeOAuth2Code(code, email, provider, platform).subscribe({
        next: (res) => {
          console.log(res);
          localStorage.setItem('access_token', res.data.token);
          this.loadUserProfile();
        },
        error: (err) => {
          console.error('OAuth2 exchange failed', err);
          this.setLoading(false);
          this.toastService.show(`${err}`, 'Failed', ToastType.Error);
        },
      });
    }
  }

  private loadUserProfile(): void {
    this.profile.getProfile().subscribe({
      next: (profile: ProfileResponse) => {
        const userProfile: UserProfile = {
          accountId: profile.data.accountId,
          profileId: profile.data.profileId,
          firstName: profile.data.firstName,
          lastName: profile.data.lastName,
          dob: profile.data.dob,
          address: profile.data.address,
          gender: profile.data.gender,
          avatarUri: profile.data.avatarUri,
        };
        console.log(profile.data.avatarUri);
        this.userContext.setUser(userProfile);
        this.toastService.show('Login success', 'Success', ToastType.Success);
        this.setLoading(false);
        this.router.navigate(['/']);
      },
      error: (err) => {
        console.error('Get profile error', err);
        this.setLoading(false);
        this.toastService.show(`Get profile error`, 'Failed', ToastType.Error);
      },
    });
  }

  private setLoading(isLoading: boolean): void {
    this.loading = isLoading;
    isLoading ? this.loaderService.show() : this.loaderService.hide();
  }

  private setError(message: string): void {
    this.error = message;
  }
}
