import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { Auth } from '../../services/auth/auth';
import { UserContextService } from '../../services/profile/profile-service';

@Component({
  selector: 'app-oauth2-callback',
  standalone: true,
  imports: [CommonModule],
  template: '<p>Processing OAuth2...</p>',
})
export class OAuth2CallbackComponent implements OnInit {
  constructor(
    private auth: Auth,
    private userContext: UserContextService,
    private router: Router,
    private secureStorage: SecureStorageService
  ) {}
  PROFILE_KEY = environment.profileKey;
  ACCESSTOKEN_KEY = environment.accessTokenKey;
  ngOnInit() {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const email = params.get('email');
    const provider = params.get('provider');
    const platform = params.get('platform');

    console.log('OAuth2 callback: ', code, email, provider, platform);

    if (code && email && provider && platform) {
      this.auth.exchangeOAuth2Code(code, email, provider, platform).subscribe({
        next: (res) => {
          console.log('OAuth2 exchange success:', res);
          this.secureStorage.setItem(this.ACCESSTOKEN_KEY, res.data.token);
          this.secureStorage.setItem(
            this.PROFILE_KEY,
            JSON.stringify(res.data)
          );
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
