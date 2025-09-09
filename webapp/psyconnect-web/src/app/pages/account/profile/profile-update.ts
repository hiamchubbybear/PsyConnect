import { CommonModule } from '@angular/common';
import {
    ChangeDetectorRef,
    Component,
    HostListener,
    OnInit,
} from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    FormsModule,
    ReactiveFormsModule,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { map, Observable, of } from 'rxjs';
import { SingleButton } from '../../../components/single-button/single-button';
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
import { ProfileModel } from './profile-model';
import { ProfileOverlayComponent } from './profile-overlay';
@Component({
  selector: 'app-profile-section',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    ProfileOverlayComponent,
    SingleButton,
  ],
  templateUrl: './profile-update.html',
  styleUrls: ['./profile-update.scss'],
})
export class ProfileSectionComponent implements OnInit {
  private initialProfile: any;
  selectedImage: File | null = null;
  imagePreview: string | null = null;
  username: string | null = '';
  isUploadImage?: boolean | false;
  avatarUriOrigin: string = '';
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
      this.updateProfile(avatarUri || this.avatarUriOrigin);
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
    private userContext: UserContextService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      firstName: [''],
      lastName: [''],
      dob: [''],
      address: [''],
      gender: [''],
      avatarUri: [''],
    });

    this.username = this.secureStorage.getItem('username');

    this.profileService.getProfile().subscribe((res) => {
      if (res?.data) {
        this.profile = res.data;

        this.form.patchValue(this.profile);

        this.initialProfile = { ...this.profile };

        this.avatarUriOrigin = res.data.avatarUri;
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

    this.editing = null;
    this.cdr.detectChanges();
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
        this.toastService.show(
          'Cập nhật profile thành công!',
          'Success',
          ToastType.Success
        );
        console.log('Cập nhật thành công', response);

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
  get canConfirm(): boolean {
    if (this.isUploadImage) return true;
    if (!this.initialProfile) return false;

    if (this.profileUpdate) {
      return Object.keys(this.profileUpdate).some((key) => {
        const newVal = (this.profileUpdate as any)[key];
        const oldVal = (this.initialProfile as any)[key];

        return newVal !== oldVal;
      });
    }
    return false;
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.cancel();
    }
  }
}
