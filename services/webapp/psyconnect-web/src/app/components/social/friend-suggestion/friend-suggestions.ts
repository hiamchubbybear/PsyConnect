
import { CommonModule } from '@angular/common';
import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { LucideAngularModule } from 'lucide-angular';
import { FriendSuggestion } from '../../../models/connection.model';
import { FriendsCard } from '../friends-card/friend-card';

@Component({
  selector: 'app-friend-suggestions',
  standalone: true,
  imports: [CommonModule, TranslateModule, FriendsCard, LucideAngularModule],
  templateUrl: './friend-suggestions.html',
  styleUrls: ['./friend-suggestions.scss'],
})
export class FriendSuggestionsComponent {
  @Input() suggestions: FriendSuggestion[] = [];
  @Output() addFriend = new EventEmitter<string>();
  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;
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

  onAddFriend(id: string): void {
    this.addFriend.emit(id);
  }
}
