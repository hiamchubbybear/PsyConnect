import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  ConsultationSession,
  SessionRequest,
} from '../../models/consultation.model';

export interface CalendarEvent {
  id: string;
  title: string;
  start: string;
  end: string;
  mode: string;
  status: string;
  color: string;
  therapist_id?: string;
  client_id?: string;
}

export interface OverviewData {
  upcoming_count: number;
  completed_count: number;
  cancelled_count: number;
  pending_payment_count: number;
  total_spent: number;
  next_session?: ConsultationSession;
  recent_sessions: ConsultationSession[];
}

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/sessions`;

  constructor(private http: HttpClient) {}

  getAllSessions(): Observable<ConsultationSession[]> {
    return this.http.get<ConsultationSession[]>(`${this.baseUrl}/me`);
  }

  getSessionById(id: string): Observable<ConsultationSession> {
    return this.http.get<ConsultationSession>(`${this.baseUrl}/${id}`);
  }

  createSession(data: SessionRequest): Observable<ConsultationSession> {
    return this.http.post<ConsultationSession>(`${this.baseUrl}/me`, data);
  }

  deleteSession(id: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }

  cancelSession(id: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }

  getCalendar(
    from?: string,
    to?: string,
  ): Observable<{ events: CalendarEvent[] }> {
    let params = new HttpParams();
    if (from) params = params.set('from', from);
    if (to) params = params.set('to', to);
    return this.http.get<{ events: CalendarEvent[] }>(
      `${this.baseUrl}/me/calendar`,
      { params },
    );
  }

  getOverview(): Observable<OverviewData> {
    return this.http.get<OverviewData>(`${this.baseUrl}/me/overview`);
  }

  getPaymentUrl(id: string): Observable<{ payment_url: string }> {
    return this.http.get<{ payment_url: string }>(
      `${this.baseUrl}/${id}/payment-url`,
    );
  }

  processRefund(id: string): Observable<{ refund_trace_id: string }> {
    return this.http.post<{ refund_trace_id: string }>(
      `${this.baseUrl}/${id}/refund`,
      {},
    );
  }

  getAllSessionsAdmin(): Observable<ConsultationSession[]> {
    return this.http.get<ConsultationSession[]>(
      `${environment.apiUrl}/${environment.apiVersion}/consultation/admin/sessions`,
    );
  }
}
