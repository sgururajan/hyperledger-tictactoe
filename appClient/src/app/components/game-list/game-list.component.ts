import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {Actions, ofActionDispatched, ofActionSuccessful, Select, Store} from "@ngxs/store";
import {GameModel} from "../../models/game.model";
import {Observable} from "rxjs/internal/Observable";
import {GameStateModel, GetAllGameList} from "../../state/game.state";
import {MatTableDataSource} from "@angular/material";
import {NetworkState} from "../../state/network.state";
import {GameViewModel} from "../../models/gameView.model";
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";


@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss']
})
export class GameListComponent implements OnInit {

  gameDataSource:MatTableDataSource<GameViewModel>;
  columnsToDisplayFields = ["id", "status", "players", "nextPlayer", "lastTxId", "action"];

  loadingDataSource$:BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false );
  loading$:Observable<boolean>;

  @Select(state=>state.gamestate.gameList) gameList$:Observable<GameStateModel>;

  constructor(private store: Store, private action$:Actions, private changeDetector: ChangeDetectorRef) {
    this.loading$ = this.loadingDataSource$.asObservable();
    this.action$.pipe(ofActionDispatched(GetAllGameList)).subscribe(()=> this.loadingDataSource$.next(true));
    this.action$.pipe(ofActionSuccessful(GetAllGameList)).subscribe(()=> {
      console.log("GetAllGameListActionCompleted");
      this.loadingDataSource$.next(false)
      console.log(this.loadingDataSource$.getValue());
      this.changeDetector.detectChanges();
    });
  }

  ngOnInit() {


    this.store.dispatch(new GetAllGameList());
    this.gameDataSource = new MatTableDataSource<GameViewModel>();
    this.store.select(state => state.gamestate.gameList).subscribe((resp:GameModel[])=>{
      console.log(resp);
      this.gameDataSource.data = resp.map(x=> this.convertToGameViewModel(x));
      console.log(this.gameDataSource);
    });
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
    }
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

  onRowClicked(row) {
    console.log(row);
  }

  onPlayClicked(row) {
    console.log(row);
  }

}

export interface TableField{
  displayName:string
  fieldId: string
}


