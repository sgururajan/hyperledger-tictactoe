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
	gridWidth        = 3
)

type GameIdCounter struct {
	CurrentValue int
}

type Cell struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
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

type TictactoeGame struct {
}

func (m *TictactoeGame) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger := shim.NewLogger("init")
	logger.Info("intializing tictactoe chaincode")
	err := m.initializeGameIdCounter(stub)
	if err != nil {
		logger.Errorf("error while intializing game id counter. Err: %#v", err)
	}
	return shim.Success(nil)
}

func (m *TictactoeGame) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger := shim.NewLogger("Invoke")
	function, args := stub.GetFunctionAndParameters()

	logger.Infof("received invode method %s", function)
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
	case "makemove":
		return m.makeMove(stub, args)
	default:
		return shim.Error("invalid invoke method " + function)
	}

	return shim.Success(nil)
}

func (m *TictactoeGame) makeMove(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger := shim.NewLogger("makeMove")
	if len(args) < 4 {
		return shim.Error("not enough arguments. at least 4 is expected")
	}

	gameId, err := strconv.Atoi(args[0])
	if err != nil || gameId < 1 {
		return shim.Error(errorWithMessage("invalid game id (args[0])", err))
	}

	orgName := args[1]
	if orgName == "" {
		return shim.Error("orgname cannot be empty (args[1])")
	}

	row, err := strconv.Atoi(args[2])
	if err != nil || row < 0 || row > 2 {
		return shim.Error(errorWithMessage("invalid row number (args[2])", err))
	}

	col, err := strconv.Atoi(args[3])
	if err != nil || col < 0 || col > 2 {
		return shim.Error(errorWithMessage("invalid col number (args[3])", err))
	}

	game, err := getGameByKey(fmt.Sprintf("%s%d", gameKeyPrefix, gameId), stub)

	if err != nil {
		return shim.Error(errorMessage(err))
	}

	if game.Id < 1 {
		return shim.Error(fmt.Sprintf("game with id- %d not found", gameId))
	}

	if game.IsCompleted {
		return shim.Error(fmt.Sprintf("game %s is completed", gameId))
	}

	if !validateMove(game, row, col) {
		return shim.Error("invalid move. The played cell is not empty")
	}

	var symbol string
	var nextPlayer int
	if game.Players[0].Name == orgName {
		symbol = game.Players[0].Symbol
		nextPlayer = 1
	}
	if game.Players[1].Name == orgName {
		symbol = game.Players[1].Symbol
		nextPlayer = 0
	}

	if symbol == "" {
		return shim.Error(fmt.Sprintf("org %s is not a player in game %d", orgName, gameId))
	}

	game.Cells[gridWidth*row+col].Value = symbol
	winner := CheckForWinner(game)

	if winner != "" {
		game.Winner = winner
		game.IsCompleted = true
	}

	if checForDraw(game) {
		game.Winner=""
		game.IsCompleted=true
	}

	game.PlayerToPlayIndex = nextPlayer

	logger.Infof("update game state: %#v", game)

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	err = stub.PutState(fmt.Sprintf("%s%d", gameKeyPrefix, gameId), gameBytes)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	response, err := GenerateResponse(stub.GetTxID(), []Game{game})
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(response)
}

func (m *TictactoeGame) getGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger := shim.NewLogger("getGame")
	if len(args) < 1 {
		logger.Error("not enough arguments. expected at least 1")
		return shim.Error("not enough arguments. expected at least 1")
	}

	gameStateKey := fmt.Sprintf("%s%s", gameKeyPrefix, args[0])
	gameBytes, err := stub.GetState(gameStateKey)
	if err != nil {
		logger.Errorf("error while getting game. Err: %#v", err)
		return shim.Error(errorMessage(err))
	}

	if gameBytes == nil || len(gameBytes) == 0 {
		logger.Infof("game with id \"%s\" does not exists", args[0])
		return shim.Error(fmt.Sprintf("game with id \"%s\" does not exists", args[0]))
	}

	//game,err:= getGameObjFromBytes(gameBytes)
	//if err != nil {
	//	return shim.Error(errorMessage(err))
	//}

	return shim.Success(gameBytes)
}

func (m *TictactoeGame) newGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger := shim.NewLogger("newGame")

	if len(args) < 1 {
		return shim.Error("not enough arguments. expected at least 1")
	}

	gameId, err := m.getNewGameId(stub)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	game := m.CreateNewGame(args[0], gameId)

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	logger.Infof("game created with id: %v", gameId)

	stub.PutState(gameKeyPrefix+strconv.Itoa(gameId), gameBytes)

	logger.Infof("updated game state with key %s", gameKeyPrefix+strconv.Itoa(gameId))

	response, err := GenerateResponse(stub.GetTxID(), []Game{game})
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

	result, err := GenerateResponse(stub.GetTxID(), gameList)
	if err != nil {
		return shim.Error(errorMessage(err))
	}
	return shim.Success(result)
}

