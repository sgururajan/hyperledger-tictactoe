import {Component, OnInit} from '@angular/core';
import {Actions, ofActionDispatched, ofActionSuccessful, Store} from "@ngxs/store";
import {GetNetworks, SelectNetwork, SetCurrentNetworkOrganizations} from "./state/network.state";
import {Action} from "@ngxs/store/src/action";
import {GetAllGameList} from "./state/game.state";
import {WebsocketService} from "./shared/services/websocket.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit{
  title = 'Tictactoe';

  constructor(private store:Store, private action$:Actions, private socketServer: WebsocketService) {
    this.action$.pipe(ofActionSuccessful(GetNetworks)).subscribe(res=>{
      this.store.dispatch(new SelectNetwork("testnetwork"));
    });
    this.action$.pipe(ofActionSuccessful(GetNetworks, SelectNetwork)).subscribe(res=>{
      console.log("setting orgs for current network");
      this.store.dispatch(new SetCurrentNetworkOrganizations());
    });
  }

  ngOnInit() {
    this.socketServer.init();
    this.socketServer.connect();
    this.store.dispatch(new GetNetworks());
  }
}
