import { Injectable } from '@angular/core';
import { Http } from "@angular/http";
import { BehaviorSubject } from 'rxjs';
import { AppConfig } from './app-config.model';
import { BaseApiService } from '../../shared/services/base-api.service';
import { map, catchError } from "rxjs/operators";
import { AppConfigFactory } from './app-config.factory';

const Default_Config = <AppConfig>{
  productName: '',
  services: {}
}

@Injectable({
  providedIn: 'root'
})
export class AppConfigService extends BaseApiService {
  config$:BehaviorSubject<AppConfig> = new BehaviorSubject<AppConfig>(undefined);

  constructor(private http: Http) { 
    super()
  }

  load():Promise<AppConfig> {
    return new Promise(resolver=> {
      this.http.get("assets/config.json")
        .pipe(map(this.extract), catchError(this.handleError))
        .subscribe((config:AppConfig)=> {
          this.config$.next(config);
          resolver(config);

          console.log("loaded configuration from global config file");
        });
    });
  }

  public get config() {
    if (!this.config$.value) {
      return Default_Config
    }
    return this.config$.value;
  }

}
