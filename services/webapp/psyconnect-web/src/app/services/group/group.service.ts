import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface Group {
  id: string;
  name: string;
  description: string;
  category: string;
  icon?: string;
  creator_id: string;
  member_count: number;
  is_public: boolean;
  created_at: string;
}

@Injectable({
  providedIn: 'root',
})
export class GroupService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/groups`;

  constructor(private http: HttpClient) {}

  getGroups(
    category?: string,
    limit: number = 20,
    skip: number = 0,
    query?: string,
  ): Observable<Group[]> {
    let params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());

    if (category) {
      params = params.set('category', category);
    }

    if (query) {
      params = params.set('q', query);
    }

    return this.http.get<Group[]>(this.apiUrl, { params });
  }

  getGroupById(id: string): Observable<Group> {
    return this.http.get<Group>(`${this.apiUrl}/${id}`);
  }

  createGroup(group: Partial<Group>): Observable<Group> {
    return this.http.post<Group>(this.apiUrl, group);
  }

  joinGroup(id: string): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.apiUrl}/${id}/join`, {});
  }
}
