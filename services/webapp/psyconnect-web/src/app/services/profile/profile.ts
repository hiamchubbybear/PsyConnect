import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, throwError } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
    ProfileResponse,
    ProfileUpdateResponse,
    UserProfileUpdateRequest,
} from '../../models/profile';
import { Auth } from '../auth/auth';

@Injectable({
  providedIn: 'root',
})
export class Profile {
  private readonly apiUrl = environment.apiUrl;
  private readonly apiVersion = environment.apiVersion;

  constructor(private http: HttpClient, private auth: Auth) {}

  getProfile(): Observable<ProfileResponse> {
    const token = this.auth.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<ProfileResponse>(`${this.apiUrl}/profile`, {
      headers,
    });
  }

  getProfileById(userId: string): Observable<ProfileResponse> {
    const token = this.auth.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<ProfileResponse>(`${this.apiUrl}/${this.apiVersion}/profile/${userId}`, {
      headers,
    });
  }

  updateProfile(
    updateProfile: UserProfileUpdateRequest
  ): Observable<ProfileUpdateResponse> {
    const token = this.auth.getToken();

    if (!token) {
      return throwError(() => new Error('Authentication required'));
    }

    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    return this.http
      .put<ProfileUpdateResponse>(
        `${this.apiUrl}/${this.apiVersion}/profile/me`,
        updateProfile,
        { headers }
      )
      .pipe(
        catchError((error) => {
          if (error.status === 401) {
            this.auth.logout();
          }
          return throwError(() => error);
        })
      );
  }
}
