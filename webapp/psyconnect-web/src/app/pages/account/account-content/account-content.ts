import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'account-content',
  templateUrl: './account-content.html',
  styleUrls: ['./account-content.scss'] ,
  standalone :true,
  imports :  [ CommonModule ,ReactiveFormsModule

  ]
})
export class AccountContentComponent {

  accountForm: FormGroup;
  personalForm: FormGroup;
  securityForm: FormGroup;
  paymentForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.accountForm = this.fb.group({
      username: [''],
      email: ['']
    });

    this.personalForm = this.fb.group({
      fullname: [''],
      phone: ['']
    });

    this.securityForm = this.fb.group({
      currentPassword: [''],
      newPassword: [''],
      confirmPassword: ['']
    });

    this.paymentForm = this.fb.group({
      cardNumber: [''],
      expiry: [''],
      cvv: ['']
    });
  }

  saveForm(form: FormGroup, type: string) {
    console.log(`Saving ${type}:`, form.value);
  }
}
