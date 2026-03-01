import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { environment } from '../../../../environments/environment';
import { SecureStorageService } from '../../../encrypt/secure';
import { CreatePostRequest, MediaAttachment } from '../../../models/post.model';
import { CloudinaryService } from '../../../services/cloudinary/cloudinary.service';
import { NewsfeedService } from '../../../services/newsfeed/newsfeed.service';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-create-post',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './create-post.html',
  styleUrls: ['./create-post.scss'],
})
export class CreatePostComponent implements OnInit {
  postForm!: FormGroup;
  loading = false;
  selectedFiles: File[] = [];
  previewUrls: string[] = [];

  postTypes = [
    { value: 'article', labelKey: 'POST.Types.Article' },
    { value: 'question', labelKey: 'POST.Types.Question' },
    { value: 'discussion', labelKey: 'POST.Types.Discussion' },
    { value: 'resource', labelKey: 'POST.Types.Resource' },
  ];

  visibilityOptions = [
    { value: 'public', labelKey: 'POST.Visibility.Public' },
    { value: 'followers', labelKey: 'POST.Visibility.Followers' },
    { value: 'private', labelKey: 'POST.Visibility.Private' },
  ];
  constructor(
    private fb: FormBuilder,
    private newsfeedService: NewsfeedService,
    private router: Router,
    private toastService: ToastService,
    private cloudinaryService: CloudinaryService,
    private secureStorage: SecureStorageService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  initForm() {
    this.postForm = this.fb.group({
      title: [
        '',
        [
          Validators.required,
          Validators.minLength(5),
          Validators.maxLength(200),
        ],
      ],
      content: ['', [Validators.required, Validators.minLength(20)]],
      tags: [''],
      categories: [''],
      post_type: ['article', Validators.required],
      visibility: ['public', Validators.required],
    });
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files) {
      const files = Array.from(input.files);

      // Limit to 5 files
      if (this.selectedFiles.length + files.length > 5) {
        this.toastService.error('Error', 'Maximum 5 files allowed');
        return;
      }

      files.forEach((file) => {
        // Validate file size (max 5MB per file)
        if (file.size > 5 * 1024 * 1024) {
          this.toastService.error(
            'Error',
            `File ${file.name} is too large. Max 5MB per file.`,
          );
          return;
        }

        this.selectedFiles.push(file);

        // Create preview for images
        if (file.type.startsWith('image/')) {
          const reader = new FileReader();
          reader.onload = (e) => {
            this.previewUrls.push(e.target?.result as string);
          };
          reader.readAsDataURL(file);
        }
      });
    }
  }

  removeFile(index: number) {
    this.selectedFiles.splice(index, 1);
    this.previewUrls.splice(index, 1);
  }

  async onSubmit() {
    if (this.postForm.invalid) {
      Object.keys(this.postForm.controls).forEach((key) => {
        this.postForm.get(key)?.markAsTouched();
      });
      return;
    }

    this.loading = true;
    const formValue = this.postForm.value;

    try {
      // 1. Upload files to Cloudinary first
      const mediaAttachments: MediaAttachment[] = [];
      const username =
        this.secureStorage.getItem<string>(environment.usernameKey) ||
        'anonymous';

      if (this.selectedFiles.length > 0) {
        console.log(
          `[CREATE POST] Uploading ${this.selectedFiles.length} files...`,
        );
        for (const file of this.selectedFiles) {
          const type = file.type.startsWith('image/')
            ? 'image'
            : file.type.startsWith('video/')
              ? 'video'
              : 'document';

          const url = await this.cloudinaryService.uploadImage(
            file,
            username,
            'posts',
            'post',
          );

          mediaAttachments.push({
            type: type as any,
            url: url,
            caption: file.name,
          });
        }
      }

      // 2. Prepare post data
      const tags = formValue.tags
        ? formValue.tags
            .split(',')
            .map((t: string) => t.trim())
            .filter((t: string) => t)
        : [];
      const categories = formValue.categories
        ? formValue.categories
            .split(',')
            .map((c: string) => c.trim())
            .filter((c: string) => c)
        : [];

      const postData: CreatePostRequest = {
        title: formValue.title,
        content: formValue.content,
        tags,
        categories,
        post_type: formValue.post_type,
        visibility: formValue.visibility,
        media: mediaAttachments,
      };

      // 3. Create post on backend
      this.newsfeedService.createPost(postData).subscribe({
        next: (post) => {
          this.toastService.success('Success', 'Post created successfully!');
          this.router.navigate(['/feature/feed']); // Redirect to feed instead of post detail (which might be complex/missing)
        },
        error: (err) => {
          console.error('Failed to create post:', err);
          this.toastService.error(
            'Error',
            'Failed to create post record. Media uploaded but post failed.',
          );
          this.loading = false;
        },
      });
    } catch (error: any) {
      console.error('Upload process failed:', error);
      this.toastService.error(
        'Error',
        error.message || 'Failed to upload images',
      );
      this.loading = false;
    }
  }

  cancel() {
    this.router.navigate(['/feed']);
  }

  getErrorMessage(fieldName: string): string {
    const control = this.postForm.get(fieldName);
    if (!control || !control.errors || !control.touched) return '';

    if (control.errors['required']) return `${fieldName} is required`;
    if (control.errors['minlength']) {
      return `${fieldName} must be at least ${control.errors['minlength'].requiredLength} characters`;
    }
    if (control.errors['maxlength']) {
      return `${fieldName} must not exceed ${control.errors['maxlength'].requiredLength} characters`;
    }
    return '';
  }
}
