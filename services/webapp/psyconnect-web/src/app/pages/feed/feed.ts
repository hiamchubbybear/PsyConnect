import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Post } from '../../models/post.model';
import {
  REACTION_IMAGES,
  REACTION_TYPES,
  ReactionType,
} from '../../models/reaction.model';
import {
  Therapist as ApiTherapist,
  mapTherapistResponse,
} from '../../models/swipe-card';
import { SocialService } from '../../services/consultation/social.service';
import { Group, GroupService } from '../../services/group/group.service';
import { NewsfeedService } from '../../services/newsfeed/newsfeed.service';
import { ReactionService } from '../../services/newsfeed/reaction.service';
import { Profile } from '../../services/profile/profile';
import { SwipeService } from '../../services/swipe/swipe.service';
import { PsyButtonComponent as AppButtonComponent } from '../../shared/ui-atoms/button/psy-button.component';
import { PsyEmptyStateComponent as EmptyStateComponent } from '../../shared/ui-atoms/empty-state/psy-empty-state.component';
import { PostCardComponent } from '../../shared/ui-atoms/post-card/post-card.component';
import { SearchComponent } from '../search/search';
import { PostModalComponent } from './post-modal/post-modal';

interface Therapist {
  name: string;
  avatar: string;
  specialtyKey: string;
  rating: number;
  statusKey: string;
}

interface SupportGroup {
  nameKey: string;
  members: number;
}

@Component({
  selector: 'app-feed',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    TranslateModule,
    PostModalComponent,
    PostCardComponent,
    EmptyStateComponent,
    AppButtonComponent,
    SearchComponent,
  ],
  templateUrl: './feed.html',
  styleUrls: ['./feed.scss'],
})
export class FeedComponent implements OnInit {
  posts: Post[] = [];
  loading = false;
  error: string | null = null;
  feedType: 'all' | 'trending' = 'all';

  // Reaction picker state
  showReactionPicker: { [postId: string]: boolean } = {};
  reactionTypes = REACTION_TYPES;
  reactionImages = REACTION_IMAGES;

  // Profile hover state
  profileCache: Map<string, any> = new Map();
  showAuthorCard: { [postId: string]: boolean } = {};
  hoveredUserId: string | null = null;
  authorCardTimeout: any = null;

  // Post modal state
  showPostModal = false;
  selectedPostId: string | null = null;

  moodWeek = ['😊', '😐', '🙂', '😕', '😁', '😢', '🙂'];

  therapists: Therapist[] = [
    {
      name: 'Dr. Emily Chen',
      avatar: 'https://randomuser.me/api/portraits/women/65.jpg',
      specialtyKey: 'FEED.Therapists.P1.Specialty',
      rating: 4.9,
      statusKey: 'FEED.Therapists.P1.Status',
    },
    {
      name: 'Dr. Michael Torres',
      avatar: 'https://randomuser.me/api/portraits/men/43.jpg',
      specialtyKey: 'FEED.Therapists.P2.Specialty',
      rating: 4.7,
      statusKey: 'FEED.Therapists.P2.Status',
    },
  ];

  supportGroups: Group[] = [];
  viewedPosts: Set<string> = new Set();
  private currentPage = 0;
  private readonly PAGE_SIZE = 20;
  private observer: IntersectionObserver | null = null;

  constructor(
    private router: Router,
    private newsfeedService: NewsfeedService,
    private reactionService: ReactionService,
    private profileService: Profile,
    private socialService: SocialService,
    private swipeService: SwipeService,
    private groupService: GroupService,
  ) {}

  ngOnInit() {
    this.loadFeed();
    this.loadTherapists();
    this.loadGroups();
    this.setupIntersectionObserver();
  }

  loadFeed() {
    this.loading = true;
    this.error = null;

    const feedObservable =
      this.feedType === 'trending'
        ? this.newsfeedService.getTrendingPosts(20)
        : this.newsfeedService.getFeed(20, 0);

    feedObservable.subscribe({
      next: (posts) => {
        this.posts = posts || [];
        this.loading = false;
        // Load profiles for all post authors
        (posts || []).forEach((post) => {
          if (post.author_id) {
            this.loadProfile(post.author_id).then((profile) => {
              post.author_name =
                `${profile.firstName} ${profile.lastName}`.trim();
              post.author_avatar = profile.avatarUri;
            });
          }
        });
        this.viewedPosts.clear();
        setTimeout(() => this.observePosts(), 100);
      },
      error: (err: any) => {
        console.error('Failed to load feed:', err);
        this.error = 'Failed to load feed';
        this.loading = false;
      },
    });
  }

