import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

@Injectable({ providedIn: 'root' })
export class PasswordService {
  private apiUrl = `${environment.apiUrl}`;

  constructor(private http: HttpClient) {}

  requestReset(payload: {
    email: string;
  }): Observable<ApiResponse<boolean>> {
    const url = `${this.apiUrl}/identity/req/reset-password`;
    return this.http.post<ApiResponse<boolean>>(url, payload, {
      withCredentials: false,
    });
  }

  confirmReset(payload: {
    username: string;
    resetToken: string;
    email: string;
    newPassword: string;
  }): Observable<ApiResponse<boolean>> {
    const url = `${this.apiUrl}/identity/password`;
    return this.http.put<ApiResponse<boolean>>(url, payload, {
      withCredentials: false,
    });
  }
}
