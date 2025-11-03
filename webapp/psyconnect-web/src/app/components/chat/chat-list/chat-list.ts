import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Friend } from '../../../models/chat.models';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { ToastType } from '../../../shared/toast/toast-type';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-chat-list',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './chat-list.html',
  styleUrls: ['./chat-list.scss'],
})
export class ChatListComponent implements OnInit {
  constructor(
    private friendService: FriendService,
    private toastService: ToastService
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

  @Output() friendSelected = new EventEmitter<Friend>();

  searchQuery = '';
  visibleCount = 6;

  ngOnInit() {}

  private loadFriends() {
    this.friendService.getMyFriends().subscribe({
      next: (friends) => (this.friends = friends),
      error: () =>
        this.toastService.show(
          'TOAST.key_failed',
          'TOAST.key_load_friends_failed',
          ToastType.Error
        ),
    });
  }

  get filteredFriends(): Friend[] {
    const q = this.searchQuery?.toLowerCase() ?? '';
    return this.friends
      .filter(
        (f) =>
          f.firstName.toLowerCase().includes(q) ||
          f.lastName.toLowerCase().includes(q)
      )
      .slice(0, this.visibleCount);
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
  }
  trackByFn(index: number, item: Friend) {
    return item.profileId;
  }
}
