import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environment';
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

@Injectable({
  providedIn: 'root',
})
export class Profile {
  private readonly apiUrl = environment.apiUrl;

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
}
