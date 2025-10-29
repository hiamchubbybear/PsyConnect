import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    EventEmitter,
    Input,
    Output,
    ViewChild,
} from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { LucideAngularModule } from 'lucide-angular';
import { Connection } from '../../../models/connection.model';
import { FriendsCard } from '../friends-card/friend-card';

@Component({
  selector: 'app-friend-list',
  standalone: true,
  imports: [
    CommonModule,
    FriendsCard,
    LucideAngularModule,
    TranslateModule,
  ],
  templateUrl: './friend-list.html',
  styleUrls: ['./friend-list.scss'],
})
export class FriendListComponent {
  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;
  @Input() friends: Connection[] = [];
  @Output() messageFriend = new EventEmitter<string>();
  @Output() unfriendUser = new EventEmitter<string>();

  scrollLeft() {
    this.scrollContainer.nativeElement.scrollBy({
      left: -300,
      behavior: 'smooth',
    });
  }

  scrollRight() {
    this.scrollContainer.nativeElement.scrollBy({
      left: 300,
      behavior: 'smooth',
    });
  }

  onMessage(friend: any) {
    console.log('Message', friend);
  }

  onUnfriend(friend: any) {
    console.log('Unfriend', friend);
  }
}
