import { Injectable } from '@angular/core';
import { SecureStorageService } from '../../encrypt/secure';

@Injectable({ providedIn: 'root' })
export class AuthService {
    constructor(private secureService : SecureStorageService) {

    }
  private readonly TOKEN_KEY = 'access_token';

  saveToken(token: string) {
    this.secureService.setItem(this.TOKEN_KEY, token);
  }

  getToken(): string | null {
    return this.secureService.getItem(this.TOKEN_KEY);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  logout() {
    this.secureService.removeItem(this.TOKEN_KEY);
  }
}
