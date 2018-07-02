package fabnetwork

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/cryptoutil"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type fabSdkImpl struct {
	organizations    map[string]fab.OrganizationConfig
	orderers         map[string]fab.OrdererConfig
	peers            map[string]fab.PeerConfig
	channels         map[string]fab.ChannelEndpointConfig
	isSystemCertPool bool
	tlsCerts         endpoint.MutualTLSConfig
	clientConfig     entities.ClientConfiguration
	rwLock           sync.RWMutex
}

func newFabSdkImpl(network *FabricNetwork, clientConfig entities.ClientConfiguration) (*fabSdkImpl, error) {
	impl := &fabSdkImpl{
		clientConfig: clientConfig,
		channels: make(map[string]fab.ChannelEndpointConfig),
	}
	err := impl.createOrdererConfig(network)
	if err != nil {
		return nil, err
	}
	err = impl.createPeersConfig(network)

	if err != nil {
		return nil, err
	}

	err = impl.createOrgsConfig(network)
	if err != nil {
		return nil, err
	}



	return impl, nil
}

func (m *fabSdkImpl) createOrgsConfig(network *FabricNetwork) error {
	orgs, err := network.providers.orgProvider.GetOrganizations()
	if err != nil {
		return err
	}

	m.organizations = make(map[string]fab.OrganizationConfig)

	for _, v := range orgs {
		m.organizations[v.Name] = fab.OrganizationConfig{
			MSPID:                  v.MSPID,
			Peers:                  v.Peers,
			CryptoPath:             v.CryptoPath,
			CertificateAuthorities: v.CertificateAuthorities,
		}
	}

	return nil
}

func (m *fabSdkImpl) createOrdererConfig(network *FabricNetwork) error {
	orderers, err := network.providers.ordererProvider.GetOrderers()
	if err != nil {
		return err
	}

	m.orderers = make(map[string]fab.OrdererConfig)

	for _, v := range orderers {
		m.orderers[v.Name] = fab.OrdererConfig{
			URL:         v.URL,
			GRPCOptions: v.GRPCOptions,
			TLSCACert:   tlsCertBytes(v.TLSCACertPath),
		}
	}

	return nil
}

func (m *fabSdkImpl) createPeersConfig(network *FabricNetwork) error {
	peers, err := network.providers.peerProvider.GetPeers()
	if err != nil {
		return err
	}

	m.peers = make(map[string]fab.PeerConfig)

	for _, v := range peers {
		m.peers[v.EndPoint] = fab.PeerConfig{
			URL:         v.URL,
			EventURL:    v.EventURL,
			GRPCOptions: v.GrpcOptions,
			TLSCACert:   tlsCertBytes(v.TLSCertPath),
		}
	}

	return nil
}

func (m *fabSdkImpl) Client() *msp.ClientConfig {
	m.tlsCerts = endpoint.MutualTLSConfig{
		Client: endpoint.TLSKeyPair{
			Key:  newTLSConfig(m.clientConfig.TLSKeyPair.KeyPath),
			Cert: newTLSConfig(m.clientConfig.TLSKeyPair.CertPath),
		},
	}
	return &msp.ClientConfig{
		Organization:    strings.ToLower(m.clientConfig.Organization),
		Logging:         getFabLoggingType(m.clientConfig.Logging),
		CryptoConfig:    msp.CCType{Path: m.clientConfig.CryptoConfigPath},
		CredentialStore: msp.CredentialStoreType{Path: m.clientConfig.CredentialStorePath},
		TLSCert:         m.tlsCerts.Client.Cert.Bytes(),
		TLSKey:          m.tlsCerts.Client.Key.Bytes(),
	}
}

func (m *fabSdkImpl) Timeout(tType fab.TimeoutType) time.Duration {
	t, ok := defaultTimeOutTypes[tType]
	if !ok {
		return time.Second * 30
	}

	return t
}

func (m *fabSdkImpl) EventServiceType() fab.EventServiceType {
	if os.Getenv("FABRIC_SDK_CLIENT_EVENTSERVICE_TYPE") == "eventhub" {
		return fab.EventHubEventServiceType
	}

	return fab.DeliverEventServiceType
}

