package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

func main() {
	err := shim.Start(new(TictactoeGame))
	if err != nil {
		fmt.Println("error starting tictactoe chaincode")
		fmt.Printf("%#v", err)
	}
}

const (
	gameIdCounterKey = "tictactoeGameIdCounterKey"
	gameKeyPrefix    = "tictactoeGame-"
	symbolX          = "X"
	symbolO          = "O"
)

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

type GameIdCounter struct {
	CurrentValue int
}

type TictactoeGameResponse struct {
	TxId string
	Payload interface{}
}

type TictactoeGame struct {
}

func (m *TictactoeGame) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (m *TictactoeGame) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "creategame" {
		return m.newGame(stub, args)
	}

	return shim.Success(nil)
}

func (m *TictactoeGame) newGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("not enough arguments. expected at least 1")
	}

	gameId, err := m.getNewGameId(stub)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	game := m.createNewGame(args[0], gameId)

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	stub.PutState(gameKeyPrefix+string(gameId), gameBytes)

	response,err:= generateResponse(stub.GetTxID(), gameId)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(response)
}

func (m *TictactoeGame) joinGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args[0] = gameid (integer), args[1] = second player name (mostly org name)
	if len(args) < 2 {
		return shim.Error("not enough arguments. expected at least 2")
	}

	gameId,err:= strconv.Atoi(args[0])
	if err != nil {
		return shim.Error(errorWithMessage("gameid should be integer", err))
	}

	game:= Game{}
	gameBytes, err:= stub.GetState(gameKeyPrefix+string(gameId))
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	err=json.Unmarshal(gameBytes, &game)
	if err != nil {
		return shim.Error(errorWithMessage("unable to unmarshal game data", err))
	}

	otherPlayerName:= args[1]
	if otherPlayerName=="" {
		return shim.Error("joining player name cannot be empty.")
	}

	game.Players[1].Name = otherPlayerName
	game.Players[1].Symbol = symbolX

	gameBytes, err = json.Marshal(game)
	if err != nil {
		return shim.Error(errorWithMessage("unable to marshal game data", err))
	}

	err = stub.PutState(gameKeyPrefix+string(gameId), gameBytes)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(nil)
}

func (m *TictactoeGame) createNewGame(initPlayer string, gameId int) Game {
	game := Game{
		Cells:       [9]Cell{},
		Players:     [2]Player{},
		Id:          gameId,
		IsCompleted: false,
	}

	game.Players[0] = Player{
		Name:   initPlayer,
		Symbol: symbolO,
	}
	game.PlayerToPlayIndex = 0

	width := 3
	for i := 0; i < width; i++ {
		for j := 0; j < width; j++ {
			game.Cells[width*i+j] = Cell{
				Row:    i,
				Column: j,
				Value:  "",
			}
		}
	}

	return game
}

func (m *TictactoeGame) getNewGameId(stub shim.ChaincodeStubInterface) (int, error) {
	gameIdCounterBytes, err := stub.GetState(gameIdCounterKey)
	if err != nil {
		return -1, err
	}

	gameIdCounter := GameIdCounter{}
	json.Unmarshal(gameIdCounterBytes, &gameIdCounter)
	gameId := gameIdCounter.CurrentValue
	gameIdCounter.CurrentValue = gameIdCounter.CurrentValue + 1

	stateBytes, err := json.Marshal(gameIdCounter)
	err = stub.PutState(gameIdCounterKey, stateBytes)
	if err != nil {
		return -1, err
	}

	return gameId, nil
}

func errorMessage(err error) string {
	return fmt.Sprintf("%#v", err)
}

func errorWithMessage(msg string, err error) string {
	return fmt.Sprintf("%s. err: %#v", msg, err)
}

func generateResponse(txId string, payload interface{}) ([]byte, error) {
	response:= TictactoeGameResponse{
		TxId: txId,
		Payload:payload,
	}

	result, err:= json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
