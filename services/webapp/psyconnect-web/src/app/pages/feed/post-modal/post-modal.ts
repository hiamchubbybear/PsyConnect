import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Comment } from '../../../models/comment.model';
import { Post } from '../../../models/post.model';
import { REACTION_TYPES, ReactionType } from '../../../models/reaction.model';
import { SocialService } from '../../../services/consultation/social.service';
import { LoaderService } from '../../../services/loader/loader';
import { CommentService } from '../../../services/newsfeed/comment.service';
import { NewsfeedService } from '../../../services/newsfeed/newsfeed.service';
import { ReactionService } from '../../../services/newsfeed/reaction.service';
import { Profile } from '../../../services/profile/profile';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-post-modal',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, ReactiveFormsModule],
  templateUrl: './post-modal.html',
  styleUrls: ['./post-modal.scss'],
})
export class PostModalComponent implements OnInit {
  @Input() postId!: string;
  @Output() closeModal = new EventEmitter<void>();

  post: Post | null = null;
  comments: Comment[] = [];
  error: string | null = null;
  commentForm!: FormGroup;
  submittingComment = false;

  showReactionPicker = false;
  reactionTypes = REACTION_TYPES;
  isBookmarked = false;

  profileCache: Map<string, any> = new Map();
  authorProfile: any = null;
  showAuthorCard = false;
  authorCardTimeout: any = null;
  hoveredUserId: string | null = null;

  constructor(
    private newsfeedService: NewsfeedService,
    private commentService: CommentService,
    private reactionService: ReactionService,
    private socialService: SocialService,
    private profileService: Profile,
    private toastService: ToastService,
    private fb: FormBuilder,
    private loaderService: LoaderService
  ) {}

  ngOnInit() {
    this.initCommentForm();
    if (this.postId) {
      this.loadPost(this.postId);
      this.loadComments(this.postId);
    }
  }

  initCommentForm() {
    this.commentForm = this.fb.group({
      content: ['', [Validators.required, Validators.minLength(1)]],
    });
  }

  loadPost(id: string) {
    this.loaderService.show();
    this.error = null;

    this.newsfeedService.getPost(id).subscribe({
      next: (post) => {
        this.post = post;
        this.loaderService.hide();
        if (post.author_id) {
          this.loadAuthorProfile(post.author_id);
        }
      },
      error: (err) => {
        console.error('Failed to load post:', err);
        this.error = 'Failed to load post';
        this.loaderService.hide();
      },
    });
  }

  loadAuthorProfile(authorId: string) {
    this.loadProfile(authorId).then((profile) => {
      this.authorProfile = profile;
    });
  }

  loadProfile(userId: string): Promise<any> {
    return new Promise((resolve) => {
      if (this.profileCache.has(userId)) {
        resolve(this.profileCache.get(userId));
        return;
      }

      this.profileService.getProfileById(userId).subscribe({
        next: (response: any) => {
          const profile = response.data || {
            profileId: userId,
            firstName: 'Unknown',
            lastName: 'User',
            username: userId,
          };
          this.profileCache.set(userId, profile);
          resolve(profile);
        },
        error: (err) => {
          console.error('Failed to load profile:', err);
          const fallback = {
            profileId: userId,
            firstName: 'Unknown',
            lastName: 'User',
            username: userId,
          };
          this.profileCache.set(userId, fallback);
          resolve(fallback);
        },
      });
    });
  }

  loadComments(postId: string) {
    this.commentService.getComments(postId).subscribe({
      next: (comments) => {
        this.comments = comments;
        comments.forEach((comment) => {
          if (comment.author_id) {
            this.loadProfile(comment.author_id);
          }
        });
      },
      error: (err) => {
        console.error('Failed to load comments:', err);
      },
    });
  }

