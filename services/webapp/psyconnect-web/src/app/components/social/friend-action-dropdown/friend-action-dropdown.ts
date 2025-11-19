import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';
import { Connection } from '../../../models/connection.model';
import { ConfirmDialogComponent } from '../../confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-friend-action-dropdown',
  standalone: true,
  imports: [CommonModule, LucideAngularModule],
  templateUrl: './friend-action-dropdown.html',
  styleUrls: ['./friend-action-dropdown.scss'],
})
export class FriendActionDropdownComponent {
  @Input() visible = false;
  @Input() cardType: 'request' | 'friend' | 'suggestion' = 'friend';
  @Input() friend!: Connection;

  @Output() unfriend = new EventEmitter<void>();
  @Output() viewProfile = new EventEmitter<void>();
  @Output() accept = new EventEmitter<void>();
  @Output() decline = new EventEmitter<void>();
  @Output() message = new EventEmitter<void>();
  @Output() addFriend = new EventEmitter<void>();

  private overlayRef?: OverlayRef;

  constructor(private overlay: Overlay) {}

  openConfirmOverlay() {
    this.overlayRef = this.overlay.create({
      hasBackdrop: true,
      backdropClass: 'cdk-overlay-dark-backdrop',
      positionStrategy: this.overlay
        .position()
        .global()
        .centerHorizontally()
        .centerVertically(),
    });

    const portal = new ComponentPortal(ConfirmDialogComponent);
    const dialogRef = this.overlayRef.attach(portal);

    dialogRef.instance.title = 'Xóa bạn bè';
    dialogRef.instance.message = `Bạn có chắc chắn muốn xóa ${this.friend.name}?`;

    dialogRef.instance.confirm.subscribe(() => {
      this.unfriend.emit();
      this.overlayRef?.dispose();
    });

    dialogRef.instance.cancel.subscribe(() => {
      this.overlayRef?.dispose();
    });
  }
}
