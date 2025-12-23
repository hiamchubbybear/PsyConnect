import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SKIP_AUTH } from '../auth/auth.interceptor';

export interface RegisterRequest {
  username: string;
  password: string;
  firstName: string;
  lastName: string;
  address: string;
  gender: string;
  email: string;
  role: string;
  avatarUri: string;
  dob: string;
}

@Injectable({ providedIn: 'root' })
export class RegisterService {
  private apiUrl = `${environment.apiUrl}/identity/create`;

  constructor(private http: HttpClient) {}

  register(request: RegisterRequest): Observable<any> {
    return this.http.post<any>(this.apiUrl, request, {
      context: new HttpContext().set(SKIP_AUTH, true),
    });
  }
}
