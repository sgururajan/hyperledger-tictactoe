package apiHandlers

import (
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiMessage"
	"encoding/json"
)

type NetworkApiHandler struct {
	repo database.NetworkRepository
	networkHandler *networkHandlers.NetworkHandler
}

func NewNetworkAPIHandler(r database.NetworkRepository, nwhandler *networkHandlers.NetworkHandler ) *NetworkApiHandler {
	return &NetworkApiHandler{
		repo: r,
		networkHandler:nwhandler,
	}
}

func (m *NetworkApiHandler) RegisterRoutes(r *mux.Router) {
	for _, route:= range m.getRoutes() {
		r.Methods(route.Method).Path(route.Pattern).Name(route.Name).HandlerFunc(route.HandlerFunc)
}
}

func (m *NetworkApiHandler) getNetworksInfo(w http.ResponseWriter, r *http.Request) {
	networks, err:= m.repo.GetNetworks()
	if err!=nil {
		HandleServerError(w, "unable to get networks from repo", err)
		return
	}

	var networkInfos []apiMessage.NetworkInfo
	for k,v:= range networks {
		fabNetwork, err:= m.networkHandler.GetNetwork(k)
		if err!= nil {
			HandleServerError(w, "", err)
			return
		}

		chList,err:= fabNetwork.GetChannelList()
		if err != nil {
			HandleServerError(w, "", err)
			return
		}

		peers, err:= m.repo.GetPeers(k)
		if err != nil {
			HandleServerError(w, "", err)
			return
		}

		nInfo:= apiMessage.NetworkInfo{
			Name: k,
			NoOfBlocks:0,
			NoOfChannels:len(chList),
			NoOfOrganizations:len(v.Organizations),
			NoOfPeers:len(peers),
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
			Method: http.MethodGet,
			Name: "NetworksInfo",
			Pattern:"/networksinfo",
			HandlerFunc:m.getNetworksInfo,
		},
		{
			Method:http.MethodGet,
			Name:"HealthCheck",
			Pattern:"/hc",
			HandlerFunc:m.healthCheck,
		},
	}
}