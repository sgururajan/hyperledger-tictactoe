package network

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	tls2 "github.com/hyperledger/fabric-sdk-go/pkg/core/config/comm/tls"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/cryptoutil"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/pkg/errors"
	"os"
	"strings"
	"sync"
	"time"
	"io/ioutil"
	"hyperledger/hyperledger-tictactoe/server/pathUtil"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
)

type FabEndpointConfiguration struct {
	orgsConfig       map[string]fab.OrganizationConfig
	orderersConfig   map[string]fab.OrdererConfig
	peersConfig      map[string]fab.PeerConfig
	channelsConfig   map[string]fab.ChannelEndpointConfig
	caConfig         map[string]FabCAConfig
	isSystemCertPool bool
	//networkConfig    fab.NetworkConfig
	tlsCertPool      fab.CertPool
	clientConfig     FabClientConfig
	rwLock           sync.Mutex
}

type FabClientConfig struct {
	Organization    string
	Logging         api.LoggingType
	CryptoConfig    msp.CCType
	CredentialStore msp.CredentialStoreType
	TLSCerts        endpoint.MutualTLSConfig
	TLSKey          []byte
	TLSCert         []byte
}

type FabCAConfig struct {
	URL       string
	TLSCACert endpoint.MutualTLSConfig
	Registrar msp.EnrollCredentials
	CAName    string
}

type FabOrdererConfiguration struct {
	URL         string
	GRPCOptions map[string]interface{}
	TLSCACert   *x509.Certificate
}

type FabPeerConfiguration struct {
	URL         string
	EventURL    string
	GRPCOptions map[string]interface{}
	TLSCACert   *x509.Certificate
}

func newFabClientConfig(clientConfig ClientConfiguration) *FabClientConfig {
	return &FabClientConfig{
		Organization:    strings.ToLower(clientConfig.Organization),
		Logging:         getFabLoggingType(clientConfig.Logging),
		CryptoConfig:    msp.CCType{Path: clientConfig.CryptoConfigPath},
		CredentialStore: msp.CredentialStoreType{Path: clientConfig.CredentialStorePath},
		TLSCerts: endpoint.MutualTLSConfig{
			Client: endpoint.TLSKeyPair{
				Key:  newTLSConfig(clientConfig.TLSKeyPair.KeyPath),
				Cert: newTLSConfig(clientConfig.TLSKeyPair.CertPath),
			},
		},
	}
}

func newFabCAConfig(caConfig CAConfiguration) FabCAConfig {
	return FabCAConfig{
		URL: caConfig.URL,
		TLSCACert: endpoint.MutualTLSConfig{
			Path: caConfig.TLSCertPath,
			Client: endpoint.TLSKeyPair{
				Key:  newTLSConfig(caConfig.TLSCertClientPaths.KeyPath),
				Cert: newTLSConfig(caConfig.TLSCertClientPaths.CertPath),
			},
		},
		Registrar: msp.EnrollCredentials{
			EnrollID:     caConfig.RegistrarCredential.ID,
			EnrollSecret: caConfig.RegistrarCredential.Secret,
		},
		CAName: caConfig.CAName,
	}
}

func newFabOrdererConfiguration(ordererConfig OrdererConfiguration) fab.OrdererConfig {
	return fab.OrdererConfig{
		URL:         ordererConfig.URL,
		GRPCOptions: getGRPCOptions(ordererConfig.GRPCOptions),
		TLSCACert:   tlsCertBytes(ordererConfig.TLSCACertPath),
	}
}

func newFabPeerConfiguration(peerConfig PeerConfiguration) fab.PeerConfig {
	return fab.PeerConfig{
		URL:         peerConfig.URL,
		EventURL:    peerConfig.EventURL,
		GRPCOptions: getGRPCOptions(peerConfig.GRPCOptions),
		TLSCACert:   tlsCertBytes(peerConfig.TLSCACertPath),
	}
}

