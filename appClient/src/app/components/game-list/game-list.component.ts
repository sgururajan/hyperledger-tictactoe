import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {Actions, ofActionDispatched, ofActionSuccessful, Select, Store} from "@ngxs/store";
import {GameModel} from "../../models/game.model";
import {Observable} from "rxjs/internal/Observable";
import {GameStateModel, GetAllGameList, AddNewGame, UpdateGameList, JoinGame, SetCurrentGame} from "../../state/game.state";
import {MatTableDataSource} from "@angular/material";
import {NetworkState} from "../../state/network.state";
import {GameViewModel} from "../../models/gameView.model";
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";

import * as _ from "lodash";
import { ModalDialogService } from '../../shared/services/modal-dialog.service';
import { GameComponent } from '../game/game.component';
import { async } from 'rxjs/internal/scheduler/async';


@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss']
})
export class GameListComponent implements OnInit {

  gameDataSource:MatTableDataSource<GameViewModel>;
  columnsToDisplayFields = ["id", "status", "players", "nextPlayer", "winner", "action"];

  loading$:BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  @Select(state=>state.gamestate.gameList) gameList$:Observable<GameStateModel>;

  constructor(private store: Store, private action$:Actions, private changeDetector: ChangeDetectorRef,
    private modalService: ModalDialogService) {
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

    this.store.select(state=>state.gamestate.gameListById).subscribe((resp:{[id:number]:GameViewModel})=> {
      //this.gameDataSource.data = resp.map(x=> this.convertToGameViewModel(x[x]))
      // console.log(`Game list hash table: `, resp);
      // const data = _.map(resp, g=> this.convertToGameViewModel(g));
      // console.log(`data: `, data);
      // this.gameDataSource.data = _.map(resp, g=>this.convertToGameViewModel(g));
      this.gameDataSource.data = _.map(resp, x=>x);
    })

    //this.columnsToDisplay.map(x=>x.fieldId)
  }

  async onRowDblClicked(row:GameViewModel) {
    console.log(row);
    await this.store.dispatch(new SetCurrentGame(row.id));    
    let dialogResult = await this.modalService.open(GameComponent, row);
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



