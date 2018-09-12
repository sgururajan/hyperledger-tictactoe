import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {MaterialUiModule} from "./modules/material-ui";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NgxsModule} from "@ngxs/store";
import {CoreModule} from "./core/core.module";
import {ApiServiceModule} from "./api-service/api-service.module";
import {NetworkState} from "./state/network.state";
import {TopnavbarComponent} from "./components/topnavbar/topnavbar.component";
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import { NetworkLoginComponent } from './components/network-login/network-login.component';
import { GameListComponent } from './components/game-list/game-list.component';
import {NetworkAuthInterceptorService} from "./shared/services/network-auth-interceptor.service";
import {GameState} from "./state/game.state";
import { CommonModule } from '@angular/common';
import { GameListActionsComponent } from './components/game-list-actions/game-list-actions.component';
import { OverlayComponent } from './components/overlay/overlay.component';
import { GameComponent } from './components/game/game.component';
import { BlockComponent } from './components/block/block.component';
import { BlockListComponent } from './components/block-list/block-list.component';
import { BlockState } from './state/block.state';

@NgModule({
  declarations: [
    AppComponent,
    TopnavbarComponent,
    NetworkLoginComponent,
    GameListComponent,
    GameListActionsComponent,
    OverlayComponent,
    GameComponent,
    BlockComponent,
    BlockListComponent
  ],
  imports: [
    CoreModule,
    CommonModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    ReactiveFormsModule,
    NgxsModule.forRoot([
      NetworkState,
      GameState,
      BlockState,
    ]),
    MaterialUiModule,
    ApiServiceModule,
    AppRoutingModule
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: NetworkAuthInterceptorService,
      multi: true,
    }
  ],
  entryComponents: [
    NetworkLoginComponent,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
