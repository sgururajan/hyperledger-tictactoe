package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
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

type GameIdCounter struct {
	CurrentValue int
}

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

type TictactoeGame struct {
}

func (m *TictactoeGame) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger:= shim.NewLogger("init")
	logger.Info("intializing tictactoe chaincode")
	err:= m.initializeGameIdCounter(stub)
	if err != nil {
		logger.Errorf("error while intializing game id counter. Err: %#v", err)
	}
	return shim.Success(nil)
}

func (m *TictactoeGame) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger := shim.NewLogger("Invoke")
	function, args := stub.GetFunctionAndParameters()

	logger.Infof("received args: %#v", args)

	function = strings.ToLower(function)
	switch function {
	case "creategame":
		return m.newGame(stub, args)
	case "getgameslist":
		return m.getGamesList(stub, args)
	case "getallgames":
		return m.getAllGames(stub)
	case "joingame":
		return m.joinGame(stub, args)
	default:
		return shim.Error("invalid invoke method " + function)
	}

	return shim.Success(nil)
}

func (m *TictactoeGame) getGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger:= shim.NewLogger("getGame")
	if len(args)<1 {
		logger.Error("not enough arguments. expected at least 1")
		return shim.Error("not enough arguments. expected at least 1")
	}

	gameStateKey:= fmt.Sprintf("%s%s", gameKeyPrefix, args[0])
	gameBytes, err:= stub.GetState(gameStateKey)
	if err != nil {
		logger.Errorf("error while getting game. Err: %#v", err)
		return shim.Error(errorMessage(err));
	}

	if gameBytes==nil || len(gameBytes)==0 {
		logger.Infof("game with id \"%s\" does not exists", args[0])
		return shim.Error(fmt.Sprintf("game with id \"%s\" does not exists", args[0]))
	}

	return shim.Success(nil)
}

func (m *TictactoeGame) newGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger:= shim.NewLogger("newGame")

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

	logger.Infof("game created with id: %v", gameId)

	stub.PutState(gameKeyPrefix+strconv.Itoa(gameId), gameBytes)

	logger.Infof("updated game state with key %s", gameKeyPrefix+strconv.Itoa(gameId))

	response, err := generateResponse(stub.GetTxID(), []Game{game})
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(response)
}

func (m *TictactoeGame) getGamesList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger := shim.NewLogger("getGamesList")
	if len(args) < 2 {
		return shim.Error("not enough arguments. expected at least 2")
	}

	pageIndex, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error(errorWithMessage("invalid page index (args[0]). got "+args[0], err))
	}

	pageSize, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(errorWithMessage("invalid page size (args[1]). got"+args[1], err))
	}

	logger.Infof("received params pageIndex: %d, pageSize: %d", pageIndex, pageSize)
	if pageIndex < 1 {
		return shim.Error("pageIndex should be a positive integer greater than 0")
	}

	startKey := fmt.Sprintf("%s%d", gameKeyPrefix, (pageIndex-1)*pageSize)
	endKey := fmt.Sprintf("%s%d", gameKeyPrefix, ((pageIndex-1)*pageSize)+pageSize)

	logger.Infof("startKey: %s", startKey)
	logger.Infof("endKey: %s", endKey)

	gameList, err := getGameListFromStartAndEndKey(startKey, endKey, stub)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	result, err := generateResponse(stub.GetTxID(), gameList)
	if err != nil {
		return shim.Error(errorMessage(err))
	}
	return shim.Success(result)
}

