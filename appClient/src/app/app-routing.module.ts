import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {NetworkLoginComponent} from "./components/network-login/network-login.component";
import {GameListComponent} from "./components/game-list/game-list.component";
import {NetworkAuthGuardService} from "./core/guard/network-auth-guard.service";

const routes: Routes = [
  {path: "", redirectTo: "/gamelist", pathMatch: "full"},
  {path:"gamelist", component: GameListComponent, canActivate:[NetworkAuthGuardService]}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
