export interface GameModel {
    id: number
    completed: boolean
    players: Player[]
    winner: string
    cells: Cell[]
    txId: string
    playerToPlay: number
}

export interface Cell {
    row: number
    col: number
    value: string
}

export interface Player {
    name: string
    symbol: string
}
