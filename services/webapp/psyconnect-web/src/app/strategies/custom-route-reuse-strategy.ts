import {
  ActivatedRouteSnapshot,
  DetachedRouteHandle,
  RouteReuseStrategy,
} from '@angular/router';


export class CustomRouteReuseStrategy implements RouteReuseStrategy {
  private storedRoutes = new Map<string, DetachedRouteHandle>();

  
  private readonly CACHED_ROUTES = new Set([
    'feature/feed',
    'feature/search',
    'feature/chat',
    'feature/chat/:id',
  ]);

  
  
  private getCacheKey(route: ActivatedRouteSnapshot): string {
    const path = route.routeConfig?.path || '';
    if (path === 'feature/chat/:id') return 'feature/chat';
    return path;
  }

  
  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    const path = route.routeConfig?.path || '';
    return this.CACHED_ROUTES.has(path);
  }

  
  store(route: ActivatedRouteSnapshot, handle: DetachedRouteHandle): void {
    if (handle) {
      this.storedRoutes.set(this.getCacheKey(route), handle);
    }
  }

  
  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    if (!route.routeConfig) return false;
    const key = this.getCacheKey(route);
    return this.storedRoutes.has(key);
  }

  
  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle | null {
    if (!route.routeConfig) return null;
    return this.storedRoutes.get(this.getCacheKey(route)) || null;
  }

  
  shouldReuseRoute(
    future: ActivatedRouteSnapshot,
    curr: ActivatedRouteSnapshot,
  ): boolean {
    return future.routeConfig === curr.routeConfig;
  }
}
