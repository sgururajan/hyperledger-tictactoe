import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ModuleWithProviders } from '@angular/compiler/src/core';
import { AppConfigFactory } from './config/app-config.factory';
import { AppConfigService } from './config/app-config.service';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [],
  providers:[
    {
      provide: APP_INITIALIZER,
      useFactory: AppConfigFactory,
      deps:[AppConfigService],
      multi: true
    }
  ]
})
export class CoreModule { 
  
}
