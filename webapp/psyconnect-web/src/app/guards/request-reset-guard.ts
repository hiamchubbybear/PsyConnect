import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

@Injectable({ providedIn: 'root' })
export class RequestResetSuccessGuard implements CanActivate {
  constructor(private router: Router) {}

  canActivate(): boolean {
    const nav = this.router.getCurrentNavigation();
    const email = (nav?.extras.state as any)?.email;

    if (!email) {
      this.router.navigate(['/auth/login']);
      return false;
    }
    return true;
  }
}
