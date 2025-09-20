import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ConsultationSidebar } from './payment-sidebar/payment-sidebar';
import { RouterModule } from "@angular/router";

@Component({
  selector: 'app-payment-management',
  standalone: true,
  imports: [ConsultationSidebar, CommonModule, RouterModule],
  templateUrl: './payment-management.html',
  styleUrl: '../account/account-update.scss',
})
export class PaymentManagement {}
