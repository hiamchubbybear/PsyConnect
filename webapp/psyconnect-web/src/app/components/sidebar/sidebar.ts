import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.html',
  styleUrls: ['./sidebar.scss'],
  standalone: true,
  imports: [RouterModule],
})
export class SidebarComponent {
  isCollapsed = true;

  constructor(public authService: AuthService) {}

  get showSidebar(): boolean {
    return this.authService.isLoggedIn();
  }
}
