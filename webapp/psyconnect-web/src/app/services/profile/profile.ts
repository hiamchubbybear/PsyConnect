import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, throwError } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Auth } from '../auth/auth';

export interface ProfileResponse {
  code: number;
  message: string;
  data: {
    accountId: string;
    profileId: string;
    firstName: string;
    lastName: string;
    dob: string;
    address: string;
    gender: string;
    avatarUri: string;
  };
}
export interface ProfileUpdateResponse {
  code: number;
  message: string;
  data: {
    username: string;
    profileId: string;
    firstName: string;
    lastName: string;
    dob: string;
    address: string;
    gender: string;
    avatarUri: string;
  };
}
export interface UserProfileUpdateRequest {
  username: string;
  firstName: string;
  lastName: string;
  dob: string;
  address: string;
  gender: string;
  avatarUri: string;
}

@Injectable({
  providedIn: 'root',
})
export class Profile {
  private readonly apiUrl = environment.apiUrl;
  private readonly apiVersion = environment.apiVersion;

  constructor(
    private http: HttpClient,
    private auth: Auth,
  ) {}

  getProfile(): Observable<ProfileResponse> {
    const token = this.auth.getToken();
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<ProfileResponse>(`${this.apiUrl}/profile`, {
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
