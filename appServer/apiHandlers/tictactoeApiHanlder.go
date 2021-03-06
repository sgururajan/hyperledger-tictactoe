package apiHandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"github.com/sgururajan/hyperledger-tictactoe/domainModel"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork"
	"net/http"
	"strconv"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiMessage"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
)

const (
	channelName       = "tictactoechannel"
	chainCodeName     = "tictactoe"
	invodeCmd         = "invoke"
	networkNameHeader = "X-hlt3-networkName"
	orgNameHeader     = "X-hlt3-orgName"
)

type TictactoeApiHandler struct {
	repo            database.NetworkRepository
	networkHandlers *networkHandlers.NetworkHandler
	broadcastMessage BroadcastMessage
}

type BroadcastMessage func(string, interface{})

func NewTictactoeApiHandler(repo database.NetworkRepository, nHandlers *networkHandlers.NetworkHandler, broadcastFunc BroadcastMessage) *TictactoeApiHandler {
	return &TictactoeApiHandler{
		repo:            repo,
		networkHandlers: nHandlers,
		broadcastMessage: broadcastFunc,
	}
}

func (m *TictactoeApiHandler) RegisterRoutes(router *mux.Router) {
	for _, r := range m.getRoutes() {
		router.Methods(r.Method).Path(r.Pattern).Name(r.Name).HandlerFunc(r.HandlerFunc)
	}
}

func getNetworkNameFromReqHeader(request *http.Request) string {
	return request.Header.Get(networkNameHeader)
}

func getOrgNameFromReqHeader(request *http.Request) string {
	return request.Header.Get(orgNameHeader);
}

func (m *TictactoeApiHandler) addGame(writer http.ResponseWriter, request *http.Request) {
	networkName := getNetworkNameFromReqHeader(request)
	if networkName == "" {
		HandleServerError(writer, errors.New("empty network name"))
		return
	}

	network, err := m.networkHandlers.GetNetwork(networkName)
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "error getting fabric network instance"))
		return
	}
	// for now get the orgname from header
	orgName := getOrgNameFromReqHeader(request)
	if orgName == "" {
		HandleServerError(writer, errors.New("orgname not found in header"))
		return
	}

	payload, err := m.getTicTacToeGameResponseFromChaincode(network, orgName, "creategame", orgName)
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	gameList := payload.Games

	m.broadcastMessage("gameadded", gameList)

	writer.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) getGameList(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	pageIndex, err := strconv.Atoi(vars["pageindex"])
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "invalid page index value. must be integer"))
		return
	}

	pageSize, err := strconv.Atoi(vars["pagesize"])
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "invalid page size value. must be integer"))
		return
	}

	networkName := getNetworkNameFromReqHeader(request)
	if networkName == "" {
		HandleServerError(writer, errors.New("empty network name"))
		return
	}

	network, err := m.networkHandlers.GetNetwork(networkName)
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "error getting fabric network instance"))
		return
	}

	repoOrgs, err := m.repo.GetOrganizations(networkName)
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "no organizations found for network "+networkName))
		return
	}
	if len(repoOrgs) == 0 {
		HandleServerError(writer, errors.New("no organizations found for network "+networkName))
		return
	}

	// for now get the orgname from header
	orgName := getOrgNameFromReqHeader(request)
	if orgName == "" {
		orgName = repoOrgs[0].Name
	}

	payload, err := m.getTicTacToeGameResponseFromChaincode(network, orgName, "getgameslist", strconv.Itoa(pageIndex), strconv.Itoa(pageSize))
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	gameList := payload.Games
	writer.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) getAllGamesList(writer http.ResponseWriter, request *http.Request) {
	logger:= utils.NewAppLogger("getAllGamesList", "")
	networkName := getNetworkNameFromReqHeader(request)
	if networkName == "" {
		HandleServerError(writer, errors.New("empty network name"))
		return
	}

	network, err := m.networkHandlers.GetNetwork(networkName)
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, "error getting fabric network instance"))
		return
	}

	orgName := getOrgNameFromReqHeader(request)
	if orgName == "" {
		HandleServerError(writer, errors.New("org name not found in header"))
		return
	}

	payload, err := m.getTicTacToeGameResponseFromChaincode(network, orgName, "getallgames")
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	gameList := payload.Games
	logger.Infof("%#v", gameList)
	writer.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) joinGame(writer http.ResponseWriter, request *http.Request) {
	networkName:= getNetworkNameFromReqHeader(request)
	if networkName=="" {
		HandleServerError(writer, errors.New("empty network name"))
		return
	}

	orgName:= getOrgNameFromReqHeader(request)
	if orgName=="" {
		HandleServerError(writer, errors.New("org name not found in header"))
		return
	}

	vars:= mux.Vars(request)
	gameIdStr:= vars["gameid"]
	if gameIdStr=="" {
		HandleServerError(writer, errors.New("not a valid game id"))
		return
	}

	network,err:= m.networkHandlers.GetNetwork(networkName)
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	payload,err:= m.getTicTacToeGameResponseFromChaincode(network, orgName, "joingame", gameIdStr, orgName)
	if err != nil {
		HandleServerError(writer, err)
		return
	}
	gameList:= payload.Games

	m.broadcastMessage("gameupdated", gameList)

	writer.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) makeMove(writer http.ResponseWriter, request *http.Request) {
	logger:= utils.NewAppLogger("makeMove","")

	logger.Infof("request body: %#v", request.Body)

	decoder:= json.NewDecoder(request.Body)
	var moveRequest apiMessage.MakeMoveRequest
	err:= decoder.Decode(&moveRequest)
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	logger.Infof("received move request: %#v", moveRequest)

	networkName:= getNetworkNameFromReqHeader(request)
	orgName:= getOrgNameFromReqHeader(request)

	if networkName=="" {
		HandleServerError(writer, errors.New("invalid network name"))
		return
	}

	if orgName == "" {
		HandleServerError(writer, errors.New("invalid org name"))
		return
	}

	network, err:= m.networkHandlers.GetNetwork(networkName)
	if err != nil {
		HandleServerError(writer, err)
		return
	}

	payload, err:= m.getTicTacToeGameResponseFromChaincode(network,
		orgName,
		"makemove",
		strconv.Itoa(moveRequest.GameId),
		orgName,
		strconv.Itoa(moveRequest.Row),
		strconv.Itoa(moveRequest.Column))

	if err != nil {
		HandleServerError(writer, err)
		return
	}

	gameList:= payload.Games

	logger.Infof("update game states: %v", gameList)

	m.broadcastMessage("gameupdated", gameList)

	writer.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) getGame(writer http.ResponseWriter, request *http.Request) {

}

