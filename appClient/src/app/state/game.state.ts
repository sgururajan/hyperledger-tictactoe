import {GameModel} from "../models/game.model";
import {Action, State, StateContext} from "@ngxs/store";
import {GameService} from "../api-service/api/game.service";
import {Injector} from "@angular/core";
import {tap} from "rxjs/operators";
import * as _ from "lodash";

export interface GameStateModel {
  gameList: GameModel[],
  gameListById:{[gameId:number]:GameModel},
  currentGame: GameModel
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

  constructor(private injector: Injector){}

  private get gameService() {
    if(!this._gameService) {
      this._gameService = this.injector.get(GameService)
    }
    return this._gameService;
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
    const gamesById = _.keyBy(action.gameList, g=>g.id);

    return context.patchState({
      gameList: _.orderBy(_.unionBy(state.gameList, action.gameList, "id"), ["id"], "[asc]"),
      gameListById: {...state.gameListById, ...gamesById},
    });
  }
  
  @Action(JoinGame)
  joinGame(context:StateContext<GameStateModel>, action:JoinGame) {
    return this.gameService.joinGame(action.gameId).pipe(tap(resp=>{
      return context.dispatch(new UpdateGameList(resp));
    }))
  }
}
