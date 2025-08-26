import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AppRoutingModule } from '../../app.routes';
import { AccountSidebarComponent } from './account-sidebar/account-sidebar';

@Component({
  selector: 'app-update-account',
  standalone: true,
  imports: [
    RouterModule,
    CommonModule,
    ReactiveFormsModule,
    AppRoutingModule,
AccountSidebarComponent,
  ],
  templateUrl: './account.html',
  styleUrls: ['./account.scss', '../login/login.scss'],
})
export class UpdateAccountComponent {
  updateForm: FormGroup;
  shakeErrors = false;

  constructor(private fb: FormBuilder) {
    this.updateForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(6)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.minLength(6)]],
    });
  }

  onSubmit() {
    if (this.updateForm.invalid) {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 300);
      return;
    }
    console.log('Updated:', this.updateForm.value);
  }

  onCancel() {
    this.updateForm.reset();
  }
}
