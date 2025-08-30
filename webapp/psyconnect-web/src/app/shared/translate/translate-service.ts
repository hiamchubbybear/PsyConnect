import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface LanguageOption {
  code: string;
  name: string;
  flag: string;
}

@Injectable({
  providedIn: 'root',
})
export class TranslationService {
  public readonly availableLanguages: LanguageOption[] = [
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'vi', name: 'Tiếng Việt', flag: '🇻🇳' },
  ];

  private currentLanguageSubject = new BehaviorSubject<string>('en');
  public currentLanguage$ = this.currentLanguageSubject.asObservable();

  constructor(private translateService: TranslateService) {
    this.initializeTranslation();
  }

  private initializeTranslation(): void {
    this.translateService.setDefaultLang('en');

    const savedLanguage = this.getSavedLanguage();
    this.setLanguage(savedLanguage);
  }

  private getSavedLanguage(): string {
    const savedLang = localStorage.getItem('app-language');
    if (savedLang && this.isLanguageSupported(savedLang)) {
      return savedLang;
    }

    const browserLang = this.translateService.getBrowserLang();
    if (browserLang && this.isLanguageSupported(browserLang)) {
      return browserLang;
    }

    return 'en';
  }

  private isLanguageSupported(lang: string): boolean {
    return this.availableLanguages.some((l) => l.code === lang);
  }

  public setLanguage(language: string): void {
    if (!this.isLanguageSupported(language)) {
      console.warn(`Language ${language} is not supported`);
      return;
    }

    this.translateService.use(language);
    this.currentLanguageSubject.next(language);

    localStorage.setItem('app-language', language);

    document.documentElement.lang = language;
  }

  public getCurrentLanguage(): string {
    return this.currentLanguageSubject.value;
  }

  public translate(key: string, params?: any): Observable<string> {
    return this.translateService.get(key, params);
  }

  public instant(key: string, params?: any): string {
    return this.translateService.instant(key, params);
  }

  public translateMultiple(keys: string[], params?: any): Observable<any> {
    return this.translateService.get(keys, params);
  }

  public switchLanguage(): void {
    const currentLang = this.getCurrentLanguage();
    const nextLang = currentLang === 'en' ? 'vi' : 'en';
    this.setLanguage(nextLang);
  }

  public getCurrentLanguageInfo(): LanguageOption {
    const currentLang = this.getCurrentLanguage();
    return (
      this.availableLanguages.find((l) => l.code === currentLang) ||
      this.availableLanguages[0]
    );
  }

  public reloadTranslations(): void {
    const currentLang = this.getCurrentLanguage();
    this.translateService.reloadLang(currentLang);
  }

  public setTranslations(lang: string, translations: any): void {
    this.translateService.setTranslation(lang, translations, true);
  }

  public waitForTranslations(): Observable<boolean> {
    return this.translateService.onTranslationChange.pipe(map(() => true));
  }
}
