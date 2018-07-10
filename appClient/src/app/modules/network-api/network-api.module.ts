import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BASE_PATH } from './variables';
import { CreateApiFactory } from '../../utils/create-api.factory';
import { AppConfigService } from '../../core/config/app-config.service';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [],
  providers:[
    {
      provide: BASE_PATH,
      useFactory: CreateApiFactory("api"),
      deps: [AppConfigService]
    }
  ] 
})
export class NetworkApiModule { }
