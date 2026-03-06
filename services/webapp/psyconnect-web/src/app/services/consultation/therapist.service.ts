import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { TherapistV1 } from '../../models/consultation.model';

interface ApiResponse<T> {
  message: string;
  status: number;
  data: T;
}

@Injectable({
  providedIn: 'root',
})
export class TherapistService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/therapists`;

  constructor(private http: HttpClient) {}

  getTherapistProfile(): Observable<TherapistV1> {
    return this.http
      .get<ApiResponse<TherapistV1>>(`${this.baseUrl}/me`)
      .pipe(map((res) => res.data));
  }

  getTherapistById(id: string): Observable<TherapistV1> {
    return this.http
      .get<ApiResponse<TherapistV1>>(`${this.baseUrl}/${id}`)
      .pipe(map((res) => res.data));
  }

  createTherapistProfile(data: TherapistV1): Observable<TherapistV1> {
    return this.http
      .post<ApiResponse<TherapistV1>>(`${this.baseUrl}/me`, data)
      .pipe(map((res) => res.data));
  }

  updateTherapistProfile(data: Partial<TherapistV1>): Observable<TherapistV1> {
    return this.http
      .put<ApiResponse<TherapistV1>>(`${this.baseUrl}/me`, data)
      .pipe(map((res) => res.data));
  }

  changeTherapistStatus(status: boolean): Observable<boolean> {
    return this.http
      .patch<ApiResponse<boolean>>(`${this.baseUrl}/me/availability`, {
        status,
      })
      .pipe(map((res) => res.data));
  }

  searchTherapists(
    query: string,
    limit: number = 20,
    skip: number = 0,
  ): Observable<any> {
    const params = { q: query, limit: limit.toString(), skip: skip.toString() };
    const searchUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/search/therapists`;
    return this.http
      .get<ApiResponse<any>>(searchUrl, { params })
      .pipe(map((res) => res.data));
  }
}
