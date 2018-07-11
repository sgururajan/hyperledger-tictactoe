package apiHandlers

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"net/http"
	"strconv"
	"github.com/sgururajan/hyperledger-tictactoe/domainModel"
	"encoding/json"
)

const (
	channelName   = "tictactoechannel"
	chainCodeName = "tictactoe"
	invodeCmd     = "invoke"
)

type TictactoeApiHandler struct {
	repo            database.NetworkRepository
	networkHandlers *networkHandlers.NetworkHandler
}

func NewTictactoeApiHandler(repo database.NetworkRepository, nHandlers *networkHandlers.NetworkHandler) *TictactoeApiHandler {
	return &TictactoeApiHandler{
		repo:            repo,
		networkHandlers: nHandlers,
	}
}

func (m *TictactoeApiHandler) RegiterRoutes(router *mux.Router) {
	for _, r := range m.getRoutes() {
		router.Methods(r.Method).Path(r.Pattern).Name(r.Name).HandlerFunc(r.HandlerFunc)
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

	networkName := vars["network"]
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
	orgName := request.Header.Get("X-hlt3-orgName")
	if orgName == "" {
		orgName = repoOrgs[0].Name
	}

	peerEndpoints, err := m.repo.GetEndrosingPeersEndpoints(networkName)
	if err != nil {
		HandleServerError(writer, errors.WithMessage(err, ""))
		return
	}

	response, err := network.ExecuteChainCode(orgName,
		channelName,
		chainCodeName,
		peerEndpoints,
		"getgameslist",
		strconv.Itoa(pageIndex),
		strconv.Itoa(pageSize))

	if err != nil {
		HandleServerError(writer, err)
		return
	}

	payload:= domainModel.TictactoeGameResponse{}
	gameList:= []domainModel.Game{}
	err = json.Unmarshal(response.Payload, &payload)
	if err != nil {
		HandleServerError(writer, err)
	}
	gameList=payload.Games
	err = json.NewEncoder(writer).Encode(gameList)
	if err != nil {
		HandleServerError(writer, err)
	}
}

func (m *TictactoeApiHandler) getRoutes() []Route {
	return []Route{
		{
			Method:      http.MethodGet,
			Name:        "GetGameList",
			Pattern:     "/getGameList/{network}/{pageindex}/{pagesize}",
			HandlerFunc: m.getGameList,
		},
	}
}