  submitComment() {
    if (this.commentForm.invalid || !this.post) return;

    this.submittingComment = true;
    const content = this.commentForm.value.content;

    this.commentService.createComment(this.post.id, content).subscribe({
      next: (comment) => {
        this.comments.unshift(comment);
        this.commentForm.reset();
        this.submittingComment = false;
        this.toastService.success('Success', 'Comment posted successfully');
        if (this.post) {
          this.post.comment_count++;
        }
        if (comment.author_id) {
          this.loadProfile(comment.author_id);
        }
      },
      error: (err) => {
        console.error('Failed to post comment:', err);
        this.toastService.error('Error', 'Failed to post comment');
        this.submittingComment = false;
      },
    });
  }

  toggleReactionPicker(event: Event) {
    event.stopPropagation();
    this.showReactionPicker = !this.showReactionPicker;
  }

  addReaction(type: ReactionType, event: Event) {
    event.stopPropagation();
    if (!this.post) return;

    this.reactionService.addReaction(this.post.id, type).subscribe({
      next: () => {
        this.showReactionPicker = false;
        this.toastService.success('Success', 'Reaction added');
      },
      error: (err) => {
        console.error('Failed to add reaction:', err);
        this.toastService.error('Error', 'Failed to add reaction');
      },
    });
  }

  toggleBookmark() {
    if (!this.post) return;

    const action = this.isBookmarked
      ? this.socialService.removeBookmark(this.post.id)
      : this.socialService.addBookmark(this.post.id);

    action.subscribe({
      next: () => {
        this.isBookmarked = !this.isBookmarked;
        this.toastService.success(
          'Success',
          this.isBookmarked ? 'Post bookmarked' : 'Bookmark removed'
        );
      },
      error: (err) => {
        console.error('Failed to toggle bookmark:', err);
        this.toastService.error('Error', 'Failed to update bookmark');
      },
    });
  }

  sharePost() {
    if (!this.post) return;

    this.socialService.sharePost(this.post.id).subscribe({
      next: () => {
        this.toastService.success('Success', 'Post shared successfully');
        if (this.post) {
          this.post.share_count++;
        }
      },
      error: (err) => {
        console.error('Failed to share post:', err);
        this.toastService.error('Error', 'Failed to share post');
      },
    });
  }

  getReactionKeys(): ReactionType[] {
    return Object.keys(this.reactionTypes) as ReactionType[];
  }

  getReactionEmoji(type: ReactionType): string {
    return this.reactionTypes[type];
  }

  getUserName(userId: string): string {
    const profile = this.profileCache.get(userId);
    if (!profile) return 'Unknown User';

    const firstName = profile.firstName || '';
    const lastName = profile.lastName || '';
    const fullName = `${firstName} ${lastName}`.trim();

    return (
      fullName ||
      profile.username ||
      profile.display_name ||
      profile.profileId ||
      'Unknown User'
    );
  }

  getAuthorName(): string {
    if (!this.post) return 'Unknown';
    return this.getUserName(this.post.author_id);
  }

  showAuthorCardWithDelay(userId: string) {
    this.hoveredUserId = userId;
    this.authorCardTimeout = setTimeout(() => {
      if (this.hoveredUserId === userId) {
        this.loadProfile(userId).then((profile) => {
          this.authorProfile = profile;
          this.showAuthorCard = this.hoveredUserId === userId;
        });
      }
    }, 500);
  }

  hideAuthorCard() {
    if (this.authorCardTimeout) {
      clearTimeout(this.authorCardTimeout);
    }
    this.showAuthorCard = false;
    this.hoveredUserId = null;
  }

  formatTime(timestamp: string): string {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 7) {
      return date.toLocaleDateString();
    } else if (days > 0) {
      return `${days}d ago`;
    } else if (hours > 0) {
      return `${hours}h ago`;
    } else if (minutes > 0) {
      return `${minutes}m ago`;
    } else {
      return 'Just now';
    }
  }

  close() {
    this.closeModal.emit();
  }

  onOverlayClick(event: MouseEvent) {
    if (
      (event.target as HTMLElement).classList.contains('post-modal-overlay')
    ) {
      this.close();
    }
  }
}
