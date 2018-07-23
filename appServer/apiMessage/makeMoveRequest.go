package apiMessage

type MakeMoveRequest struct {
	GameId int `json:"gameId"`
	Row    int `json:"row"`
	Column int `json:"column"`
}
