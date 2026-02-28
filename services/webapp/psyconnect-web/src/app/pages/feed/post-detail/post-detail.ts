import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Comment } from '../../../models/comment.model';
import { Post } from '../../../models/post.model';
import {
  REACTION_IMAGES,
  REACTION_TYPES,
  ReactionType,
} from '../../../models/reaction.model';
import { SocialService } from '../../../services/consultation/social.service';
import { LoaderService } from '../../../services/loader/loader';
import { CommentService } from '../../../services/newsfeed/comment.service';
import { NewsfeedService } from '../../../services/newsfeed/newsfeed.service';
import { ReactionService } from '../../../services/newsfeed/reaction.service';
import { Profile } from '../../../services/profile/profile';
import { ToastService } from '../../../shared/toast/toast.service';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'app-post-detail',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, ReactiveFormsModule, AvatarFallbackPipe],
  templateUrl: './post-detail.html',
  styleUrls: ['./post-detail.scss'],
})
export class PostDetailComponent implements OnInit {
  post: Post | null = null;
  comments: Comment[] = [];
  // loading = false; // Removed manual loading state
  error: string | null = null;
  commentForm!: FormGroup;
  submittingComment = false;

  // Reply functionality
  replyingToCommentId: string | null = null;
  replyingToUsername: string | null = null;
  replyForm!: FormGroup;

  showReactionPicker = false;
  reactionTypes = REACTION_TYPES;
  isBookmarked = false;

  // Author profiles cache
  profileCache: Map<string, any> = new Map();
  authorProfile: any = null;
  showAuthorCard = false;
  authorCardTimeout: any = null;
  hoveredUserId: string | null = null;

