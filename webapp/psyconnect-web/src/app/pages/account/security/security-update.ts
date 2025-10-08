import { CommonModule } from '@angular/common';
import { Component, HostListener, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { environment } from '../../../../environments/environment';
import { ProfileFieldDropdownComponent } from '../../../components/field-row/field-row';
import { SecureStorageService } from '../../../encrypt/secure';
import { PasswordService } from '../../../services/auth/password.service';
import { Profile } from '../../../services/profile/profile';
import { ToastType } from '../../../shared/toast/toast.model';
import { ToastService } from '../../../shared/toast/toast.service';

interface PasswordData {
  current: string;
  new: string;
  confirm: string;
}

interface SecurityMethods {
  email: string;
  phone: string;
  key: string;
}

@Component({
  selector: 'app-security-section',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    ProfileFieldDropdownComponent,
    TranslateModule,
  ],
  templateUrl: './security-update.html',
  styleUrls: ['./security-update.scss'],
})
export class SecuritySectionComponent implements OnInit {
  PROFILE_KEY = environment.profileKey;
  saveName() {
    throw new Error('Method not implemented.');
  }
  confirmChanges() {
    throw new Error('Method not implemented.');
  }
  securityItems = [
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.ChangePassword',
      key: 'changePassword',
    },
    { label: 'ACCOUNT_MANAGEMENT.Security.Items.TwoFA', key: '2fa' },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.RecentLogins',
      key: 'recentLogins',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.SecurityMethods',
      key: 'securityMethods',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.LoggedDevices',
      key: 'loggedDevices',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.SuspiciousActivity',
      key: 'suspiciousActivity',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.LinkedEmails',
      key: 'linkedEmails',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.ThirdPartyApps',
      key: 'thirdPartyApps',
    },
  ];

  securityCheckItems = [
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.LoggedDevices',
      key: 'loggedDevices',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.LinkedEmails',
      key: 'linkedEmails',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.ThirdPartyApps',
      key: 'thirdPartyApps',
    },
  ];

  editing: string | null = null;
  passwordData: PasswordData = {
    current: '',
    new: '',
    confirm: '',
  };

  securityMethods: SecurityMethods = {
    email: '',
    phone: '',
    key: '',
  };

  newEmail = '';
  newPhone = '';
  newKey = '';

  recentLogins: Array<{ device: string; time: Date; location: string }> = [];
  loggedDevices: Array<{ name: string; lastActive: Date }> = [];
  suspiciousActivities: Array<{ message: string; time: Date }> = [];
  linkedEmails: Array<{ address: string }> = [];
  thirdPartyApps: Array<{ name: string; permissions: string[] }> = [];

  form!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private toastService: ToastService,
    private passwordService: PasswordService,
    private secureStorage: SecureStorageService
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      passwordCurrent: [''],
      passwordNew: [''],
      passwordConfirm: [''],
      email: [''],
      phone: [''],
      key: [''],
    });

    this.recentLogins = [
      { device: 'Chrome - Windows', time: new Date(), location: 'Hà Nội' },
      { device: 'Safari - iPhone', time: new Date(), location: 'HCM' },
    ];

    this.loggedDevices = [
      { name: 'Laptop Huy', lastActive: new Date() },
      { name: 'iPhone Huy', lastActive: new Date() },
    ];

    this.suspiciousActivities = [
      { message: 'Đăng nhập từ IP lạ', time: new Date() },
    ];

    this.linkedEmails = [{ address: 'huy@example.com' }];
    this.thirdPartyApps = [{ name: 'App XYZ', permissions: ['Đọc', 'Ghi'] }];
  }

  openEdit(field: string) {
    this.editing = field;
    switch (field) {
      case 'changePassword':
        this.passwordResetEmail = this.userEmail;
        break;
      case 'securityMethods':
        this.securityMethods = { ...this.securityMethods };
        break;
    }
  }

  passwordResetEmail = '';
  submittingReset = false;

  get userEmail(): string {
    return (
      this.secureStorage.getItem<Profile>(this.PROFILE_KEY)?.getEmail() ||
      'user@example.com'
    );
  }

  sendPasswordReset() {
    if (!this.passwordResetEmail) {
      this.toastService.show(
        'Không có email để gửi reset!',
        'Error',
        ToastType.Error
      );
      return;
    }

    this.submittingReset = true;

    this.passwordService
      .requestReset({ email: this.passwordResetEmail })
      .subscribe({
        next: (res) => {
          this.submittingReset = false;
          if (res.code === 200) {
            this.toastService.show(
              `Email reset mật khẩu đã được gửi tới ${this.passwordResetEmail}`,
              'Success',
              ToastType.Success
            );
            this.editing = null;
          } else {
            this.toastService.show(
              res.message || 'Có lỗi xảy ra!',
              'Error',
              ToastType.Error
            );
          }
        },
        error: (err) => {
          this.submittingReset = false;
          this.toastService.show(
            err.error?.message || 'Có lỗi xảy ra!',
            'Error',
            ToastType.Error
          );
        },
      });
  }

  saveSecurityMethods() {
    this.toastService.show(
      'Phương thức bảo mật đã được cập nhật!',
      'Success',
      ToastType.Success
    );
    this.editing = null;
  }

  enable2FA(method: 'authApp' | 'sms' | 'securityKey') {
    this.toastService.show(
      `2FA được bật bằng ${method}`,
      'Success',
      ToastType.Success
    );
    this.editing = null;
  }

  logoutDevice(device: { name: string; lastActive: Date }) {
    this.toastService.show(
      `Đã đăng xuất thiết bị ${device.name}`,
      'Success',
      ToastType.Success
    );
    this.loggedDevices = this.loggedDevices.filter((d) => d !== device);
  }

  resolveAlert(alert: { message: string; time: Date }) {
    this.toastService.show(
      `Đã xác minh: ${alert.message}`,
      'Success',
      ToastType.Success
    );
    this.suspiciousActivities = this.suspiciousActivities.filter(
      (a) => a !== alert
    );
  }

  addEmail() {
    if (this.newEmail.trim() !== '') {
      this.linkedEmails.push({ address: this.newEmail });
      this.newEmail = '';
      this.toastService.show(
        'Email đã được thêm!',
        'Success',
        ToastType.Success
      );
    }
  }

  removeEmail(email: { address: string }) {
    this.linkedEmails = this.linkedEmails.filter((e) => e !== email);
    this.toastService.show('Email đã bị xóa!', 'Success', ToastType.Success);
  }

  revokeApp(app: { name: string; permissions: string[] }) {
    this.thirdPartyApps = this.thirdPartyApps.filter((a) => a !== app);
    this.toastService.show(
      `Quyền của ${app.name} đã bị thu hồi`,
      'Success',
      ToastType.Success
    );
  }

  cancel() {
    this.editing = null;
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.cancel();
    }
  }
  getFieldDisplay(key: string): string {
    switch (key) {
      case 'changePassword':
        return '••••••••';
      case '2fa':
        return this.editing === '2fa' ? 'Enabled' : 'Disabled';
      case 'recentLogins':
        return `${this.recentLogins.length} logins`;
      case 'loggedDevices':
        return `${this.loggedDevices.length} devices`;
      case 'suspiciousActivity':
        return `${this.suspiciousActivities.length} alerts`;
      case 'linkedEmails':
        return this.linkedEmails.map((e) => e.address).join(', ');
      case 'thirdPartyApps':
        return this.thirdPartyApps.map((a) => a.name).join(', ');
      case 'securityMethods':
        return '••••';
      default:
        return '-';
    }
  }
  closeOverlay() {
    this.editing = null;
  }
  saveOverlay() {
    switch (this.editing) {
      case 'changePassword':
        this.savePassword();
        break;
      case 'securityMethods':
        this.saveSecurityMethods();
        break;
      case '2fa':
        this.toastService.show(
          '2FA đã được cập nhật!',
          'Success',
          ToastType.Success
        );
        break;
      default:
        this.toastService.show(
          'Cập nhật thành công!',
          'Success',
          ToastType.Success
        );
        break;
    }
    this.closeOverlay();
  }
}