func (m *TictactoeApiHandler) getGameFromChannel(networkName, orgName string, gameId int) {

}

func (m *TictactoeApiHandler) getTicTacToeGameResponseFromChaincode(network *fabnetwork.FabricNetwork, orgName, cmd string, args ...string) (domainModel.TictactoeGameResponse, error) {
	result := domainModel.TictactoeGameResponse{}
	peerEndpoints, err := m.repo.GetEndrosingPeersEndpoints(network.Name)
	if err != nil {
		return result, err
	}
	response, err := network.ExecuteChainCode(orgName,
		channelName,
		chainCodeName,
		peerEndpoints,
		cmd,
		args...)

	if err != nil {
		return result, err
	}

	payload := domainModel.TictactoeGameResponse{}
	err = json.Unmarshal(response.Payload, &payload)
	if err != nil {
		return result, err
	}

	payload.TxId = response.TxId

	result = payload
	return result, nil
}

func (m *TictactoeApiHandler) getRoutes() []Route {
	return []Route{
		{
			Method:      http.MethodGet,
			Name:        "GetGameList",
			Pattern:     "/getGameList/{pageindex}/{pagesize}",
			HandlerFunc: m.getGameList,
		},
		{
			Method:      http.MethodGet,
			Name:        "GetAllGameList",
			Pattern:     "/getAllGameList",
			HandlerFunc: m.getAllGamesList,
		},
		{
			Method:      http.MethodPost,
			Name:        "AddGame",
			Pattern:     "/addgame",
			HandlerFunc: m.addGame,
		},
		{
			Method:      http.MethodPost,
			Name:        "JoinGame",
			Pattern:     "/joingame/{gameid}",
			HandlerFunc: m.joinGame,
		},
		{
			Method:      http.MethodPost,
			Name:        "MakeMove",
			Pattern:     "/makemove",
			HandlerFunc: m.makeMove,
		},
	}
}
