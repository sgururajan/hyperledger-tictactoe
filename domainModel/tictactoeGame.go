package domainModel

type Cell struct {
	Row    int    `json:"row"`
	Column int    `json:"col"`
	Value  string `json:"value"`
}

type Player struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Game struct {
	Id                int       `json:"id"`
	IsCompleted       bool      `json:"completed"`
	Players           [2]Player `json:"players"`
	PlayerToPlayIndex int       `json:"playerToPlay"`
	Winner            string    `json:"winner"`
	Cells             [9]Cell   `json:"cells"`
}

type TictactoeGameResponse struct {
	TxId  string `json:"txid,omitempty"`
	Games []Game `json:"games"`
}
