import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DriveSidebar } from './sidebar/sidebar';

@Component({
  selector: 'app-consultation',
  imports: [DriveSidebar, CommonModule, RouterModule],
  templateUrl: './drive-management.html',
  styleUrl: '../account/account-update.scss',
  standalone: true,
})
export class DriveManagement {}
