import { Overlay } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { Injectable, Injector } from '@angular/core';
import { ConfirmDialogComponent } from '../../components/confirm-dialog/confirm-dialog';

@Injectable({ providedIn: 'root' })
export class ConfirmDialogService {
  constructor(private overlay: Overlay, private injector: Injector) {}

  open(title: string, message: string): Promise<boolean> {
    return new Promise((resolve) => {
      const overlayRef = this.overlay.create({
        hasBackdrop: true,
        backdropClass: 'cdk-overlay-dark-backdrop',
        positionStrategy: this.overlay
          .position()
          .global()
          .centerHorizontally()
          .centerVertically(),
      });

      const dialogPortal = new ComponentPortal(ConfirmDialogComponent);
      const componentRef = overlayRef.attach(dialogPortal);

      componentRef.instance.title = title;
      componentRef.instance.message = message;

      const sub1 = componentRef.instance.confirm.subscribe(() => {
        resolve(true);
        overlayRef.dispose();
        sub1.unsubscribe();
        sub2.unsubscribe();
      });

      const sub2 = componentRef.instance.cancel.subscribe(() => {
        resolve(false);
        overlayRef.dispose();
        sub1.unsubscribe();
        sub2.unsubscribe();
      });

      overlayRef.backdropClick().subscribe(() => {
        resolve(false);
        overlayRef.dispose();
        sub1.unsubscribe();
        sub2.unsubscribe();
      });
    });
  }
}
