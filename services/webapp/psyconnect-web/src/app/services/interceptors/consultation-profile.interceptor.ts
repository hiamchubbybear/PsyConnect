import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../auth/auth.service';

/**
 * Interceptor to handle 404 errors from consultation API
 * Redirects to create profile page instead of 404 page
 */
export const consultationProfileInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const authService = inject(AuthService);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      // Check if it's a 404 from consultation API
      if (error.status === 404 && isConsultationProfileRequest(req.url)) {
        console.log('🔍 Consultation profile not found (404). Redirecting to create profile...');

        const role = authService.getRole();

        if (role) {
          // Redirect to create profile page with role
          router.navigate(['/feature/consultation/create-profile'], {
            queryParams: { role }
          });
        } else {
          console.warn('No role found. Cannot redirect to create profile.');
        }

        // Return null instead of error to prevent error propagation
        return throwError(() => new Error('Profile not found. Redirecting to create profile.'));
      }

      // For other errors, pass through
      return throwError(() => error);
    })
  );
};

/**
 * Check if the request is for consultation profile
 */
function isConsultationProfileRequest(url: string): boolean {
  const consultationProfilePatterns = [
    '/consultation/clients/me',
    '/consultation/therapists/me',
    '/consultation/clients/',
    '/consultation/therapists/'
  ];

  return consultationProfilePatterns.some(pattern => url.includes(pattern));
}
