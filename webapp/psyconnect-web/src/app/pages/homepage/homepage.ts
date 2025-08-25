import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ThemeService } from '../../services/theme/theme-service';

@Component({
  selector: 'app-homepage',
  imports: [RouterModule],
  standalone: true,
  templateUrl: './homepage.html',
  styleUrl: './homepage.scss',
})
export class Homepage {
  constructor(private themeService: ThemeService) {}
  bannerUrl = './assets/icon/banner-psyconnect-light.svg';
  ngOnInit(): void {
    if (this.themeService.getTheme() == 'light') {
      console.log('Light banner');

      this.bannerUrl = './assets/icon/banner-psyconnect-dark-meme.svg';
    } else {
      this.bannerUrl = './assets/icon/banner-psyconnect-lightmeme.svg';
    }
  }
}
