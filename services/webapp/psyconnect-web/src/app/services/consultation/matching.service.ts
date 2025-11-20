import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { MatchRequest, MatchResponse } from '../../models/consultation.model';

@Injectable({
  providedIn: 'root',
})
export class MatchingService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation/therapist/match`;

  constructor(private http: HttpClient) {}

  getAllMatches(): Observable<MatchResponse[]> {
    return this.http.get<MatchResponse[]>(this.baseUrl);
  }

  requestMatch(data: MatchRequest): Observable<MatchResponse> {
    return this.http.post<MatchResponse>(this.baseUrl, data);
  }
}
