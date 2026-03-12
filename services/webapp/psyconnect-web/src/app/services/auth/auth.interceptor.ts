import {
  HttpContextToken,
  HttpErrorResponse,
  HttpEvent,
  HttpHandlerFn,
  HttpInterceptorFn,
  HttpRequest,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, throwError, timer } from 'rxjs';
import {
  catchError,
  delayWhen,
  filter,
  finalize,
  retryWhen,
  scan,
  switchMap,
  take,
} from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { ToastService } from '../../shared/toast/toast.service';
import { ToastType } from '../../shared/toast/toast.model';
import { LoaderService } from '../loader/loader';
import { Auth } from './auth';

const usernameKey = environment.usernameKey;
const ACCESSTOKEN_KEY = environment.accessTokenKey;
export const SKIP_AUTH = new HttpContextToken(() => false);

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<any>,
  next: HttpHandlerFn,
): Observable<HttpEvent<any>> => {
  const authService = inject(Auth);
  const secureStorage = inject(SecureStorageService);
  const loaderService = inject(LoaderService);
  const toastService = inject(ToastService);
  const isRefreshing = { value: false };
  const tokenSubject = new BehaviorSubject<string | null>(null);
  const maxRetry = 3;
  const router = inject(Router);
  if (req.context.get(SKIP_AUTH)) {
    return next(req);
  }
  const token = secureStorage.getItem<string>(ACCESSTOKEN_KEY);
  let authReq = req;
  if (token) {
    
    authReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`,
        roles: 'client', 
      },
    });
  }

  return handleRequest(authReq, next);

  function handleRequest(
    req: HttpRequest<any>,
    next: HttpHandlerFn,
  ): Observable<HttpEvent<any>> {
    return next(req).pipe(
      retryWhen((errors) =>
        errors.pipe(
          scan((count, err) => {
            if (!(err instanceof HttpErrorResponse) || err.status !== 500)
              throw err;
            if (count >= maxRetry) throw err;
            return count + 1;
          }, 0),
          delayWhen(() => timer(1000)),
        ),
      ),
      catchError((err) => {
        if (!(err instanceof HttpErrorResponse)) {
          return throwError(() => err);
        }
        const message = (err.error?.message || '').toLowerCase();
        const username = secureStorage.getItem<string>(usernameKey);
        if (!username) {
          if (err.status === 401 || err.status === 403) {
            toastService.show({
              title: 'Truy cập bị từ chối',
              message: 'Bạn không có quyền hoặc phiên đăng nhập đã hết hạn.',
              type: ToastType.Error,
            });
          }
          return throwError(() => new Error('No username in storage'));
        }
        if (err.status === 401) {
          return refreshTokenAndRetry(req, next, username);
        }

        
        
        

        if (err.status === 500) {
          console.error('[AuthInterceptor] 500 error (non-server):', message);
          return throwError(() => err);
        }
        return throwError(() => err);
      }),
    );
  }

  function refreshTokenAndRetry(
    req: HttpRequest<any>,
    next: HttpHandlerFn,
    username: string,
  ): Observable<HttpEvent<any>> {
    if (!isRefreshing.value) {
      console.log('[AuthInterceptor] Starting token refresh...');
      isRefreshing.value = true;
      tokenSubject.next(null);
      loaderService.show();

      return authService.refreshToken(username).pipe(
        switchMap((newToken) => {
          isRefreshing.value = false;
          secureStorage.setItem(ACCESSTOKEN_KEY, newToken);
          tokenSubject.next(newToken);

          toastService.show({
            title: 'TOAST.key_token_refreshed',
            message: 'TOAST.key_success',
            type: ToastType.Success,
          });

          return next(
            req.clone({ setHeaders: { Authorization: `Bearer ${newToken}` } }),
          ).pipe(
            catchError((err) => {
              
              
              if (err instanceof HttpErrorResponse && err.status === 401) {
                toastService.show({
                  title: 'Truy cập bị hạn chế',
                  message: 'Bạn không có quyền thực hiện hành động này.',
                  type: ToastType.Warning,
                });
              } else {
                toastService.show({
                  title: 'Lỗi không xác định',
                  message: 'Đã có lỗi xảy ra, vui lòng thử lại sau.',
                  type: ToastType.Error,
                });
              }
              
              return throwError(() => err);
            }),
          );
        }),
        catchError((err) => {
          isRefreshing.value = false;
          toastService.show({
            title: 'TOAST.key_session_expired',
            message: 'TOAST.key_failed',
            type: ToastType.Error,
          });
          authService.logout();
          router.navigate(['/auth/login']);
          return throwError(() => err);
        }),
        finalize(() => {
          loaderService.hide();
        }),
      );
    } else {
      return tokenSubject.pipe(
        filter((t) => t != null),
        take(1),
        switchMap((t) =>
          next(req.clone({ setHeaders: { Authorization: `Bearer ${t!}` } })),
        ),
      );
    }
  }
};
