import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { Therapist } from '../../models/swipe-card';
export interface SwipeProfile {
  clientId: string;
  therapistId: string;
  points: number;
  reasons: string[];
  status: 'pending' | 'accepted' | 'rejected';
  createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class SwipeService {
  private readonly apiUrl = environment.apiUrl;
  private readonly version = environment.apiVersion;

  constructor(
    private http: HttpClient,
    private secureStorage: SecureStorageService
  ) {}

  triggerUpdate(): Observable<{
    message: string;
    status: number;
    data: SwipeProfile[];
  }> {
    const token = this.secureStorage.getItem('access_token');
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });

    return this.http.post<{
      message: string;
      status: number;
      data: SwipeProfile[];
    }>(
      `${this.apiUrl}/${this.version}/consultation/clients/me/recommend`,
      {},
      { headers }
    );
  }

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
    }>(`${this.apiUrl}/${this.version}/consultation/clients/me/recommend/top`, {
      headers,
    });
  }
}
