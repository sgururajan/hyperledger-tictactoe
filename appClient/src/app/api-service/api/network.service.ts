import { Injectable, Inject } from '@angular/core';
import { BASE_PATH } from '../variables';
import { Observable } from 'rxjs';
import { HttpClient, HttpEvent, HttpResponse } from '@angular/common/http';
import { Network } from '../../models/network.model';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class NetworkService {

  private basepath = "http://localhost:4300";

  constructor(private httpClient: HttpClient, @Inject(BASE_PATH) basepath: string) {
    if (basepath) {
      this.basepath = basepath;
    }
  }

  public getNetworks(observer?: "body", reportProgress?: true): Observable<Network[]>;
  public getNetworks(observer?: "events", reportProgress?: true): Observable<HttpEvent<Network[]>>;
  public getNetworks(observer?: "response", reportProgress?: true): Observable<HttpResponse<Network[]>>;
  public getNetworks(observer: any = "body", reportProgress: boolean = true): Observable<any> {
    const url = `${this.basepath}/api/networks`;
    return this.httpClient.get<Network[]>(url, { observe: observer, reportProgress: reportProgress });
  }
}
