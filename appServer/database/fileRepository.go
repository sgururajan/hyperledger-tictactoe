package database

import (
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"io/ioutil"
	"encoding/json"
)

type NetworkFileRepository struct {
	fileName string
	networks map[string]Network
}

func NewNetworkFileRepository(fName string) NetworkRepository {
	repo:= &NetworkFileRepository{
		fileName:fName,
	}

	var networks map[string]Network
	cBytes, err:= ioutil.ReadFile(fName)
	if err!=nil {
		panic(err)
	}

	err = json.Unmarshal(cBytes, &networks)
	if err!=nil {
		panic(err)
	}

	repo.networks = networks
	return repo
}

func (m *NetworkFileRepository) GetNetworks() (map[string]Network, error) {
	return m.networks, nil
}

func (m *NetworkFileRepository) GetOrganizations(networkName string) ([]entities.Organization, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	orgs:= []entities.Organization{}
	for k,v:= range network.Organizations {
		orgPeers, _:= m.GetPeersForOrgId(networkName, k)
		peers:= getPeerNameSliceFromPeers(orgPeers)
		org:= entities.Organization{
			CertificateAuthorities: v.CertificateAuthorities,
			CryptoPath:v.CryptoPath,
			Peers: peers,
			MSPID: v.MSPID,
			Name: v.Name,
			ID: v.ID,
			MSPDir:v.MSPDir,
			AdminUser: v.AdminUser,
		}

		orgs = append(orgs, org)
	}

	return orgs, nil

}

func (m *NetworkFileRepository) GetPeersForOrgId(networkName, orgId string) ([]entities.Peer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Peer{}
	for _,v:= range network.Peers {
		if v.Organization == orgId {
			p:= getFabNetworkPeer(v)
			result = append(result, p)
		}
	}

	return result, nil
}

func (m *NetworkFileRepository) GetEndrosingPeersForOrgId(networkName, orgId string) ([]entities.Peer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Peer{}
	for _,v:= range network.Peers {
		if v.Organization == orgId && v.IsEndrosingPeer {
			result = append(result, getFabNetworkPeer(v))
		}
	}

	return result, nil
}

func (m *NetworkFileRepository) GetChainCodePeersForOrgId(networkName, orgId string) ([]entities.Peer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Peer{}
	for _,v:= range network.Peers {
		if v.Organization == orgId && v.IsChainCodePeer {
			result = append(result, getFabNetworkPeer(v))
		}
	}

	return result, nil
}

func (m *NetworkFileRepository) GetPeers(networkName string)([]entities.Peer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Peer{}
	for _,v:= range network.Peers {
		result = append(result, getFabNetworkPeer(v))
	}

	return result, nil
}

func (m *NetworkFileRepository) GetOrderers(networkName string) ([]entities.Orderer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Orderer{}
	for _,v:= range network.Orderers {
		result= append(result, getFabNetworkOrderer(v))
	}

	return result, nil
}

func (m *NetworkFileRepository) GetOrderersForOrgId(networkName, orgId string) ([]entities.Orderer, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.Orderer{}
	for _,v:= range network.Orderers {
		if v.Organization==orgId {
			result= append(result, getFabNetworkOrderer(v))
		}
	}

	return result, nil
}

func (m *NetworkFileRepository) GetSecurityConfig(networkName string)(entities.SecurityConfiguration, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return entities.SecurityConfiguration{}, errors.Errorf("network %s does not exists", networkName)
	}

	return entities.SecurityConfiguration{
		KeyStoragePath: network.SecurityConfiguration.KeyStoragePath,
		ProviderLabel:network.SecurityConfiguration.ProviderLabel,
		ProviderPin: network.SecurityConfiguration.ProviderPin,
		ProviderLibPath: network.SecurityConfiguration.ProviderLibPath,
		IsSoftVerify: network.SecurityConfiguration.IsSoftVerify,
		Provider: network.SecurityConfiguration.Provider,
		Level: network.SecurityConfiguration.Level,
		Algorithm: network.SecurityConfiguration.Algorithm,
		IsEnabled: network.SecurityConfiguration.IsEnabled,
	}, nil
}

func (m *NetworkFileRepository) GetCertificateAuthorities(networkName string)([]entities.CertificateAuthority, error){
	network, ok:= m.networks[networkName]
	if !ok {
		return nil, errors.Errorf("network %s does not exists", networkName)
	}

	result:= []entities.CertificateAuthority{}
	for _,v:= range network.CertificateAuthority {
		result = append(result, getFabNetworkCertificateAutority(v))
	}

	return result, nil
}


func getPeerNameSliceFromPeers(peers []entities.Peer) []string {
	result:= []string{}
	for _,v:= range peers {
		result = append(result, v.EndPoint)
	}

	return result
}

func getFabNetworkPeer(peer Peer) entities.Peer {
	return entities.Peer{
		TLSCertPath: peer.TLSCertPath,
		URL: peer.URL,
		GrpcOptions: peer.GrpcOptions,
		EventURL: peer.EventURL,
		EndPoint: peer.EndPoint,
		IsEndrosingPeer: peer.IsEndrosingPeer,
		CanQueryLedger: peer.CanQueryLedger,
		IsChainCodePeer: peer.IsChainCodePeer,
		EventSource: peer.EventSource,
	}
}

func getFabNetworkOrderer(o Orderer) entities.Orderer {
	return entities.Orderer{
		MSPDir: o.MSPDir,
		Name: o.Name,
		MSPID: o.MSPID,
		URL: o.URL,
		CryptoPath:o.CryptoPath,
		GRPCOptions:o.GRPCOptions,
		TLSCACertPath:o.TLSCACertPath,
	}
}

func getFabNetworkCertificateAutority(v CertificateAuthority) entities.CertificateAuthority {
	return entities.CertificateAuthority{
		URL: v.URL,
		TLSCertPath:v.TLSCertPath,
		CryptoConfigPath:v.CryptoConfigPath,
		CAName: v.CAName,
		RegistrarCredential: entities.Credential{
			ID: v.RegistrarCredential.ID,
			Secret: v.RegistrarCredential.Secret,
		},
		TLSCertClientPaths: entities.TLSKeyPathPair{
			CertPath: v.TLSCertClientPaths.CertPath,
			KeyPath: v.TLSCertClientPaths.KeyPath,
		},
	}
}
