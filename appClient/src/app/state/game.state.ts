import {GameModel} from "../models/game.model";
import {Action, State, StateContext, Store, Selector} from "@ngxs/store";
import {GameService} from "../api-service/api/game.service";
import {Injector} from "@angular/core";
import {tap} from "rxjs/operators";
import * as _ from "lodash";
import { GameViewModel } from "../models/gameView.model";
import { NetworkState } from "./network.state";

export interface GameStateModel {
  gameListById:{[gameId:number]:GameViewModel},
  currentGame: GameViewModel
}

export class GetGameList {
  static readonly type="[GameState] GetGameList";
  constructor(public pageIndex:number, public pageSize:number){}
}

export class GetAllGameList {
  static readonly type="[GameState] GetAllGameList";
}

export class AddNewGame {
  static readonly type="[GameState] AddNewGame";
}

export class JoinGame {
  static readonly type="[GameState] Joingame";
  constructor(public gameId: number){}
}

export class MakeMove {
  static readonly type="[GameState] MakeMove";
  constructor(public gameId:number, public row:number, public col:number){}
}

export class SetCurrentGame {
  static readonly type="[GameState] SetCurrentGame";
  constructor(public gameId:number){}
}

export class UpdateGameList {
  static readonly type="[GameState] UpdateGameList";
  constructor (public  gameList:GameModel[]){}
}

@State({
  name: "gamestate",
  defaults: {
    gameList: [],
    gameListById: {},
    currentGame: undefined
  }
})
export class GameState {

  private _gameService:GameService;

  constructor(private injector: Injector, private store:Store){}

  private get gameService() {
    if(!this._gameService) {
      this._gameService = this.injector.get(GameService)
    }
    return this._gameService;
  }

  @Selector() static currentGame(state: GameStateModel): GameViewModel {
    return state.currentGame;
  }

  @Action(GetGameList)
  getGameList(context:StateContext<GameStateModel>, action: GetGameList) {
    return this.gameService.getGameList(action.pageIndex, action.pageSize).pipe(tap(resp=>{      
      return context.dispatch(new UpdateGameList(resp))
    }));
  }

  @Action(GetAllGameList)
  getAllGameList(context: StateContext<GameStateModel>, action: GetAllGameList) {
    return this.gameService.getAllGameList().pipe(tap(resp=>{
      console.log(`all games list: `, resp);
      return context.dispatch(new UpdateGameList(resp));
    }));
  }

  @Action(AddNewGame)
  addNewGame(context: StateContext<GameStateModel>, action:AddNewGame) {
    return this.gameService.addGame().pipe(tap(resp=>{
      return context.dispatch(new UpdateGameList(resp));
    }));
  }

  @Action(UpdateGameList)
  updateGameList(context: StateContext<GameStateModel>, action: UpdateGameList) {
    let state = context.getState();    
    const gameViewModels = _.map(action.gameList, g=>this.convertToGameViewModel(g));
    const gamesById = _.keyBy(gameViewModels, g=>g.id);
    const currentGameId = state.currentGame?state.currentGame.id:-1;
    const combindedGameList = {...state.gameListById,...gamesById};
    const currentGame = currentGameId>-1 ?combindedGameList[currentGameId]:undefined;
    return context.patchState({
      gameListById: combindedGameList,
      currentGame: currentGame      
    });
  }
  
  @Action(JoinGame)
  joinGame(context:StateContext<GameStateModel>, action:JoinGame) {
    return this.gameService.joinGame(action.gameId).pipe(tap(resp=>{
      console.log(`joing game response: `, resp);
      return context.dispatch(new UpdateGameList(resp));
    }))
  }

  @Action(MakeMove)
  makeMove(context:StateContext<GameStateModel>, action:MakeMove) {
    console.log(`received in action, gameid: ${action.gameId}, row: ${action.row}, col: ${action.col}`);
    return this.gameService.makeMove(action.gameId, action.row, action.col).pipe(tap(resp=>{
      console.log(`make move response: `, resp);
      return context.dispatch(new UpdateGameList(resp));
    }));
  }

  @Action(SetCurrentGame)
  setCurrentGame(context:StateContext<GameStateModel>, action:SetCurrentGame) {
    let state = context.getState()
    return context.patchState({
      currentGame: state.gameListById[action.gameId]
    });
  }

  convertToGameViewModel(game: GameModel):GameViewModel {
    console.log(`game model: `, game);
    return <GameViewModel>{
      id: game.id,
      status: this.getGameStatus(game),
      players: game.players.map(x=>x.name).join(","),
      lastTxId: game.txId,
      nextPlayer: game.players[game.playerToPlay].name,
      canJoin: this.getCanJoin(game),
      canPlay: this.getCanPlay(game),
      awaitingOtherPlayer: this.getAwaitingOtherPlayer(game),
      cells: game.cells,
      completed: game.completed,
      winner: game.winner,
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
}
