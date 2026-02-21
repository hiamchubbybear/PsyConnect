import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { LucideAngularModule } from 'lucide-angular';
import { FriendRequest } from '../../../models/connection.model';
import { HorizontalListComponent } from '../../../shared/ui-atoms/horizontal-list/horizontal-list.component';
import { FriendsCard } from '../friends-card/friend-card';

@Component({
  selector: 'app-friend-requests',
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule,
    FriendsCard,
    LucideAngularModule,
    HorizontalListComponent,
  ],
  templateUrl: './friend-requests.html',
  styleUrls: ['./friend-requests.scss'],
})
export class FriendRequestsComponent {
  @Input() requests: FriendRequest[] = [];
  @Output() acceptRequest = new EventEmitter<string>();
  @Output() declineRequest = new EventEmitter<string>();

  onAccept(id: string): void {
    this.acceptRequest.emit(id);
  }

  onDecline(id: string): void {
    this.declineRequest.emit(id);
  }
}
