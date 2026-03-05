import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../auth/auth.service';


export const consultationProfileInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const authService = inject(AuthService);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      
      if (error.status === 404 && isConsultationProfileRequest(req.url)) {
        console.log('🔍 Consultation profile not found (404). Redirecting to create profile...');

        const role = authService.getRole();

        if (role) {
          
          router.navigate(['/feature/consultation/create-profile'], {
            queryParams: { role }
          });
        } else {
          console.warn('No role found. Cannot redirect to create profile.');
        }

        
        return throwError(() => new Error('Profile not found. Redirecting to create profile.'));
      }

      
      return throwError(() => error);
    })
  );
};


function isConsultationProfileRequest(url: string): boolean {
  const consultationProfilePatterns = [
    '/consultation/clients/me',
    '/consultation/therapists/me',
    '/consultation/clients/',
    '/consultation/therapists/'
  ];

  return consultationProfilePatterns.some(pattern => url.includes(pattern));
}
