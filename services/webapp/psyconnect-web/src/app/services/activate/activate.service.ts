import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SKIP_AUTH } from '../auth/auth.interceptor';

export interface ActivateRequest {
  email: string;
  token: string;
  verifedTime: string;
}

export interface ResendActivationRequest {
  email: string;
}

@Injectable({ providedIn: 'root' })
export class ActivateService {
  private activateUrl = `${environment.apiUrl}/identity/activate`;
  private resendUrl = `${environment.apiUrl}/identity/req/activate`;

  constructor(private http: HttpClient) {}

  activate(request: ActivateRequest): Observable<any> {
    return this.http.post<any>(this.activateUrl, request, {
      context: new HttpContext().set(SKIP_AUTH, true),
    });
  }

  resendActivation(email: string): Observable<any> {
    return this.http.post<any>(
      this.resendUrl,
      { email },
      {
        context: new HttpContext().set(SKIP_AUTH, true),
      }
    );
  }
}
