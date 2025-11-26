import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Post } from '../../../../models/post.model';
import { LoaderService } from '../../../../services/loader/loader';
import { NewsfeedService } from '../../../../services/newsfeed/newsfeed.service';

@Component({
  selector: 'app-all-posts',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, FormsModule],
  templateUrl: './all-posts.html',
  styleUrl: './all-posts.scss'
})
export class AllPosts implements OnInit {
  posts: Post[] = [];
  filteredPosts: Post[] = [];
  loading = false;
  searchQuery = '';
  selectedFilter: 'all' | 'published' | 'draft' = 'all';
  selectedSort: 'newest' | 'oldest' | 'popular' = 'newest';

  stats = {
    total: 0,
    published: 0,
    draft: 0,
    views: 0
  };

  constructor(
    private newsfeedService: NewsfeedService,
    private loaderService: LoaderService
  ) {}

  ngOnInit() {
    this.loadPosts();
  }

  loadPosts() {
    this.loading = true;
    this.loaderService.show();

    this.newsfeedService.getFeed(0).subscribe({
      next: (posts) => {
        this.posts = posts;
        this.filteredPosts = posts;
        this.calculateStats();
        this.applyFilters();
        this.loading = false;
        this.loaderService.hide();
      },
      error: (err) => {
        console.error('Failed to load posts:', err);
        this.loading = false;
        this.loaderService.hide();
      }
    });
  }

  calculateStats() {
    this.stats.total = this.posts.length;
    this.stats.published = this.posts.filter(p => p.status === 'published').length;
    this.stats.draft = this.posts.filter(p => p.status === 'draft').length;
    this.stats.views = this.posts.reduce((sum, p) => sum + (p.view_count || 0), 0);
  }

  applyFilters() {
    let filtered = [...this.posts];

    // Filter by status
    if (this.selectedFilter !== 'all') {
      filtered = filtered.filter(p => p.status === this.selectedFilter);
    }

    // Search
    if (this.searchQuery) {
      const query = this.searchQuery.toLowerCase();
      filtered = filtered.filter(p =>
        p.title.toLowerCase().includes(query) ||
        p.content.toLowerCase().includes(query)
      );
    }

    // Sort
    filtered.sort((a, b) => {
      switch (this.selectedSort) {
        case 'newest':
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        case 'oldest':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        case 'popular':
          return (b.view_count || 0) - (a.view_count || 0);
        default:
          return 0;
      }
    });

    this.filteredPosts = filtered;
  }

  onSearchChange(query: string) {
    this.searchQuery = query;
    this.applyFilters();
  }

  onFilterChange(filter: 'all' | 'published' | 'draft') {
    this.selectedFilter = filter;
    this.applyFilters();
  }

  onSortChange(sort: 'newest' | 'oldest' | 'popular') {
    this.selectedSort = sort;
    this.applyFilters();
  }

  deletePost(postId: string) {
    if (confirm('Are you sure you want to delete this post?')) {
      // TODO: Implement delete
      console.log('Delete post:', postId);
    }
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }
}
