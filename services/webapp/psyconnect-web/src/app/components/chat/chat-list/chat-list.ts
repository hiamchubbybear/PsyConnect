import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Subject, Subscription, of } from 'rxjs';
import {
  catchError,
  debounceTime,
  distinctUntilChanged,
  switchMap,
} from 'rxjs/operators';
import { Friend } from '../../../models/chat.models';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { ImgFallbackDirective } from '../../../shared/directives/img-fallback.directive';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-chat-list',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    AvatarFallbackPipe,
    ImgFallbackDirective,
  ],
  templateUrl: './chat-list.html',
  styleUrls: ['./chat-list.scss'],
})
export class ChatListComponent implements OnInit, OnDestroy {
  constructor(
    private friendService: FriendService,
    private toastService: ToastService,
  ) {}
  @Input() selectedFriendId: string | null = null;

  @Input() currentUser!: {
    profileId: string;
    name: string;
    avatar: string;
    isOnline: boolean;
    status: string;
  };

  @Input() friends: Friend[] = [];
  @Input() isLoading = false;

  @Output() friendSelected = new EventEmitter<Friend>();

  searchQuery = '';
  searchSubject = new Subject<string>();
  private searchSubscription?: Subscription;

  isSearching = false;
  searchResults: Friend[] = [];
  visibleCount = 50;

  ngOnInit() {
    this.searchSubscription = this.searchSubject
      .pipe(
        debounceTime(300),
        distinctUntilChanged(),
        switchMap((query) => {
          if (!query.trim()) {
            this.isSearching = false;
            this.searchResults = [];
            return of([]);
          }
          this.isSearching = true;
          return this.friendService.searchProfiles(query).pipe(
            catchError(() => {
              return of([]);
            }),
          );
        }),
      )
      .subscribe((results) => {
        this.searchResults = results;
        this.isSearching = false;
      });
  }

  ngOnDestroy() {
    if (this.searchSubscription) {
      this.searchSubscription.unsubscribe();
    }
  }

  onSearchChange(query: string) {
    this.searchQuery = query;
    this.searchSubject.next(query);
  }

  get filteredFriends(): Friend[] {
    const q = this.searchQuery?.toLowerCase() ?? '';

    const localMatches = this.friends.filter(
      (f) =>
        f.firstName.toLowerCase().includes(q) ||
        f.lastName.toLowerCase().includes(q),
    );

    if (!q) {
      return localMatches.slice(0, this.visibleCount);
    }

    const combined = [...localMatches];
    const localIds = new Set(localMatches.map((f) => f.profileId));

    for (const res of this.searchResults) {
      if (!localIds.has(res.profileId)) {
        combined.push(res);
      }
    }

    return combined.slice(0, Math.max(this.visibleCount, combined.length));
  }

  loadMoreFriends() {
    if (this.visibleCount < this.friends.length) {
      this.visibleCount += 6;
    }
  }
  selectedUserId: string | null = null;

  onSelect(friend: Friend) {
    this.selectedFriendId = friend.profileId ?? null;
    this.friendSelected.emit(friend);

    // Clear search and reset view
    this.searchQuery = '';
    this.searchSubject.next('');
    this.searchResults = [];
    this.isSearching = false;
  }
  trackByFn(index: number, item: Friend) {
    return item.profileId;
  }
}
