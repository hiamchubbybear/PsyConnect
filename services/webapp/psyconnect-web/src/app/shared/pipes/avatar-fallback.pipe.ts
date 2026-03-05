import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'avatarFallback',
  standalone: true
})
export class AvatarFallbackPipe implements PipeTransform {
  
  private readonly fallbacks = [
    '/default_avatar/7249ed41cf4ab50c5d37c169786ff768.jpg',
    '/default_avatar/8e4cd15170685b5cef2ae84488a491b9.jpg',
    '/default_avatar/95c890dba4a37b89d0bf046693540b3e.jpg'
  ];

  
  transform(value: string | null | undefined, identifier?: string | null): string {
    
    if (!value || value.trim() === '' || this.isInvalidUrl(value)) {
      return this.getFallback(identifier);
    }
    return value;
  }

  private isInvalidUrl(url: string): boolean {
    const trimmed = url.trim().toLowerCase();
    
    if (['null', 'undefined', 'none', 'default'].includes(trimmed)) {
      return true;
    }
    
    if (!trimmed.startsWith('http') && !trimmed.startsWith('/') && !trimmed.startsWith('data:')) {
      return true;
    }
    
    const deadDomains = ['placeimg.com', 'placeholder.com', 'placekitten.com', 'lorempixel.com', 'example.com'];
    if (deadDomains.some(domain => trimmed.includes(domain))) {
      return true;
    }
    return false;
  }

  private getFallback(identifier?: string | null): string {

    
    if (!identifier) {
      return this.fallbacks[0];
    }

    
    let hash = 0;
    for (let i = 0; i < identifier.length; i++) {
        const char = identifier.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; 
    }

    const index = Math.abs(hash) % this.fallbacks.length;
    return this.fallbacks[index];
  }
}
