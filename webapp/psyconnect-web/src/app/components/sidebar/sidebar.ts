import { Component, ElementRef } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Auth } from '../../services/auth/auth';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.html',
  styleUrls: ['./sidebar.scss'],
  standalone: true,
  imports: [RouterModule],
})
export class SidebarComponent {
  isCollapsed = true;

  private observer?: IntersectionObserver;

  constructor(private elementRef: ElementRef, public authService: Auth) {}

  get showSidebar(): boolean {
    return this.authService.isLoggedIn();
  }

  //   ngOnInit(): void {
  //     const footer = document.querySelector('app-footer, footer');

  //     if (footer) {
  //       this.observer = new IntersectionObserver(
  //         (entries) => {
  //           entries.forEach((entry) => {
  //             this.isCollapsed = !entry.isIntersecting;
  //           });
  //         },
  //         { threshold: 0.1 }
  //       );
  //       this.observer.observe(footer);
  //     }
  //   }

  //   ngOnDestroy(): void {
  //     this.observer?.disconnect();
  //   }
}
