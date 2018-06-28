package networkconfig

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/sgururajan/hyperledger-tictactoe/server/pathUtil"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	tls2 "github.com/hyperledger/fabric-sdk-go/pkg/core/config/comm/tls"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/cryptoutil"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"github.com/pkg/errors"
	"strconv"
)

// FabNetworkConfiguration - FabNetworkConfiguration
type FabNetworkConfiguration struct {
	Name                  string
	Organizations         map[string]fab.OrganizationConfig
	Orderers              map[string]fab.OrdererConfig
	Peers                 map[string]fab.PeerConfig
	Channels              map[string]fab.ChannelEndpointConfig
	CAs                   map[string]FabCAConfig
	SecurityConfiguration *FabSecurityConfiguration
	OrgsByName            map[string]*OrgInfo
	PeersByEndpoint       map[string]*PeerInfo
	OrderersByEndpoint    map[string]*OrdererInfo
	OrgsInfo              []*OrgInfo
	PeersInfo             []*PeerInfo
	OrderersInfo          []*OrdererInfo
	isSystemCertPool      bool
	//networkConfig    fab.NetworkConfig
	tlsCertPool  fab.CertPool
	ClientConfig *FabClientConfig
	rwLock       sync.Mutex
}

type OrdererInfo struct {
	Endpoint    string
	IsPrimary   bool
	OrdererType string
}

type OrgInfo struct {
	Name          string
	Endpoint      string
	AdminUserName string
	IsPrimary     bool
	MSPID         string
	MSPDir        string
	IsOrderer     bool
}

type PeerInfo struct {
	Endpoint        string
	OrgName         string
	IsEndrosingPeer bool
	Port            uint32
}

// FabClientConfig - FabClientConfig
type FabClientConfig struct {
	UserName        string
	Organization    string
	Logging         api.LoggingType
	CryptoConfig    msp.CCType
	CredentialStore msp.CredentialStoreType
	TLSCerts        endpoint.MutualTLSConfig
	TLSKey          []byte
	TLSCert         []byte
}

// FabCAConfig - FabCAConfig
type FabCAConfig struct {
	URL       string
	TLSCACert endpoint.MutualTLSConfig
	Registrar msp.EnrollCredentials
	CAName    string
}

// FabOrdererConfiguration - FabOrdererConfiguration
type FabOrdererConfiguration struct {
	URL         string
	GRPCOptions map[string]interface{}
	TLSCACert   *x509.Certificate
}

// FabPeerConfiguration - FabPeerConfiguration
type FabPeerConfiguration struct {
	URL         string
	EventURL    string
	GRPCOptions map[string]interface{}
	TLSCACert   *x509.Certificate
}

