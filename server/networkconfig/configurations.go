package networkconfig

import "time"

//type ClientConfig struct {
//	Organization    string
//	Logging         api.LoggingType
//	CryptoConfig    msp.CCType
//	CredentialStore msp.CredentialStoreType
//	TLSCerts        endpoint.MutualTLSConfig
//	TLSKey          []byte
//	TLSCert         []byte
//}

// NetworkConfiguration - NetworkConfiguration
type NetworkConfiguration struct {
	Name                       string                               `json:"name"`
	OrganizationsConfiguration map[string]OrganizationConfiguration `json:"organizationsConfiguration"`
	OrderersConfiguration      map[string]OrdererConfiguration      `json:"orderersConfiguration"`
	PeersConfiguration         map[string]PeerConfiguration         `json:"peersConfiguration"`
	ChannelsConfiguration      map[string]ChannelConfiguration      `json:"channelsConfiguration"`
	CAConfiguration            map[string]CAConfiguration           `json:"caConfiguration"`
	ClientConfiguration        ClientConfiguration                  `json:"clientConfiguration"`
	IsSystemCertPool           bool                                 `json:"isSystemCertPool"`
	SecurityConfiguration      SecurityConfiguration                `json:"securityConfiguration"`
}

// TLSKeyPathPair - TLSKeyPathPair
type TLSKeyPathPair struct {
	KeyPath  string `json:"keyPath"`
	CertPath string `json:"certPath"`
}

// Credential - Credential
type Credential struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// LoggingLevel - LoggingLevel
type LoggingLevel int

const (
	// DEBUG ...
	DEBUG LoggingLevel = iota + 1
	// INFO ...
	INFO
	// WARNING ...
	WARNING
	// ERROR ...
	ERROR
	// FATAL ...
	FATAL
)

// ClientConfiguration - ClientConfiguration
type ClientConfiguration struct {
	Organization        string         `json:"organization"`
	Logging             LoggingLevel   `json:"logging"`
	CryptoConfigPath    string         `json:"cryptoConfigPath"`
	CredentialStorePath string         `json:"credentialStorePath"`
	TLSKeyPair          TLSKeyPathPair `json:"tlsKeyPair"`
}

// GRPCOptions - GRPCOptions
type GRPCOptions struct {
	SSLTargetNameOveride string        `json:"sslTargetNameOverride"`
	KeepAliveTime        time.Duration `json:"keepAliveTime"`
	KeepAliveTimeOut     time.Duration `json:"keepAliveTimeout"`
	KeepAlivePermit      bool          `json:"keepAlivePermit"`
	FailFast             bool          `json:"failFast"`
	AllowInsecure        bool          `json:"allowInsecure"`
}

// PeerConfiguration - PeerConfiguration
type PeerConfiguration struct {
	URL           string      `json:"url"`
	EventURL      string      `json:"eventUrl"`
	GRPCOptions   GRPCOptions `json:"grpcOptions"`
	TLSCACertPath string      `json:"tlsCACertPath"`
}

// CAConfiguration - CAConfiguration
type CAConfiguration struct {
	URL                 string         `json:"url"`
	TLSCertPath         string         `json:"tlsCertPath"`
	TLSCertClientPaths  TLSKeyPathPair `json:"tlsCertClientPaths"`
	CAName              string         `json:"caName"`
	RegistrarCredential Credential     `json:"registrarCredential"`
	CryptoConfigPath    string         `json:"cryptoConfigPath"`
}

// OrdererConfiguration - OrdererConfiguration
type OrdererConfiguration struct {
	URL           string      `json:"url"`
	GRPCOptions   GRPCOptions `json:"grpcOptions"`
	TLSCACertPath string      `json:"tlsCACertPath"`
}

// OrganizationConfiguration - OrganizationConfiguration
type OrganizationConfiguration struct {
	MSPID                  string   `json:"mspid"`
	CryptoPath             string   `json:"cryptoPath"`
	Peers                  []string `json:"peers"`
	CertificateAuthorities []string `json:"certificateAuthorities"`
	UserName               string   `json:"userName"`
}

// ChannelPeerConfiguration - ChannelPeerConfiguration
type ChannelPeerConfiguration struct {
	IsEndrosingPeer bool `json:"isEndrosingPeer"`
	IsChainCodePeer bool `json:"isChainCodePeer"`
	CanQueryLedger  bool `json:"canQueryLedger"`
	EventSource     bool `json:"eventSource"`
}

// ConfigurationPolicy - ConfigurationPolicy
type ConfigurationPolicy struct {
	MinResponses        int           `json:"minResponses"`
	MaxTargets          int           `json:"maxTargets"`
	RetryAttempts       int           `json:"retryAttempts"`
	RetryInitialBackoff time.Duration `json:"retryInitialBackoff"`
	RetryMaxBackoff     time.Duration `json:"retryMaxBackoff"`
	RetryBackoffFactor  float64       `json:"retryBackoffFactor"`
}

// ChannelConfiguration - ChannelConfiguration
type ChannelConfiguration struct {
	Orderers           []string                            `json:"orderers"`
	Peers              map[string]ChannelPeerConfiguration `json:"peers"`
	QueryChannelPolicy ConfigurationPolicy                 `json:"queryChannelPolicy"`
}

// SecurityConfiguration - SecurityConfiguration
type SecurityConfiguration struct {
	IsEnabled       bool   `json:"isEnabled"`
	Algorithm       string `json:"algorithm"`
	Level           int    `json:"level"`
	IsSoftVerify    bool   `json:"isSoftVerify"`
	Provider        string `json:"provider"`
	ProviderLibPath string `json:"providerLibPath"`
	ProviderPin     string `json:"providerPin"`
	ProviderLabel   string `json:"providerLabel"`
	KeyStoragePath  string `json:"keyStoragePath"`
}
