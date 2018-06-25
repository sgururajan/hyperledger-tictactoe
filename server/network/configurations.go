package network

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

type TLSKeyPathPair struct {
	KeyPath  string
	CertPath string
}

type Credential struct {
	ID     string
	Secret string
}

type LoggingLevel int

const (
	DEBUG LoggingLevel = iota + 1
	INFO
	WARNING
	ERROR
	FATAL
)

type ClientConfiguration struct {
	Organization        string
	Logging             LoggingLevel
	CryptoConfigPath    string
	CredentialStorePath string
	TLSKeyPair          TLSKeyPathPair
}

type GRPCOptions struct {
	SSLTargetNameOveride string
	KeepAliveTime        time.Duration
	KeepAliveTimeOut     time.Duration
	KeepAlivePermit      bool
	FailFast             bool
	AllowInsecure        bool
}

type PeerConfiguration struct {
	URL           string
	EventURL      string
	GRPCOptions   GRPCOptions
	TLSCACertPath string
}

type CAConfiguration struct {
	URL                 string
	TLSCertPath         string
	TLSCertClientPaths  TLSKeyPathPair
	CAName              string
	RegistrarCredential Credential
	CryptoConfigPath string
}

type OrdererConfiguration struct {
	URL           string
	GRPCOptions   GRPCOptions
	TLSCACertPath string
}

type OrganizationConfiguration struct {
	MSPID                  string
	CryptoPath             string
	Peers                  []string
	CertificateAuthorities []string
}

type ChannelPeerConfiguration struct {
	IsEndrosingPeer bool
	IsChainCodePeer bool
	CanQueryLedger  bool
	EventSource     bool
}

type ConfigurationPolicy struct {
	MinResponses        int
	MaxTargets          int
	RetryAttempts       int
	RetryInitialBackoff time.Duration
	RetryMaxBackoff     time.Duration
	RetryBackoffFactor  float64
}

type ChannelConfiguration struct {
	Orderers           []string
	Peers              map[string]ChannelPeerConfiguration
	QueryChannelPolicy ConfigurationPolicy
}

type SecurityConfiguration struct {
	IsEnabled       bool
	Algorithm       string
	Level           int
	IsSoftVerify    bool
	Provider        string
	ProviderLibPath string
	ProviderPin     string
	ProviderLabel   string
	KeyStoragePath  string
}
