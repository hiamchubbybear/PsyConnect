import { CommonModule } from '@angular/common';
import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import emotesData from '../../../assets/emotes-svg/emotes.json';

interface Emote {
  id: string;
  name: string;
  file: string;
  category: string;
  keywords: string[];
  shortcode: string;
}

interface EmoteCategory {
  id: string;
  name: string;
  icon: string;
  description: string;
}

@Component({
  selector: 'app-emote-picker',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './emote-picker.component.html',
  styleUrls: ['./emote-picker.component.scss'],
})
export class EmotePickerComponent implements OnInit {
  @Output() emoteSelected = new EventEmitter<Emote>();
  @Output() closed = new EventEmitter<void>();

  emotes: Emote[] = [];
  categories: EmoteCategory[] = [];
  selectedCategory: string = 'all';
  searchQuery: string = '';
  recentEmotes: Emote[] = [];

  ngOnInit() {
    this.emotes = (emotesData as any).emotes;
    this.categories = (emotesData as any).categories;
    this.loadRecentEmotes();
  }

  get filteredEmotes(): Emote[] {
    let filtered = this.emotes;

    
    if (this.selectedCategory !== 'all') {
      filtered = filtered.filter((e) => e.category === this.selectedCategory);
    }

    
    if (this.searchQuery) {
      const query = this.searchQuery.toLowerCase();
      filtered = filtered.filter(
        (e) =>
          e.name.toLowerCase().includes(query) ||
          e.keywords.some((k) => k.toLowerCase().includes(query)) ||
          e.shortcode.toLowerCase().includes(query)
      );
    }

    return filtered;
  }

  get emotesByCategory(): { [key: string]: Emote[] } {
    const grouped: { [key: string]: Emote[] } = {};

    this.categories.forEach((cat) => {
      grouped[cat.id] = this.emotes.filter((e) => e.category === cat.id);
    });

    return grouped;
  }

  selectEmote(emote: Emote) {
    this.addToRecent(emote);
    this.emoteSelected.emit(emote);
  }

  selectCategory(categoryId: string) {
    this.selectedCategory = categoryId;
  }

  close() {
    this.closed.emit();
  }

  private loadRecentEmotes() {
    const recent = localStorage.getItem('recentEmotes');
    if (recent) {
      const ids = JSON.parse(recent);
      this.recentEmotes = ids
        .map((id: string) => this.emotes.find((e) => e.id === id))
        .filter((e: Emote | undefined): e is Emote => e !== undefined)
        .slice(0, 8);
    }
  }

  private addToRecent(emote: Emote) {
    const recent = localStorage.getItem('recentEmotes');
    let ids: string[] = recent ? JSON.parse(recent) : [];

    
    ids = ids.filter((id) => id !== emote.id);

    
    ids.unshift(emote.id);

    
    ids = ids.slice(0, 8);

    localStorage.setItem('recentEmotes', JSON.stringify(ids));
    this.loadRecentEmotes();
  }

  getEmotePath(emote: Emote): string {
    return `assets/emotes-svg/${emote.file}`;
  }
}
