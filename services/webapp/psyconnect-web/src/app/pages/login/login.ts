import { CommonModule } from '@angular/common';
import { Component, HostListener, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { ProfileResponse } from '../../models/profile';
import { Auth } from '../../services/auth/auth';
import { AuthStateService } from '../../services/auth/auth-state.service';
import { LoaderService } from '../../services/loader/loader';
import { Profile } from '../../services/profile/profile';
import {
  UserContextService,
  UserProfile,
} from '../../services/profile/profile-service';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';

import { PsyButtonComponent } from '../../shared/ui-atoms/button/psy-button.component';
import { InputComponent } from '../../shared/ui-atoms/input/input.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    RouterModule,
    TranslateModule,
    InputComponent,
    PsyButtonComponent,
  ],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login implements OnInit {
  USERNAME_KEY = environment.usernameKey;
  ACCESSTOKEN_KEY = environment.accessTokenKey;
  loginForm: FormGroup;
  error: string | null = null;
  loading = false;
  apiUrl = environment.apiUrl;
  shakeErrors = false;
  showEmailForm = false; 

  constructor(
    private fb: FormBuilder,
    private auth: Auth,
    private profile: Profile,
    private router: Router,
    private userContext: UserContextService,
    private loaderService: LoaderService,
    private toastService: ToastService,
    private authState: AuthStateService,
    private secureStorage: SecureStorageService,
  ) {
    this.loginForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(6)]],
      password: ['', [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.checkOAuth2Callback();
  }

  toggleEmailForm(): void {
    this.showEmailForm = !this.showEmailForm;
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      this.triggerShake();
      return;
    }

    this.setLoading(true);

    this.auth.login(this.loginForm.value).subscribe({
      next: () => {
        this.loadUserProfile();
        this.secureStorage.setItem(
          this.USERNAME_KEY,
          this.loginForm.value.username,
        );

        const username = this.secureStorage.getItem(this.USERNAME_KEY);
      },
      error: () => {
        this.setError('Invalid email or password.');
        this.setLoading(false);
        this.toastService.show('', `${this.error}`, ToastType.Error);
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
          this.secureStorage.setItem(this.ACCESSTOKEN_KEY, res.data.token);
          this.secureStorage.setItem(this.USERNAME_KEY, email);
          this.loadUserProfile();
        },
        error: (err) => {
          console.error('OAuth2 exchange failed', err);
          this.setLoading(false);
          this.toastService.show('', `${err}`, ToastType.Error);
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

        this.userContext.setUser(userProfile);
        this.toastService.show('', 'TOAST.login_success', ToastType.Success);
        this.setLoading(false);
        this.router.navigate(['/feature/feed']);
        this.authState.setLoggedIn(true);
      },
      error: (err) => {
        console.error('Get profile error', err);
        this.setLoading(false);
        this.toastService.show('', 'TOAST.error_generic', ToastType.Error);
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
  private triggerShake() {
    this.shakeErrors = false;
    requestAnimationFrame(() => {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 320);
    });
  }
  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      this.onSubmit();
    }
  }
}
