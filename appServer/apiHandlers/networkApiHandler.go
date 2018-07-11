package apiHandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiMessage"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"net/http"
	"github.com/pkg/errors"
)

type NetworkApiHandler struct {
	repo           database.NetworkRepository
	networkHandler *networkHandlers.NetworkHandler
}

func NewNetworkAPIHandler(r database.NetworkRepository, nwhandler *networkHandlers.NetworkHandler) *NetworkApiHandler {
	return &NetworkApiHandler{
		repo:           r,
		networkHandler: nwhandler,
	}
}

func (m *NetworkApiHandler) RegisterRoutes(r *mux.Router) {
	for _, route := range m.getRoutes() {
		r.Methods(route.Method).Path(route.Pattern).Name(route.Name).HandlerFunc(route.HandlerFunc)
	}
}

func (m *NetworkApiHandler) getNetworksInfo(w http.ResponseWriter, r *http.Request) {
	networks, err := m.repo.GetNetworks()
	if err != nil {
		HandleServerError(w, errors.WithMessage(err, "unable to get networks from repo"))
		return
	}

	var networkInfos []apiMessage.NetworkInfo
	for k, v := range networks {
		fabNetwork, err := m.networkHandler.GetNetwork(k)
		if err != nil {
			HandleServerError(w,errors.WithMessage(err, "error getting network"))
			return
		}

		chList, err := fabNetwork.GetChannelList()
		if err != nil {
			HandleServerError(w, errors.WithMessage(err, "error getting channel list"))
			return
		}

		peers, err := m.repo.GetPeers(k)
		if err != nil {
			HandleServerError(w, errors.WithMessage(err, "error getting peer list"))
			return
		}

		nInfo := apiMessage.NetworkInfo{
			Name:              k,
			NoOfBlocks:        0,
			NoOfChannels:      len(chList),
			NoOfOrganizations: len(v.Organizations),
			NoOfPeers:         len(peers),
		}

		networkInfos = append(networkInfos, nInfo)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(networkInfos)
}

func (m *NetworkApiHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy. I am"))
}

func (m *NetworkApiHandler) getRoutes() []Route {
	return []Route{
		{
			Method:      http.MethodGet,
			Name:        "NetworksInfo",
			Pattern:     "/networksinfo",
			HandlerFunc: m.getNetworksInfo,
		},
		{
			Method:      http.MethodGet,
			Name:        "HealthCheck",
			Pattern:     "/hc",
			HandlerFunc: m.healthCheck,
		},
	}
}
