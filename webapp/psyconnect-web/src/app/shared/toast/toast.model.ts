export enum ToastType {
  Success = 'success',
  Info = 'info',
  Warning = 'warning',
  Error = 'error',
}

export interface ToastData {
  id: string;
  title: string;
  message: string;
  type: ToastType;
  duration?: number;
}
