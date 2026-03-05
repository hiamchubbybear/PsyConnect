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

  
  decodeToken(token: string): JWTPayload | null {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) {
        return null;
      }

      let payload = parts[1].replace(/-/g, '+').replace(/_/g, '/');
      const pad = payload.length % 4 === 0 ? '' : '='.repeat(4 - (payload.length % 4));
      const decoded = atob(payload + pad);
      return JSON.parse(decoded) as JWTPayload;
    } catch (error) {
      console.error('Failed to decode JWT:', error);
      return null;
    }
  }

  
  private extractRoleFromScope(scope: string): 'therapist' | 'client' | null {
    if (!scope) return null;
    const s = scope.toLowerCase();
    if (s.includes('therapist')) return 'therapist';
    if (s.includes('client') || s.includes('user')) return 'client';
    return null;
  }

  
  saveRole(role: 'therapist' | 'client'): void {
    this.secureService.setItem(this.ROLE_KEY, role);
  }

  getRole(): 'therapist' | 'client' | null {
    const role = this.secureService.getItem(this.ROLE_KEY) as string | null;
    if (role) {
      const extracted = this.extractRoleFromScope(role);
      if (extracted) return extracted;
    }

    
    const payload = this.getCurrentUser();
    if (payload && payload.scope) {
      const extractedRole = this.extractRoleFromScope(payload.scope);
      if (extractedRole) {
        this.saveRole(extractedRole);
        return extractedRole;
      }
    }
    return null;
  }

  
  isTherapist(): boolean {
    return this.getRole() === 'therapist';
  }

  
  isClient(): boolean {
    return this.getRole() === 'client';
  }

  getCurrentUser(): JWTPayload | null {
    const token = this.getToken();
    if (!token) {
      return null;
    }
    return this.decodeToken(token);
  }

  
  isTokenExpired(): boolean {
    const payload = this.getCurrentUser();
    if (!payload) {
      return true;
    }
    const now = Math.floor(Date.now() / 1000);
    return payload.exp < now;
  }
}
