import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Connection } from '../../../models/connection.model';
import { ConfirmDialogService } from '../../../services/dialog/ConfirmDialog.service';
import { FriendActionDropdownComponent } from '../friend-action-dropdown/friend-action-dropdown';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'app-friend-card',
  standalone: true,
  imports: [CommonModule, TranslateModule, FriendActionDropdownComponent, AvatarFallbackPipe],
  templateUrl: './friend-card.html',
  styleUrls: ['./friend-card.scss'],
})
export class FriendsCard {
  @Input() friend!: Connection;
  @Input() cardType: 'request' | 'friend' | 'suggestion' = 'friend';
  @Input() mutualCount?: number;
  @Input() mutualAvatars?: string[];

  @Output() accept = new EventEmitter<string>();
  @Output() decline = new EventEmitter<string>();
  @Output() message = new EventEmitter<string>();
  @Output() unfriend = new EventEmitter<string>();
  @Output() addFriend = new EventEmitter<string>();
  @Output() viewProfile = new EventEmitter<string>();

  isHovered = false;

  onAccept() {
    this.accept.emit(this.friend.id);
  }

  onDecline() {
    this.decline.emit(this.friend.id);
  }

  onMessage() {
    this.message.emit(this.friend.id);
  }

  onAddFriend() {
    this.addFriend.emit(this.friend.id);
  }

  onViewProfile() {
    this.viewProfile.emit(this.friend.id);
  }

  constructor(private confirmDialog: ConfirmDialogService) {}

  async onUnfriend() {
    const confirmed = await this.confirmDialog.open(
      'dialog.unfriend_title',
      'dialog.unfriend_message'
    );

    if (confirmed) {
      this.unfriend.emit(this.friend.id);
    }
  }
}
