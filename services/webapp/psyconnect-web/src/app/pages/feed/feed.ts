import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Post } from '../../models/post.model';
import { REACTION_TYPES, ReactionType } from '../../models/reaction.model';
import { NewsfeedService } from '../../services/newsfeed/newsfeed.service';
import { ReactionService } from '../../services/newsfeed/reaction.service';
import { Profile } from '../../services/profile/profile';

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
  imports: [CommonModule, RouterModule, TranslateModule],
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

  // Profile hover state
  profileCache: Map<string, any> = new Map();
  showAuthorCard: { [postId: string]: boolean } = {};
  hoveredUserId: string | null = null;
  authorCardTimeout: any = null;

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

  supportGroups: SupportGroup[] = [
    { nameKey: 'FEED.SupportGroups.P1', members: 127 },
    { nameKey: 'FEED.SupportGroups.P2', members: 89 },
  ];

  constructor(
    private router: Router,
    private newsfeedService: NewsfeedService,
    private reactionService: ReactionService,
    private profileService: Profile
  ) {}

  ngOnInit() {
    this.loadFeed();
  }

  loadFeed() {
    this.loading = true;
    this.error = null;

    const feedObservable = this.feedType === 'trending'
      ? this.newsfeedService.getTrendingPosts(20)
      : this.newsfeedService.getFeed(20, 0);

    feedObservable.subscribe({
      next: (posts) => {
        this.posts = posts;
        this.loading = false;
        // Load profiles for all post authors
        posts.forEach(post => {
          if (post.author_id && !this.profileCache.has(post.author_id)) {
            this.loadProfile(post.author_id);
          }
        });
      },
      error: (err: any) => {
        console.error('Failed to load feed:', err);
        this.error = 'Failed to load feed';
        this.loading = false;
      }
    });
  }

  switchFeed(type: 'all' | 'trending') {
    this.feedType = type;
    this.loadFeed();
  }

  toggleMoodLog() {
    alert('Open mood log (implement modal)');
  }

  bookTherapist(therapist: any) {
    alert(`Book ${therapist.name} (implement booking)`);
  }

  newPost() {
    this.router.navigate(['/feature/article/create']);
  }

  // Reaction methods
  toggleReactionPicker(postId: string, event: Event) {
    event.stopPropagation();
    this.showReactionPicker[postId] = !this.showReactionPicker[postId];
  }

  addReaction(post: Post, reactionType: ReactionType, event: Event) {
    event.stopPropagation();

    this.reactionService.addReaction(post.id, reactionType).subscribe({
      next: () => {
        post.like_count++;
        this.showReactionPicker[post.id] = false;
      },
      error: (err) => {
        console.error('Failed to add reaction:', err);
      }
    });
  }

  removeReaction(post: Post, event: Event) {
    event.stopPropagation();

    this.reactionService.removeReaction(post.id).subscribe({
      next: () => {
        post.like_count = Math.max(0, post.like_count - 1);
      },
      error: (err) => {
        console.error('Failed to remove reaction:', err);
      }
    });
  }

  openPost(post: Post) {
    this.router.navigate(['/feed/post', post.id]);
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
          const fallback = { firstName: userId, lastName: '', profileId: userId };
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
}
