import { CommonModule } from '@angular/common';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs/operators';

interface FileItem {
  id: string;
  name: string;
  size: number;
  folderId?: string;
}
interface FilesResponse {
  files: FileItem[];
}
interface FolderItem {
  id: string;
  name: string;
}

@Component({
  selector: 'app-files',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './files.html',
  styleUrls: ['./files.scss'],
})
export class FilesComponent implements OnInit {
  private baseUrl = 'http://localhost:8080';

  files: FileItem[] = [];
  folders: FolderItem[] = [];
  selectedFolderId!: string;
  uploadFile!: File;
  loading = false;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadFolders();
  }

  loadFolders() {
    this.http
      .get<any>(`${this.baseUrl}/folders`, this.getHttpOptions())
      .subscribe({
        next: (res) => {
          this.folders = res.folders || [];
          if (this.folders.length) {
            this.selectedFolderId = this.folders[0].id;
            this.loadFilesByFolder(this.selectedFolderId); // load file folder đầu tiên
          }
        },
        error: (err) => console.error('Load folders failed', err),
      });
  }

  private getHttpOptions() {
    return { withCredentials: true };
  }

  loadFiles() {
    this.loading = true;
    this.http
      .get<any>(`${this.baseUrl}/files`, this.getHttpOptions())
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res) => {
          // Nếu backend trả trực tiếp mảng hoặc object { files: [] }
          this.files = res.files || res || [];
        },
        error: (err) => console.error('Load files failed', err),
      });
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) this.uploadFile = file;
  }

  upload() {
    if (!this.uploadFile || !this.selectedFolderId) return;

    const formData = new FormData();
    formData.append('file', this.uploadFile);
    formData.append('folderId', this.selectedFolderId);

    this.http
      .post<any>(`${this.baseUrl}/upload`, formData, this.getHttpOptions())
      .subscribe({
        next: () => {
          this.uploadFile = undefined!;
          this.loadFiles();
        },
        error: (err) => console.error('Upload failed', err),
      });
  }

  downloadFile(file: FileItem) {
    this.http
      .get(`${this.baseUrl}/files/${file.id}`, {
        responseType: 'blob',
        ...this.getHttpOptions(),
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
      .delete(`${this.baseUrl}/files/${file.id}`, this.getHttpOptions())
      .subscribe({
        next: () => this.loadFiles(),
        error: (err) => console.error('Delete failed', err),
      });
  }

  onFolderChange() {
    this.loadFilesByFolder(this.selectedFolderId || undefined);
  }

  //   loadFilesByFolder(folderId?: string) {
  //     this.loading = true;
  //     let params = new HttpParams();
  //     if (folderId) {
  //       params = params.set('folderId', folderId);
  //     }

  //     this.http
  //       .get<FileItem[]>(`${this.baseUrl}/files/folder`, {
  //         ...this.getHttpOptions(),
  //         params,
  //       })
  //       .pipe(finalize(() => (this.loading = false)))
  //       .subscribe({
  //         next: (res) => {
  //           this.files = res; // backend trả mảng file
  //         },
  //         error: (err) => console.error('Load files failed', err),
  //       });
  //   }

  loadFilesByFolder(folderId?: string) {
    this.loading = true;
    let params = new HttpParams();
    if (folderId) {
      params = params.set('folderId', folderId);
    }

    this.http
      .get<FileItem[]>(`${this.baseUrl}/files/folder`, {
        ...this.getHttpOptions(),
        params,
      })
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res) => {
          this.files = res || []; // nếu folderId không truyền => backend trả tất cả file
        },
        error: (err) => console.error('Load files failed', err),
      });
  }
}
