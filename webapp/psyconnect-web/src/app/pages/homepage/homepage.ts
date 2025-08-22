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
  bannerUrl = 'assets/images/ndbanner-light.png';
  ngOnInit(): void {
    if (this.themeService.getTheme() == 'light') {
        console.log("Light banner");

      this.bannerUrl = 'assets/images/white-second-poster.avif';
    } else {
      this.bannerUrl = 'assets/images/black-second-poster.avif';
    }
  }
}
