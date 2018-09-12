import { Component, OnInit } from '@angular/core';
import {Store} from "@ngxs/store";
import {AddNewGame} from "../../state/game.state";

@Component({
  selector: 'app-game-list-actions',
  templateUrl: './game-list-actions.component.html',
  styleUrls: ['./game-list-actions.component.scss']
})
export class GameListActionsComponent implements OnInit {

  constructor(private store:Store) { }

  ngOnInit() {
  }

  onAddGame() {
    this.store.dispatch(new AddNewGame());
  }

}
