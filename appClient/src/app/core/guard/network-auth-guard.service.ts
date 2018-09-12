import { Injectable } from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot} from "@angular/router";
import {ModalDialogService} from "../../shared/services/modal-dialog.service";
import {Store} from "@ngxs/store";
import {NetworkState} from "../../state/network.state";
import {NetworkLoginComponent} from "../../components/network-login/network-login.component";
import {GetAllGameList} from "../../state/game.state";

@Injectable({
  providedIn: 'root'
})
export class NetworkAuthGuardService {

  constructor(private router: Router, private modalService: ModalDialogService, private store:Store) { }

  async canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    const currentOrg = this.store.selectSnapshot(NetworkState.currentOrganization);
    if (!currentOrg || currentOrg.name==="") {
      let dialogResult = await this.modalService.open(NetworkLoginComponent);
      if(dialogResult) {
        //this.store.dispatch(GetAllGameList);
        return dialogResult.Success;
      }
      return false;
    }
    return true;
  }
}
