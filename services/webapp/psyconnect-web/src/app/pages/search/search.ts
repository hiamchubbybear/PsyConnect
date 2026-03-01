import { CommonModule } from '@angular/common';
import {
  Component,
  ElementRef,
  HostListener,
  OnInit,
  ViewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { forkJoin, map } from 'rxjs';
import { TherapistService } from '../../services/consultation/therapist.service';
import { GroupService } from '../../services/group/group.service';
import { NewsfeedService } from '../../services/newsfeed/newsfeed.service';
import { Profile } from '../../services/profile/profile';

type Category = 'posts' | 'people' | 'therapists' | 'groups' | 'tags';

interface SearchResult {
  type: 'post' | 'user' | 'therapist' | 'group' | 'tag';
  id: string;
  title: string;
  subtitle?: string;
  image?: string;
  metadata?: any;
}

interface RecentSearch {
  query: string;
  category?: string;
  icon?: string;
}

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [FormsModule, CommonModule, RouterModule],
  templateUrl: './search.html',
  styleUrls: ['./search.scss'],
})
export class SearchComponent implements OnInit {
  @ViewChild('searchWrapper') searchWrapper!: ElementRef;

  query = '';
  isDropdownOpen = false;
  selectedCategory: Category | null = null;
  results: SearchResult[] = [];
  loading = false;
  recentSearches: RecentSearch[] = [];
  popularTags: string[] = [];
  hasSearched = false;

  readonly categories: { id: Category; label: string; icon: string }[] = [
    { id: 'posts', label: 'Posts', icon: 'post' },
    { id: 'people', label: 'People', icon: 'people' },
    { id: 'therapists', label: 'Therapists', icon: 'therapist' },
    { id: 'groups', label: 'Groups', icon: 'group' },
    { id: 'tags', label: 'Tags', icon: 'tag' },
  ];

  constructor(
    private newsfeedService: NewsfeedService,
    private profileService: Profile,
    private therapistService: TherapistService,
    private groupService: GroupService,
    private router: Router,
  ) {}

  ngOnInit() {
    this.loadRecentSearches();
    this.loadPopularTags();
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent) {
    if (
      this.searchWrapper &&
      !this.searchWrapper.nativeElement.contains(event.target)
    ) {
      this.isDropdownOpen = false;
    }
  }

  loadRecentSearches() {
    const saved = localStorage.getItem('psy_recentSearches');
    if (saved) {
      this.recentSearches = JSON.parse(saved).slice(0, 5);
    }
  }

  loadPopularTags() {
    this.newsfeedService.getPopularTags(8).subscribe({
      next: (tags) => (this.popularTags = tags || []),
      error: () => (this.popularTags = []),
    });
  }

  openDropdown() {
    this.isDropdownOpen = true;
  }

  selectCategory(category: Category) {
    if (this.selectedCategory === category) {
      this.selectedCategory = null;
    } else {
      this.selectedCategory = category;
    }
    if (this.query.trim()) {
      this.performSearch();
    }
  }

  removeCategory() {
    this.selectedCategory = null;
    if (this.query.trim()) {
      this.performSearch();
    }
  }

  onQueryChange() {
    this.isDropdownOpen = true;
    if (this.query.trim().length >= 1) {
      this.performSearch();
    } else {
      this.results = [];
      this.hasSearched = false;
    }
  }

