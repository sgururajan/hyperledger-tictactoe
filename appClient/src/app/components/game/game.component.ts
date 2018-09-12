import { Component, OnInit, Input, Inject, Optional, ChangeDetectorRef } from '@angular/core';
import { GameViewModel } from '../../models/gameView.model';
import { Store, Select, Actions, ofActionDispatched, ofActionSuccessful } from '@ngxs/store';
import { NetworkState, NetworkModel } from '../../state/network.state';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { Cell } from '../../models/game.model';
import { Observable, BehaviorSubject } from 'rxjs';
import { GameStateModel, JoinGame, MakeMove } from '../../state/game.state';

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.scss']
})
export class GameComponent implements OnInit {
  
  disabled:boolean=false;
  game: GameViewModel;
  currentOrg:string;

  loading$:BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  @Select(state=> state.gamestate) games$:Observable<GameStateModel>;
  @Select(state=> state.networkstate) network$:Observable<NetworkModel>;


  constructor(private store:Store, private actions$:Actions,
  @Optional() private dialogRef: MatDialogRef<GameComponent>,
  private changeDetector: ChangeDetectorRef) { 

    this.actions$.pipe(ofActionDispatched(MakeMove)).subscribe(()=>this.loading$.next(true));
    this.actions$.pipe(ofActionSuccessful(MakeMove)).subscribe(()=>{
      this.loading$.next(false);
      this.changeDetector.detectChanges();
    });

    this.network$.subscribe(x=>{
      this.currentOrg = x.currentOrganization.name;
    });

    this.games$.subscribe(x=>{
      this.game = x.currentGame;
      this.disabled = this.game.nextPlayer!=this.currentOrg || this.game.completed || !this.game.canPlay;
    });
  }

  ngOnInit() {
    if (!this.game) {
      this.disabled=true;
      return;
    }   
  }

  onCellClick(cell:Cell) {
    console.log(`clicked cell: `, cell);  
    this.store.dispatch(new MakeMove(this.game.id, cell.row, cell.column));
  }

  onCloseButtonClick() {
    if (this.dialogRef && this.dialogRef.componentInstance) {
      this.dialogRef.close();
    }
  }

  createMockGame():GameViewModel {
    console.log("creating mock game");
    const gameCells:Cell[] = [
      { column:0, row: 0, value:"x", },
      { column:1, row: 0, value:"o", },
      { column:2, row: 0, value:"x", },
      { column:0, row: 1, value:"o", },
      { column:1, row: 1, value:"x", },
      { column:2, row: 1, value:"o", },
      { column:0, row: 2, value:"x", },
      { column:1, row: 2, value:"o", },
      { column:2, row: 2, value:"x", }
    ]
    return <GameViewModel>{
      id: 1,
      status: "In Progress",
      players: "org1, org2",
      lastTxId: "",
      nextPlayer: "org1",
      canJoin: false,
      canPlay: true,
      awaitingOtherPlayer: false,
      cells: gameCells,
    }
  }

}
