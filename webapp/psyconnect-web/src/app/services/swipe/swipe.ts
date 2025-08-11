import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';

class Swipe {
  constructor(
    private http: HttpClient,
    private securityStorage: SecureStorageService
  ) {}
  private readonly apiUrl = environment.apiUrl;
  getSwipeData(): Observable<any> {
    const token = this.securityStorage.getItem('access_token');

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
