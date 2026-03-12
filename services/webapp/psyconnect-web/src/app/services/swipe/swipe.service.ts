import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
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
interface TherapistRecommendResponse {
  message: string;
  status: number;
  data: {
    client_id: string;
    therapist_id: string;
    points: number;
    reasons: string[];
    status: string;
    created_at: string;
  }[];
}

@Injectable({ providedIn: 'root' })
export class SwipeService {
  private readonly apiUrl = environment.apiUrl;
  private readonly version = environment.apiVersion;
  PROFILE_KEY = environment.profileKey;
  USERNAME_KEY = environment.usernameKey;
  ACCESSTOKEN_KEY = environment.accessTokenKey;
  THERAPIST_KEY = environment.therapistsKey;
  constructor(
    private http: HttpClient,
    private secureStorage: SecureStorageService,
  ) {}

  triggerUpdate(): Observable<{
    message: string;
    status: number;
    data: SwipeProfile[];
  }> {
    const token = this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });

    return this.http.post<{
      message: string;
      status: number;
      data: SwipeProfile[];
    }>(
      `${this.apiUrl}/${this.version}/consultation/clients/me/recommend`,
      {},
      { headers },
    );
  }

  getSwipeData(): Observable<{
    message: string;
    status: number;
    data: Therapist[];
  }> {
    const token = this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });

    return this.http.get<{
      message: string;
      status: number;
      data: Therapist[];
    }>(`${this.apiUrl}/${this.version}/consultation/clients/me/recommend/top`, {
      headers,
    });
  }
  getUpdateTherapistHandler(): Observable<boolean> {
    const token = this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
    const url = `${this.apiUrl}/${this.version}/consultation/clients/me/recommend`;
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<TherapistRecommendResponse>(url, { headers }).pipe(
      map((res) => {
        return Array.isArray(res.data) && res.data.length > 0;
      }),
    );
  }
  postSwipe(therapistId: string, status: string = 'swiped'): Observable<any> {
    const token = this.secureStorage.getItem(this.ACCESSTOKEN_KEY);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    });

    return this.http.post(
      `${this.apiUrl}/${this.version}/consultation/clients/me/swipe`,
      {
        therapist_id: therapistId,
        status: status,
      },
      { headers },
    );
  }
}
