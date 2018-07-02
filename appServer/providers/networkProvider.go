package providers

import (
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
)

type FabNetworkProvider struct {
	networkName string
	repository database.NetworkRepository
}

func NewFabNetworkProvider(name string, repo database.NetworkRepository) *FabNetworkProvider {
	return &FabNetworkProvider{
		networkName:name,
		repository: repo,
	}
}

func (m *FabNetworkProvider) GetCertificateAuthorities() ([]entities.CertificateAuthority, error){
	return m.repository.GetCertificateAuthorities(m.networkName)
}

func (m *FabNetworkProvider) GetOrderers() ([]entities.Orderer, error){
	return m.repository.GetOrderers(m.networkName)
}

func (m *FabNetworkProvider) GetOrderersForOrgId(orgId string) ([]entities.Orderer, error){
	return m.repository.GetOrderersForOrgId(m.networkName, orgId)
}

func (m *FabNetworkProvider) GetOrganizations() ([]entities.Organization, error){
	return m.repository.GetOrganizations(m.networkName)
}

func (m *FabNetworkProvider) GetPeersForOrgId(orgId string)([]entities.Peer, error){
	return m.repository.GetPeersForOrgId(m.networkName, orgId)
}

func (m *FabNetworkProvider) GetEndrosingPeersForOrgId(orgId string) ([]entities.Peer, error){
	return m.repository.GetEndrosingPeersForOrgId(m.networkName, orgId)
}

func (m *FabNetworkProvider) GetChainCodePeersForOrgId(orgId string) ([]entities.Peer, error){
	return m.repository.GetChainCodePeersForOrgId(m.networkName, orgId)
}

func (m *FabNetworkProvider) GetPeers() ([]entities.Peer, error){
	return m.repository.GetPeers(m.networkName)
}

func (m *FabNetworkProvider) IsSecurityEnabled() bool{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return true
	}
	return secConfig.IsEnabled
}

func (m *FabNetworkProvider) SecurityAlgorithm() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return "SHA2"
	}
	return secConfig.Algorithm
}

func (m *FabNetworkProvider) SecurityLevel() int{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return 256
	}
	return secConfig.Level
}

func (m *FabNetworkProvider) SecurityProvider() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return "sw"
	}
	return secConfig.Provider
}

func (m *FabNetworkProvider) SoftVerify() bool{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return true
	}
	return secConfig.IsSoftVerify
}

func (m *FabNetworkProvider) SecurityProviderLibPath() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return ""
	}
	return secConfig.ProviderLibPath
}

func (m *FabNetworkProvider) SecurityProviderPin() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return ""
	}
	return secConfig.ProviderPin
}

func (m *FabNetworkProvider) SecurityProviderLabel() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return ""
	}
	return secConfig.ProviderLabel
}

func (m *FabNetworkProvider) KeyStorePath() string{
	secConfig, err:= m.repository.GetSecurityConfig(m.networkName)
	if err!=nil {
		return "/tmp/msp/store"
	}
	return secConfig.KeyStoragePath
}
