import {
  ApplicationRef,
  ComponentRef,
  EnvironmentInjector,
  Injectable,
  createComponent,
} from '@angular/core';
import { ToastContainerComponent } from './toast-container';
import { ToastData, ToastType } from './toast.model';

@Injectable({ providedIn: 'root' })
export class ToastService {
  private containerRef?: ComponentRef<ToastContainerComponent>;

  constructor(
    private appRef: ApplicationRef,
    private environmentInjector: EnvironmentInjector
  ) {}

  show(message: string, title: string, type: ToastType, duration = 4000) {
    if (!this.containerRef) {
      this.containerRef = createComponent(ToastContainerComponent, {
        environmentInjector: this.environmentInjector,
      });
      this.appRef.attachView(this.containerRef.hostView);
      const domElem = (this.containerRef.hostView as any)
        .rootNodes[0] as HTMLElement;
      document.body.appendChild(domElem);
    }

    const { bgColor, accentColor } = this.resolveColors(type);
    const id = crypto.randomUUID();

    const toast: ToastData = {
      id,
      title,
      message,
      type,
      duration,
      bgColor,
      accentColor,
    };
    this.containerRef.instance.addToast(toast);
  }

  private resolveColors(type: ToastType) {
    switch (type) {
      case ToastType.Success:
        return { bgColor: '#f0fdf4', accentColor: '#22c55e' };
      case ToastType.Info:
        return { bgColor: '#f0f9ff', accentColor: '#3b82f6' };
      case ToastType.Warning:
        return { bgColor: '#fff7ed', accentColor: '#f97316' };
      case ToastType.Error:
        return { bgColor: '#fef2f2', accentColor: '#ef4444' };
    }
  }
}
