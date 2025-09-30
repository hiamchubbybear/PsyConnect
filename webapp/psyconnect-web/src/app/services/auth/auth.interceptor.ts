import {
    HttpErrorResponse,
    HttpEvent,
    HttpHandlerFn,
    HttpInterceptorFn,
    HttpRequest,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { BehaviorSubject, Observable, throwError, timer } from 'rxjs';
import {
    catchError,
    delayWhen,
    filter,
    retryWhen,
    scan,
    switchMap,
    take,
} from 'rxjs/operators';
import { SecureStorageService } from '../../encrypt/secure';
import { Auth } from './auth';

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<any>,
  next: HttpHandlerFn
): Observable<HttpEvent<any>> => {
  const authService = inject(Auth);
  const secureStorage = inject(SecureStorageService);

  const isRefreshing = { value: false };
  const tokenSubject = new BehaviorSubject<string | null>(null);
  const maxRetry = 3;

  const token = secureStorage.getItem<string>('access_token');
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

        console.log(
          `[AuthInterceptor] Caught error after retries: ${err.status} ${err.url}`
        );

        if (err.status === 500) {
          const username = secureStorage.getItem<string>('username');
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

      return authService.refreshToken(username).pipe(
        switchMap((newToken) => {
          console.log('[AuthInterceptor] New token received:', newToken);
          isRefreshing.value = false;
          secureStorage.setItem('access_token', newToken);
          tokenSubject.next(newToken);

          return next(
            req.clone({ setHeaders: { Authorization: `Bearer ${newToken}` } })
          ).pipe(
            retryWhen((errors) =>
              errors.pipe(
                scan((count, err) => {
                  if (count >= maxRetry) {
                    console.error(
                      '[AuthInterceptor] Max retries reached after refresh, logout'
                    );
                    authService.logout();
                    throw err;
                  }
                  console.warn(
                    `[AuthInterceptor] Retry after refresh #${count + 1}`
                  );
                  return count + 1;
                }, 0),
                delayWhen(() => timer(1000))
              )
            )
          );
        }),
        catchError((err) => {
          console.error(
            '[AuthInterceptor] Refresh token failed, logging out',
            err
          );
          isRefreshing.value = false;
          authService.logout();
          return throwError(() => err);
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
