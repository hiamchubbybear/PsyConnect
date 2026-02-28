import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'avatarFallback',
  standalone: true
})
export class AvatarFallbackPipe implements PipeTransform {
  // Constant array pointing directly to the public/default_avatar images
  private readonly fallbacks = [
    '/default_avatar/7249ed41cf4ab50c5d37c169786ff768.jpg',
    '/default_avatar/8e4cd15170685b5cef2ae84488a491b9.jpg',
    '/default_avatar/95c890dba4a37b89d0bf046693540b3e.jpg'
  ];

  /**
   * Transforms a given avatar URI.
   * If it's empty, null, or undefined, it falls back to one of the default avatars.
   * Which default avatar is chosen depends on hashing the provided identifier string.
   *
   * @param value Optional current avatar URI string
   * @param identifier Optional identifier (User ID, username) to generate a stable, deterministic fallback
   * @returns Resolvable image URI path
   */
  transform(value: string | null | undefined, identifier?: string | null): string {
    // Check for empty, null, undefined, or obviously invalid values
    if (!value || value.trim() === '' || this.isInvalidUrl(value)) {
      return this.getFallback(identifier);
    }
    return value;
  }

  private isInvalidUrl(url: string): boolean {
    const trimmed = url.trim().toLowerCase();
    // Reject strings like "null", "undefined", "default", etc.
    if (['null', 'undefined', 'none', 'default'].includes(trimmed)) {
      return true;
    }
    // Must start with http, https, / or data:
    if (!trimmed.startsWith('http') && !trimmed.startsWith('/') && !trimmed.startsWith('data:')) {
      return true;
    }
    // Block known dead placeholder image services
    const deadDomains = ['placeimg.com', 'placeholder.com', 'placekitten.com', 'lorempixel.com', 'example.com'];
    if (deadDomains.some(domain => trimmed.includes(domain))) {
      return true;
    }
    return false;
  }

  private getFallback(identifier?: string | null): string {

    // If there's no identifier to hash against, safely return the first default
    if (!identifier) {
      return this.fallbacks[0];
    }

    // Hash the identifier string to derive an index stably tied to that specific user
    let hash = 0;
    for (let i = 0; i < identifier.length; i++) {
        const char = identifier.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32bit integer
    }

    const index = Math.abs(hash) % this.fallbacks.length;
    return this.fallbacks[index];
  }
}