func newFabClientConfig(clientConfig ClientConfiguration) *FabClientConfig {
	return &FabClientConfig{
		UserName:        clientConfig.UserName,
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

func NewFabNetworkConfigurationFromConfig(config NetworkConfiguration) *FabNetworkConfiguration {
	return NewFabNetworkConfiguration(config.Name,
		config.CAConfiguration,
		config.OrderersConfiguration,
		config.PeersConfiguration,
		config.OrganizationsConfiguration,
		config.ChannelsConfiguration,
		config.ClientConfiguration,
		config.SecurityConfiguration)
}

// NewFabNetworkConfiguration - creates new instance of FabNetworkConfiguration
func NewFabNetworkConfiguration(name string,
	caConfig map[string]CAConfiguration,
	ordererConfig map[string]OrdererConfiguration,
	peerConfig map[string]PeerConfiguration,
	orgConfig map[string]OrganizationConfiguration,
	chConfig map[string]ChannelConfiguration,
	clientConfig ClientConfiguration,
	securityConfig SecurityConfiguration) *FabNetworkConfiguration {

	result := FabNetworkConfiguration{Name: name}

	result.isSystemCertPool = false
	result.tlsCertPool = tls2.NewCertPool(result.isSystemCertPool)

	result.Orderers = make(map[string]fab.OrdererConfig)
	result.Peers = make(map[string]fab.PeerConfig)
	result.Organizations = make(map[string]fab.OrganizationConfig)
	result.Channels = make(map[string]fab.ChannelEndpointConfig)
	result.CAs = make(map[string]FabCAConfig)

	result.OrgsInfo = []*OrgInfo{}
	result.PeersInfo = []*PeerInfo{}
	result.OrderersInfo = []*OrdererInfo{}
	result.OrgsByName = make(map[string]*OrgInfo)
	result.PeersByEndpoint=make(map[string]*PeerInfo)
	result.OrderersByEndpoint=make(map[string]*OrdererInfo)

	for k, v := range ordererConfig {
		result.Orderers[k] = newFabOrdererConfiguration(v)
		result.OrderersInfo = append(result.OrderersInfo, &OrdererInfo{
			Endpoint:  k,
			IsPrimary: v.IsPrimary,
		})
	}

	for k, v := range orgConfig {
		result.Organizations[k] = newFabOrgsConfiguration(v)

		orgInfo:= &OrgInfo{
			Name:          k,
			Endpoint:      k,
			MSPID:         v.MSPID,
			AdminUserName: "Admin",
			IsPrimary:     v.IsPrimary,
			MSPDir:        v.MSPDir,
			IsOrderer:     v.IsOrderer,
		}
		result.OrgsInfo = append(result.OrgsInfo, orgInfo)
		result.OrgsByName[k]=orgInfo

		for _, pv := range v.Peers {
			peerInfo:= &PeerInfo{
				Endpoint: pv,
				OrgName:  k,
			}
			result.PeersInfo = append(result.PeersInfo, peerInfo)
			result.PeersByEndpoint[pv]=peerInfo
		}
	}

	for k, v := range peerConfig {
		result.Peers[k] = newFabPeerConfiguration(v)
		port := getPortFromUrl(v.URL)
		for _, pv := range result.PeersInfo {
			if strings.Contains(v.URL, pv.Endpoint) {
				pv.Port = port
			}
		}
	}

	for k, v := range chConfig {
		result.Channels[k] = newFabChannelConfig(v)
		for _, pv := range result.PeersInfo {
			cpv, ok := v.Peers[pv.Endpoint]
			if !ok {
				break
			}
			pv.IsEndrosingPeer = cpv.IsEndrosingPeer
		}
	}

	for k, v := range caConfig {
		result.CAs[k] = newFabCAConfig(v)
	}

	result.SecurityConfiguration = newFabSecurityConfig(securityConfig)

	result.ClientConfig = newFabClientConfig(clientConfig)

	return &result
}

func getPortFromUrl(url string) uint32 {
	index := strings.LastIndex(url, ":")
	if index > 0 {
		res, err := strconv.ParseUint(url[index:], 10, 32)
		if err != nil {
			return 0
		}

		return uint32(res)
	}

	return 0
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

	config.Peers = make(map[string]fab.PeerChannelConfig)

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
		"ssl-target-Name-override": opts.SSLTargetNameOveride,
		"keep-alive-time":          opts.KeepAliveTime,
		"keep-alive-timeout":       opts.KeepAliveTimeOut,
		"keep-alive-permit":        opts.KeepAlivePermit,
		"fail-fast":                opts.FailFast,
		"allow-insecure":           opts.AllowInsecure,
	}
}

func (m *FabNetworkConfiguration) Client() *msp.ClientConfig {
	return &msp.ClientConfig{
		Organization:    strings.ToLower(m.ClientConfig.Organization),
		Logging:         m.ClientConfig.Logging,
		CryptoConfig:    m.ClientConfig.CryptoConfig,
		CredentialStore: m.ClientConfig.CredentialStore,
		TLSKey:          m.ClientConfig.TLSCerts.Client.Key.Bytes(),
		TLSCert:         m.ClientConfig.TLSCerts.Client.Cert.Bytes(),
	}
}

func (m *FabNetworkConfiguration) Timeout(tType fab.TimeoutType) time.Duration {
	t, ok := defaultTimeOutTypes[tType]
	if !ok {
		return time.Second * 30
	}

	return t
}

func (m *FabNetworkConfiguration) EventServiceType() fab.EventServiceType {
	if os.Getenv("FABRIC_SDK_CLIENT_EVENTSERVICE_TYPE") == "eventhub" {
		return fab.EventHubEventServiceType
	}

	return fab.DeliverEventServiceType
}

func (m *FabNetworkConfiguration) OrderersConfig() []fab.OrdererConfig {
	orderers := []fab.OrdererConfig{}

	for _, o := range m.Orderers {
		if o.TLSCACert == nil && !m.isSystemCertPool {
			return nil
		}

		orderers = append(orderers, o)
	}

	return orderers
}

func (m *FabNetworkConfiguration) OrdererConfig(nameOrUrl string) (*fab.OrdererConfig, bool) {
	orderer, ok := m.Orderers[strings.ToLower(nameOrUrl)]
	if !ok {
		return nil, false
	}

	return &orderer, true
}

func (m *FabNetworkConfiguration) PeersConfig(org string) ([]fab.PeerConfig, bool) {
	orgPeers := m.Organizations[strings.ToLower(org)].Peers
	peers := []fab.PeerConfig{}

	for _, peerName := range orgPeers {
		p := m.Peers[strings.ToLower(peerName)]
		if err := m.verifyPeerConfig(p, peerName, endpoint.IsTLSEnabled(p.URL)); err != nil {
			return nil, false
		}

		peers = append(peers, p)
	}

	return peers, true
}

func (m *FabNetworkConfiguration) PeerConfig(nameOrUrl string) (*fab.PeerConfig, bool) {
	pConfig, ok := m.Peers[strings.ToLower(nameOrUrl)]
	if ok {
		return &pConfig, true
	}

	i := strings.Index(nameOrUrl, ":")
	if i > 0 {
		return m.PeerConfig(nameOrUrl[0:i])
	}

	return nil, false
}

func (m *FabNetworkConfiguration) verifyPeerConfig(p fab.PeerConfig, peerName string, tlsEnabled bool) error {
	if p.URL == "" {
		return errors.Errorf("URL does not exists or empty for peer %s", peerName)
	}

	if tlsEnabled && p.TLSCACert == nil && !m.isSystemCertPool {
		return errors.Errorf("tls certificates does not exists or empty for peer %s", peerName)
	}

	return nil
}

func (m *FabNetworkConfiguration) NetworkConfig() *fab.NetworkConfig {
	//return &m.networkConfig
	return &fab.NetworkConfig{
		Channels:      m.Channels,
		Organizations: m.Organizations,
		Orderers:      m.Orderers,
		Peers:         m.Peers,
	}
}

func (m *FabNetworkConfiguration) NetworkPeers() []fab.NetworkPeer {
	netPeers := []fab.NetworkPeer{}

	for name, p := range m.Peers {
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

func (m *FabNetworkConfiguration) PeerMSPID(name string) (string, bool) {
	for _, org := range m.Organizations {
		for i := 0; i < len(org.Peers); i++ {
			if strings.EqualFold(org.Peers[i], name) {
				return org.MSPID, true
			}
		}
	}

	return "", false
}

func (m *FabNetworkConfiguration) ChannelConfig(channelName string) (*fab.ChannelEndpointConfig, bool) {
	ch, ok := m.Channels[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}
	return &ch, true
}

func (m *FabNetworkConfiguration) ChannelPeers(channelName string) ([]fab.ChannelPeer, bool) {
	peers := []fab.ChannelPeer{}

	chConfig, ok := m.Channels[strings.ToLower(channelName)]
	if !ok {
		return nil, false
	}

	for peerName, chPeerConfig := range chConfig.Peers {
		p, ok := m.Peers[strings.ToLower(peerName)]
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

func (m *FabNetworkConfiguration) ChannelOrderers(channelName string) ([]fab.OrdererConfig, bool) {
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

func (m *FabNetworkConfiguration) TLSCACertPool() fab.CertPool {
	return m.tlsCertPool
}

func (m *FabNetworkConfiguration) TLSClientCerts() []tls.Certificate {
	var clientCerts tls.Certificate
	cb := m.ClientConfig.TLSCerts.Client.Cert.Bytes()

	if len(cb) == 0 {
		return []tls.Certificate{clientCerts}
	}

	cs := cryptosuite.GetDefault()
	pk, err := cryptoutil.GetPrivateKeyFromCert(cb, cs)

	if err != nil || pk == nil {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
		ccs, err := m.loadPrivateKeyFromConfig(m.ClientConfig, clientCerts, cb)
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

func (m *FabNetworkConfiguration) CryptoConfigPath() string {
	return m.ClientConfig.CryptoConfig.Path
}

func (m *FabNetworkConfiguration) loadPrivateKeyFromConfig(clientConfig *FabClientConfig, clientCert tls.Certificate, cb []byte) ([]tls.Certificate, error) {
	kb := clientConfig.TLSCerts.Client.Key.Bytes()
	clientCerts, err := tls.X509KeyPair(cb, kb)
	if err != nil {
		return nil, errors.Errorf("error loading cert/key pair as TLS client credential: %s", err)
	}

	return []tls.Certificate{clientCerts}, nil
}

/* Methods for Identity configuration */
func (m *FabNetworkConfiguration) CAConfig(orgName string) (*msp.CAConfig, bool) {
	return m.getCAConfig(orgName)
}

func (m *FabNetworkConfiguration) CAServerCerts(orgName string) ([][]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAServerCerts, true
}

func (m *FabNetworkConfiguration) CAClientKey(orgName string) ([]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientKey, true
}

func (m *FabNetworkConfiguration) CAClientCert(orgName string) ([]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientCert, true
}

func (m *FabNetworkConfiguration) CAKeyStorePath() string {
	return "/tmp/msp"
}

func (m *FabNetworkConfiguration) CredentialStorePath() string {
	return "/tmp/state-store"
}

func (m *FabNetworkConfiguration) getCAConfig(orgName string) (*msp.CAConfig, bool) {
	if len(m.Organizations[strings.ToLower(orgName)].CertificateAuthorities) == 0 {
		return nil, false
	}

	org := m.Organizations[strings.ToLower(orgName)]
	certAuthName := org.CertificateAuthorities[0]
	if certAuthName == "" {
		return nil, false
	}

	caConfig, ok := m.CAs[strings.ToLower(certAuthName)]
	if !ok {
		return nil, false
	}

	mspCAConfig, err := caConfig.getMSPCAConfig()
	if err != nil {
		return nil, false
	}

	return mspCAConfig, true
}

func (m *FabCAConfig) getMSPCAConfig() (*msp.CAConfig, error) {
	serverCerts, err := m.getServerCerts()
	if err != nil {
		return nil, err
	}

	return &msp.CAConfig{
		URL:              m.URL,
		Registrar:        m.Registrar,
		CAName:           m.CAName,
		TLSCAClientCert:  m.TLSCACert.Client.Cert.Bytes(),
		TLSCAClientKey:   m.TLSCACert.Client.Key.Bytes(),
		TLSCAServerCerts: serverCerts,
	}, nil
}

func (m *FabCAConfig) getServerCerts() ([][]byte, error) {
	var serverCerts [][]byte

	pems := m.TLSCACert.Pem
	if len(pems) > 0 {
		serverCerts := make([][]byte, len(pems))
		for i, pem := range pems {
			serverCerts[i] = []byte(pem)
		}

		return serverCerts, nil
	}

	certFiles := strings.Split(m.TLSCACert.Path, ",")
	serverCerts = make([][]byte, len(certFiles))
	for i, certPath := range certFiles {
		bytes, err := ioutil.ReadFile(pathUtil.Substitute(certPath))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to load server certificates")
		}
		serverCerts[i] = bytes
	}

	return serverCerts, nil
}
