import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SocialService {
  private baseUrl = `${environment.apiUrl}/v1/consultation`;

  constructor(private http: HttpClient) {}

  
  followUser(userId: string): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/users/${userId}/follow`, {});
  }

  unfollowUser(userId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/users/${userId}/follow`);
  }

  
  addBookmark(postId: string): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/posts/${postId}/bookmark`, {});
  }

  removeBookmark(postId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/posts/${postId}/bookmark`);
  }

  
  sharePost(postId: string): Observable<void> {
    return this.http.post<void>(`${this.baseUrl}/posts/${postId}/share`, {});
  }
}
