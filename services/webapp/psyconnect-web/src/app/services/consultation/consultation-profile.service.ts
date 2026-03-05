import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, of, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { AuthService } from '../auth/auth.service';

export interface ConsultationProfile {
  id: string;
  userId: string;
  role: 'therapist' | 'client';
  
}

@Injectable({
  providedIn: 'root'
})
export class ConsultationProfileService {
  private apiUrl = `${environment.apiUrl}/v1/consultation`;

  constructor(
    private http: HttpClient,
    private authService: AuthService,
    private router: Router
  ) {}

  
  getClientProfile(): Observable<ConsultationProfile | null> {
    return this.http.get<ConsultationProfile>(`${this.apiUrl}/clients/me`).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          console.warn('Client profile not found (401). User needs to create consultation profile.');
          return of(null);
        }
        return throwError(() => error);
      })
    );
  }

  
  getTherapistProfile(): Observable<ConsultationProfile | null> {
    return this.http.get<ConsultationProfile>(`${this.apiUrl}/therapists/me`).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          console.warn('Therapist profile not found (401). User needs to create consultation profile.');
          return of(null);
        }
        return throwError(() => error);
      })
    );
  }

  
  getCurrentProfile(): Observable<ConsultationProfile | null> {
    const role = this.authService.getRole();

    if (!role) {
      console.error('No role found. Cannot fetch consultation profile.');
      return of(null);
    }

    if (role === 'therapist') {
      return this.getTherapistProfile();
    } else {
      return this.getClientProfile();
    }
  }

  
  hasProfile(): Observable<boolean> {
    return this.getCurrentProfile().pipe(
      map(profile => profile !== null)
    );
  }

  
  ensureProfileExists(): Observable<boolean> {
    return this.hasProfile().pipe(
      map(hasProfile => {
        if (!hasProfile) {
          const role = this.authService.getRole();
          console.log(`No consultation profile found for role: ${role}. Redirecting to create profile...`);
          this.router.navigate(['/consultation/create-profile'], {
            queryParams: { role }
          });
          return false;
        }
        return true;
      })
    );
  }

  
  createClientProfile(data: any): Observable<ConsultationProfile> {
    return this.http.post<ConsultationProfile>(`${this.apiUrl}/clients/me`, data);
  }

  
  createTherapistProfile(data: any): Observable<ConsultationProfile> {
    return this.http.post<ConsultationProfile>(`${this.apiUrl}/therapists/me`, data);
  }

  
  updateClientProfile(data: any): Observable<ConsultationProfile> {
    return this.http.put<ConsultationProfile>(`${this.apiUrl}/clients/me`, data);
  }

  
  updateTherapistProfile(data: any): Observable<ConsultationProfile> {
    return this.http.put<ConsultationProfile>(`${this.apiUrl}/therapists/me`, data);
  }
}