  switchFeed(type: 'all' | 'trending') {
    this.feedType = type;
    this.loadFeed();
  }

  loadTherapists() {
    this.swipeService.getSwipeData().subscribe({
      next: (res) => {
        const therapists = mapTherapistResponse(res);
        const slice = therapists.slice(0, 3);
        if (slice.length === 0) return;

        // Seed list with initial data from swipe API
        this.therapists = slice.map((t: ApiTherapist) => ({
          name: t.name || 'Therapist',
          avatar: t.avatarOverride || 'assets/images/default-avatar.png',
          specialtyKey: t.specialization?.[0] || 'Counseling',
          rating: t.rating || 5.0,
          profileId: t.profileId,
        })) as any[];

        // Enrich each entry with real profile data
        slice.forEach((t: ApiTherapist, i: number) => {
          if (!t.profileId) return;
          this.profileService.getProfileById(t.profileId).subscribe({
            next: (response: any) => {
              const profile = response?.data || response;
              const fullName =
                `${profile.firstName || ''} ${profile.lastName || ''}`.trim();
              this.therapists[i] = {
                ...this.therapists[i],
                name: fullName || this.therapists[i].name,
                avatar: profile.avatarUri || this.therapists[i].avatar,
              };
              this.therapists = [...this.therapists]; // trigger change detection
            },
            error: () => {},
          });
        });
      },
      error: (err: any) => {
        console.error('Failed to load therapists:', err);
      },
    });
  }

  loadGroups() {
    this.groupService.getGroups().subscribe({
      next: (groups) => {
        this.supportGroups = groups || [];
      },
      error: (err) => {
        console.error('Failed to load groups:', err);
        this.supportGroups = [];
      },
    });
  }

