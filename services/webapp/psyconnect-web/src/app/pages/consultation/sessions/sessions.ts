import { Component } from '@angular/core';

import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'consultation-sessions',
  imports: [AvatarFallbackPipe],
  templateUrl: './sessions.html',
  styleUrl: './sessions.scss',
})
export class Sessions {}
