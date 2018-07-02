package entities

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
	Logging             LoggingLevel   `json:"serverlog"`
	CryptoConfigPath    string         `json:"cryptoConfigPath"`
	CredentialStorePath string         `json:"credentialStorePath"`
	TLSKeyPair          TLSKeyPathPair `json:"tlsKeyPair"`
	UserName            string         `json:"userName"`
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

type Credential struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// TLSKeyPathPair - TLSKeyPathPair
type TLSKeyPathPair struct {
	KeyPath  string `json:"keyPath"`
	CertPath string `json:"certPath"`
}

type CertificateAuthority struct {
	URL                 string         `json:"url"`
	TLSCertPath         string         `json:"tlsCertPath"`
	TLSCertClientPaths  TLSKeyPathPair `json:"tlsCertClientPaths"`
	CAName              string         `json:"caName"`
	RegistrarCredential Credential     `json:"registrarCredential"`
	CryptoConfigPath    string         `json:"cryptoConfigPath"`
}

type Orderer struct {
	Name          string                 `json:"name"`
	URL           string                 `json:"url"`
	GRPCOptions   map[string]interface{} `json:"grpcOptions"`
	TLSCACertPath string                 `json:"tlsCACertPath"`
	MSPID         string                 `json:"mspid"`
	CryptoPath    string                 `json:"cryptoPath"`
	MSPDir        string                 `json:"mspDir"`
}

type Organization struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name" yaml:"name"`
	MSPID                  string   `json:"mspid" yaml:"name"`
	Peers                  []string `json:"peers"`
	CryptoPath             string   `json:"cryptoPath"`
	MSPDir                 string   `json:"mspDir"`
	CertificateAuthorities []string `json:"certificateAuthorities"`
	AdminUser              string   `json:"adminUser"`
}

type Peer struct {
	URL             string                 `json:"url"`
	EventURL        string                 `json:"eventUrl"`
	EndPoint        string                 `json:"endPoint"`
	GrpcOptions     map[string]interface{} `json:"grpcOptions"`
	TLSCertPath     string                 `json:"tlsCertPath"`
	IsEndrosingPeer bool                   `json:"isEndrosingPeer"`
	IsChainCodePeer bool                   `json:"isChainCodePeer"`
	CanQueryLedger  bool                   `json:"canQueryLedger"`
	EventSource     bool                   `json:"eventSource"`
}

type NetWork struct {
	Name                  string                          `json:"name"`
	Organizations         []string                        `json:"organizations"`
	Orderers              []string                        `json:"orderers"`
	Peers                 map[string]Peer                 `json:"peers"`
	CertificateAuthority  map[string]CertificateAuthority `json:"certificateAuthority"`
	IsSystemCertPool      bool                            `json:"isSystemCertPool"`
	SecurityConfiguration SecurityConfiguration           `json:"securityConfiguration"`
	Consortiums           map[string][]string             `json:"consortiums"` // consortiums and its participating organization groups
}

type CreateChannelRequest struct {
	ChannelName       string              `json:"channelName"`
	OrganizationNames []string            `json:"organizationNames"`
	ConsortiumName    string              `json:"consortiumName"`
	AnchorPeers       map[string][]string `json:"anchorPeers"` //anchor peers for each organization
}

type InstallChainCodeRequest struct {
	ChannelName string `json:"channelName"`
	ChainCodeName string `json:"chainCodeName"`
	ChainCodePath string `json:"chainCodePath"`
	ChainCodeVersion string `json:"chainCodeVersion"`
}