func (m *TictactoeGame) getAllGames(stub shim.ChaincodeStubInterface) peer.Response {
	logger := shim.NewLogger("getAllGames")
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

	result, err := GenerateResponse(stub.GetTxID(), gameList)
	if err != nil {
		logger.Error(errorMessage(err))
		return shim.Error(errorMessage(err))
	}

	return shim.Success(result)
}

func (m *TictactoeGame) joinGame(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger := shim.NewLogger("joinGame")
	// args[0] = gameid (integer), args[1] = second player name (mostly org name)
	if len(args) < 2 {
		return shim.Error("not enough arguments. expected at least 2")
	}

	logger.Infof("arguments recieved: %s, %s", args[0], args[1])

	if args[0] == "" {
		return shim.Error("game id cannot be empty (args[0])")
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

	if game.Players[0].Name == otherPlayerName {
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

	response, err := GenerateResponse(stub.GetTxID(), []Game{game})
	if err != nil {
		return shim.Error(errorMessage(err))
	}

	return shim.Success(response)
}

func (m *TictactoeGame) CreateNewGame(initPlayer string, gameId int) Game {
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

	fmt.Println(game)

	return game
}

func (m *TictactoeGame) initializeGameIdCounter(stub shim.ChaincodeStubInterface) error {
	gameIdCounterBytes, err := stub.GetState(gameIdCounterKey)
	if err != nil {
		return err
	}

	if gameIdCounterBytes == nil {
		gameIdCounter := GameIdCounter{
			CurrentValue: 1,
		}

		stateBytes, err := json.Marshal(gameIdCounter)
		if err != nil {
			return err
		}

		err = stub.PutState(gameIdCounterKey, stateBytes)
		if err != nil {
			return err
		}
	}

	return nil
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

func getGameByKey(key string, stub shim.ChaincodeStubInterface) (Game, error) {
	logger := shim.NewLogger("getGameByKey")
	logger.Infof("getting game for key: %s", key)
	gameBytes, err := stub.GetState(key)
	if err != nil {
		return Game{}, err
	}

	game, err := getGameObjFromBytes(gameBytes)
	if err != nil {
		return Game{}, err
	}

	logger.Infof("found game with id: %d", game.Id)

	return game, nil
}

func getGameListFromStartAndEndKey(startKey, endKey string, stub shim.ChaincodeStubInterface) ([]Game, error) {
	logger := shim.NewLogger("getGameListFromStartAndEndKey")
	logger.Infof("getting games for start key: %s and end key: %s", startKey, endKey)

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

	logger.Infof("%d games found", len(gameList))

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

func GenerateResponse(txId string, payload []Game) ([]byte, error) {
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

func validateMove(game Game, row, col int) bool {
	valid := false
	valid = game.Cells[gridWidth*row+col].Value == ""
	return valid
}

func CheckForWinner(game Game) string {
	hasWinner := false
	var a, b, c string
	for i := 0; i < gridWidth; i++ {
		a, b, c = game.Cells[gridWidth*i+0].Value, game.Cells[gridWidth*i+1].Value, game.Cells[gridWidth*i+2].Value
		//hasWinner = (game.Cells[gridWidth*i+0].Value == game.Cells[gridWidth*i+1].Value) && (game.Cells[gridWidth*i+0].Value == game.Cells[gridWidth*i+2].Value)
		hasWinner, winner := findWinner(game, a, b, c)
		if hasWinner {
			return winner
		}
	}

	for j := 0; j < gridWidth; j++ {
		//hasWinner = (game.Cells[gridWidth*0+j].Value == game.Cells[gridWidth*1+j].Value) && (game.Cells[gridWidth*0+j].Value == game.Cells[gridWidth*2+j].Value)
		a, b, c = game.Cells[gridWidth*0+j].Value, game.Cells[gridWidth*1+j].Value, game.Cells[gridWidth*2+j].Value
		hasWinner, winner := findWinner(game, a, b, c)
		if hasWinner {
			return winner
		}
	}

	a, b, c = game.Cells[gridWidth*0+0].Value, game.Cells[gridWidth*1+1].Value, game.Cells[gridWidth*2+2].Value
	hasWinner, winner := findWinner(game, a, b, c)
	if hasWinner {
		return winner
	}

	a, b, c = game.Cells[gridWidth*0+2].Value, game.Cells[gridWidth*1+1].Value, game.Cells[gridWidth*2+0].Value
	hasWinner, winner = findWinner(game, a, b, c)
	if hasWinner {
		return winner
	}

	return ""
}

func checForDraw(game Game) bool {
	isDraw:= true

	for _,c:= range game.Cells {
		isDraw = isDraw && c.Value!=""
	}

	return isDraw
}

func findWinner(game Game, a, b, c string) (bool, string) {
	if a == "" && b == "" && c == "" {
		return false, ""
	}

	hasWinner := (a == b) && (a == c)
	if hasWinner {
		return hasWinner, findWinningPlayerNameBySymbol(game, a)
	}
	return false, ""
}

func findWinningPlayerNameBySymbol(game Game, symbol string) string {
	if strings.EqualFold(game.Players[0].Symbol, symbol) {
		return game.Players[0].Name
	}

	if strings.EqualFold(game.Players[1].Symbol, symbol) {
		return game.Players[1].Name
	}

	return ""
}
