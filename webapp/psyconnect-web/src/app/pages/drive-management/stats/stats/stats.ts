import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { finalize } from 'rxjs/operators';

interface KeyValue {
  key: string;
  value: number;
}
interface FolderItem {
  id: string;
  name: string;
}
@Component({
  standalone: true,
  imports: [CommonModule],
  selector: 'app-stats',
  templateUrl: './stats.html',
  styleUrls: ['./stats.scss'],
})
export class StatsComponent implements OnInit {
  private baseUrl = 'http://localhost:8080';

  totalFiles = 0;
  totalSize = 0;
  filesPerFolder: KeyValue[] = [];
  filesByType: KeyValue[] = [];
  folderSizes: KeyValue[] = [];

  folderMap: Record<string, string> = {};
  loading = false;

  constructor(private http: HttpClient) {}

  loadFilesPerFolder() {
    this.http
      .get<Record<string, number>>(
        `${this.baseUrl}/stats/files-per-folder`,
        this.getHttpOptions()
      )
      .subscribe({
        next: (res) => {
          this.filesPerFolder = Object.entries(res).map(([key, value]) => ({
            key,
            value,
          }));
        },
        error: (err) => console.error(err),
      });
  }

  ngOnInit() {
    this.loadFolders();
    this.loadStats();
    this.loadFilesPerFolder();
  }

  private getHttpOptions() {
    return { withCredentials: true };
  }

  loadFolders() {
    this.http
      .get<any>(`${this.baseUrl}/folders`, this.getHttpOptions())
      .subscribe({
        next: (res) => {
          const folders = res.folders || [];

          folders.forEach(({ id, name }: FolderItem) => {
            this.folderMap[id] = name;
          });

          this.loadStats();
        },
        error: (err) => console.error('Load folders failed', err),
      });
  }

  loadStats() {
    this.loading = true;

    this.http
      .get<any>(`${this.baseUrl}/stats/total-files`, this.getHttpOptions())
      .subscribe({ next: (res) => (this.totalFiles = res.totalFiles || 0) });

    this.http
      .get<any>(`${this.baseUrl}/stats/total-size`, this.getHttpOptions())
      .subscribe({ next: (res) => (this.totalSize = res.totalSize || 0) });

    this.http
      .get<Record<string, number>>(
        `${this.baseUrl}/stats/files-per-folder`,
        this.getHttpOptions()
      )
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res) => {
          this.filesPerFolder = Object.entries(res).map(([key, value]) => ({
            key: this.folderMap[key] || key,
            value,
          }));
        },
      });

    this.http
      .get<Record<string, number>>(
        `${this.baseUrl}/stats/files-by-type`,
        this.getHttpOptions()
      )
      .subscribe({
        next: (res) => {
          this.filesByType = Object.entries(res).map(([type, count]) => ({
            key: type,
            value: count,
          }));
        },
        error: (err) => console.error('Load files by type failed', err),
      });

    this.http
      .get<{ sizePerFolder: Record<string, number> }>(
        `${this.baseUrl}/stats/folder-size`,
        this.getHttpOptions()
      )
      .subscribe({
        next: (res) => {
          this.folderSizes = Object.entries(res.sizePerFolder).map(
            ([folderId, size]) => ({
              key: this.folderMap[folderId] || folderId,
              value: +(size / (1024 * 1024)).toFixed(2),
            })
          );
          console.log(this.folderSizes);
        },
        error: (err) => console.error('Load folder sizes failed', err),
      });
  }
  formatSize(bytes: number | string): string {
    let size = typeof bytes === 'string' ? parseInt(bytes, 10) : bytes;

    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
    if (size < 1024 * 1024 * 1024)
      return `${(size / (1024 * 1024)).toFixed(2)} MB`;
    return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  }
}
