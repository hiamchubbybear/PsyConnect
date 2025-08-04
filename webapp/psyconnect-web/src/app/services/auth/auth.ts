import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class Auth {
  apiUrl = `${environment.apiUrl}`;

  constructor(private http: HttpClient) {}

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
            localStorage.setItem('access_token', res?.data.token);
          }
        })
      );
  }
  logout(): void {
    localStorage.removeItem('access_token');
  }

  isLoggedIn(): boolean {
    return !!localStorage.getItem('access_token');
  }
  getToken(): string | null {
    return localStorage.getItem('access_token');
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
