import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {NetworkLoginComponent} from "./components/network-login/network-login.component";
import {GameListComponent} from "./components/game-list/game-list.component";
import {NetworkAuthGuardService} from "./core/guard/network-auth-guard.service";
import { GameComponent } from './components/game/game.component';

const routes: Routes = [
  {path: "", redirectTo: "/game", pathMatch: "full"},
  {path:"game", component: GameComponent},
  {path:"gamelist", component: GameListComponent, canActivate:[NetworkAuthGuardService]}, 
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
