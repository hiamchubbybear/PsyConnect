import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
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
    private router: Router
  ) {}

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
          localStorage.setItem('access_token', res.data.token);
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
