import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { Therapist } from '../../models/swipe-card';

@Injectable({ providedIn: 'root' })
export class SwipeService {
  private readonly apiUrl = environment.apiUrl;

  constructor(
    private http: HttpClient,
    private secureStorage: SecureStorageService
  ) {}

  getSwipeData(): Observable<{
    message: string;
    status: number;
    data: Therapist[];
  }> {
    const token = this.secureStorage.getItem('access_token');
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });

    return this.http.get<{
      message: string;
      status: number;
      data: Therapist[];
    }>(`${this.apiUrl}/consultation/client/recommend/top`, { headers });
  }
}
