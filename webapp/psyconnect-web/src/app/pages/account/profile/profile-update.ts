import { CommonModule } from '@angular/common';
import { Component, HostListener, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { map, Observable, of } from 'rxjs';
import { ButtonGroupComponent } from '../../../components/button-group/button-group';
import { SecureStorageService } from '../../../encrypt/secure';
import { CloudinaryService } from '../../../services/cloudinary/cloudinary.service';
import { LoaderService } from '../../../services/loader/loader';
import {
  Profile,
  UserProfileUpdateRequest,
} from '../../../services/profile/profile';
import {
  UserContextService,
  UserProfile,
} from '../../../services/profile/profile-service';
import { ToastType } from '../../../shared/toast/toast.model';
import { ToastService } from '../../../shared/toast/toast.service';
import { ProfileModel } from './profile-model-update';

@Component({
  selector: 'app-profile-section',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    ButtonGroupComponent,
  ],
  templateUrl: './profile-update.html',
  styleUrls: ['./profile-update.scss'],
})
export class ProfileSectionComponent implements OnInit {
  selectedImage: File | null = null;
  imagePreview: string | null = null;
  username: string | null = '';
  isUploadImage?: boolean | false;
  profileUpdate: UserProfileUpdateRequest | undefined;
  onImageSelected(event: any) {
    const file = event.target.files[0];
    if (file) {
      if (!file.type.startsWith('image/')) {
        this.toastService.show(
          'Vui lòng chọn file ảnh hợp lệ!',
          'Error',
          ToastType.Error
        );
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        this.toastService.show(
          'File ảnh không được vượt quá 5MB!',
          'Error',
          ToastType.Error
        );
        return;
      }
      this.selectedImage = file;
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.imagePreview = e.target.result as string;
      };
      reader.readAsDataURL(file);
      this.isUploadImage = true;
    }
  }

  async confirmChanges() {
    this.loaderService.show();
    try {
      let avatarUri = this.profile?.avatarUri;

      if (this.isUploadImage && this.selectedImage) {
        const uploadedUrl = await this.cloudinaryService.uploadImage(
          this.selectedImage,
          this.username || ''
        );

        if (!uploadedUrl || uploadedUrl.trim() === '') {
          throw new Error('Upload failed - no URL returned');
        }

        avatarUri = uploadedUrl;
        this.selectedImage = null;
        this.imagePreview = null;
        this.isUploadImage = false;
      }
      if (avatarUri) this.updateProfile(avatarUri);
    } catch (err) {
      console.error('Upload error:', err);
      this.toastService.show(
        'Upload ảnh thất bại! Vui lòng thử lại.',
        'Error',
        ToastType.Error
      );
      this.selectedImage = null;
      this.imagePreview = null;
      this.isUploadImage = false;
    } finally {
      this.loaderService.hide();
    }
  }
  form!: FormGroup;
  profile: ProfileModel | null = null;

  editing: string | null = null;
  editValue: any = '';
  isOverlayOpen: boolean = false;
  draftData = {
    firstName: '',
    middleName: '',
    lastName: '',
  };
  constructor(
    private fb: FormBuilder,
    private secureStorage: SecureStorageService,
    private profileService: Profile,
    private cloudinaryService: CloudinaryService,
    private toastService: ToastService,
    private loaderService: LoaderService,
    private userContext: UserContextService
  ) {}

  ngOnInit(): void {
    this.getProfile().subscribe((profile) => {
      if (profile) {
        this.profile = profile;
        this.form.patchValue(profile);
      }
    });
    this.username = this.secureStorage.getItem('username');
    this.form = this.fb.group({
      firstName: [this.profile?.firstName || ''],
      lastName: [this.profile?.lastName || ''],
      dob: [this.profile?.dob || ''],
      address: [this.profile?.address || ''],
      gender: [this.profile?.gender || ''],
      avatarUri: [this.profile?.avatarUri || ''],
    });

    this.profileService.getProfile().subscribe((res) => {
      if (res?.data) {
        this.form.patchValue({
          firstName: res.data.firstName,
          lastName: res.data.lastName,
          dob: res.data.dob,
          address: res.data.address,
          gender: res.data.gender,
          avatarUri: res.data.avatarUri,
        });
      }
    });
  }

  getProfile(): Observable<ProfileModel | null> {
    const profileJson = this.secureStorage.getItem('profile');
    if (profileJson) {
      return of(profileJson as ProfileModel);
    } else {
      return this.profileService
        .getProfile()
        .pipe(map((res) => res?.data || null));
    }
  }

  openOverlay() {
    this.isOverlayOpen = true;
  }

  closeOverlay() {
    this.isOverlayOpen = false;
  }

  openEdit(field: string) {
    this.editing = field;
    if (field === 'name') {
      this.draftData = {
        firstName: this.form.get('firstName')?.value || 'Huy',
        middleName: this.form.get('middleName')?.value || '',
        lastName: this.form.get('lastName')?.value || 'Tran',
      };
    } else {
      this.editValue = this.form.get(field)?.value || '';
    }
  }
  saveName() {
    this.form.patchValue({
      firstName: this.draftData.firstName,
      middleName: this.draftData.middleName,
      lastName: this.draftData.lastName,
    });

    if (!this.profileUpdate) {
      this.profileUpdate = {
        username: this.username || '',
      } as UserProfileUpdateRequest;
    }
    this.profileUpdate.firstName = this.draftData.firstName;
    this.profileUpdate.lastName =
      this.draftData.lastName + ' ' + this.draftData.middleName;
    this.editing = null;
  }
  saveField(fieldName: string) {
    if (!this.editing) return;
    this.form.patchValue({ [this.editing]: this.editValue });
    if (!this.profileUpdate) {
      this.profileUpdate = {
        username: this.username || '',
      } as UserProfileUpdateRequest;
    }
    (this.profileUpdate as any)[this.editing] = this.editValue;
    console.log(`Updated ${this.editing}:`, this.editValue);
    this.editing = null;
  }

  cancel() {
    this.editing = null;
  }

  private updateProfile(avatarUri: string) {
    this.profileUpdate = {
      username: this.username || '',
      firstName: this.form.value.firstName,
      lastName: this.form.value.lastName,
      dob: this.form.value.dob,
      address: this.form.value.address,
      gender: this.form.value.gender,
      avatarUri,
    };

    this.profileService.updateProfile(this.profileUpdate).subscribe({
      next: (response) => {
        // this.toastService.show(
        //   'Cập nhật profile thành công!',
        //   'Success',
        //   ToastType.Success
        // );

        this.form.patchValue({
          ...this.profileUpdate,
          avatarUri: avatarUri,
        });

        this.profile = {
          ...this.profile,
          ...this.profileUpdate,
        } as ProfileModel;
        if (this.profileUpdate) {
          this.profileService.updateProfile(this.profileUpdate).subscribe({
            next: () => {
              this.toastService.show(
                'Profile updated successfully!',
                'Success',
                ToastType.Success
              );

              this.secureStorage.setItem('profile', this.profileUpdate);
              if (this.profileUpdate) {
                const updatedUser: UserProfile = {
                  firstName: this.profileUpdate.firstName,
                  lastName: this.profileUpdate.lastName,
                  dob: this.profileUpdate.dob,
                  address: this.profileUpdate.address,
                  gender: this.profileUpdate.gender,
                  avatarUri: this.profileUpdate.avatarUri,
                  accountId: '',
                  profileId: '',
                };

                this.userContext.setUser(updatedUser);
              }
              this.isUploadImage = false;
              this.selectedImage = null;
              this.imagePreview = null;
            },
            error: (err) => {
              this.toastService.show(
                'Profile update failed!',
                'Error',
                ToastType.Error
              );
              console.error(err);
            },
          });
        }
      },
      error: (err) => {
        console.error('Profile update error:', err);
        this.toastService.show(
          'Cập nhật profile thất bại!',
          'Error',
          ToastType.Error
        );
      },
    });
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.cancel();
    }
  }
}
