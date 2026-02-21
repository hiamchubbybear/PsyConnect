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
  private observer?: IntersectionObserver;

  constructor(
    private elementRef: ElementRef,
    public authService: Auth,
  ) {}

  get showSidebar(): boolean {
    return this.authService.isLoggedIn();
  }
}