  performSearch() {
    if (!this.query.trim()) return;
    this.loading = true;
    this.hasSearched = true;
    this.saveRecentSearch(this.query);

    const searches: any = {};

    if (!this.selectedCategory || this.selectedCategory === 'posts') {
      searches.posts = this.newsfeedService.searchPosts(this.query, 8, 0).pipe(
        map((posts) =>
          (posts || []).map((p) => ({
            type: 'post' as const,
            id: p.id,
            title: p.title || 'Post',
            subtitle: p.content?.slice(0, 80),
            image: p.media?.[0]?.url || p.image_url,
            metadata: { tags: p.hashtags },
          })),
        ),
      );
    }

    if (!this.selectedCategory || this.selectedCategory === 'people') {
      searches.users = this.profileService
        .searchProfiles(this.query, 0, 8)
        .pipe(
          map((res) =>
            (res?.data || []).map((p: any) => ({
              type: 'user' as const,
              id: p.profileId,
              title: `${p.firstName} ${p.lastName}`.trim() || 'User',
              subtitle: p.address || 'Member',
              image: p.avatarUri,
            })),
          ),
        );
    }

    if (!this.selectedCategory || this.selectedCategory === 'therapists') {
      searches.therapists = this.therapistService
        .searchTherapists(this.query, 8, 0)
        .pipe(
          map((res) => {
            const list = res?.data || res || [];
            return list.map((t: any) => ({
              type: 'therapist' as const,
              id: t.profile_id,
              title: t.name || 'Therapist',
              subtitle: t.specialization?.join(', ') || 'Licensed Professional',
              image: t.avatar_override,
              metadata: { rating: t.rating },
            }));
          }),
        );
    }

    if (!this.selectedCategory || this.selectedCategory === 'groups') {
      searches.groups = this.groupService
        .getGroups(undefined, 8, 0, this.query)
        .pipe(
          map((groups) =>
            (groups || []).map((g) => ({
              type: 'group' as const,
              id: g.id,
              title: g.name,
              subtitle: g.description,
              metadata: { memberCount: g.member_count },
            })),
          ),
        );
    }

    forkJoin(searches).subscribe({
      next: (res: any) => {
        this.results = [
          ...(res.users || []),
          ...(res.therapists || []),
          ...(res.groups || []),
          ...(res.posts || []),
        ];
        this.loading = false;
      },
      error: () => {
        this.loading = false;
      },
    });
  }

  saveRecentSearch(query: string) {
    if (!query.trim()) return;
    const entry: RecentSearch = {
      query,
      category: this.selectedCategory
        ? this.categories.find((c) => c.id === this.selectedCategory)?.label
        : undefined,
    };
    this.recentSearches = [
      entry,
      ...this.recentSearches.filter((r) => r.query !== query),
    ].slice(0, 5);
    localStorage.setItem(
      'psy_recentSearches',
      JSON.stringify(this.recentSearches),
    );
  }

  clearRecentSearches() {
    this.recentSearches = [];
    localStorage.removeItem('psy_recentSearches');
  }

  selectRecentSearch(recent: RecentSearch) {
    this.query = recent.query;
    if (recent.category) {
      const cat = this.categories.find((c) => c.label === recent.category);
      if (cat) this.selectedCategory = cat.id;
    }
    this.performSearch();
  }

  selectTag(tag: string) {
    this.selectedCategory = 'tags';
    this.query = '#' + tag;
    this.performSearch();
  }

  navigateToResult(result: SearchResult) {
    this.isDropdownOpen = false;
    switch (result.type) {
      case 'user':
        this.router.navigate(['/profile', result.id]);
        break;
      case 'post':
        this.router.navigate(['/feature/feed/post', result.id]);
        break;
      case 'therapist':
        this.router.navigate(['/feature/consultation/therapist', result.id]);
        break;
      case 'group':
        this.router.navigate(['/feature/groups', result.id]);
        break;
    }
  }

  clearQuery() {
    this.query = '';
    this.results = [];
    this.hasSearched = false;
  }

  get showTagSuggestions(): boolean {
    return (
      this.selectedCategory === 'tags' &&
      !this.hasSearched &&
      this.popularTags.length > 0
    );
  }

  get showResults(): boolean {
    return this.hasSearched && !this.loading;
  }

  get showSuggestions(): boolean {
    return !this.hasSearched && !this.loading;
  }

  get selectedCategoryLabel(): string {
    return (
      this.categories.find((c) => c.id === this.selectedCategory)?.label || ''
    );
  }

  trackById(_: number, result: SearchResult): string {
    return `${result.type}-${result.id}`;
  }
}
