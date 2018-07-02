package networkHandlers

import (
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/providers"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/pkg/errors"
)

var cryptoConfigPath = "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/crypto-config"

type NetworkHandler struct {
	repo                database.NetworkRepository
	fabNetworks         map[string]*fabnetwork.FabricNetwork
	fabNetworkProviders map[string]*providers.FabNetworkProvider
}

func NewNetworkHandler(repo database.NetworkRepository) (*NetworkHandler, error) {
	handler := &NetworkHandler{
		repo:                repo,
		fabNetworks:         make(map[string]*fabnetwork.FabricNetwork),
		fabNetworkProviders: make(map[string]*providers.FabNetworkProvider),
	}

	err:= handler.bootFabricNetworks()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (m *NetworkHandler) Close() {
	for _,v:= range m.fabNetworks {
		v.Close()
	}
}

func (m *NetworkHandler) GetNetwork(networkName string) (*fabnetwork.FabricNetwork, error) {
	if n,ok:= m.fabNetworks[networkName];ok {
		return n, nil
	}

	return nil, errors.Errorf("network with name \"%s\" does not exists", networkName)
}

func (m *NetworkHandler) bootFabricNetworks() error {
	networks, err := m.repo.GetNetworks()
	if err != nil {
		return err
	}

	for k, _ := range networks {
		m.fabNetworkProviders[k] = providers.NewFabNetworkProvider(k, m.repo)
		m.fabNetworks[k], err = fabnetwork.NewFabricNetwork(getClientConfigurtion(),
				fabnetwork.WithOrdererProvider(m.fabNetworkProviders[k]),
				fabnetwork.WithOrganizationProvider(m.fabNetworkProviders[k]),
				fabnetwork.WithCAProvider(m.fabNetworkProviders[k]),
				fabnetwork.WithPeerProvider(m.fabNetworkProviders[k]),
				fabnetwork.WithSecurityConfig(m.fabNetworkProviders[k]))
		if err != nil {
			return err
		}
	}

	return nil
}

func getClientConfigurtion() entities.ClientConfiguration {
	return entities.ClientConfiguration{
		Logging:             entities.DEBUG,
		CryptoConfigPath:    cryptoConfigPath,
		CredentialStorePath: "/tmp/store",
		TLSKeyPair: entities.TLSKeyPathPair{
			KeyPath:  "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-key.pem",
			CertPath: "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
		},
	}
}
