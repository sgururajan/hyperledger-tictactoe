export interface GameViewModel {
  id: number
  status: string
  players: string
  nextPlayer: string
  lastTxId: string
  canPlay:boolean
  awaitingOtherPlayer:boolean
  canJoin:boolean
}
