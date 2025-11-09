import { Component } from '@angular/core';
import { SharedTranslateModule } from '../../shared/translate/translate.module';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-footer',
  imports: [TranslateModule],
  templateUrl: './footer.html',
  styleUrl: './footer.scss'
})
export class Footer {

}
