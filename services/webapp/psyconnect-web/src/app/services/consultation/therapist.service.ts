import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { TherapistV1 } from '../../models/consultation.model';

@Injectable({
  providedIn: 'root',
})
export class TherapistService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/therapists`;

  constructor(private http: HttpClient) {}

  getTherapistProfile(): Observable<TherapistV1> {
    return this.http.get<TherapistV1>(`${this.baseUrl}/me`);
  }

  getTherapistById(id: string): Observable<TherapistV1> {
    return this.http.get<TherapistV1>(`${this.baseUrl}/${id}`);
  }

  createTherapistProfile(data: TherapistV1): Observable<TherapistV1> {
    return this.http.post<TherapistV1>(`${this.baseUrl}/me`, data);
  }

  updateTherapistProfile(data: Partial<TherapistV1>): Observable<TherapistV1> {
    return this.http.put<TherapistV1>(`${this.baseUrl}/me`, data);
  }

  changeTherapistStatus(status: boolean): Observable<boolean> {
    return this.http.patch<boolean>(`${this.baseUrl}/me/availability`, { status });
  }
}
