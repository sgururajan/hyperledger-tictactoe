package database

type Network struct {
	Name                  string                          `json:"Name"`
	Organizations         map[string]Organization         `json:"organizations"`
	Orderers              map[string]Orderer              `json:"orderers"`
	Peers                 map[string]Peer                 `json:"peers"`
	CertificateAuthority  map[string]CertificateAuthority `json:"certificateAuthority"`
	IsSystemCertPool      bool                            `json:"isSystemCertPool"`
	SecurityConfiguration SecurityConfiguration           `json:"securityConfiguration"`
	Consortiums           map[string][]string             `json:"consortiums"` // consortiums and its participating organization groups
	//OrganizationsConfiguration map[string]OrganizationConfiguration `json:"organizationsConfiguration"`
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
	Organization        string         `json:"organization"`
}

type Orderer struct {
	Name          string                 `json:"name"`
	URL           string                 `json:"url"`
	GRPCOptions   map[string]interface{} `json:"grpcOptions"`
	TLSCACertPath string                 `json:"tlsCACertPath"`
	MSPID         string                 `json:"mspid"`
	CryptoPath    string                 `json:"cryptoPath"`
	MSPDir        string                 `json:"mspDir"`
	Organization  string                 `json:"organization"`
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
	IsOrderer              bool     `json:"isOrderer"`
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
	Organization    string                 `json:"organization"`
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
