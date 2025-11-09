import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { FriendListComponent } from '../../../components/social/friend-list/friend-list';
import { FriendSuggestionsComponent } from '../../../components/social/friend-suggestion/friend-suggestions';
import { FriendRequestsComponent } from '../../../components/social/friends-request/friend-requests';
import {
    Connection,
    FriendRequest,
    FriendSuggestion,
} from '../../../models/connection.model';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { ToastType } from '../../../shared/toast/toast.model';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-friends-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    FriendRequestsComponent,
    FriendListComponent,
    FriendSuggestionsComponent,
  ],
  templateUrl: './friends-page.html',
  styleUrls: ['./friends-page.scss'],
})
export class FriendsPage {
  searchQuery = '';

  private allFriends: Connection[] = [];
  private allSuggestions: FriendSuggestion[] = [];

  friends: Connection[] = [];
  friendRequests: FriendRequest[] = [];
  suggestions: FriendSuggestion[] = [];

  constructor(
    private friendService: FriendService,
    private toast: ToastService
  ) {}

  ngOnInit() {
    this.loadFriends();
    this.loadRequests();
    this.loadSuggestions();
  }

  private loadFriends() {
    this.friendService.getMyFriends().subscribe((friends) => {
      this.friends = friends.map((f) => ({
        id: f.profileId,
        name: `${f.firstName} ${f.lastName}`.trim(),
        avatarUrl: f.avatarUri,
      }));
    });
  }

  private loadRequests() {
    this.friendService.getReceivedRequests().subscribe((requests) => {
      this.friendRequests = requests.map((r) => ({
        id: r.profileId,
        name: `${r.firstName} ${r.lastName}`.trim(),
        avatarUrl: r.avatarUri,
        requestedAt: new Date(),
      }));
    });
  }

  private loadSuggestions() {
    this.friendService.getFriendSuggestions().subscribe((suggestions) => {
      this.suggestions = suggestions.map((s) => ({
        id: s.profileId,
        name: `${s.firstName} ${s.lastName}`.trim(),
        avatarUrl: s.avatarUri,
        mutualCount: 0,
        mutualAvatars: [],
      }));
    });
  }
  onAcceptRequest(id: string): void {
    console.log('[UI] Accept friend request:', id);

    this.friendService.acceptFriendRequest(id).subscribe({
      next: (res) => {
        console.log('[API] Accept response:', res);
        this.friendRequests = this.friendRequests.filter(
          (req) => req.id !== id
        );
        this.toast.show(
          'Thành công',
          'Đã chấp nhận lời mời kết bạn',
          ToastType.Success
        );
      },
      error: (err) => {
        console.error('[API] Accept failed:', err);
        this.toast.show('Lỗi', 'Không thể chấp nhận lời mời', ToastType.Error);
      },
    });
  }

  onDeclineRequest(id: string): void {
    console.log('[UI] Decline friend request:', id);
    this.friendRequests = this.friendRequests.filter((req) => req.id !== id);
    this.toast.show(
      'Đã từ chối',
      'Đã từ chối lời mời kết bạn',
      ToastType.Success
    );
  }

  onMessageFriend(id: string): void {
    console.log('[UI] Message friend:', id);

    this.toast.show(
      'Nhắn tin',
      `Đang mở khung chat với bạn ${id}`,
      ToastType.Success
    );
  }

  onUnfriend(id: string): void {
    console.log('[UI] Unfriend:', id);

    this.friendService.unfriend(id).subscribe({
      next: (res) => {
        console.log('[API] Unfriend response:', res);
        this.friends = this.friends.filter((f) => f.id !== id);
        this.allFriends = this.allFriends.filter((f) => f.id !== id);
        this.toast.show(
          'Đã hủy bạn bè',
          'Bạn đã hủy kết bạn thành công',
          ToastType.Success
        );
      },
      error: (err) => {
        console.error('[API] Unfriend failed:', err);
        this.toast.show('Lỗi', 'Không thể hủy kết bạn', ToastType.Error);
      },
    });
  }

  onAddFriend(id: string): void {
    console.log('[UI] Add friend:', id);

    this.friendService.sendFriendRequest(id).subscribe({
      next: (res) => {
        console.log('[API] Send request response:', res);
        this.suggestions = this.suggestions.filter((s) => s.id !== id);
        this.allSuggestions = this.allSuggestions.filter((s) => s.id !== id);
        this.toast.show(
          'Đã gửi',
          'Lời mời kết bạn đã được gửi',
          ToastType.Success
        );
      },
      error: (err) => {
        console.error('[API] Send request failed:', err);
        this.toast.show(
          'Lỗi',
          'Không thể gửi lời mời kết bạn',
          ToastType.Error
        );
      },
    });
  }
  onSearchFriends(): void {
    const query = this.searchQuery.trim().toLowerCase();
    if (!query) {
      this.friends = [...this.allFriends];
      this.suggestions = [...this.allSuggestions];
      return;
    }

    this.friends = this.allFriends.filter((f) =>
      f.name.toLowerCase().includes(query)
    );

    this.suggestions = this.allSuggestions.filter((s) =>
      s.name.toLowerCase().includes(query)
    );
  }
}
