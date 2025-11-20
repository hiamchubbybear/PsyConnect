import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Post } from '../../models/post.model';

@Injectable({
  providedIn: 'root',
})
export class TagsCategoriesService {
  private readonly baseUrl = `${environment.apiUrl}/${environment.apiVersion}/consultation`;

  constructor(private http: HttpClient) {}

  getPostsByTag(tag: string, limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.baseUrl}/tags/${tag}/posts`, { params });
  }

  getPostsByCategory(category: string, limit: number = 20, skip: number = 0): Observable<Post[]> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString());
    return this.http.get<Post[]>(`${this.baseUrl}/categories/${category}/posts`, { params });
  }
}
