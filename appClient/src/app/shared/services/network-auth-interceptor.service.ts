import { Injectable } from '@angular/core';
import {HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from "@angular/common/http";
import {Observable} from "rxjs/internal/Observable";
import {Store} from "@ngxs/store";
import {NetworkState} from "../../state/network.state";

@Injectable({
  providedIn: 'root'
})
export class NetworkAuthInterceptorService implements HttpInterceptor{

  constructor(private store: Store) { }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const currNetwork=this.store.selectSnapshot(NetworkState.currentNetwork);
    const currOrg = this.store.selectSnapshot(NetworkState.currentOrganization);

    if(currNetwork && currOrg) {
      console.log("intercepting http call and adding headers");
      const reqWithNetworkInfo = req.clone({
        headers: req.headers.append("X-hlt3-networkName", currNetwork.name).append("X-hlt3-orgName", currOrg.name)
      });
      return next.handle(reqWithNetworkInfo);
    }

    return next.handle(req);
  }


}
