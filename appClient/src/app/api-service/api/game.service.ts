import {Inject, Injectable} from '@angular/core';
import {HttpClient, HttpEvent, HttpResponse} from "@angular/common/http";
import {BASE_PATH} from "../variables";
import {Observable} from "rxjs/internal/Observable";
import {GameModel} from "../../models/game.model";

@Injectable({
  providedIn: 'root'
})
export class GameService {

  private basepath:string="http://localhost:4300";

  constructor(private httpClient: HttpClient, @Inject(BASE_PATH) basepath:string) {
    if (basepath) {
      this.basepath = basepath;
    }
  }

  public getGameList(pageIndex:number, pageSize:number, observer?:"body", reportProgress?:true):Observable<GameModel[]>;
  public getGameList(pageIndex:number, pageSize:number, observer?:"response", reportProgress?:true):Observable<HttpResponse<GameModel[]>>;
  public getGameList(pageIndex:number, pageSize:number, observer?:"events", reportProgress?:true):Observable<HttpEvent<GameModel[]>>;
  public getGameList(pageIndex:number=1, pageSize:number=10, observer:any ="body", reportProgress:boolean=true):Observable<any> {
    const url = `${this.basepath}/api/innetwork/getGameList/${pageIndex}/${pageSize}`;
    return this.httpClient.get<GameModel[]>(url, {observe: observer, reportProgress:reportProgress});
  }

  public getAllGameList(observer?:"body", reportProgress?:true):Observable<GameModel[]>;
  public getAllGameList(observer?:"response", reportProgress?:true):Observable<HttpResponse<GameModel[]>>;
  public getAllGameList(observer?:"events", reportProgress?:true):Observable<HttpEvent<GameModel[]>>;
  public getAllGameList(observer:any ="body", reportProgress:boolean=true):Observable<any> {
    const url = `${this.basepath}/api/innetwork/getAllGameList`;
    return this.httpClient.get<GameModel[]>(url, {observe: observer, reportProgress:reportProgress});
  }

  public addGame(observer?:"body", reportProgress?:true):Observable<GameModel[]>;
  public addGame(observer?:"response", reportProgress?:true):Observable<HttpResponse<GameModel[]>>;
  public addGame(observer?:"events", reportProgress?:true):Observable<HttpEvent<GameModel[]>>;
  public addGame(observer:any ="body", reportProgress:boolean=true):Observable<any> {
    const url = `${this.basepath}/api/innetwork/addgame`;
    return this.httpClient.post<GameModel[]>(url, {}, {observe: observer, reportProgress:reportProgress});
  }

  public joinGame(gameId?,observer?:"body", reportProgress?:true):Observable<GameModel[]>;
  public joinGame(gameId?,observer?:"response", reportProgress?:true):Observable<HttpResponse<GameModel[]>>;
  public joinGame(gameId?, observer?:"events", reportProgress?:true):Observable<HttpEvent<GameModel[]>>;
  public joinGame(gameId:number, observer:any ="body", reportProgress:boolean=true):Observable<any> {
    const url = `${this.basepath}/api/innetwork/joingame/${gameId}`;
    return this.httpClient.post<GameModel[]>(url, {}, {observe: observer, reportProgress:reportProgress});
  }

  public makeMove(gameId?,row?,col?,observer?:"body", reportProgress?:true):Observable<GameModel[]>;
  public makeMove(gameId?,row?,col?,observer?:"response", reportProgress?:true):Observable<HttpResponse<GameModel[]>>;
  public makeMove(gameId?,row?,col?,observer?:"events", reportProgress?:true):Observable<HttpEvent< GameModel[]>>;
  public makeMove(gameId:number,row:number,col:number,observer:any="body", reportProgress:boolean=true):Observable<any>{
    console.log(`received gameid: ${gameId}, row: ${row}, col: ${col}`);
    const payload = {
      gameId: gameId,
      row: row,
      column: col
    };

    console.log(`payload sending to server: `, payload);

    const url = `${this.basepath}/api/innetwork/makemove`;
    return this.httpClient.post<GameModel[]>(url, payload, {observe: observer, reportProgress: reportProgress});
  }
}
