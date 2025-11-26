import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';

interface JWTPayload {
  sub: string;
  accountId: string;
  profileId: string;
  scope: string;
  iss: string;
  exp: number;
  type: string;
  iat: number;
  jti: string;
  platform: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  constructor(private secureService: SecureStorageService) {}
  private readonly TOKEN_KEY = environment.accessTokenKey;
  private readonly ROLE_KEY = environment.roleKey;

  saveToken(token: string) {
    this.secureService.setItem(this.TOKEN_KEY, token);

    // Decode JWT and save role
    const payload = this.decodeToken(token);
    if (payload) {
      const role = this.extractRoleFromScope(payload.scope);
      if (role) {
        this.saveRole(role);
      }
    }
  }

  getToken(): string | null {
    return this.secureService.getItem(this.TOKEN_KEY);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  logout() {
    this.secureService.removeItem(this.TOKEN_KEY);
    this.secureService.removeItem(this.ROLE_KEY);
  }

  /**
   * Decode JWT token to get payload
   */
  decodeToken(token: string): JWTPayload | null {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) {
        return null;
      }

      const payload = parts[1];
      const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
      return JSON.parse(decoded) as JWTPayload;
    } catch (error) {
      console.error('Failed to decode JWT:', error);
      return null;
    }
  }

  /**
   * Extract role from scope string
   * Example: "role.therapist:permission ..." -> "therapist"
   */
  private extractRoleFromScope(scope: string): 'therapist' | 'client' | null {
    const roleMatch = scope.match(/role\.(therapist|client)/);
    if (roleMatch && roleMatch[1]) {
      return roleMatch[1] as 'therapist' | 'client';
    }
    return null;
  }

  /**
   * Save user role to localStorage
   */
  saveRole(role: 'therapist' | 'client'): void {
    this.secureService.setItem(this.ROLE_KEY, role);
  }

  /**
   * Get user role from localStorage
   */
  getRole(): 'therapist' | 'client' | null {
    return this.secureService.getItem(this.ROLE_KEY) as 'therapist' | 'client' | null;
  }

  /**
   * Check if user is therapist
   */
  isTherapist(): boolean {
    return this.getRole() === 'therapist';
  }

  /**
   * Check if user is client
   */
  isClient(): boolean {
    return this.getRole() === 'client';
  }

  /**
   * Get current JWT payload
   */
  getCurrentUser(): JWTPayload | null {
    const token = this.getToken();
    if (!token) {
      return null;
    }
    return this.decodeToken(token);
  }

  /**
   * Check if token is expired
   */
  isTokenExpired(): boolean {
    const payload = this.getCurrentUser();
    if (!payload) {
      return true;
    }
    const now = Math.floor(Date.now() / 1000);
    return payload.exp < now;
  }
}
