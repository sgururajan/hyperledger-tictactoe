import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NetworkListComponent } from './components/network-list/network-list.component';

const routes: Routes = [
  {path: '', redirectTo: "/networklist", pathMatch: "full"},
  {path: "networklist", component: NetworkListComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
