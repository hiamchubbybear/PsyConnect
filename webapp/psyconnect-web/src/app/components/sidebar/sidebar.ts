import { Component } from '@angular/core';
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

  constructor(public authService: Auth) {}

  get showSidebar(): boolean {
    return this.authService.isLoggedIn();
  }
}
