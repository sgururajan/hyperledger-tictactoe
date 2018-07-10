import { BrowserModule } from '@angular/platform-browser';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA, APP_INITIALIZER } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { FlexLayoutModule } from "@angular/flex-layout";
import { PolymerModule } from '@codebakery/origami';

import './polymer'

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TopnavbarComponent } from './components/topnavbar/topnavbar.component';
import { NetworkListComponent } from './components/network-list/network-list.component';
import { CoreModule } from './core/core.module';
import { AppConfigFactory } from './core/config/app-config.factory';
import { AppConfigService } from './core/config/app-config.service';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { LayoutModule } from '@angular/cdk/layout';
import { MatToolbarModule, MatButtonModule, MatSidenavModule, MatIconModule, MatListModule } from '@angular/material';
import { NgxsModule } from '@ngxs/store';
import { NetworksInfoState } from './state/network-info.state';
import { HttpClientModule } from '@angular/common/http';
import { NetworkApiModule } from './modules/network-api/network-api.module';
import { NetworkService } from './modules/network-api/api/network.service';

@NgModule({
  declarations: [
    AppComponent,
    TopnavbarComponent,
    NetworkListComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    HttpClientModule,
    PolymerModule.forRoot(),
    FlexLayoutModule,
    CoreModule,    
    NetworkApiModule,
    BrowserAnimationsModule,
    LayoutModule,
    MatToolbarModule,
    MatButtonModule,
    MatSidenavModule,
    MatIconModule,
    MatListModule,
    NgxsModule.forRoot([
      NetworksInfoState,      
    ]),    
    AppRoutingModule,
  ],
  providers: [    
  ],
  bootstrap: [AppComponent],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ]
})
export class AppModule { }
