import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs/operators';

interface FolderItem {
  id: string;
  name: string;
}

@Component({
  selector: 'app-folders',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './folders.html',
  styleUrls: ['./folders.scss'],
})
export class FoldersComponent implements OnInit {
  private baseUrl = 'http://localhost:8080';

  folders: FolderItem[] = [];
  newFolderName = '';
  renameFolderId = '';
  renameFolderName = '';
  loading = false;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadFolders();
  }

  private getHttpOptions() {
    return { withCredentials: true };
  }

  loadFolders() {
    this.loading = true;
    this.http
      .get<{ folders: FolderItem[] }>(
        `${this.baseUrl}/folders`,
        this.getHttpOptions()
      )
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res) => (this.folders = res.folders),
        error: (err) => console.error('Load folders failed', err),
      });
  }

  createFolder() {
    if (!this.newFolderName) return;
    this.http
      .post(
        `${this.baseUrl}/folders`,
        { name: this.newFolderName },
        this.getHttpOptions()
      )
      .subscribe({
        next: () => {
          this.newFolderName = '';
          this.loadFolders();
        },
      });
  }

  startRename(folder: FolderItem) {
    this.renameFolderId = folder.id;
    this.renameFolderName = folder.name;
  }

  renameFolder() {
    if (!this.renameFolderId || !this.renameFolderName) return;
    this.http
      .put(
        `${this.baseUrl}/folders/${this.renameFolderId}`,
        { newName: this.renameFolderName },
        this.getHttpOptions()
      )
      .subscribe({
        next: () => {
          this.renameFolderId = '';
          this.renameFolderName = '';
          this.loadFolders();
        },
      });
  }

  deleteFolder(id: string) {
    if (!confirm('Delete this folder?')) return;
    this.http
      .delete(`${this.baseUrl}/folders/${id}`, this.getHttpOptions())
      .subscribe({ next: () => this.loadFolders() });
  }
}
