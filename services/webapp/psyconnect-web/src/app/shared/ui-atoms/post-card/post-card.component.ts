import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { LucideAngularModule } from 'lucide-angular'; // Added this import
import { Post } from '../../../models/post.model';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { PsyButtonComponent } from '../button/psy-button.component';

@Component({
  selector: 'app-post-card',
  standalone: true,
  imports: [
    CommonModule,
    LucideAngularModule,
    PsyButtonComponent,
    AvatarFallbackPipe,
    TranslateModule,
  ],
  templateUrl: './post-card.component.html',
  styleUrls: ['./post-card.component.scss'],
})
export class PostCardComponent {
  @Input() post!: Post;
  @Input() currentUserId?: string;

  @Output() upvote = new EventEmitter<Event>();
  @Output() downvote = new EventEmitter<Event>();
  @Output() share = new EventEmitter<Event>();
  @Output() bookmark = new EventEmitter<Event>();
  @Output() openPost = new EventEmitter<Event>();
  @Output() openProfile = new EventEmitter<string>();

  // Helper for reaction keys if needed, but logic should be simple

  formatTime(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
  }

  onUpvote(event: Event) {
    event.stopPropagation();
    this.upvote.emit(event);
  }

  onDownvote(event: Event) {
    event.stopPropagation();
    this.downvote.emit(event);
  }

  onShare(event: Event) {
    event.stopPropagation();
    this.share.emit(event);
  }

  onBookmark(event: Event) {
    event.stopPropagation();
    this.bookmark.emit(event);
  }
}
