import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Discover } from './discover/discover';
import { ConsultationSidebar } from './sidebar/sidebar';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-consultation',
  imports: [ConsultationSidebar, CommonModule, RouterModule],
  templateUrl: './consultation.html',
  styleUrl: '../account/account-update.scss',
  standalone: true,
})
export class Consultation {}