  currentUserProfile: any = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
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
    const postId = this.route.snapshot.paramMap.get('id');
    if (postId) {
      this.loadPost(postId);
      this.loadComments(postId);
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
      error: (err) =>
        console.error('Failed to load current user profile:', err),
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

  loadPost(id: string) {
    this.loaderService.show(); // Show global loader
    this.error = null;

    this.newsfeedService.getPost(id).subscribe({
      next: (post) => {
        this.post = post;
        this.loaderService.hide(); // Hide global loader
        // Load author profile
        if (post.author_id) {
          this.loadAuthorProfile(post.author_id);
        }
      },
      error: (err) => {
        console.error('Failed to load post:', err);
        this.error = 'Failed to load post';
        this.loaderService.hide(); // Hide global loader
      },
    });
  }

  loadAuthorProfile(authorId: string) {
    this.loadProfile(authorId).then((profile) => {
      this.authorProfile = profile;
    });
  }

  // Generic profile loading with caching
  private ensureProfileLoaded(comment: any) {
    if (comment.author_id && !this.profileCache.has(comment.author_id)) {
      this.loadProfile(comment.author_id);
    }
    if (comment.replies && comment.replies.length > 0) {
      comment.replies.forEach((reply: any) => this.ensureProfileLoaded(reply));
    }
  }

  loadProfile(userId: string): Promise<any> {
    console.log('📥 loadProfile called for:', userId);
    return new Promise((resolve) => {
      // Check cache first
      if (this.profileCache.has(userId)) {
        console.log(
          '💾 Profile found in cache:',
          this.profileCache.get(userId)
        );
        resolve(this.profileCache.get(userId));
        return;
      }

      // Fetch from API
      console.log('🌐 Fetching profile from API...');
      this.profileService.getProfileById(userId).subscribe({
        next: (response: any) => {
          console.log('✅ Profile fetched successfully:', response);
          const profileData = response.data || response;
          this.profileCache.set(userId, profileData);
          resolve(profileData);
        },
        error: (err: any) => {
          console.error('❌ Failed to load profile:', err);
          // Cache a fallback profile
          const fallback = {
            firstName: 'Unknown',
            lastName: 'User',
            profileId: userId,
            username: userId,
          };
          this.profileCache.set(userId, fallback);
          resolve(fallback);
        },
      });
    });
  }

  loadComments(postId: string) {
    this.commentService.getComments(postId, 20, 0, 2).subscribe({
      next: (response: any) => {
        // Handle both response formats: array or {comments, total}
        if (Array.isArray(response)) {
          // Old format: direct array
          this.comments = response;
        } else if (response && response.comments) {
          // New format: {comments: [], total: N}
          this.comments = response.comments;
        } else {
          this.comments = [];
        }

        // Add mock comments for demo if needed
        if (this.comments.length === 0) {
          const now = new Date().toISOString();
          this.comments = [
            {
              id: 'mock-1',
              post_id: postId,
              user_id: 'user-1',
              author_id: 'user-1',
              content:
                'I think for our second campaign we can try to target a different audience. How does it sound for you?',
              like_count: 2,
              is_deleted: false,
              created_at: now,
              updated_at: now,
              isLiked: true,
              replies: [
                {
                  id: 'mock-1-1',
                  post_id: postId,
                  user_id: 'user-2',
                  author_id: 'user-2',
                  content:
                    'Yes, that sounds good! I can think about this tomorrow. When do we plan to start that campaign?',
                  like_count: 3,
                  is_deleted: false,
                  created_at: now,
                  updated_at: now,
                  isLiked: true,
                  parent_comment_id: 'mock-1',
                },
              ],
            },
          ];
        }

        // Load profiles for all comment authors
        this.comments.forEach((comment) => {
          this.ensureProfileLoaded(comment);
        });
      },
      error: (err: any) => {
        console.error('Failed to load comments:', err);
        this.comments = [];
      },
    });
  }

  submitComment() {
    if (this.commentForm.invalid || !this.post) return;

    this.submittingComment = true;
    const content = this.commentForm.value.content;

    // DEBUG: Log what we're sending
    console.log('=== SUBMIT COMMENT DEBUG ===');
    console.log('Form value:', this.commentForm.value);
    console.log('Content:', content);
    console.log('Payload:', { content });
    console.log('JSON:', JSON.stringify({ content }));
    console.log('============================');

    this.commentService.createComment(this.post.id, { content }).subscribe({
      next: (comment) => {
        if (!this.comments) {
          this.comments = [];
        }
        this.comments.unshift(comment);

        // Load profile for the new comment (current user)
        if (comment.author_id) {
          this.loadProfile(comment.author_id);
        }

        this.commentForm.reset();
        this.submittingComment = false;
        if (this.post) {
          this.post.comment_count++;
        }
      },
      error: (err: any) => {
        console.error('Failed to submit comment:', err);
        console.error('Error status:', err.status);
        console.error('Error message:', err.message);
        console.error('Error body:', err.error);
        console.error('Full error object:', JSON.stringify(err, null, 2));

        // Show user-friendly error
        const errorMsg =
          err.error?.error || err.message || 'Failed to submit comment';
        this.toastService.error('Error', errorMsg);
        this.submittingComment = false;
      },
    });
  }

  // Reply functionality methods
  startReply(comment: Comment) {
    console.log('🔵 START REPLY CALLED!', comment);
    this.replyingToCommentId = comment.id;
    // Get author name from profile cache
    const profile = this.profileCache.get(comment.user_id);
    this.replyingToUsername = profile
      ? `${profile.first_name} ${profile.last_name}`
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
          // Add reply to parent's replies array
          this.addReplyToComment(this.comments, parentCommentId, reply);

          // Load profile for reply author
          if (reply.user_id) {
            this.loadProfile(reply.user_id);
          }

          this.cancelReply();
          this.toastService.success('Success', 'Reply posted');
        },
        error: (err) => {
          console.error('Failed to post reply:', err);
          const errorMsg = err.error?.error || 'Failed to post reply';
          this.toastService.error('Error', errorMsg);
        },
      });
  }

  addReplyToComment(
    comments: Comment[],
    parentId: string,
    reply: Comment
  ): boolean {
    for (const comment of comments) {
      if (comment.id === parentId) {
        if (!comment.replies) {
          comment.replies = [];
        }
        comment.replies.unshift(reply);
        return true;
      }
      if (comment.replies && comment.replies.length > 0) {
        if (this.addReplyToComment(comment.replies, parentId, reply)) {
          return true;
        }
      }
    }
    return false;
  }

  canReply(comment: Comment): boolean {
    const depth = comment.depth ?? 0; // Default to 0 if undefined
    return depth < 2; // Max 3 levels (0, 1, 2)
  }

  getReplyCharCount(): number {
    return this.replyForm.get('content')?.value?.length || 0;
  }

  toggleReactionPicker(event: Event) {
    event.stopPropagation();
    this.showReactionPicker = !this.showReactionPicker;
  }

  // New upvote/downvote methods
  handleVote(voteType: ReactionType, event: Event) {
    event.stopPropagation();

    if (!this.post) return;

    this.reactionService.toggleVote(this.post.id, voteType).subscribe({
      next: () => {
        // Update local vote state
        if (this.post) {
          const currentVote = this.reactionService.getUserVote(this.post.id);
          this.post.user_vote = currentVote;

          // Reload post to get updated counts from server
          this.loadPost(this.post.id);
        }
        this.showReactionPicker = false;
      },
      error: (err) => {
        console.error('Failed to vote:', err);
        this.toastService.error('Error', 'Failed to vote');
      },
    });
  }

  hasUpvoted(): boolean {
    return this.post ? this.post.user_vote === 'up' : false;
  }

  hasDownvoted(): boolean {
    return this.post ? this.post.user_vote === 'down' : false;
  }

  // Legacy method - kept for backward compatibility
  addReaction(reactionType: ReactionType, event: Event) {
    this.handleVote(reactionType, event);
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
        if (this.post) {
          this.post.share_count++;
        }
        this.toastService.success('Success', 'Post shared successfully');
      },
      error: (err) => {
        console.error('Failed to share post:', err);
        this.toastService.error('Error', 'Failed to share post');
      },
    });
  }

  getReactionEmoji(type: ReactionType): string {
    return this.reactionTypes[type] || '';
  }

  getReactionImage(type: ReactionType): string {
    return REACTION_IMAGES[type] || '';
  }

  getReactionKeys(): ReactionType[] {
    return Object.keys(this.reactionTypes) as ReactionType[];
  }

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

  // Author card hover methods
  showAuthorCardWithDelay(userId: string) {
    console.log('🔍 Hover started for userId:', userId);
    this.hoveredUserId = userId;
    this.authorCardTimeout = setTimeout(() => {
      console.log('⏰ Timeout triggered, loading profile...');
      this.loadProfile(userId).then((profile) => {
        console.log('✅ Profile loaded:', profile);
        console.log('🎯 hoveredUserId:', this.hoveredUserId, 'userId:', userId);
        if (this.hoveredUserId === userId) {
          this.authorProfile = profile;
          this.showAuthorCard = true;
          console.log('✨ Showing author card!', this.showAuthorCard);
        }
      });
    }, 300);
  }

  hideAuthorCard() {
    console.log('👋 Hiding author card');
    if (this.authorCardTimeout) {
      clearTimeout(this.authorCardTimeout);
    }
    this.showAuthorCard = false;
    this.hoveredUserId = null;
  }

  goBack() {
    this.router.navigate(['/feature/feed']);
  }

  getAuthorName(): string {
    if (!this.authorProfile) return 'Loading...';
    const firstName = this.authorProfile.firstName;
    const lastName = this.authorProfile.lastName;

    if (firstName && lastName) {
      return `${firstName} ${lastName}`;
    } else if (firstName) {
      return firstName;
    } else if (lastName) {
      return lastName;
    }

    return (
      this.authorProfile.username ||
      this.authorProfile.display_name ||
      'Unknown User'
    );
  }

  getUserName(userId: string): string {
    const profile = this.profileCache.get(userId);
    if (!profile) return userId;
    const firstName = profile.firstName || '';
    const lastName = profile.lastName || '';
    return `${firstName} ${lastName}`.trim() || profile.profileId || userId;
  }

  getAvatarUrl(profile: any): string | null {
    return profile?.avatarUri || null;
  }

  // Auto-resize textarea as user types
  autoResizeTextarea(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    textarea.style.height = 'auto';
    textarea.style.height = Math.min(textarea.scrollHeight, 150) + 'px';
  }

  // Get current character count
  getCharCount(): number {
    return this.commentForm.get('content')?.value?.length || 0;
  }
}