func (m *TictactoeGame) getAllGames(stub shim.ChaincodeStubInterface) peer.Response {
	logger:= shim.NewLogger("getAllGames")
	logger.Info("getAllGames method invoked")
	gameIdBytes, err := stub.GetState(gameIdCounterKey)
	if err != nil {
		return shim.Error(errorMessage(err))
	}
	logger.Infof("gameIdCounter bytes: %#v", gameIdBytes)

	counter, err := getGameIdCounterObjFromBytes(gameIdBytes)
	if err != nil {
		logger.Error(errorMessage(err))
		return shim.Error(errorMessage(err))
	}

	startKey := gameKeyPrefix + "1"
	endKey := gameKeyPrefix + strconv.Itoa(counter.CurrentValue)

	logger.Infof("getting games list with starting key: %s and ending key: %s", startKey, endKey)

	gameList, err := getGameListFromStartAndEndKey(startKey, endKey, stub)
	if err != nil {
		logger.Error(errorMessage(err))
		return shim.Error(errorMessage(err))
	}

	result, err := generateResponse(stub.GetTxID(), gameList)
	if err != nil {
		logger.Error(errorMessage(err))
		return shim.Error(errorMessage(err))
	}

	return shim.Success(result)
}

func (m *TictactoeGame) joinGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger:= shim.NewLogger("joinGame")
	// args[0] = gameid (integer), args[1] = second player name (mostly org name)
	if len(args) < 2 {
		return shim.Error("not enough arguments. expected at least 2")
	}

	logger.Infof("arguments recieved: %s, %s", args[0], args[1])

	if args[0]=="" {
		return shim.Error("game id cannot be empty (args[0])");
	}

	_, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error(errorWithMessage("gameid should be integer", err))
	}

	game := Game{}
	gameBytes, err := stub.GetState(gameKeyPrefix + args[0])
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	err = json.Unmarshal(gameBytes, &game)
	if err != nil {
		return shim.Error(errorWithMessage("unable to unmarshal game data", err))
	}

	otherPlayerName := args[1]
	if otherPlayerName == "" {
		return shim.Error("joining player name cannot be empty.")
	}

	if game.Players[0].Name==otherPlayerName {
		return shim.Error("this player already joined the game")
	}

	game.Players[1].Name = otherPlayerName
	game.Players[1].Symbol = symbolX

	gameBytes, err = json.Marshal(game)
	if err != nil {
		return shim.Error(errorWithMessage("unable to marshal game data", err))
	}

	err = stub.PutState(gameKeyPrefix+args[0], gameBytes)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	logger.Info("updated game with other player")

	response, err:= generateResponse(stub.GetTxID(), []Game{game})
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(response)
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

func (m *TictactoeGame) initializeGameIdCounter(stub shim.ChaincodeStubInterface) error {
	gameIdCounterBytes, err := stub.GetState(gameIdCounterKey)
	if err != nil {
		return err
	}

	if gameIdCounterBytes==nil {
		gameIdCounter:= GameIdCounter{
			CurrentValue: 1,
		}

		stateBytes,err:= json.Marshal(gameIdCounter)
		if err != nil {
			return err
		}

		err=stub.PutState(gameIdCounterKey, stateBytes)
		if err != nil {
			return err
		}
	}

	return  nil
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

func getGameListFromStartAndEndKey(startKey, endKey string, stub shim.ChaincodeStubInterface) ([]Game, error) {
	gameList := []Game{}
	query, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	for query.HasNext() {
		v, err := query.Next()
		if err != nil {
			return nil, err
		}
		game, err := getGameObjFromBytes(v.Value)
		if err != nil {
			return nil, err
		}
		gameList = append(gameList, game)
	}

	return gameList, nil
}

func getGameIdCounterObjFromBytes(input []byte) (GameIdCounter, error) {
	counter := GameIdCounter{}
	err := json.Unmarshal(input, &counter)
	if err != nil {
		return GameIdCounter{}, err
	}
	return counter, nil
}

func getGameObjFromBytes(input []byte) (Game, error) {
	game := Game{}
	err := json.Unmarshal(input, &game)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func errorMessage(err error) string {
	return fmt.Sprintf("%#v", err)
}

func errorWithMessage(msg string, err error) string {
	return fmt.Sprintf("%s. err: %#v", msg, err)
}

func generateResponse(txId string, payload []Game) ([]byte, error) {
	response := TictactoeGameResponse{
		TxId:  txId,
		Games: payload,
	}

	result, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
