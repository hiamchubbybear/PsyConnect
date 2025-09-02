import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';

import { HttpClient, provideHttpClient, withFetch } from '@angular/common/http';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { routes } from './app.routes';

// Custom loader class đơn giản
export class CustomTranslateLoader implements TranslateLoader {
  constructor(private http: HttpClient) {}

  getTranslation(lang: string): Observable<any> {
    return this.http.get(`./assets/i18n/${lang}.json`);
  }
}

// Factory function
export function createTranslateLoader(http: HttpClient) {
  return new CustomTranslateLoader(http);
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes), // Router
    provideHttpClient(withFetch()), // HTTP Client
    provideAnimations(), // Angular animations
    importProvidersFrom(
      TranslateModule.forRoot({
        defaultLanguage: 'en', // hoặc 'vi'
        loader: {
          provide: TranslateLoader,
          useFactory: createTranslateLoader,
          deps: [HttpClient],
        },
      })
    ),
    // Nếu cần ReactiveFormsModule hoặc FormsModule thì importProvidersFrom tương tự
    // importProvidersFrom(ReactiveFormsModule),
    // importProvidersFrom(FormsModule),
  ],
};
