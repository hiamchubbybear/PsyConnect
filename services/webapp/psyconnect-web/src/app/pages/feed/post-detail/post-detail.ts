import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
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
  selector: 'app-post-detail',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, ReactiveFormsModule],
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

  showReactionPicker = false;
  reactionTypes = REACTION_TYPES;
  isBookmarked = false;

  // Author profiles cache
  profileCache: Map<string, any> = new Map();
  authorProfile: any = null;
  showAuthorCard = false;
  authorCardTimeout: any = null;
  hoveredUserId: string | null = null;

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
    const postId = this.route.snapshot.paramMap.get('id');
    if (postId) {
      this.loadPost(postId);
      this.loadComments(postId);
    }
  }

  initCommentForm() {
    this.commentForm = this.fb.group({
      content: ['', [Validators.required, Validators.minLength(1)]],
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
    this.loadProfile(authorId).then(profile => {
      this.authorProfile = profile;
    });
  }

  // Generic profile loading with caching
  loadProfile(userId: string): Promise<any> {
    console.log('📥 loadProfile called for:', userId);
    return new Promise((resolve) => {
      // Check cache first
      if (this.profileCache.has(userId)) {
        console.log('💾 Profile found in cache:', this.profileCache.get(userId));
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
            username: userId
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
        this.comments = comments || []; // Ensure array is initialized
        // Load profiles for all comment authors
        comments?.forEach(comment => {
          if (comment.author_id && !this.profileCache.has(comment.author_id)) {
            this.loadProfile(comment.author_id);
          }
        });
      },
      error: (err: any) => {
        console.error('Failed to load comments:', err);
        this.comments = []; // Initialize empty array on error
      },
    });
  }

  submitComment() {
    if (this.commentForm.invalid || !this.post) return;

    this.submittingComment = true;
    const content = this.commentForm.value.content;

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
        this.toastService.error('Error', 'Failed to submit comment');
        this.submittingComment = false;
      },
    });
  }

  toggleReactionPicker(event: Event) {
    event.stopPropagation();
    this.showReactionPicker = !this.showReactionPicker;
  }

  addReaction(reactionType: ReactionType, event: Event) {
    event.stopPropagation();
    if (!this.post) return;

    this.reactionService.addReaction(this.post.id, reactionType).subscribe({
      next: () => {
        if (this.post) {
          this.post.like_count++;
        }
        this.showReactionPicker = false;
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
    return this.reactionTypes[type];
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
      this.loadProfile(userId).then(profile => {
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

    return this.authorProfile.username || this.authorProfile.display_name || 'Unknown User';
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
}
