export enum ToastType {
  Success = 'success',
  Info = 'info',
  Warning = 'warning',
  Error = 'error',
  Schedule = 'schedule',
  Session = 'session',
}

export interface ToastAction {
  label: string;
  action: () => void;
  primary?: boolean;
}

export interface ToastData {
  id: string;
  title?: string;
  message: string;
  description?: string;
  type: ToastType;
  duration?: number;
  showCountdown?: boolean;
  actions?: ToastAction[];
}
