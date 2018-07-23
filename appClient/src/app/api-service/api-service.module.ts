import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BASE_PATH } from './variables';
import { AppConfigService } from '../core/config/app-config.service';
import { ApiFactory } from '../shared/api.factory';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [],
  providers: [
    {
      provide: BASE_PATH,
      useFactory: ApiFactory("api"),
      deps: [AppConfigService],
      multi: true,
    }
  ]
})
export class ApiServiceModule { }
