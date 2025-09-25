import { CommonModule } from '@angular/common';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs/operators';

interface FileItem {
  id: string;
  name: string;
  mimeType: string;
  size: string;
}

interface FilesResponse {
  files: FileItem[];
}

@Component({
  selector: 'app-search-page',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './search.html',
  styleUrls: ['./search.scss'],
})
export class SearchComponent implements OnInit {
  private baseUrl = 'http://localhost:8080';

  files: FileItem[] = [];
  keyword = '';
  loading = false;

  constructor(private http: HttpClient) {}

  ngOnInit(): void {}

  search() {
    this.loading = true;

    const params = new HttpParams().set('q', this.keyword || '');

    this.http
      .get<FilesResponse>(`${this.baseUrl}/search`, {
        params,
        withCredentials: true,
      })
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res) => {
          this.files = res.files || [];
        },
        error: (err) => console.error('Search failed', err),
      });
  }

  downloadFile(file: FileItem) {
    this.http
      .get(`${this.baseUrl}/files/${file.id}`, {
        responseType: 'blob',
        withCredentials: true,
      })
      .subscribe((blob) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = file.name;
        a.click();
        window.URL.revokeObjectURL(url);
      });
  }

  deleteFile(file: FileItem) {
    if (!confirm(`Delete file ${file.name}?`)) return;
    this.http
      .delete(`${this.baseUrl}/files/${file.id}`, { withCredentials: true })
      .subscribe({
        next: () => this.search(),
        error: (err) => console.error('Delete failed', err),
      });
  }
}
