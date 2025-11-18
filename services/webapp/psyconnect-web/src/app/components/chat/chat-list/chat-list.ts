import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Friend } from '../../../models/chat.models';
import { FriendService } from '../../../services/chat/profile.chat.service';

@Component({
  selector: 'app-chat-list',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './chat-list.html',
  styleUrls: ['./chat-list.scss'],
})
export class ChatListComponent implements OnInit {
  constructor(private friendService: FriendService) {}

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
      error: (err) => console.error('Error loading friends:', err),
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

  onSelect(friend: Friend) {
    this.friendSelected.emit(friend);
  }
  trackByFn(index: number, item: Friend) {
    return item.profileId;
  }
}
