import {GameModel} from "../models/game.model";
import {Action, State, StateContext} from "@ngxs/store";
import {GameService} from "../api-service/api/game.service";
import {Injector} from "@angular/core";
import {tap} from "rxjs/operators";
import * as _ from "lodash";

export interface GameStateModel {
  gameList: GameModel[],
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

@State({
  name: "gamestate",
  defaults: {
    gameList: [],
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
      let state = context.getState();
      return context.patchState({
        gameList: _.orderBy(_.unionBy(state.gameList, resp, "id"), ["id"],["asc"])
      });
    }));
  }

  @Action(GetAllGameList)
  getAllGameList(context: StateContext<GameStateModel>, action: GetAllGameList) {
    return this.gameService.getAllGameList().pipe(tap(resp=>{
      let state = context.getState()
      return context.patchState({
        gameList: _.orderBy(_.unionBy(state.gameList, resp, "id"), ["id"], ["asc"])
      });
    }));
  }

  @Action(AddNewGame)
  addNewGame(context: StateContext<GameStateModel>, action:AddNewGame) {
    return this.gameService.addGame().pipe(tap(resp=>{
      let state = context.getState();
      return context.patchState({
        gameList: _.orderBy(_.unionBy(state.gameList, resp, "id"), ["id"], ["asc"])
      });
    }));
  }
}
