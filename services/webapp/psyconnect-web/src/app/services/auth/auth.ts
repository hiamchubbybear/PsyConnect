import { HttpClient, HttpContext, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { SKIP_AUTH } from './auth.interceptor';
@Injectable({
  providedIn: 'root',
})
export class Auth {
  apiUrl = `${environment.apiUrl}`;
  apiVersion = `${environment.apiVersion}`;
  refreshKey = `${environment.refreshKey}`;
  PROFILE_KEY = environment.profileKey;
  USERNAME_KEY = environment.usernameKey;
  ACCESSTOKEN_KEY = environment.accessTokenKey;
  THERAPIST_KEY = environment.therapistsKey;
  constructor(
    private http: HttpClient,
    private secureStorage: SecureStorageService,
  ) {}

  login(credentials: {
    username: string;
    password: string;
  }): Observable<{ data: { token: string; refreshToken: string } }> {
    const url = `${this.apiUrl}/${this.apiVersion}/auth/login?provider=NORMAL&platform=web`;
    return this.http
      .post<{ data: { token: string; refreshToken: string } }>(
        url,
        credentials,
        {
          withCredentials: false,
        },
      )
      .pipe(
        tap((res) => {
          console.log(res);
          let resToken = res?.data.token;
          let refreshToken = res?.data.refreshToken;
          if (resToken != null && resToken != '') {
            console.log('Token: ' + resToken);
            console.log('Refresh Token: ' + refreshToken);
            this.secureStorage.setItem(this.ACCESSTOKEN_KEY, resToken);
            this.secureStorage.setItem(this.refreshKey, refreshToken);
          }
        }),
      );
  }

  logout(): void {
    this.secureStorage.clear();
  }

  isLoggedIn(): boolean {
    const token = this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
    return !!(token && token !== 'null' && token !== 'undefined');
  }

  getToken(): string | null {
    return this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
  }

  loginWithGoogle(token: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/oauth2/authorization/google`, {
      token,
    });
  }

  loginWithFacebook(token: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/oauth2/authorization/facebook`, {
      token,
    });
  }
  isTokenExpired(token: string): boolean {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) return true;
      let payload = parts[1].replace(/-/g, '+').replace(/_/g, '/');
      const pad = payload.length % 4 === 0 ? '' : '='.repeat(4 - (payload.length % 4));
      const decodedInfo = JSON.parse(atob(payload + pad));
      const exp = decodedInfo.exp * 1000;
      return Date.now() > exp;
    } catch (e) {
      return true;
    }
  }
  checkTokenOnStartup() {
    const token = this.secureStorage.getItem<string>(this.ACCESSTOKEN_KEY);
    const username = this.secureStorage.getItem<string>(this.USERNAME_KEY);
    if (!token || !username) {
      this.logout();
      return;
    }
    if (token && this.isTokenExpired(token)) {
      this.refreshToken(username).subscribe({
        next: () => console.log('Token refreshed on startup'),
        error: () => this.logout(),
      });
    }
  }

  refreshToken(username: string): Observable<string> {
    const refreshToken = this.secureStorage.getItem<string>(this.refreshKey);
    console.log(refreshToken);
    console.log('User name' + username);
    const params = new HttpParams()
      .set('provider', 'NORMAL')
      .set('platform', 'web')
      .set('token', refreshToken || '')
      .set('username', username);

    return this.http
      .post<{
        code: number;
        message: string;
        data: { token: string; refreshToken: string; successful: boolean };
      }>(
        `${this.apiUrl}/${this.apiVersion}/auth/refresh`,
        {},
        { params, context: new HttpContext().set(SKIP_AUTH, true) },
      )
      .pipe(
        tap((res) => {
          this.secureStorage.setItem(this.ACCESSTOKEN_KEY, res.data.token);
          this.secureStorage.setItem(this.refreshKey, res.data.refreshToken);
        }),
        map((res) => res.data.token),
      );
  }

  exchangeOAuth2Code(
    code: string,
    email: string,
    provider: string,
    plateform: string,
  ): Observable<any> {
    const params = {
      code: code,
      email: email,
      provider: provider,
      platform: plateform,
    };
    return this.http.get(`${this.apiUrl}/auth/oauth2/callback/exchange-code`, {
      params,
    });
  }
}
