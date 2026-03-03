import {
  ActivatedRouteSnapshot,
  DetachedRouteHandle,
  RouteReuseStrategy,
} from '@angular/router';

/**
 * Custom route reuse strategy that caches certain routes (like chat and feed)
 * so they are never destroyed when navigating away — like Messenger.
 *
 * Key behaviors:
 * - `feature/chat` and `feature/chat/:id` are stored under the SAME key (`feature/chat`)
 *   so navigating between different chat users never recreates the component.
 * - `feature/feed` and `feature/search` are also cached.
 * - All other routes are destroyed and recreated normally.
 */
export class CustomRouteReuseStrategy implements RouteReuseStrategy {
  private storedRoutes = new Map<string, DetachedRouteHandle>();

  // Route paths that should be permanently cached (never destroyed)
  private readonly CACHED_ROUTES = new Set([
    'feature/feed',
    'feature/search',
    'feature/chat',
    'feature/chat/:id',
  ]);

  // Normalize route key: treat `feature/chat/:id` same as `feature/chat`
  // so switching contacts doesn't destroy the component
  private getCacheKey(route: ActivatedRouteSnapshot): string {
    const path = route.routeConfig?.path || '';
    if (path === 'feature/chat/:id') return 'feature/chat';
    return path;
  }

  /** Should we save this route when navigating away? */
  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    const path = route.routeConfig?.path || '';
    return this.CACHED_ROUTES.has(path);
  }

  /** Save the detached route component tree */
  store(route: ActivatedRouteSnapshot, handle: DetachedRouteHandle): void {
    if (handle) {
      this.storedRoutes.set(this.getCacheKey(route), handle);
    }
  }

  /** Should we restore a previously cached route? */
  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    if (!route.routeConfig) return false;
    const key = this.getCacheKey(route);
    return this.storedRoutes.has(key);
  }

  /** Restore the cached component tree */
  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle | null {
    if (!route.routeConfig) return null;
    return this.storedRoutes.get(this.getCacheKey(route)) || null;
  }

  /**
   * Should Angular reuse the existing component instance for the next route?
   * - Same component class → reuse (e.g. /chat/user1 → /chat/user2)
   * - Different component → don't reuse
   */
  shouldReuseRoute(
    future: ActivatedRouteSnapshot,
    curr: ActivatedRouteSnapshot,
  ): boolean {
    return future.routeConfig === curr.routeConfig;
  }
}