func NewFabEndpointConfiguration(caConfig map[string]CAConfiguration,
	ordererConfig map[string]OrdererConfiguration,
	peerConfig map[string]PeerConfiguration,
	orgConfig map[string]OrganizationConfiguration,
	chConfig map[string]ChannelConfiguration,
	clientConfig ClientConfiguration) *FabEndpointConfiguration {

	result := FabEndpointConfiguration{}

	result.isSystemCertPool = false
	result.tlsCertPool = tls2.NewCertPool(result.isSystemCertPool)

	for k, v := range ordererConfig {
		result.orderersConfig[k] = newFabOrdererConfiguration(v)
	}

	for k, v := range peerConfig {
		result.peersConfig[k] = newFabPeerConfiguration(v)
	}

	for k, v := range orgConfig {
		result.orgsConfig[k] = newFabOrgsConfiguration(v)
	}

	for k, v := range chConfig {
		result.channelsConfig[k] = newFabChannelConfig(v)
	}

	for k, v := range caConfig {
		result.caConfig[k] = newFabCAConfig(v)
	}

	result.clientConfig = *newFabClientConfig(clientConfig)

	return &result
}

func newFabOrgsConfiguration(orgConfig OrganizationConfiguration) fab.OrganizationConfig {
	return fab.OrganizationConfig{
		MSPID:      orgConfig.MSPID,
		CryptoPath: orgConfig.CryptoPath,
		Peers:      orgConfig.Peers,
		CertificateAuthorities: orgConfig.CertificateAuthorities,
	}
}

func newFabChannelConfig(chConfig ChannelConfiguration) fab.ChannelEndpointConfig {
	config := fab.ChannelEndpointConfig{
		Orderers: chConfig.Orderers,
		Policies: fab.ChannelPolicies{
			QueryChannelConfig: newFabQueryChannelConfig(chConfig.QueryChannelPolicy),
		},
	}

	for k, v := range chConfig.Peers {
		pConfig := newFabChannelPeerConfiguration(v)
		config.Peers[k] = pConfig
	}

	return config
}

func newFabQueryChannelConfig(pConfig ConfigurationPolicy) fab.QueryChannelConfigPolicy {
	return fab.QueryChannelConfigPolicy{
		MinResponses: pConfig.MinResponses,
		MaxTargets:   pConfig.MaxTargets,
		RetryOpts: retry.Opts{
			Attempts:       pConfig.RetryAttempts,
			InitialBackoff: pConfig.RetryInitialBackoff,
			MaxBackoff:     pConfig.RetryMaxBackoff,
			BackoffFactor:  pConfig.RetryBackoffFactor,
		},
	}
}

func newFabChannelPeerConfiguration(cpConfig ChannelPeerConfiguration) fab.PeerChannelConfig {
	return fab.PeerChannelConfig{
		EndorsingPeer:  cpConfig.IsEndrosingPeer,
		ChaincodeQuery: cpConfig.IsChainCodePeer,
		LedgerQuery:    cpConfig.CanQueryLedger,
		EventSource:    cpConfig.EventSource,
	}
}

