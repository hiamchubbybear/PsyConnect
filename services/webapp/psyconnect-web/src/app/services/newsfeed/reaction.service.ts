import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Reaction, ReactionType } from '../../models/reaction.model';

@Injectable({
  providedIn: 'root'
})
export class ReactionService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/posts`;

  constructor(private http: HttpClient) {}

  // Add reaction to post
  addReaction(postId: string, reactionType: ReactionType): Observable<Reaction> {
    return this.http.post<Reaction>(`${this.apiUrl}/${postId}/react`, {
      reaction_type: reactionType
    });
  }

  // Remove reaction from post
  removeReaction(postId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${postId}/react`);
  }

  // Get all reactions for a post
  getReactions(postId: string): Observable<Reaction[]> {
    return this.http.get<Reaction[]>(`${this.apiUrl}/${postId}/reactions`);
  }
}
