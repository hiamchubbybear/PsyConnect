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
import { ToastType } from '../../shared/toast/toast-type';
import { ToastService } from '../../shared/toast/toast.service';
import { LoaderService } from '../loader/loader';
import { Auth } from './auth';

const usernameKey = environment.usernameKey;
const ACCESSTOKEN_KEY = environment.accessTokenKey;
export const SKIP_AUTH = new HttpContextToken(() => false);

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<any>,
  next: HttpHandlerFn
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
    console.log('[AuthInterceptor] Skipping auth for:', req.url);
    return next(req);
  }

  const token = secureStorage.getItem<string>(ACCESSTOKEN_KEY);
  let authReq = req;
  if (token) {
    authReq = req.clone({ setHeaders: { Authorization: `Bearer ${token}` } });
    console.log('[AuthInterceptor] Adding token to request header');
  }

  return handleRequest(authReq, next);

  function handleRequest(
    req: HttpRequest<any>,
    next: HttpHandlerFn
  ): Observable<HttpEvent<any>> {
    return next(req).pipe(
      retryWhen((errors) =>
        errors.pipe(
          scan((count, err) => {
            if (!(err instanceof HttpErrorResponse) || err.status !== 500)
              throw err;
            if (count >= maxRetry) throw err;
            console.warn(`[AuthInterceptor] Retry #${count + 1} due to 500`);
            return count + 1;
          }, 0),
          delayWhen(() => timer(1000))
        )
      ),
      catchError((err) => {
        if (!(err instanceof HttpErrorResponse)) return throwError(() => err);
        if (err.status === 500) {
          const username = secureStorage.getItem<string>(usernameKey);
          if (!username) {
            console.warn('[AuthInterceptor] No username in storage, logout');
            authService.logout();
            return throwError(() => new Error('No username in storage'));
          }
          return refreshTokenAndRetry(req, next, username);
        }

        return throwError(() => err);
      })
    );
  }

  function refreshTokenAndRetry(
    req: HttpRequest<any>,
    next: HttpHandlerFn,
    username: string
  ): Observable<HttpEvent<any>> {
    if (!isRefreshing.value) {
      console.log('[AuthInterceptor] Starting token refresh...');
      isRefreshing.value = true;
      tokenSubject.next(null);
      loaderService.show();

      return authService.refreshToken(username).pipe(
        switchMap((newToken) => {
          console.log('[AuthInterceptor] New token received:', newToken);
          isRefreshing.value = false;
          secureStorage.setItem(ACCESSTOKEN_KEY, newToken);
          tokenSubject.next(newToken);

          toastService.show('Token refreshed', 'Success', ToastType.Success);

          return next(
            req.clone({ setHeaders: { Authorization: `Bearer ${newToken}` } })
          ).pipe(
            catchError((err) => {
              console.error(
                '[AuthInterceptor] Request failed after refresh, logout'
              );
              toastService.show(
                'Phiên đăng nhập của bạn đã hết hạn',
                'Failed',
                ToastType.Error
              );
              authService.logout();
              return throwError(() => err);
            })
          );
        }),
        catchError((err) => {
          console.error(
            '[AuthInterceptor] Refresh token failed, logging out',
            err
          );
          isRefreshing.value = false;
          toastService.show(
            'Phiên đăng nhập của bạn đã hết hạn',
            'Failed',
            ToastType.Error
          );
          authService.logout();
          router.navigate(['/auth/login']);
          return throwError(() => err);
        }),
        finalize(() => {
          loaderService.hide();
        })
      );
    } else {
      console.log('[AuthInterceptor] Waiting for ongoing token refresh...');
      return tokenSubject.pipe(
        filter((t) => t != null),
        take(1),
        switchMap((t) =>
          next(req.clone({ setHeaders: { Authorization: `Bearer ${t!}` } }))
        )
      );
    }
  }
};
