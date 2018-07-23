import { Injectable } from '@angular/core';
import {AppConfigService} from "../../core/config/app-config.service";
import {WebSocketSubject} from "rxjs/webSocket";
import {SocketMessage} from "../../models/socketMessage.model";
import {Store} from "@ngxs/store";
import {UpdateGameList} from "../../state/game.state";

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  private _host:string = "http://localhost:4300"
  private _socket$:WebSocketSubject<SocketMessage>

  constructor(private appConfig: AppConfigService, private store:Store) {

  }

  init() {
    this.appConfig.config$.subscribe(()=>{
      this._host = this.appConfig.config.services["api"].replace("http://", "ws://")
      console.log(`socker server host set to: `, this._host);
    });
  }

  connect() {
    // this._socket$ = WebSocketSubject.create(`${this._host}/ws`);
    //
    // this._socket$.subscribe(
    //   (message:SocketMessage)=>console.log(`socket message: `,message),
    //   (err)=> console.error(`Error: `, err),
    //   ()=> console.warn("socket connection completed")
    // );
    // console.log(`socket obj created: `, this._socket$);

    var socket = new WebSocket("ws://localhost:4300/ws");
    socket.onopen = e=>console.log(`socket opened: `,  e);
    socket.onerror= err=>console.error(`socket error: `, err);
    socket.onmessage=e=> this.handleSocketMessage(e.data); //console.log(`socket message: `, );


  }

  handleSocketMessage(data) {
    var msg = JSON.parse(data);
    console.log(`received msg with type `, msg.type);
    console.log(`received msg with payload `, msg.payload);
    switch (msg.type) {
      case "gameadded": {
        this.handleGameUpdate(msg);
        break;
      }
      case "gameupdated": {
        this.handleGameUpdate(msg);
        break;
      }
      default: {
        console.log("no handler definded for socket msg type: " + msg.type);
        break;
      }
    }
  }

  // handleGameAdded(data:SocketMessage) {
  //   this.store.dispatch(new UpdateGameList(data.payload))
  // }

  handleGameUpdate(data:SocketMessage) {
    this.store.dispatch(new UpdateGameList(data.payload));
  }


}