func (m *fabSdkImpl) OrderersConfig() []fab.OrdererConfig {
	orderers := []fab.OrdererConfig{}

	for _, o := range m.orderers {
		if o.TLSCACert == nil && !m.isSystemCertPool {
			return nil
		}

		orderers = append(orderers, o)
	}

	return orderers
}

func (m *fabSdkImpl) OrdererConfig(nameOrUrl string) (*fab.OrdererConfig, bool) {
	orderer, ok := m.orderers[strings.ToLower(nameOrUrl)]
	if !ok {
		return nil, false
	}

	return &orderer, true
}

func (m *fabSdkImpl) PeersConfig(org string) ([]fab.PeerConfig, bool) {
	orgPeers := m.organizations[strings.ToLower(org)].Peers
	peers := []fab.PeerConfig{}

	for _, peerName := range orgPeers {
		p := m.peers[strings.ToLower(peerName)]
		if err := m.verifyPeerConfig(p, peerName, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil, false
		}

		peers = append(peers, p)
	}

	return peers, true
}

func (m *fabSdkImpl) PeerConfig(nameOrUrl string) (*fab.PeerConfig, bool) {
	pConfig, ok := m.peers[strings.ToLower(nameOrUrl)]
	if ok {
		return &pConfig, true
	}

	i := strings.Index(nameOrUrl, ":")
	if i > 0 {
		return m.PeerConfig(nameOrUrl[0:i])
	}

	return nil, false
}

func (m *fabSdkImpl) verifyPeerConfig(p fab.PeerConfig, peerName string, tlsEnabled bool) error {
	if p.URL == "" {
		return errors.Errorf("URL does not exists or empty for peer %s", peerName)
	}

	if tlsEnabled && p.TLSCACert == nil && !m.isSystemCertPool {
		return errors.Errorf("tls certificates does not exists or empty for peer %s", peerName)
	}

	return nil
}

func (m *fabSdkImpl) NetworkConfig() *fab.NetworkConfig {
	//return &m.networkConfig
	return &fab.NetworkConfig{
		Organizations: m.organizations,
		Orderers:      m.orderers,
		Peers:         m.peers,
		Channels:      m.channels,
	}
}

func (m *fabSdkImpl) NetworkPeers() []fab.NetworkPeer {
	netPeers := []fab.NetworkPeer{}

	for name, p := range m.peers {
		if err := m.verifyPeerConfig(p, name, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil
		}

		mspID, ok := m.PeerMSPID(name)
		if !ok {
			return nil
		}

		netPeer := fab.NetworkPeer{MSPID: mspID, PeerConfig: p}
		netPeers = append(netPeers, netPeer)
	}

	return netPeers
}

func (m *fabSdkImpl) PeerMSPID(name string) (string, bool) {
	for _, org := range m.organizations {
		for i := 0; i < len(org.Peers); i++ {
			if strings.EqualFold(org.Peers[i], name) {
				return org.MSPID, true
			}
		}
	}

	return "", false
}

func (m *fabSdkImpl) TLSClientCerts() []tls.Certificate {
	var clientCerts tls.Certificate
	cb := m.tlsCerts.Client.Cert.Bytes()

	if len(cb) == 0 {
		return []tls.Certificate{clientCerts}
	}

	cs := cryptosuite.GetDefault()
	pk, err := cryptoutil.GetPrivateKeyFromCert(cb, cs)

	if err != nil || pk == nil {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
		ccs, err := m.loadPrivateKeyFromConfig(cb)
		if err != nil {
			return nil
		}

		return ccs
	}

	clientCerts, err = cryptoutil.X509KeyPair(cb, pk, cs)
	if err != nil {
		return nil
	}

	return []tls.Certificate{clientCerts}
}

func (m *fabSdkImpl) CryptoConfigPath() string {
	return m.clientConfig.CryptoConfigPath
}

func (m *fabSdkImpl) loadPrivateKeyFromConfig(cb []byte) ([]tls.Certificate, error) {
	kb := m.tlsCerts.Client.Key.Bytes()
	clientCerts, err := tls.X509KeyPair(cb, kb)
	if err != nil {
		return nil, errors.Errorf("error loading cert/key pair as TLS client credential: %s", err)
	}

	return []tls.Certificate{clientCerts}, nil
}

