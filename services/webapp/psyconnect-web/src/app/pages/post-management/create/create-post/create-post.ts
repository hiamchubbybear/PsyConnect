import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-create-post',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './create-post.html',
  styleUrls: ['./create-post.scss'],
})
export class CreatePost {
  title = 'How does stress affect health?';
  description = '';
  anonymous = false;

  examples = [
    'I just need someone to listen...',
    "I feel lonely even when I'm not alone",
  ];

  saveDraft() {
    console.log('Save as draft:', this.title, this.description);
  }

  post() {
    console.log('Post submitted:', this.title, this.description);
  }
}
