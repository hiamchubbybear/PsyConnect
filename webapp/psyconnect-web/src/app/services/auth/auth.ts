import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { AuthStateService } from './auth-state.service';

@Injectable({
  providedIn: 'root',
})
export class Auth {
  apiUrl = `${environment.apiUrl}`;

  constructor(
    private http: HttpClient,
    private secureStorage: SecureStorageService,
    private authState: AuthStateService
  ) {}
  

  login(credentials: {
    username: string;
    password: string;
  }): Observable<{ data: { token: string } }> {
    const url = `${this.apiUrl}/auth/login?provider=NORMAL&platform=web`;
    return this.http
      .post<{ data: { token: string } }>(url, credentials, {
        withCredentials: false,
      })
      .pipe(
        tap((res) => {
          let resToken = res?.data.token;
          if (resToken != null && resToken != '') {
            console.log(resToken);
            this.secureStorage.setItem('access_token', res?.data.token);
          }
        })
      );
  }
  logout(): void {
    this.secureStorage.removeItem('access_token');
  }

  isLoggedIn(): boolean {
    return !!this.secureStorage.getItem('access_token');
  }
  getToken(): string | null {
    return this.secureStorage.getItem('access_token');
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
  exchangeOAuth2Code(
    code: string,
    email: string,
    provider: string,
    plateform: string
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
