// src/app/guards/auth.guard.ts
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth/auth.service';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // Check if user is logged in
  if (!authService.isLoggedIn()) {
    console.log('❌ Not logged in. Redirecting to login...');
    router.navigate(['/auth/login']);
    return false;
  }

  // Check if token is expired
  if (authService.isTokenExpired()) {
    console.log('❌ Token expired. Redirecting to login...');
    authService.logout();
    router.navigate(['/auth/login']);
    return false;
  }

  console.log('✅ Auth guard passed');
  return true;
};
