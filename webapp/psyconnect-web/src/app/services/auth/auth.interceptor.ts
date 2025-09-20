import {
    HttpErrorResponse,
    HttpEvent,
    HttpHandler,
    HttpInterceptor,
    HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { catchError, filter, switchMap, take, tap } from 'rxjs/operators';
import { SecureStorageService } from '../../encrypt/secure';
import { Auth } from './auth';
import { AuthService } from './auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private isRefreshing = false;
  private tokenSubject: BehaviorSubject<string | null> = new BehaviorSubject<
    string | null
  >(null);

  constructor(
    private auth: AuthService,
    private secureStorage: SecureStorageService,
    private authService: Auth
  ) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    const token = this.auth.getToken();
    let cloned = req;
    if (token) {
      cloned = this.addTokenHeader(req, token);
    }

    return next.handle(cloned).pipe(
      tap({
        next: () => console.log('Request success'),
        error: (err) => console.log('Request error tapped', err),
      }),
      catchError((error: HttpErrorResponse) => {
        console.log('CatchError triggered', error);
        if (error.status === 500) {
          return this.handle500Error(cloned, next);
        }
        return throwError(() => error);
      })
    );
  }

  private addTokenHeader(
    request: HttpRequest<any>,
    token: string
  ): HttpRequest<any> {
    // clone đầy đủ các options để Angular nhận đúng lỗi
    return request.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
      reportProgress: request.reportProgress,
      responseType: request.responseType,
      withCredentials: request.withCredentials,
      params: request.params,
      body: request.body,
    });
  }

  private handle500Error(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    if (!this.isRefreshing) {
      this.isRefreshing = true;
      this.tokenSubject.next(null);

      const username = this.secureStorage.getItem<string>('username');
      if (!username) {
        this.isRefreshing = false;
        this.auth.logout();
        return throwError(() => new Error('No username in storage'));
      }

      return this.authService.refreshToken(username).pipe(
        switchMap((newToken: string) => {
          this.isRefreshing = false;
          this.secureStorage.setItem('access_token', newToken);
          this.tokenSubject.next(newToken);
          return next.handle(this.addTokenHeader(request, newToken));
        }),
        catchError((err) => {
          this.isRefreshing = false;
          this.auth.logout();
          return throwError(() => err);
        })
      );
    } else {
      // Nếu token đang refresh, chờ token mới
      return this.tokenSubject.pipe(
        filter((token) => token != null),
        take(1),
        switchMap((token) => next.handle(this.addTokenHeader(request, token!)))
      );
    }
  }
}
