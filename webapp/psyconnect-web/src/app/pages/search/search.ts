import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './search.html',
  styleUrls: ['./search.scss'],
})
export class SearchComponent {
  selectedCategory = 'All categories';
  query = '';
  posts = [
    {
      title: 'Healing Begins Within',
      author: '@phanhthegoodgirl',
      image:
        'https://i.pinimg.com/originals/4b/84/5c/4b845c9d4ae4b78530e8ea6eceab0970.jpg',
    },
    {
      title: 'Mind Over Matter',
      author: '@dr_mindful',
      image: 'https://i.pinimg.com/originals/bb/61/a6/bb61a69f0830c39a1b0cfb701c880898.jpg',
    },
    {
      title: 'Breaking Barriers',
      author: '@therapy_today',
      image: 'https://tse1.mm.bing.net/th/id/OIP.xarYcfZK0OzfFGzZGYQJ8gHaHa?rs=1&pid=ImgDetMain&o=7&rm=3',
    },
  ];

  get filtered() {
    if (!this.query) return this.posts;
    const q = this.query.toLowerCase();
    return this.posts.filter(
      (p) =>
        p.title.toLowerCase().includes(q) || p.author.toLowerCase().includes(q)
    );
  }
  get resultsCount() {
    return this.filtered.length;
  }

  clearTag() {
    this.query = '';
  }

  loadMore() {
    console.log('load more…');
  }

  trackByTitle(index: number, post: any) {
    return post.title;
  }
}
