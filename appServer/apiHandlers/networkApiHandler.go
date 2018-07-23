package apiHandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiHandlers/viewModel"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiMessage"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"net/http"
)

const (
	contentTypeKey       = "Content-Type"
	contentTypeJsonValue = "application/json;charset=UTF-8"
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
			HandleServerError(w, errors.WithMessage(err, "error getting network"))
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

func (m *NetworkApiHandler) getOrganizations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	networkName := vars["network"]

	var orgsList []string
	orgs, err := m.repo.GetOrganizations(networkName)
	if err != nil {
		HandleServerError(w, err)
		return
	}

	for _, o := range orgs {
		if len(o.Peers) > 0 {
			orgsList = append(orgsList, o.Name)
		}
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(orgsList)
	if err != nil {
		HandleServerError(w, err)
	}
}

func (m *NetworkApiHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Write([]byte("Yoda says: Healthy. I am " + vars["network"]))
}

func (m *NetworkApiHandler) getNetworks(w http.ResponseWriter, r *http.Request) {
	var result []viewModel.NetworkViewModel
	networks, err := m.repo.GetNetworks()
	if err != nil {
		HandleServerError(w, err)
		return
	}

	for nk, _ := range networks {
		nmodel, err := m.getNetworkViewModel(nk)
		if err != nil {
			HandleServerError(w, err)
			return
		}

		result = append(result, nmodel)
	}

	w.Header().Set(contentTypeKey, contentTypeJsonValue)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleServerError(w, err)
	}
}

func (m *NetworkApiHandler) getNetworkViewModel(networkName string) (viewModel.NetworkViewModel, error) {
	orgs, err := m.getOrganizationViewModel(networkName)
	if err != nil {
		return viewModel.NetworkViewModel{}, err
	}
	channels, err := m.getChannelsViewModel(networkName)
	if err != nil {
		return viewModel.NetworkViewModel{}, err
	}

	vm := viewModel.NetworkViewModel{
		Name:          networkName,
		Organizations: orgs,
		Channels:      channels,
	}

	return vm, nil
}

func (m *NetworkApiHandler) getChannelsViewModel(networkName string) ([]viewModel.ChannelViewModel, error) {
	var result []viewModel.ChannelViewModel
	fnetwork, err := m.networkHandler.GetNetwork(networkName)
	if err != nil {
		return result, err
	}

	chList, err := fnetwork.GetChannelList()
	if err != nil {
		return result, err
	}

	for _, ch := range chList {
		cvm := viewModel.ChannelViewModel{
			Name: ch,
		}

		result = append(result, cvm)
	}

	return result, nil
}

func (m *NetworkApiHandler) getOrganizationViewModel(networkName string) ([]viewModel.OrganizationViewModel, error) {
	var result []viewModel.OrganizationViewModel
	orgs, err := m.repo.GetOrganizations(networkName)
	if err != nil {
		return result, err
	}

	for _, o := range orgs {
		if len(o.Peers) < 1 {
			continue
		}
		ovm := viewModel.OrganizationViewModel{
			Name: o.Name,
		}
		peers, err := m.getPeerViewModel(networkName, o.Name)
		if err != nil {
			return result, err
		}
		ovm.Peers = peers

		result = append(result, ovm)
	}

	return result, nil
}

func (m *NetworkApiHandler) getPeerViewModel(networkName, orgName string) ([]viewModel.PeerViewModel, error) {
	var result []viewModel.PeerViewModel
	peers, err := m.repo.GetPeersForOrgId(networkName, orgName)
	if err != nil {
		return result, err
	}

	for _, p := range peers {
		pvm := viewModel.PeerViewModel{
			Name: p.EndPoint,
			Url:  p.URL,
		}
		result = append(result, pvm)
	}

	return result, nil
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
			Name:        "Networks",
			Pattern:     "/networks",
			HandlerFunc: m.getNetworks,
		},
		{
			Method:      http.MethodGet,
			Name:        "HealthCheck",
			Pattern:     "/hc",
			HandlerFunc: m.healthCheck,
		},
		{
			Method:      http.MethodGet,
			Name:        "OrgsList",
			Pattern:     "/getorgslist/{network}",
			HandlerFunc: m.getOrganizations,
		},
	}
}
