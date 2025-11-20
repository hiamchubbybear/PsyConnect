import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { CreatePostRequest, Post } from '../../models/post.model';

@Injectable({
  providedIn: 'root'
})
export class NewsfeedService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/posts`;

  constructor(private http: HttpClient) {}

  // Get personalized feed
  getFeed(limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(this.apiUrl, { params });
  }

  // Get single post
  getPost(id: string): Observable<Post> {
    return this.http.get<Post>(`${this.apiUrl}/${id}`);
  }

  // Create new post
  createPost(post: CreatePostRequest): Observable<Post> {
    return this.http.post<Post>(this.apiUrl, post);
  }

  // Update post
  updatePost(id: string, post: Partial<CreatePostRequest>): Observable<Post> {
    return this.http.put<Post>(`${this.apiUrl}/${id}`, post);
  }

  // Delete post
  deletePost(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  // Get user's posts
  getUserPosts(userId: string, limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/user/${userId}`, { params });
  }

  // Get trending posts
  getTrendingPosts(limit: number = 20): Observable<Post[]> {
    const params = new HttpParams().set('limit', limit.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/trending`, { params });
  }

  // Search posts
  searchPosts(query: string, limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('q', query)
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.apiUrl}/search`, { params });
  }

  // Get posts by tag
  getPostsByTag(tag: string, limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${environment.apiUrl}/v1/consultation/tags/${tag}/posts`, { params });
  }
}
