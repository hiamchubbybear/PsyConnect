import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';

class Swipe {
  constructor(
    private http: HttpClient,
    private securityStorage: SecureStorageService
  ) {}
  ACCESSTOKEN_KEY = environment.accessTokenKey;
  private readonly apiUrl = environment.apiUrl;
  getSwipeData(): Observable<any> {
    const token = this.securityStorage.getItem(this.ACCESSTOKEN_KEY);
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });
    return this.http.post<any>(
      `${this.apiUrl}/consultation/client/recommend/top`,
      {
        headers,
      }
    );
  }
}
