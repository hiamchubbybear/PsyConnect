import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Homepage } from './pages/homepage/homepage';
import { Notfound } from './pages/notfound/notfound';

export const routes: Routes = [
    {
        path : '' , component : Homepage,
    },
{
        path : '**' , component : Notfound,
    }

];
@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
