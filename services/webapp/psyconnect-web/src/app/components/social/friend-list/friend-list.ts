import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { LucideAngularModule } from 'lucide-angular';
import { Connection } from '../../../models/connection.model';
import { HorizontalListComponent } from '../../../shared/ui-atoms/horizontal-list/horizontal-list.component';
import { FriendsCard } from '../friends-card/friend-card';

@Component({
  selector: 'app-friend-list',
  standalone: true,
  imports: [
    CommonModule,
    FriendsCard,
    LucideAngularModule,
    TranslateModule,
    HorizontalListComponent,
  ],
  templateUrl: './friend-list.html',
  styleUrls: ['./friend-list.scss'],
})
export class FriendListComponent {
  @Input() friends: Connection[] = [];
  @Output() messageFriend = new EventEmitter<string>();
  @Output() unfriendUser = new EventEmitter<string>();
  constructor(private router: Router) {}
  onMessage(friendId: string) {
    this.router.navigate(['/chat', friendId]);
  }
  onUnfriend(friend: any) {
    console.log('Unfriend', friend);
  }
}
