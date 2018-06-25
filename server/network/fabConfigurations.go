package network

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"strings"
)

type FabClientConfig struct {
	Organization    string
	Logging         api.LoggingType
	CryptoConfig    msp.CCType
	CredentialStore msp.CredentialStoreType
	TLSCerts        endpoint.MutualTLSConfig
	TLSKey          []byte
	TLSCert         []byte
}

func (m *FabClientConfig) Client() *msp.ClientConfig {
	return &msp.ClientConfig{
		Organization: strings.ToLower(m.Organization),
		Logging: m.Logging,
		CryptoConfig: m.CryptoConfig,
		CredentialStore:m.CredentialStore,
		TLSKey: m.TLSCerts.Client.Key.Bytes(),
		TLSCert: m.TLSCerts.Client.Cert.Bytes(),
	}
}


func NewFabClientConfig(clientConfig ClientConfiguration) *FabClientConfig {
	return &FabClientConfig{
		Organization:    strings.ToLower(clientConfig.Organization),
		Logging:         GetFabLoggingType(clientConfig.Logging),
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

func GetFabLoggingType(level LoggingLevel) api.LoggingType {
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
