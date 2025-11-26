import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Reaction, ReactionType } from '../../models/reaction.model';

interface VoteState {
  postId: string;
  voteType: ReactionType | null; // null means no vote
  upvoteCount: number;
  downvoteCount: number;
}

@Injectable({
  providedIn: 'root',
})
export class ReactionService {
  private apiUrl = `${environment.apiUrl}/v1/consultation/posts`;

  // Track user votes locally to prevent multiple votes
  private userVotes = new Map<string, ReactionType | null>();
  private voteStates = new BehaviorSubject<Map<string, VoteState>>(new Map());

  constructor(private http: HttpClient) {
    this.loadVotesFromStorage();
  }

  /**
   * Toggle vote - if same type clicked, remove vote. If different type, change vote.
   */
  toggleVote(postId: string, voteType: ReactionType): Observable<Reaction | void> {
    const currentVote = this.userVotes.get(postId);

    // If clicking same vote type, remove it
    if (currentVote === voteType) {
      return this.removeReaction(postId).pipe(
        tap(() => {
          this.userVotes.set(postId, null);
          this.saveVotesToStorage();
        })
      );
    }

    // Otherwise, add/change vote
    return this.addReaction(postId, voteType).pipe(
      tap(() => {
        this.userVotes.set(postId, voteType);
        this.saveVotesToStorage();
      })
    );
  }

  /**
   * Get current user's vote for a post
   */
  getUserVote(postId: string): ReactionType | null {
    return this.userVotes.get(postId) || null;
  }

  /**
   * Check if user has voted on a post
   */
  hasVoted(postId: string): boolean {
    return this.userVotes.has(postId) && this.userVotes.get(postId) !== null;
  }

  /**
   * Check if user upvoted
   */
  hasUpvoted(postId: string): boolean {
    return this.userVotes.get(postId) === 'upvote';
  }

  /**
   * Check if user downvoted
   */
  hasDownvoted(postId: string): boolean {
    return this.userVotes.get(postId) === 'downvote';
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

  /**
   * Save votes to localStorage for persistence
   */
  private saveVotesToStorage(): void {
    const votesObj = Object.fromEntries(this.userVotes);
    localStorage.setItem('user_votes', JSON.stringify(votesObj));
  }

  /**
   * Load votes from localStorage
   */
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

  /**
   * Clear all votes (for logout)
   */
  clearVotes(): void {
    this.userVotes.clear();
    localStorage.removeItem('user_votes');
  }
}
