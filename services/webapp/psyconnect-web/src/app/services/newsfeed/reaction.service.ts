import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Reaction, ReactionType } from '../../models/reaction.model';

interface VoteState {
  postId: string;
  voteType: ReactionType | null; 
  upvoteCount: number;
  downvoteCount: number;
}

@Injectable({
  providedIn: 'root',
})
export class ReactionService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/posts`;

  
  private userVotes = new Map<string, ReactionType | null>();
  private voteStates = new BehaviorSubject<Map<string, VoteState>>(new Map());

  constructor(private http: HttpClient) {
    this.loadVotesFromStorage();
  }

  
  toggleVote(
    postId: string,
    voteType: ReactionType
  ): Observable<Reaction | void> {
    const currentVote = this.userVotes.get(postId);

    
    if (currentVote === voteType) {
      return this.removeReaction(postId).pipe(
        tap(() => {
          this.userVotes.set(postId, null);
          this.saveVotesToStorage();
        })
      );
    }

    
    return this.addReaction(postId, voteType).pipe(
      tap(() => {
        this.userVotes.set(postId, voteType);
        this.saveVotesToStorage();
      })
    );
  }

  
  getUserVote(postId: string): ReactionType | null {
    return this.userVotes.get(postId) || null;
  }

  
  hasVoted(postId: string): boolean {
    return this.userVotes.has(postId) && this.userVotes.get(postId) !== null;
  }

  
  hasUpvoted(postId: string): boolean {
    return this.userVotes.get(postId) === 'up';
  }

  
  hasDownvoted(postId: string): boolean {
    return this.userVotes.get(postId) === 'down';
  }

  addReaction(
    postId: string,
    reactionType: ReactionType
  ): Observable<Reaction> {
    return this.http.post<Reaction>(`${this.apiUrl}/${postId}/react`, {
      reaction_type: reactionType,
    });
  }

  removeReaction(postId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${postId}/react`);
  }

  getReactions(postId: string): Observable<Reaction[]> {
    return this.http.get<Reaction[]>(`${this.apiUrl}/${postId}/reactions`);
  }

  
  private saveVotesToStorage(): void {
    const votesObj = Object.fromEntries(this.userVotes);
    localStorage.setItem('user_votes', JSON.stringify(votesObj));
  }

  
  private loadVotesFromStorage(): void {
    const stored = localStorage.getItem('user_votes');
    if (stored) {
      try {
        const votesObj = JSON.parse(stored);
        this.userVotes = new Map(Object.entries(votesObj));
      } catch (e) {
        console.error('Failed to load votes from storage', e);
      }
    }
  }

  
  clearVotes(): void {
    this.userVotes.clear();
    localStorage.removeItem('user_votes');
  }
}
