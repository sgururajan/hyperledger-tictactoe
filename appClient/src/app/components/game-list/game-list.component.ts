import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {Actions, ofActionDispatched, ofActionSuccessful, Select, Store} from "@ngxs/store";
import {GameModel} from "../../models/game.model";
import {Observable} from "rxjs/internal/Observable";
import {GameStateModel, GetAllGameList, AddNewGame, UpdateGameList, JoinGame} from "../../state/game.state";
import {MatTableDataSource} from "@angular/material";
import {NetworkState} from "../../state/network.state";
import {GameViewModel} from "../../models/gameView.model";
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";

import * as _ from "lodash";


@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss']
})
export class GameListComponent implements OnInit {

  gameDataSource:MatTableDataSource<GameViewModel>;
  columnsToDisplayFields = ["id", "status", "players", "nextPlayer", "lastTxId", "action"];

  loading$:BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false );

  @Select(state=>state.gamestate.gameList) gameList$:Observable<GameStateModel>;

  constructor(private store: Store, private action$:Actions, private changeDetector: ChangeDetectorRef) {
    //this.loading$ = this.loadingDataSource$.asObservable();
    this.action$.pipe(ofActionDispatched(GetAllGameList)).subscribe(()=> this.loading$.next(true));
    this.action$.pipe(ofActionSuccessful(GetAllGameList)).subscribe(()=> {
      this.loading$.next(false)
      this.changeDetector.detectChanges();
    });

    this.action$.pipe(ofActionDispatched(UpdateGameList)).subscribe(()=> this.loading$.next(true));
    this.action$.pipe(ofActionSuccessful(UpdateGameList)).subscribe(()=> {
      this.loading$.next(false)
      this.changeDetector.detectChanges();
    });

    this.action$.pipe(ofActionDispatched(AddNewGame)).subscribe(()=> this.loading$.next(true));
    this.action$.pipe(ofActionSuccessful(AddNewGame)).subscribe(()=> {
      this.loading$.next(false);
      this.changeDetector.detectChanges();
    });

    this.action$.pipe(ofActionDispatched(JoinGame)).subscribe(()=> this.loading$.next(true));
    this.action$.pipe(ofActionSuccessful(JoinGame)).subscribe(()=> {
      this.loading$.next(false);
      this.changeDetector.detectChanges();
    });
  }

  ngOnInit() {


    this.store.dispatch(new GetAllGameList());
    this.gameDataSource = new MatTableDataSource<GameViewModel>();
    this.store.select(state => state.gamestate.gameList).subscribe((resp:GameModel[])=>{
      //this.gameDataSource.data = resp.map(x=> this.convertToGameViewModel(x));
    });

    this.store.select(state=>state.gamestate.gameListById).subscribe((resp:{[id:number]:GameModel})=> {
      //this.gameDataSource.data = resp.map(x=> this.convertToGameViewModel(x[x]))
      console.log(`Game list hash table: `, resp);
      const data = _.map(resp, g=> this.convertToGameViewModel(g));
      console.log(`data: `, data);
      this.gameDataSource.data = _.map(resp, g=>this.convertToGameViewModel(g));
    })

    //this.columnsToDisplay.map(x=>x.fieldId)
  }

  convertToGameViewModel(game: GameModel):GameViewModel {

    return <GameViewModel>{
      id: game.id,
      status: this.getGameStatus(game),
      players: game.players.map(x=>x.name).join(","),
      lastTxId: game.txId,
      nextPlayer: game.players[game.playerToPlay].name,
      canJoin: this.getCanJoin(game),
      canPlay: this.getCanPlay(game),
      awaitingOtherPlayer: this.getAwaitingOtherPlayer(game),
    }
  }

  getAwaitingOtherPlayer(game:GameModel):boolean {
    const currentOrg = this.store.selectSnapshot(NetworkState.currentOrganization).name;
    return !game.completed && game.players.every(x=>x.name!=="") && game.players[game.playerToPlay].name!= currentOrg;
  }

  getCanJoin(game:GameModel): boolean {
    const currentOrg = this.store.selectSnapshot(NetworkState.currentOrganization).name;
    return !game.completed && !game.players.some(x=>x.name==currentOrg) && game.players.some(x=>x.name==="");
  }

  getCanPlay(game:GameModel): boolean {
    const currentOrg = this.store.selectSnapshot(NetworkState.currentOrganization).name;
    return !game.completed && game.players[game.playerToPlay].name===currentOrg;
  }

  getGameStatus(game:GameModel):string {
    let gameStatus = "In Progress";

    if (game.players.some(x=>x.name==="")) {
      gameStatus = "Awaiting...";
    } else if(game.completed) {
      gameStatus = "Completed";
    }

    return gameStatus;
  }

  onRowDblClicked(row) {
    console.log(row);
  }

  onPlayClicked(row) {
    console.log(row);
  }

  onJoinClicked(game:GameViewModel) {
    // console.log(game);
    // console.log(game.id);
    this.store.dispatch(new JoinGame(game.id));
  }

}



