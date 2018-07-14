package domainModel

type Cell struct {
	Row    int
	Column int
	Value  string
}

type Player struct {
	Name   string
	Symbol string
}

type Game struct {
	Id                int
	IsCompleted       bool
	Players           [2]Player
	PlayerToPlayIndex int
	Winner            string
	Cells             [9]Cell
}

type TictactoeGameResponse struct {
	TxId  string `json:"txid,omitempty"`
	Games []Game `json:"games"`
}
