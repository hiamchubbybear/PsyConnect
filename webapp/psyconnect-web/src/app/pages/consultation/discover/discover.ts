import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { SwipeDeckComponent } from '../../swipe/swipe-deck';

@Component({
  selector: 'consultation-discover',
  imports: [CommonModule, SwipeDeckComponent],
  templateUrl: './discover.html',
  styleUrl: './discover.scss',
})
export class Discover {}
