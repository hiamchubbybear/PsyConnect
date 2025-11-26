import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { NewsfeedService } from '../../services/newsfeed/newsfeed.service';

interface SearchResult {
  type: 'post' | 'user' | 'tag';
  id: string;
  title: string;
  subtitle?: string;
  image?: string;
  metadata?: any;
}

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [FormsModule, CommonModule, RouterModule, TranslateModule],
  templateUrl: './search.html',
  styleUrls: ['./search.scss'],
})
export class SearchComponent implements OnInit {
  query = '';
  selectedCategory: 'all' | 'posts' | 'users' | 'tags' = 'all';
  results: SearchResult[] = [];
  loading = false;
  recentSearches: string[] = [];

  constructor(private newsfeedService: NewsfeedService) {}

  ngOnInit() {
    this.loadRecentSearches();
  }

  loadRecentSearches() {
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      this.recentSearches = JSON.parse(saved);
    }
  }

  saveRecentSearch(query: string) {
    if (!query.trim()) return;

    this.recentSearches = [query, ...this.recentSearches.filter(q => q !== query)].slice(0, 5);
    localStorage.setItem('recentSearches', JSON.stringify(this.recentSearches));
  }

  clearRecentSearches() {
    this.recentSearches = [];
    localStorage.removeItem('recentSearches');
  }

  onSearch() {
    if (!this.query.trim()) {
      this.results = [];
      return;
    }

    this.loading = true;
    this.saveRecentSearch(this.query);

    // Mock search - replace with actual API call
    setTimeout(() => {
      this.results = this.mockSearch(this.query);
      this.loading = false;
    }, 500);
  }

  mockSearch(query: string): SearchResult[] {
    const q = query.toLowerCase();
    const mockResults: SearchResult[] = [
      {
        type: 'post',
        id: '1',
        title: 'Understanding Anxiety and Depression',
        subtitle: 'A comprehensive guide to mental health',
        image: 'https://via.placeholder.com/100'
      },
      {
        type: 'post',
        id: '2',
        title: 'Mindfulness Meditation Techniques',
        subtitle: 'Learn to calm your mind',
        image: 'https://via.placeholder.com/100'
      },
      {
        type: 'user',
        id: '1',
        title: 'Dr. Sarah Johnson',
        subtitle: 'Clinical Psychologist',
        image: 'https://via.placeholder.com/100'
      },
      {
        type: 'tag',
        id: '1',
        title: '#mentalhealth',
        subtitle: '1,234 posts'
      }
    ];

    return mockResults.filter(r =>
      r.title.toLowerCase().includes(q) ||
      (r.subtitle && r.subtitle.toLowerCase().includes(q))
    );
  }

  onCategoryChange(category: 'all' | 'posts' | 'users' | 'tags') {
    this.selectedCategory = category;
    if (this.query) {
      this.onSearch();
    }
  }

  clearSearch() {
    this.query = '';
    this.results = [];
  }

  selectRecentSearch(query: string) {
    this.query = query;
    this.onSearch();
  }

  getResultIcon(type: string): string {
    switch (type) {
      case 'post': return 'fa-file-alt';
      case 'user': return 'fa-user';
      case 'tag': return 'fa-hashtag';
      default: return 'fa-circle';
    }
  }
}
