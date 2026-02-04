// Rebuild trigger
import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-start',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './start.html',
  styleUrl: './start.scss',
})
export class Start {}
