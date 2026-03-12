import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
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
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { map, Observable, of } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { SingleButton } from '../../../components/single-button/single-button';
import { SecureStorageService } from '../../../encrypt/secure';
import { UserProfileUpdateRequest } from '../../../models/profile';
import { CloudinaryService } from '../../../services/cloudinary/cloudinary.service';
import { LoaderService } from '../../../services/loader/loader';
import { Profile } from '../../../services/profile/profile';
import {
  UserContextService,
  UserProfile,
} from '../../../services/profile/profile-service';
import { ToastType } from '../../../shared/toast/toast.model';
import { ToastService } from '../../../shared/toast/toast.service';
import { FieldRowComponent } from '../../../shared/ui-atoms/field-row/field-row.component';
import { ProfileModel } from './profile-model';
import { ProfileOverlayComponent } from './profile-overlay';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
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
    FieldRowComponent,
    AvatarFallbackPipe,
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
  editValue: string = '';
  addressSuggestions: string[] = [];
  profileUpdate: UserProfileUpdateRequest | undefined;
  PROFILE_KEY = environment.profileKey;
  USERNAME_KEY = environment.usernameKey;
  expandedField: string | null = null;

  toggleExpand(field: string): void {
    if (this.expandedField === field) {
      this.expandedField = null;
    } else {
      this.expandedField = field;
    }
  }

  onExpandClick(event: Event, field: string): void {
    event.stopPropagation();
    this.expandedField = null;
    this.openEdit(field);
  }

  onImageSelected(event: any) {
    const file = event.target.files[0];
    if (file) {
      if (!file.type.startsWith('image/')) {
        this.toastService.show({
          message: 'TOAST.invalid_image',
          title: 'TOAST.error',
          type: ToastType.Error,
        });
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        this.toastService.show({
          message: 'TOAST.image_too_large',
          title: 'TOAST.error',
          type: ToastType.Error,
        });

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

  async confirmChanges(): Promise<void> {
    this.loaderService.show();

    try {
      let avatarUri: string | undefined;

      if (this.isUploadImage && this.selectedImage) {
        const uploadedUrl = await this.cloudinaryService.uploadImage(
          this.selectedImage,
          this.username || '',
        );

        if (!uploadedUrl || uploadedUrl.trim() === '') {
          throw new Error('Upload failed - no URL returned');
        }

        avatarUri = uploadedUrl;
      }

      this.updateProfile(avatarUri);
    } catch (err) {
      console.error('Upload error:', err);

      this.toastService.show({
        message: 'TOAST.upload_failed',
        title: 'TOAST.error',
        type: ToastType.Error,
      });

      this.resetUploadStates();
    } finally {
      this.loaderService.hide();
    }
  }
  private resetUploadStates(): void {
    this.isUploadImage = false;
    this.selectedImage = null;
    this.imagePreview = null;
    this.editing = null;
  }
  form!: FormGroup;
  profile: ProfileModel | null = null;

  editing: string | null = null;
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
    private cdr: ChangeDetectorRef,
    private http: HttpClient,
    private translate: TranslateService,
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

    this.username = this.secureStorage.getItem(this.USERNAME_KEY);

    this.profileService.getProfile().subscribe((res) => {
      if (res?.data) {
        this.profile = res.data;
        console.log(res.data);

        this.form.patchValue(this.profile);

        this.initialProfile = {
          ...this.profile,
          profile_id: this.profile.profileId,
        };

        this.avatarUriOrigin = res.data.avatarUri;
      }
    });
  }
  private prepareCompleteProfileUpdate(
    avatarUri?: string,
  ): UserProfileUpdateRequest {
    const currentFormValues = this.form.value;

    const completeProfileData: UserProfileUpdateRequest = {
      username: this.username || '',
      firstName: currentFormValues.firstName || this.profile?.firstName || '',
      lastName: currentFormValues.lastName || this.profile?.lastName || '',
      dob: currentFormValues.dob || this.profile?.dob || '',
      address: currentFormValues.address || this.profile?.address || '',
      gender: currentFormValues.gender || this.profile?.gender || '',
      avatarUri:
        avatarUri ||
        currentFormValues.avatarUri ||
        this.profile?.avatarUri ||
        '',
    };

    return completeProfileData;
  }

  private updateProfile(avatarUri?: string): void {
    const profileUpdateData = this.prepareCompleteProfileUpdate(avatarUri);

    this.profileService.updateProfile(profileUpdateData).subscribe({
      next: (response) => {
        this.toastService.show({
          message: 'TOAST.profile_update_success',
          title: 'TOAST.success',
          type: ToastType.Success,
        });
        this.profile = {
          ...this.profile,
          ...profileUpdateData,
        } as ProfileModel;
        this.form.patchValue(profileUpdateData);
        this.initialProfile = { ...profileUpdateData };
        this.secureStorage.setItem(this.PROFILE_KEY, profileUpdateData);
        const updatedUser: UserProfile = {
          firstName: profileUpdateData.firstName,
          lastName: profileUpdateData.lastName,
          dob: profileUpdateData.dob,
          address: profileUpdateData.address,
          gender: profileUpdateData.gender,
          avatarUri: profileUpdateData.avatarUri,
          accountId: this.profile?.accountId || '',
          profileId: this.profile?.profileId || '',
        };
        this.userContext.setUser(updatedUser);
        this.resetUploadStates();
        this.profileUpdate = undefined;
      },
      error: (err) => {
        this.toastService.show({
          message: 'TOAST.profile_update_failed',
          title: 'TOAST.error',
          type: ToastType.Error,
        });
      },
    });
  }

  getProfile(): Observable<ProfileModel | null> {
    const profileJson = this.secureStorage.getItem(this.PROFILE_KEY);
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
        firstName: this.form.get('firstName')?.value || '',
        middleName: this.form.get('middleName')?.value || '',
        lastName: this.form.get('lastName')?.value || '',
      };
    } else {
      this.editValue = this.form.get(field)?.value || '';
    }
  }

  onAddressInput(query: string) {
    if (query && query.length > 2) {
      const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
        query,
      )}&format=json&limit=5`;
      this.http.get<any[]>(url).subscribe((data) => {
        this.addressSuggestions = data.map((item) => item.display_name);
      });
    } else {
      this.addressSuggestions = [];
    }
  }

  selectAddress(suggestion: string) {
    this.editValue = suggestion;
    this.addressSuggestions = [];
  }

  saveName(): void {
    this.form.patchValue({
      firstName: this.draftData.firstName,
      middleName: this.draftData.middleName,
      lastName: this.draftData.lastName,
    });
    if (!this.profileUpdate) {
      this.profileUpdate = this.prepareCompleteProfileUpdate();
    }
    this.profileUpdate.firstName = this.draftData.firstName;
    this.profileUpdate.lastName =
      this.draftData.lastName + ' ' + this.draftData.middleName;
    this.editing = null;
  }

  saveField(event: { field: string; value: string }): void {
    if (!this.editing) return;
    if (this.editing === 'name') {
      this.form.patchValue({
        firstName: this.draftData.firstName,
        middleName: this.draftData.middleName,
        lastName: this.draftData.lastName,
      });
    } else {
      this.form.patchValue({ [this.editing]: event.value });
    }
    if (!this.profileUpdate) {
      this.profileUpdate = this.prepareCompleteProfileUpdate();
    }
    if (this.editing === 'name') {
      this.profileUpdate.firstName = this.draftData.firstName;
      this.profileUpdate.lastName = this.draftData.lastName;
    } else {
      const field = this.editing as keyof UserProfileUpdateRequest;
      (this.profileUpdate as any)[field] = event.value;
    }

    this.editing = null;
    this.cdr.detectChanges();
  }

  private initProfileUpdate() {}
  cancel() {
    this.editing = null;
  }

  get canConfirm(): boolean {
    if (this.isUploadImage) return true;

    if (!this.initialProfile) return false;

    const currentFormValues = this.form.value;

    return Object.keys(currentFormValues).some((key) => {
      const currentValue = currentFormValues[key];
      const initialValue = (this.initialProfile as any)[key];
      return (
        currentValue !== initialValue &&
        currentValue !== null &&
        currentValue !== ''
      );
    });
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.cancel();
    }
  }
  copyToClipboard(value: string) {
    if (!value) return;
    navigator.clipboard.writeText(value).then(() => {
      console.log('Copied Profile ID:', value);
      this.toastService.show({
        message: 'TOAST.id_copied',
        title: 'TOAST.success',
        type: ToastType.Success,
      });
    });
  }

  getDisplayValue(field: string): string {
    if (!this.form?.value) return '';
    switch (field) {
      case 'name':
        return `${this.form.value.firstName || ''} ${
          this.form.value.middleName || ''
        } ${this.form.value.lastName || ''}`.trim();
      case 'dob':
        return this.form.value.dob
          ? new Date(this.form.value.dob).toLocaleDateString('vi-VN')
          : '-';
      case 'address':
        return this.form.value.address || '-';
      case 'gender':
        if (this.form.value.gender === 'male')
          return this.translate.instant('ACCOUNT_MANAGEMENT.Profile.Male');
        if (this.form.value.gender === 'female')
          return this.translate.instant('ACCOUNT_MANAGEMENT.Profile.Female');
        return this.translate.instant('ACCOUNT_MANAGEMENT.Profile.Other');
      default:
        return '';
    }
  }
}
