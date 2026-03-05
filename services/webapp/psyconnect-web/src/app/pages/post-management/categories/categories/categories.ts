import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

interface Category {
  id: string;
  name: string;
  description: string;
  postCount: number;
  color: string;
  icon: string;
}

@Component({
  selector: 'app-categories',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, FormsModule],
  templateUrl: './categories.html',
  styleUrl: './categories.scss'
})
export class Categories implements OnInit {
  categories: Category[] = [];
  searchQuery = '';
  showAddModal = false;
  newCategory: Partial<Category> = {};

  ngOnInit() {
    this.loadCategories();
  }

  loadCategories() {
    
    this.categories = [
      { id: '1', name: 'Mental Health', description: 'Posts about mental wellness', postCount: 45, color: '#3b82f6', icon: 'fa-brain' },
      { id: '2', name: 'Therapy', description: 'Therapy techniques and tips', postCount: 32, color: '#10b981', icon: 'fa-user-md' },
      { id: '3', name: 'Self Care', description: 'Self care practices', postCount: 28, color: '#f59e0b', icon: 'fa-heart' },
      { id: '4', name: 'Mindfulness', description: 'Mindfulness and meditation', postCount: 21, color: '#8b5cf6', icon: 'fa-spa' },
      { id: '5', name: 'Relationships', description: 'Healthy relationships', postCount: 19, color: '#ec4899', icon: 'fa-users' },
    ];
  }

  get filteredCategories() {
    if (!this.searchQuery) return this.categories;
    const q = this.searchQuery.toLowerCase();
    return this.categories.filter(c =>
      c.name.toLowerCase().includes(q) ||
      c.description.toLowerCase().includes(q)
    );
  }

  openAddModal() {
    this.showAddModal = true;
    this.newCategory = { color: '#3b82f6', icon: 'fa-folder' };
  }

  closeAddModal() {
    this.showAddModal = false;
    this.newCategory = {};
  }

  saveCategory() {
    if (this.newCategory.name) {
      this.categories.push({
        id: Date.now().toString(),
        name: this.newCategory.name!,
        description: this.newCategory.description || '',
        postCount: 0,
        color: this.newCategory.color || '#3b82f6',
        icon: this.newCategory.icon || 'fa-folder'
      });
      this.closeAddModal();
    }
  }

  deleteCategory(id: string) {
    if (confirm('Delete this category?')) {
      this.categories = this.categories.filter(c => c.id !== id);
    }
  }
}
