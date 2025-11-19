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
import { FriendRequest } from '../../../models/connection.model';
import { FriendsCard } from '../friends-card/friend-card';

@Component({
  selector: 'app-friend-requests',
  standalone: true,
  imports: [CommonModule, TranslateModule, FriendsCard, LucideAngularModule],
  templateUrl: './friend-requests.html',
  styleUrls: ['./friend-requests.scss'],
})
export class FriendRequestsComponent {
  @Input() requests: FriendRequest[] = [];
  @Output() acceptRequest = new EventEmitter<string>();
  @Output() declineRequest = new EventEmitter<string>();
  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;
  onAccept(id: string): void {
    this.acceptRequest.emit(id);
  }

  onDecline(id: string): void {
    this.declineRequest.emit(id);
  }

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
}