func (m *fabSdkImpl) ChannelConfig(channelName string) (*fab.ChannelEndpointConfig, bool) {
	ch, ok := m.channels[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}
	return &ch, true
}

func (m *fabSdkImpl) ChannelPeers(channelName string) ([]fab.ChannelPeer, bool) {
	peers := []fab.ChannelPeer{}

	chConfig, ok := m.channels[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}

	for peerName, chPeerConfig := range chConfig.Peers {
		p, ok := m.PeerConfig(strings.ToLower(peerName))
		if !ok {
			return nil, false
		}

		if err := m.verifyPeerConfig(*p, peerName, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil, false
		}

		mspID, ok := m.PeerMSPID(peerName)
		if !ok {
			return nil, false
		}

		networkPeer := fab.NetworkPeer{PeerConfig: *p, MSPID: mspID}
		peer := fab.ChannelPeer{PeerChannelConfig: chPeerConfig, NetworkPeer: networkPeer}

		peers = append(peers, peer)
	}

	return peers, true
}

func (m *fabSdkImpl) ChannelOrderers(channelName string) ([]fab.OrdererConfig, bool) {
	orderers := []fab.OrdererConfig{}

	channel, ok := m.ChannelConfig(channelName)
	if !ok || channel == nil {
		return nil, false
	}

	for _, chOrderer := range channel.Orderers {
		orderer, ok := m.OrdererConfig(chOrderer)
		if !ok || orderer == nil {
			return nil, false
		}

		orderers = append(orderers, *orderer)
	}

	return orderers, true
}

func defaultQueryChannelPolicy() fab.QueryChannelConfigPolicy {
	return fab.QueryChannelConfigPolicy{
		MinResponses: 1,
		MaxTargets:   1,
		RetryOpts: retry.Opts{
			Attempts:       5,
			InitialBackoff: 500 * time.Millisecond,
			MaxBackoff:     5 * time.Second,
			BackoffFactor:  2.0,
		},
	}
}

var defaultTimeOutTypes = map[fab.TimeoutType]time.Duration{
	fab.EndorserConnection:       time.Second * 10,
	fab.PeerResponse:             time.Minute * 3,
	fab.DiscoveryGreylistExpiry:  time.Second * 10,
	fab.EventHubConnection:       time.Second * 15,
	fab.EventReg:                 time.Second * 15,
	fab.OrdererConnection:        time.Second * 15,
	fab.OrdererResponse:          time.Minute * 2,
	fab.DiscoveryConnection:      time.Second * 15,
	fab.DiscoveryResponse:        time.Second * 15,
	fab.Query:                    time.Minute * 3,
	fab.Execute:                  time.Minute * 3,
	fab.ResMgmt:                  time.Minute * 3,
	fab.ConnectionIdle:           time.Second * 30,
	fab.EventServiceIdle:         time.Minute * 2,
	fab.ChannelConfigRefresh:     time.Minute * 90,
	fab.ChannelMembershipRefresh: time.Second * 60,
	fab.DiscoveryServiceRefresh:  time.Second * 10,
	fab.SelectionServiceRefresh:  time.Minute * 15,
	fab.CacheSweepInterval:       time.Second * 15,
}

func newTLSConfig(path string) endpoint.TLSConfig {
	config := endpoint.TLSConfig{Path: utils.Substitute(path)}
	if err := config.LoadBytes(); err != nil {
		panic(errors.Errorf("error loading bytes: %s", err))
	}

	return config
}

func tlsCertBytes(path string) *x509.Certificate {
	bytes, err := ioutil.ReadFile(utils.Substitute(path))
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(bytes)
	if block != nil {
		pub, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil
		}

		return pub
	}

	return nil
}

func getFabLoggingType(level entities.LoggingLevel) api.LoggingType {
	switch level {
	case entities.INFO:
		return api.LoggingType{Level: "INFO"}
	case entities.WARNING:
		return api.LoggingType{Level: "WARNING"}
	case entities.ERROR:
		return api.LoggingType{Level: "ERROR"}
	case entities.DEBUG:
		return api.LoggingType{Level: "DEBUG"}
	case entities.FATAL:
		return api.LoggingType{Level: "FATAL"}
	default:
		return api.LoggingType{Level: "INFO"}
	}
}
