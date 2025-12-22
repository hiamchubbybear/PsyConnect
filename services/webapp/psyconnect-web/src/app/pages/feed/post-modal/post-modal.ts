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

  // Loading state
  isInitialized = false;
  currentUserProfile: any = null;

  // Reply functionality
  replyingToCommentId: string | null = null;
  replyingToUsername: string | null = null;
  replyForm!: FormGroup;

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
    this.initReplyForm();
    this.loadCurrentUserProfile();
    if (this.postId) {
      this.loadInitialData();
    }
  }

  loadCurrentUserProfile() {
    this.profileService.getProfile().subscribe({
      next: (response: any) => {
        this.currentUserProfile = response.data || response;
        if (this.currentUserProfile?.userId) {
          this.profileCache.set(
            this.currentUserProfile.userId,
            this.currentUserProfile
          );
        }
      },
    });
  }

  loadInitialData() {
    this.isInitialized = false;
    this.error = null;

    this.newsfeedService.getPost(this.postId).subscribe({
      next: (post) => {
        this.post = post;
        this.loadProfile(post.author_id).then((profile) => {
          this.authorProfile = profile;
        });
        this.loadComments(this.postId, () => {
          this.isInitialized = true;
        });
      },
      error: (err) => {
        console.error('Failed to load post:', err);
        this.error = 'Failed to load post';
      },
    });
  }

  initCommentForm() {
    this.commentForm = this.fb.group({
      content: [
        '',
        [
          Validators.required,
          Validators.minLength(1),
          Validators.maxLength(500),
        ],
      ],
    });
  }

  initReplyForm() {
    this.replyForm = this.fb.group({
      content: [
        '',
        [
          Validators.required,
          Validators.minLength(1),
          Validators.maxLength(500),
        ],
      ],
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
          const profile = response.data || response;
          this.profileCache.set(userId, profile);
          resolve(profile);
        },
        error: (err) => {
          const fallback = {
            profileId: userId,
            firstName: 'User',
            lastName: '',
          };
          this.profileCache.set(userId, fallback);
          resolve(fallback);
        },
      });
    });
  }

  loadComments(postId: string, done?: () => void) {
    this.commentService.getComments(postId, 20, 0, 2).subscribe({
      next: (response: any) => {
        const comments = Array.isArray(response)
          ? response
          : response.comments || [];
        this.comments = comments;

        // Load profiles for all authors
        const profilePromises = comments.map((c: any) =>
          this.ensureProfileLoadedPromise(c)
        );
        Promise.all(profilePromises).then(() => {
          if (done) done();
        });
      },
      error: (err) => {
        console.error('Failed to load comments:', err);
        if (done) done();
      },
    });
  }

  private ensureProfileLoadedPromise(comment: any): Promise<any> {
    const promises: Promise<any>[] = [];
    if (comment.author_id && !this.profileCache.has(comment.author_id)) {
      promises.push(this.loadProfile(comment.author_id));
    }
    if (comment.replies && comment.replies.length > 0) {
      comment.replies.forEach((reply: any) =>
        promises.push(this.ensureProfileLoadedPromise(reply))
      );
    }
    return Promise.all(promises);
  }

  submitComment() {
    if (this.commentForm.invalid || !this.post) return;
    this.submittingComment = true;
    const content = this.commentForm.value.content;

    this.commentService.createComment(this.post.id, { content }).subscribe({
      next: (comment) => {
        this.comments.unshift(comment);
        this.commentForm.reset();
        this.submittingComment = false;
        if (this.post) this.post.comment_count++;
        this.loadProfile(comment.author_id);
      },
      error: (err) => {
        this.toastService.error('Error', 'Failed to post comment');
        this.submittingComment = false;
      },
    });
  }

  startReply(comment: Comment) {
    this.replyingToCommentId = comment.id;
    const userId = comment.author_id || comment.user_id;
    const profile = this.profileCache.get(userId);
    this.replyingToUsername = profile
      ? `${profile.firstName || ''} ${profile.lastName || ''}`.trim()
      : 'User';
    this.replyForm.reset();
  }

  cancelReply() {
    this.replyingToCommentId = null;
    this.replyingToUsername = null;
    this.replyForm.reset();
  }

  submitReply(parentCommentId: string) {
    if (this.replyForm.invalid || !this.post) return;
    const content = this.replyForm.value.content;

    this.commentService
      .createComment(this.post.id, {
        content,
        parent_comment_id: parentCommentId,
      })
      .subscribe({
        next: (reply) => {
          this.addReplyToComment(this.comments, parentCommentId, reply);
          this.cancelReply();
          if (this.post) this.post.comment_count++;
        },
        error: (err) => {
          this.toastService.error('Error', 'Failed to post reply');
        },
      });
  }

  private addReplyToComment(
    comments: Comment[],
    parentId: string,
    reply: Comment
  ): boolean {
    for (const comment of comments) {
      if (comment.id === parentId) {
        if (!comment.replies) comment.replies = [];
        comment.replies.unshift(reply);
        return true;
      }
      if (
        comment.replies &&
        this.addReplyToComment(comment.replies, parentId, reply)
      )
        return true;
    }
    return false;
  }

  getReplyCharCount(): number {
    return this.replyForm.get('content')?.value?.length || 0;
  }

  // Voting and Engagement
  upvotePost(event: Event) {
    event.stopPropagation();
    if (!this.post) return;

    // If already upvoted, remove the vote
    if (this.post.user_vote === 'up') {
      this.post.upvote_count--;
      this.post.user_vote = null;

      this.reactionService.removeReaction(this.post.id).subscribe({
        error: (err) => console.error('Failed to remove upvote:', err),
      });
      return;
    }

    // Otherwise, add/switch to upvote
    if (this.post.user_vote === 'down') this.post.downvote_count--;
    this.post.upvote_count++;
    this.post.user_vote = 'up';

    this.reactionService.addReaction(this.post.id, 'up').subscribe({
      error: (err) => console.error('Failed to upvote:', err),
    });
  }

  downvotePost(event: Event) {
    event.stopPropagation();
    if (!this.post) return;

    // If already downvoted, remove the vote
    if (this.post.user_vote === 'down') {
      this.post.downvote_count--;
      this.post.user_vote = null;

      this.reactionService.removeReaction(this.post.id).subscribe({
        error: (err) => console.error('Failed to remove downvote:', err),
      });
      return;
    }

    // Otherwise, add/switch to downvote
    if (this.post.user_vote === 'up') this.post.upvote_count--;
    this.post.downvote_count++;
    this.post.user_vote = 'down';

    this.reactionService.addReaction(this.post.id, 'down').subscribe({
      error: (err) => console.error('Failed to downvote:', err),
    });
  }

  toggleBookmark() {
    if (!this.post) return;
    const isBookmarked = this.post.user_bookmark;

    const action = isBookmarked
      ? this.socialService.removeBookmark(this.post.id)
      : this.socialService.addBookmark(this.post.id);

    action.subscribe({
      next: () => {
        if (this.post) this.post.user_bookmark = !isBookmarked;
        this.toastService.success(
          'Success',
          isBookmarked ? 'Bookmark removed' : 'Post bookmarked'
        );
      },
      error: (err) => console.error('Failed to toggle bookmark:', err),
    });
  }

  sharePost() {
    if (!this.post) return;
    this.socialService.sharePost(this.post.id).subscribe({
      next: () => {
        if (this.post) this.post.share_count++;
        this.toastService.success('Success', 'Post shared');
      },
      error: (err) => console.error('Failed to share:', err),
    });
  }

  // Utils
  getUserName(userId: string): string {
    const profile = this.profileCache.get(userId);
    if (!profile) return 'User';
    return (
      `${profile.firstName || ''} ${profile.lastName || ''}`.trim() ||
      profile.username ||
      'User'
    );
  }

  getAuthorName(): string {
    return this.post ? this.getUserName(this.post.author_id) : 'Author';
  }

  showAuthorCardWithDelay(userId: string) {
    this.hoveredUserId = userId;
    this.authorCardTimeout = setTimeout(() => {
      if (this.hoveredUserId === userId) {
        this.loadProfile(userId).then((p) => {
          this.authorProfile = p;
          this.showAuthorCard = true;
        });
      }
    }, 500);
  }

  hideAuthorCard() {
    clearTimeout(this.authorCardTimeout);
    this.showAuthorCard = false;
    this.hoveredUserId = null;
  }

  formatTime(timestamp: string): string {
    const date = new Date(timestamp);
    const diff = new Date().getTime() - date.getTime();
    const mins = Math.floor(diff / 60000);
    const hrs = Math.floor(mins / 60);
    const days = Math.floor(hrs / 24);

    if (days > 0) return `${days}d ago`;
    if (hrs > 0) return `${hrs}h ago`;
    if (mins > 0) return `${mins}m ago`;
    return 'Just now';
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