func getFabLoggingType(level LoggingLevel) api.LoggingType {
	switch level {
	case INFO:
		return api.LoggingType{Level: "INFO"}
	case WARNING:
		return api.LoggingType{Level: "WARNING"}
	case ERROR:
		return api.LoggingType{Level: "ERROR"}
	case DEBUG:
		return api.LoggingType{Level: "DEBUG"}
	case FATAL:
		return api.LoggingType{Level: "FATAL"}
	default:
		return api.LoggingType{Level: "INFO"}
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

func getGRPCOptions(opts GRPCOptions) map[string]interface{} {
	return map[string]interface{}{
		"ssl-target-name-override": opts.SSLTargetNameOveride,
		"keep-alive-time":          opts.KeepAliveTime,
		"keep-alive-timeout":       opts.KeepAliveTimeOut,
		"keep-alive-permit":        opts.KeepAlivePermit,
		"fail-fast":                opts.FailFast,
		"allow-insecure":           opts.AllowInsecure,
	}
}

func (m *FabEndpointConfiguration) Client() *msp.ClientConfig {
	return &msp.ClientConfig{
		Organization: strings.ToLower(m.clientConfig.Organization),
		Logging: m.clientConfig.Logging,
		CryptoConfig: m.clientConfig.CryptoConfig,
		CredentialStore:m.clientConfig.CredentialStore,
		TLSKey: m.clientConfig.TLSCerts.Client.Key.Bytes(),
		TLSCert: m.clientConfig.TLSCerts.Client.Cert.Bytes(),
	}
}

func (m *FabEndpointConfiguration) Timeout(tType fab.TimeoutType) time.Duration {
	t, ok := defaultTimeOutTypes[tType]
	if !ok {
		return time.Second * 30
	}

	return t
}

func (m *FabEndpointConfiguration) EventServiceType() fab.EventServiceType {
	if os.Getenv("FABRIC_SDK_CLIENT_EVENTSERVICE_TYPE") == "eventhub" {
		return fab.EventHubEventServiceType
	}

	return fab.DeliverEventServiceType
}

func (m *FabEndpointConfiguration) OrderersConfig() []fab.OrdererConfig {
	orderers := []fab.OrdererConfig{}

	for _, o := range m.orderersConfig {
		if o.TLSCACert == nil && !m.isSystemCertPool {
			return nil
		}

		orderers = append(orderers, o)
	}

	return orderers
}

func (m *FabEndpointConfiguration) OrdererConfig(nameOrUrl string) (*fab.OrdererConfig, bool) {
	orderer, ok := m.orderersConfig[strings.ToLower(nameOrUrl)]
	if !ok {
		return nil, false
	}

	return &orderer, true
}

func (m *FabEndpointConfiguration) PeersConfig(org string) ([]fab.PeerConfig, bool) {
	orgPeers := m.orgsConfig[strings.ToLower(org)].Peers
	peers := []fab.PeerConfig{}

	for _, peerName := range orgPeers {
		p := m.peersConfig[strings.ToLower(peerName)]
		if err := m.verifyPeerConfig(p, peerName, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil, false
		}

		peers = append(peers, p)
	}

	return peers, true
}

func (m *FabEndpointConfiguration) PeerConfig(nameOrUrl string) (*fab.PeerConfig, bool) {
	pConfig, ok := m.peersConfig[strings.ToLower(nameOrUrl)]
	if ok {
		return &pConfig, true
	}

	i := strings.Index(nameOrUrl, ":")
	if i > 0 {
		return m.PeerConfig(nameOrUrl[0:i])
	}

	return nil, false
}

func (m *FabEndpointConfiguration) verifyPeerConfig(p fab.PeerConfig, peerName string, tlsEnabled bool) error {
	if p.URL == "" {
		return errors.Errorf("URL does not exists or empty for peer %s", peerName)
	}

	if tlsEnabled && p.TLSCACert == nil && !m.isSystemCertPool {
		return errors.Errorf("tls certificates does not exists or empty for peer %s", peerName)
	}

	return nil
}

func (m *FabEndpointConfiguration) NetworkConfig() *fab.NetworkConfig {
	//return &m.networkConfig
	return &fab.NetworkConfig{
		Channels: m.channelsConfig,
		Organizations: m.orgsConfig,
		Orderers: m.orderersConfig,
		Peers: m.peersConfig,
	}
}

func (m *FabEndpointConfiguration) NetworkPeers() []fab.NetworkPeer {
	netPeers := []fab.NetworkPeer{}

	for name, p := range m.peersConfig {
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

func (m *FabEndpointConfiguration) PeerMSPID(name string) (string, bool) {
	for _, org := range m.orgsConfig {
		for i := 0; i < len(org.Peers); i++ {
			if strings.EqualFold(org.Peers[i], name) {
				return org.MSPID, true
			}
		}
	}

	return "", false
}

func (m *FabEndpointConfiguration) ChannelConfig(channelName string) (*fab.ChannelEndpointConfig, bool) {
	ch, ok := m.channelsConfig[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}
	return &ch, true
}

func (m *FabEndpointConfiguration) ChannelPeers(channelName string) ([]fab.ChannelPeer, bool) {
	peers := []fab.ChannelPeer{}

	chConfig, ok := m.channelsConfig[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}

	for peerName, chPeerConfig := range chConfig.Peers {
		p, ok := m.peersConfig[strings.ToLower(peerName)]
		if !ok {
			return nil, false
		}

		if err := m.verifyPeerConfig(p, peerName, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil, false
		}

		mspID, ok := m.PeerMSPID(peerName)
		if !ok {
			return nil, false
		}

		networkPeer := fab.NetworkPeer{PeerConfig: p, MSPID: mspID}
		peer := fab.ChannelPeer{PeerChannelConfig: chPeerConfig, NetworkPeer: networkPeer}

		peers = append(peers, peer)
	}

	return peers, true
}

func (m *FabEndpointConfiguration) ChannelOrderers(channelName string) ([]fab.OrdererConfig, bool) {
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

func (m *FabEndpointConfiguration) TLSCACertPool() fab.CertPool {
	return m.tlsCertPool
}

func (m *FabEndpointConfiguration) TLSClientCerts() []tls.Certificate {
	var clientCerts tls.Certificate
	cb := m.clientConfig.TLSCerts.Client.Cert.Bytes()

	if len(cb) == 0 {
		return []tls.Certificate{clientCerts}
	}

	cs := cryptosuite.GetDefault()
	pk, err := cryptoutil.GetPrivateKeyFromCert(cb, cs)

	if err != nil || pk == nil {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
		ccs, err := m.loadPrivateKeyFromConfig(&m.clientConfig, clientCerts, cb)
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

func (m *FabEndpointConfiguration) CryptoConfigPath() string {
	return m.clientConfig.CryptoConfig.Path
}

func (m *FabEndpointConfiguration) loadPrivateKeyFromConfig(clientConfig *FabClientConfig, clientCert tls.Certificate, cb []byte) ([]tls.Certificate, error) {
	kb := clientConfig.TLSCerts.Client.Key.Bytes()
	clientCerts, err := tls.X509KeyPair(cb, kb)
	if err != nil {
		return nil, errors.Errorf("error loading cert/key pair as TLS client credential: %s", err)
	}

	return []tls.Certificate{clientCerts}, nil
}


/* Methods for Identity configuration */
func (m *FabEndpointConfiguration) CAConfig(orgName string) (*msp.CAConfig, bool) {
	return m.getCAConfig(orgName)
}

func (m *FabEndpointConfiguration) CAServerCerts(orgName string) ([][]byte, bool) {
	caConfig, ok:= m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAServerCerts, true
}

func (m *FabEndpointConfiguration) CAClientKey(orgName string) ([]byte, bool) {
	caConfig, ok:= m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientKey, true
}

func (m *FabEndpointConfiguration) CAClientCert(orgName string) ([]byte, bool) {
	caConfig, ok:= m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientCert, true
}

func (m *FabEndpointConfiguration) CAKeyStorePath() string {
	return "/tmp/msp"
}

func (m *FabEndpointConfiguration) CredentialStorePath() string {
	return "/tmp/state-store"
}

func (m *FabEndpointConfiguration) getCAConfig(orgName string) (*msp.CAConfig, bool) {
	if len(m.orgsConfig[strings.ToLower(orgName)].CertificateAuthorities)==0 {
		return nil, false
	}

	org:= m.orgsConfig[strings.ToLower(orgName)]
	certAuthName:= org.CertificateAuthorities[0]
	if certAuthName=="" {
		return nil, false
	}

	caConfig, ok:= m.caConfig[strings.ToLower(certAuthName)]
	if !ok {
		return nil, false
	}

	mspCAConfig, err:= caConfig.getMSPCAConfig()
	if err != nil {
		return nil, false
	}

	return mspCAConfig, true
}

func (m *FabCAConfig) getMSPCAConfig() (*msp.CAConfig, error) {
	serverCerts, err:= m.getServerCerts()
	if err!=nil {
		return nil, err
	}

	return &msp.CAConfig{
		URL: m.URL,
		Registrar: m.Registrar,
		CAName: m.CAName,
		TLSCAClientCert: m.TLSCACert.Client.Cert.Bytes(),
		TLSCAClientKey: m.TLSCACert.Client.Key.Bytes(),
		TLSCAServerCerts: serverCerts,
	}, nil
}

func (m *FabCAConfig) getServerCerts() ([][]byte, error) {
	var serverCerts [][]byte

	pems:= m.TLSCACert.Pem
	if len(pems) > 0 {
		serverCerts:= make([][]byte, len(pems))
		for i, pem:= range pems {
			serverCerts[i]=[]byte(pem)
		}

		return serverCerts, nil
	}

	certFiles:= strings.Split(m.TLSCACert.Path, ",")
	serverCerts = make([][]byte, len(certFiles))
	for i, certPath:= range certFiles{
		bytes, err:= ioutil.ReadFile(pathUtil.Substitute(certPath))
		if err!=nil {
			return nil, errors.WithMessage(err,"failed to load server certificates")
		}
		serverCerts[i]=bytes
	}

	return serverCerts, nil
}