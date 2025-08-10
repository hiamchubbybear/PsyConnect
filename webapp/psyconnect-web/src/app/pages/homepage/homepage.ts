import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-homepage',
  imports: [RouterModule],
  standalone: true,
  templateUrl: './homepage.html',
  styleUrl: './homepage.scss',
})
export class Homepage {}
