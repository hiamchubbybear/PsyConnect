import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router } from '@angular/router';

@Injectable({ providedIn: 'root' })
export class PasswordResetGuard implements CanActivate {
  constructor(private router: Router) {}

  canActivate(route: ActivatedRouteSnapshot): boolean {
    const username = route.queryParamMap.get('username');
    const email = route.queryParamMap.get('email');
    const token = route.queryParamMap.get('token');

    if (!username || !email || !token) {
      this.router.navigate(['/auth/login']);
      return false;
    }
    return true;
  }
}
