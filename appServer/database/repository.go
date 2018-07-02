package database

import (
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
)

type NetworkRepository interface {
	GetNetworks() (map[string]Network, error)
	GetOrganizations(networkName string) ([]entities.Organization, error)
	GetPeersForOrgId(networkName, orgId string) ([]entities.Peer, error)
	GetEndrosingPeersForOrgId(networkName, orgId string) ([]entities.Peer, error)
	GetChainCodePeersForOrgId(networkName, orgId string) ([]entities.Peer, error)
	GetPeers(networkName string)([]entities.Peer, error)
	GetOrderers(networkName string) ([]entities.Orderer, error)
	GetOrderersForOrgId(networkName, orgId string) ([]entities.Orderer, error)
	GetSecurityConfig(networkName string)(entities.SecurityConfiguration, error)
	GetCertificateAuthorities(networkName string)([]entities.CertificateAuthority, error)
}

