import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { DeleteSessionRequest, Session, SessionRequest } from '../../models/consultation.model';

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/session`;

  constructor(private http: HttpClient) {}

  getAllSessions(): Observable<Session[]> {
    return this.http.get<Session[]>(`${this.baseUrl}/all`);
  }

  getSessionById(id: string): Observable<Session> {
    return this.http.get<Session>(`${this.baseUrl}/${id}`);
  }

  createSession(data: SessionRequest): Observable<Session> {
    return this.http.post<Session>(this.baseUrl, data);
  }

  deleteSession(data: DeleteSessionRequest): Observable<void> {
    return this.http.delete<void>(this.baseUrl, { body: data });
  }

  getAllSessionsAdmin(): Observable<Session[]> {
    return this.http.get<Session[]>(
      `${environment.apiUrl}/${environment.apiVersion}/consultation/admin/session`
    );
  }
}
