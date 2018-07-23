import { Injectable } from '@angular/core';
import { BehaviorSubject, config } from 'rxjs';
import { AppConfig } from './app-config.model';
import { BaseApiService } from  '../../shared/base-api.service';
import { Http } from '@angular/http';
import { map, catchError } from "rxjs/operators";
import {HttpClient} from "@angular/common/http";

const Default_Config = <AppConfig>{
  productName: '',
  services: {}
}

@Injectable({
  providedIn: 'root'
})
export class AppConfigService extends BaseApiService {
  config$: BehaviorSubject<AppConfig> = new BehaviorSubject<AppConfig>(undefined);

  constructor(private httpClient:HttpClient) {
    super()
   }

   load():Promise<AppConfig> {
     return new Promise(resolve => {
      this.httpClient.get<AppConfig>("assets/config.json")
        //.pipe(map(this.extract), catchError(this.handleError))
        .subscribe((config:AppConfig)=>{
          this.config$.next(config);
          resolve(config);
          console.log("loaded config from global config file");
        });
     });
   }

   public get config() {
     if (!this.config$.value) {
       return Default_Config;
     }
     return this.config$.value;
   }
}
