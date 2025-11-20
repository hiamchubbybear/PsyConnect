import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class SocialService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation`;

  constructor(private http: HttpClient) {}

  // Follow/Unfollow
  followUser(userId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/users/${userId}/follow`, {});
  }

  unfollowUser(userId: string): Observable<any> {
    return this.http.delete(`${this.baseUrl}/users/${userId}/follow`);
  }

  getFollowers(userId: string, limit: number = 20, skip: number = 0): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/users/${userId}/followers`, {
      params: { limit: limit.toString(), skip: skip.toString() },
    });
  }

  getFollowing(userId: string, limit: number = 20, skip: number = 0): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/users/${userId}/following`, {
      params: { limit: limit.toString(), skip: skip.toString() },
    });
  }

  // Bookmarks
  addBookmark(postId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/posts/${postId}/bookmark`, {});
  }

  removeBookmark(postId: string): Observable<any> {
    return this.http.delete(`${this.baseUrl}/posts/${postId}/bookmark`);
  }

  getBookmarks(limit: number = 20, skip: number = 0): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/bookmarks`, {
      params: { limit: limit.toString(), skip: skip.toString() },
    });
  }

  // Share
  sharePost(postId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/posts/${postId}/share`, {});
  }
}
