// src/app/features/profile-update/services/consultation-client.service.ts
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, of } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Client } from '../../models/consultation.model';

@Injectable({
  providedIn: 'root',
})
export class ConsultationClientService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/clients/me`;

  constructor(private http: HttpClient) {}
  getClientProfile(): Observable<Client> {
    return this.http.get<Client>(this.baseUrl).pipe(
      catchError((err) => {
        console.error('Error loading client profile:', err);
        return of({} as Client);
      })
    );
  }

  getClientConsultationProfile() {
    return this.http.get(`${this.baseUrl}`);
  }
  updateClientConsultationProfile(data: any) {
    return this.http.patch(`${this.baseUrl}`, data);
  }
  saveClientProfile(data: Client): Observable<Client> {
    return this.http.post<Client>(this.baseUrl, data).pipe(
      catchError((err) => {
        console.error('Error saving client profile:', err);
        throw err;
      })
);
  }

  deleteClientProfile(): Observable<void> {
    return this.http.delete<void>(this.baseUrl).pipe(
      catchError((err) => {
        console.error('Error deleting client profile:', err);
        throw err;
      })
    );
  }
}
