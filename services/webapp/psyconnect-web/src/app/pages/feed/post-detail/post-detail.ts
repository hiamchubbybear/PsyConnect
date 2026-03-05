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
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-post-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    TranslateModule,
    ReactiveFormsModule,
    AvatarFallbackPipe,
  ],
  templateUrl: './post-detail.html',
  styleUrls: ['./post-detail.scss'],
})
export class PostDetailComponent implements OnInit {
  post: Post | null = null;
  comments: Comment[] = [];
  
  error: string | null = null;
  commentForm!: FormGroup;
  submittingComment = false;

  
  replyingToCommentId: string | null = null;
  replyingToUsername: string | null = null;
  replyForm!: FormGroup;

  showReactionPicker = false;
  reactionTypes = REACTION_TYPES;
  isBookmarked = false;

  
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
    private loaderService: LoaderService,
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
            this.currentUserProfile,
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
      
      if (this.profileCache.has(userId)) {
        console.log(
          '💾 Profile found in cache:',
          this.profileCache.get(userId),
        );
        resolve(this.profileCache.get(userId));
        return;
      }

      
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
        
        if (Array.isArray(response)) {
          
          this.comments = response;
        } else if (response && response.comments) {
          
          this.comments = response.comments;
        } else {
          this.comments = [];
        }

        
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

        
        if (this.currentUserProfile) {
          this.profileCache.set(
            comment.author_id || comment.user_id,
            this.currentUserProfile,
          );
        }

        this.comments.unshift(comment);

        
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

        
        const errorMsg =
          err.error?.error || err.message || 'Failed to submit comment';
        this.toastService.error('Error', errorMsg);
        this.submittingComment = false;
      },
    });
  }

  
  startReply(comment: Comment) {
    console.log('🔵 START REPLY CALLED!', comment);
    this.replyingToCommentId = comment.id;
    
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
          
          this.addReplyToComment(this.comments, parentCommentId, reply);

          
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
    reply: Comment,
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
    const depth = comment.depth ?? 0; 
    return depth < 2; 
  }

  getReplyCharCount(): number {
    return this.replyForm.get('content')?.value?.length || 0;
  }

  toggleReactionPicker(event: Event) {
    event.stopPropagation();
    this.showReactionPicker = !this.showReactionPicker;
  }

  
  handleVote(voteType: ReactionType, event: Event) {
    event.stopPropagation();

    if (!this.post) return;

    const previousVote = this.post.user_vote;
    const previousUpvotes = this.post.upvote_count || 0;
    const previousDownvotes = this.post.downvote_count || 0;

    
    if (previousVote === voteType) {
      
      this.post.user_vote = null;
      if (voteType === 'up')
        this.post.upvote_count = Math.max(0, previousUpvotes - 1);
      else this.post.downvote_count = Math.max(0, previousDownvotes - 1);
    } else {
      
      this.post.user_vote = voteType;

      if (voteType === 'up') {
        this.post.upvote_count = previousUpvotes + 1;
        if (previousVote === 'down')
          this.post.downvote_count = Math.max(0, previousDownvotes - 1);
      } else {
        this.post.downvote_count = previousDownvotes + 1;
        if (previousVote === 'up')
          this.post.upvote_count = Math.max(0, previousUpvotes - 1);
      }
    }

    this.reactionService.toggleVote(this.post.id, voteType).subscribe({
      next: () => {
        
        this.newsfeedService.getPost(this.post!.id).subscribe((post) => {
          if (this.post) {
            this.post.upvote_count = post.upvote_count;
            this.post.downvote_count = post.downvote_count;
            this.post.user_vote = post.user_vote;
          }
        });
        this.showReactionPicker = false;
      },
      error: (err) => {
        console.error('Failed to vote:', err);
        
        if (this.post) {
          this.post.user_vote = previousVote;
          this.post.upvote_count = previousUpvotes;
          this.post.downvote_count = previousDownvotes;
        }
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
          this.isBookmarked ? 'Post bookmarked' : 'Bookmark removed',
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

    
    const firstName =
      this.authorProfile.firstName || this.authorProfile.first_name;
    const lastName =
      this.authorProfile.lastName || this.authorProfile.last_name;

    if (firstName || lastName) {
      return `${firstName || ''} ${lastName || ''}`.trim();
    }

    return (
      this.authorProfile.username || this.authorProfile.display_name || 'User'
    );
  }

  getUserName(userId: string): string {
    const profile = this.profileCache.get(userId);
    if (!profile) return 'User';

    const firstName = profile.firstName || profile.first_name || '';
    const lastName = profile.lastName || profile.last_name || '';

    const fullName = `${firstName} ${lastName}`.trim();
    if (fullName) return fullName;

    return profile.username || profile.display_name || 'User';
  }

  getAvatarUrl(profile: any): string | null {
    return profile?.avatarUri || null;
  }

  
  autoResizeTextarea(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    textarea.style.height = 'auto';
    textarea.style.height = Math.min(textarea.scrollHeight, 150) + 'px';
  }

  
  getCharCount(): number {
    return this.commentForm.get('content')?.value?.length || 0;
  }
}