  setupIntersectionObserver() {
    this.observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const postId = entry.target.getAttribute('data-post-id');
            if (postId && !this.viewedPosts.has(postId)) {
              this.markAsViewed(postId);
            }
          }
        });
      },
      { threshold: 0.5 }, // 50% of the post must be visible
    );
  }

  observePosts() {
    if (!this.observer) return;
    const postElements = document.querySelectorAll('app-post-card');
    postElements.forEach((el) => this.observer?.observe(el));
  }

  markAsViewed(postId: string) {
    this.viewedPosts.add(postId);
    this.newsfeedService.recordView(postId).subscribe({
      error: (err: any) => console.error('Failed to record view:', err),
    });

    // When all posts have been seen, reload the feed
    if (this.posts.length > 0 && this.viewedPosts.size >= this.posts.length) {
      setTimeout(() => {
        this.loadFeed();
      }, 800);
    }
  }

  joinGroup(group: Group) {
    this.groupService.joinGroup(group.id).subscribe({
      next: () => {
        group.member_count++;
        // Show success toast or feedback
      },
      error: (err) => console.error('Failed to join group:', err),
    });
  }

  toggleMoodLog() {
    // Removed
  }

  bookTherapist(therapist: any) {
    this.router.navigate(['/feature/consultation/smart-match']);
  }

  contactTherapist(therapist: any) {
    if (therapist?.profileId) {
      this.router.navigate(['/feature/chat', therapist.profileId]);
    } else {
      this.router.navigate(['/feature/consultation/smart-match']);
    }
  }

  viewAllTherapists() {
    this.router.navigate(['/feature/consultation/smart-match']);
  }

  newPost() {
    this.router.navigate(['/feature/article/create']);
  }

  // Social Engagement methods
  upvotePost(post: Post, event: Event) {
    event.stopPropagation();

    // If already upvoted, remove the vote
    if (post.user_vote === 'up') {
      post.upvote_count--;
      post.user_vote = null;

      this.reactionService.removeReaction(post.id).subscribe({
        error: (err: any) => {
          console.error('Failed to remove upvote:', err);
          this.loadFeed(); // Revert on error
        },
      });
      return;
    }

    // Otherwise, add/switch to upvote
    if (post.user_vote === 'down') post.downvote_count--;
    post.upvote_count++;
    post.user_vote = 'up';

    this.reactionService.addReaction(post.id, 'up').subscribe({
      error: (err: any) => {
        console.error('Failed to upvote:', err);
        this.loadFeed(); // Revert on error
      },
    });
  }

  downvotePost(post: Post, event: Event) {
    event.stopPropagation();

    // If already downvoted, remove the vote
    if (post.user_vote === 'down') {
      post.downvote_count--;
      post.user_vote = null;

      this.reactionService.removeReaction(post.id).subscribe({
        error: (err: any) => {
          console.error('Failed to remove downvote:', err);
          this.loadFeed(); // Revert on error
        },
      });
      return;
    }

    // Otherwise, add/switch to downvote
    if (post.user_vote === 'up') post.upvote_count--;
    post.downvote_count++;
    post.user_vote = 'down';

    this.reactionService.addReaction(post.id, 'down').subscribe({
      error: (err: any) => {
        console.error('Failed to downvote:', err);
        this.loadFeed(); // Revert on error
      },
    });
  }

  sharePost(post: Post, event: Event) {
    event.stopPropagation();
    this.socialService.sharePost(post.id).subscribe({
      next: () => {
        post.share_count++;
      },
      error: (err: any) => console.error('Failed to share:', err),
    });
  }

  toggleBookmark(post: Post, event: Event) {
    event.stopPropagation();
    const isBookmarked = post.user_bookmark; // Assuming we add this to model or track it

    const action = isBookmarked
      ? this.socialService.removeBookmark(post.id)
      : this.socialService.addBookmark(post.id);

    action.subscribe({
      next: () => {
        post.user_bookmark = !isBookmarked;
      },
      error: (err: any) => console.error('Failed to toggle bookmark:', err),
    });
  }

  removeReaction(post: Post, event: Event) {
    event.stopPropagation();

    this.reactionService.removeReaction(post.id).subscribe({
      next: () => {
        if (post.user_vote === 'up') post.upvote_count--;
        else if (post.user_vote === 'down') post.downvote_count--;
        post.user_vote = null;
      },
      error: (err: any) => {
        console.error('Failed to remove reaction:', err);
      },
    });
  }

  // Post modal methods
  openPostModal(postId: string, event?: Event) {
    if (event) {
      event.stopPropagation();
    }
    this.selectedPostId = postId;
    this.showPostModal = true;
    // Prevent body scroll when modal is open
    document.body.style.overflow = 'hidden';
  }

  closePostModal() {
    this.showPostModal = false;
    this.selectedPostId = null;
    // Restore body scroll
    document.body.style.overflow = '';
  }

  openPost(post: Post) {
    this.router.navigate(['/feed/post', post.id]);
  }

  getReactionEmoji(type: ReactionType): string {
    return this.reactionTypes[type];
  }

  getReactionImage(type: ReactionType): string {
    return this.reactionImages[type];
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

  // Profile methods
  loadProfile(userId: string): Promise<any> {
    return new Promise((resolve) => {
      if (this.profileCache.has(userId)) {
        resolve(this.profileCache.get(userId));
        return;
      }

      this.profileService.getProfileById(userId).subscribe({
        next: (response: any) => {
          const profileData = response.data || response;
          this.profileCache.set(userId, profileData);
          resolve(profileData);
        },
        error: (err: any) => {
          const fallback = {
            firstName: userId,
            lastName: '',
            profileId: userId,
          };
          this.profileCache.set(userId, fallback);
          resolve(fallback);
        },
      });
    });
  }

  getUserName(userId: string): string {
    const profile = this.profileCache.get(userId);
    if (!profile) return userId;
    const firstName = profile.firstName || '';
    const lastName = profile.lastName || '';
    return `${firstName} ${lastName}`.trim() || profile.profileId || userId;
  }

  showAuthorCardWithDelay(postId: string, userId: string) {
    this.hoveredUserId = userId;
    this.authorCardTimeout = setTimeout(() => {
      this.loadProfile(userId).then(() => {
        if (this.hoveredUserId === userId) {
          this.showAuthorCard[postId] = true;
        }
      });
    }, 300);
  }

  hideAuthorCard(postId: string) {
    if (this.authorCardTimeout) {
      clearTimeout(this.authorCardTimeout);
    }
    this.showAuthorCard[postId] = false;
    this.hoveredUserId = null;
  }

  viewProfile(userId: string) {
    this.router.navigate(['/profile', userId]);
  }

  trackByPostId(index: number, post: Post): string {
    return post.id;
  }

  trackByTherapistId(index: number, therapist: any): string {
    return therapist.profileId || index.toString();
  }

  trackByGroupId(index: number, group: Group): string {
    return group.id;
  }
}
