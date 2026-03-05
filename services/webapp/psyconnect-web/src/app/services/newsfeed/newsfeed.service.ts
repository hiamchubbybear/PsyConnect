import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { CreatePostRequest, Post } from '../../models/post.model';

@Injectable({
  providedIn: 'root',
})
export class NewsfeedService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/posts`;

  constructor(private http: HttpClient) {}

  
  getFeed(limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(this.apiUrl, { params });
  }

  
  getPost(id: string): Observable<Post> {
    return this.http.get<Post>(`${this.apiUrl}/${id}`);
  }

  
  createPost(post: CreatePostRequest): Observable<Post> {
    return this.http.post<Post>(this.apiUrl, post);
  }

  
  updatePost(id: string, post: Partial<CreatePostRequest>): Observable<Post> {
    return this.http.put<Post>(`${this.apiUrl}/${id}`, post);
  }

  
  deletePost(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  
  getUserPosts(
    userId: string,
    limit: number = 20,
    skip: number = 0,
  ): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/user/${userId}`, { params });
  }

  
  getTrendingPosts(limit: number = 20): Observable<Post[]> {
    const params = new HttpParams().set('limit', limit.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/trending`, { params });
  }

  
  searchPosts(
    query: string,
    limit: number = 20,
    skip: number = 0,
  ): Observable<Post[]> {
    const params = new HttpParams()
      .set('q', query)
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/search`, { params });
  }

  
  recordView(id: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/${id}/view`, {});
  }

  
  getPopularTags(limit: number = 10): Observable<string[]> {
    const params = new HttpParams().set('limit', limit.toString());
    return this.http.get<string[]>(
      `${this.apiUrl.replace('/posts', '/tags')}/popular`,
      { params },
    );
  }
}
