import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from '@angular/router';
import { Footer } from "./components/footer/footer";
import { Header } from "./components/header/header";
import { Homepage } from './pages/homepage/homepage';
import { Notfound } from './pages/notfound/notfound';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterModule, Header, Footer, Notfound, Homepage],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {
  protected title = 'psyconnect-web';
}
