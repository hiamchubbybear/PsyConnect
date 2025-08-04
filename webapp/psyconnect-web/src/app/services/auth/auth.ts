import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environment';

@Injectable({
  providedIn: 'root',
})
export class Auth {
  private apiUrl = `${environment.apiUrl}/auth`;

  constructor(private http: HttpClient) {}

  login(credentials: {
    username: string;
    password: string;
  }): Observable<{ data: { token: string } }> {
    const url = `${this.apiUrl}/login?provider=NORMAL&platform=web`;
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
}
