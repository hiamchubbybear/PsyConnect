import {
  ActivatedRouteSnapshot,
  DetachedRouteHandle,
  RouteReuseStrategy,
} from '@angular/router';

export class CustomRouteReuseStrategy implements RouteReuseStrategy {
  private storedRoutes = new Map<string, DetachedRouteHandle>();

  // Routes that should be cached
  private routesToCache = [
    'feature/feed',
    'feature/search',
    'feature/chat',
    'feature/chat/:id',
  ];

  /**
   * Determines if this route (and its subtree) should be detached to be reused later
   */
  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    const path = this.getFullPath(route);
    return this.routesToCache.includes(path);
  }

  /**
   * Stores the detached route.
   */
  store(route: ActivatedRouteSnapshot, handle: DetachedRouteHandle): void {
    if (handle) {
      this.storedRoutes.set(this.getFullPath(route), handle);
    }
  }

  /**
   * Determines if this route (and its subtree) should be reattached
   */
  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    const path = this.getFullPath(route);
    return !!route.routeConfig && !!this.storedRoutes.get(path);
  }

  /**
   * Retrieves the previously stored route
   */
  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle | null {
    if (!route.routeConfig) {
      return null;
    }
    const path = this.getFullPath(route);
    return this.storedRoutes.get(path) || null;
  }

  /**
   * Determines if a route should be reused
   */
  shouldReuseRoute(
    future: ActivatedRouteSnapshot,
    curr: ActivatedRouteSnapshot,
  ): boolean {
    // Reuse the component if it's the same class (e.g. navigating between chat IDs)
    return future.component === curr.component;
  }

  private getFullPath(route: ActivatedRouteSnapshot): string {
    return route.routeConfig?.path || '';
  }
}
