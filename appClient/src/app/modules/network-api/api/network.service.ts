import { Injectable, Optional, Inject } from '@angular/core';
import { HttpClient, HttpResponse, HttpEvent } from '@angular/common/http';
import { BASE_PATH } from '../variables';
import { NetworksInfo } from '../../../models/network.model';
import { Observable } from 'rxjs';
import { CreateApiFactory } from '../../../utils/create-api.factory';
import { AppConfigService } from '../../../core/config/app-config.service';

@Injectable({
  providedIn: 'root'
})
export class NetworkService {

  protected basePath = "http://something"

  constructor(private httpClient: HttpClient, @Inject(BASE_PATH) basePath: string) {
    if (basePath) {
      this.basePath = basePath;
    }
  }

  public getNetworksInfo(observer?: "body", reportProgress?: true): Observable<NetworksInfo[]>;
  public getNetworksInfo(observer?: "events", reportProgress?: true): Observable<HttpEvent<NetworksInfo[]>>;
  public getNetworksInfo(observer?: "response", reportProgress?: true): Observable<HttpResponse<NetworksInfo[]>>;
  public getNetworksInfo(observer: any = "body", reportProgress: boolean = true): Observable<any> {
    console.log("calling service getNetworkInfos");
    console.log(`url trying: ${this.basePath}/api/networksinfo`);
    return this.httpClient.get<NetworksInfo[]>(`${this.basePath}/api/networksinfo`, {
      observe: observer,
      reportProgress: reportProgress,
    });
  }
}
