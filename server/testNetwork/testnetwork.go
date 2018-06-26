package testNetwork

import (
	"fmt"
	"hyperledger/hyperledger-tictactoe/server/networkconfig"
	"path/filepath"
	"time"
)

var defaultHost = "localhost"
var cryptoConfigPath = "${GOPATH}/src/hyperledger/hyperledger-tictactoe/networkconfig/crypto-config"

func DefaultNetworkConfiguration() networkconfig.NetworkConfiguration {
	clientConfig := getDefaultClientConfig()
	channelsConfig := getDefaultChannelConfig()
	organizationsConfig := getDefaultOrganizationConfig()
	orderersConfig := getDefaultOrderersConfig()
	peersConfig := getDefaultPeersConfig()
	caConfig := getDefaultCAConfig()
	securityConfig := getSecurityConfiguration()

	return networkconfig.NetworkConfiguration{
		Name:"testNetwork",
		SecurityConfiguration:      securityConfig,
		CAConfiguration:            caConfig,
		ChannelsConfiguration:      channelsConfig,
		IsSystemCertPool:           false,
		OrderersConfiguration:      orderersConfig,
		OrganizationsConfiguration: organizationsConfig,
		PeersConfiguration:         peersConfig,
		ClientConfiguration:        clientConfig,
	}
}

func LoadDefaultNetwork() *networkconfig.FabNetworkConfiguration {
	//clientConfig := getDefaultClientConfig()
	//channelsConfig := getDefaultChannelConfig()
	//organizationsConfig := getDefaultOrganizationConfig()
	//orderersConfig := getDefaultOrderersConfig()
	//peersConfig := getDefaultPeersConfig()
	//caConfig := getDefaultCAConfig()
	//securityConfig := getSecurityConfiguration()
	//
	//return networkconfig.NewFabNetworkConfiguration("testNetwork",
	//	caConfig,
	//	orderersConfig,
	//	peersConfig,
	//	organizationsConfig,
	//	channelsConfig,
	//	clientConfig,
	//	securityConfig)

	return networkconfig.NewFabNetworkConfigurationFromConfig(DefaultNetworkConfiguration())
}

func getSecurityConfiguration() networkconfig.SecurityConfiguration {
	return networkconfig.SecurityConfiguration{
		Level:           256,
		IsEnabled:       true,
		KeyStoragePath:  "/tmp/msp/keystore",
		ProviderLabel:   "",
		ProviderPin:     "",
		ProviderLibPath: "",
		IsSoftVerify:    true,
		Provider:        "sw",
		Algorithm:       "SHA2",
	}
}

func getDefaultCAConfig() map[string]networkconfig.CAConfiguration {
	return map[string]networkconfig.CAConfiguration{
		"ca.sivatech.com": {
			URL: fmt.Sprintf("%s:7054", defaultHost),
			TLSCertClientPaths: networkconfig.TLSKeyPathPair{
				KeyPath:  "${GOPATH}/src/hyperledger/hyperledger-tictactoe/network/client-crypto/client-key.pem",
				CertPath: "${GOPATH}/src/hyperledger/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
			},
			TLSCertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/ca/ca.tictactoe.sivatech.com-cert.pem"),
			RegistrarCredential: networkconfig.Credential{
				ID:     "admin",
				Secret: "adminpw",
			},
			CAName: "ca.sivatech.com",
		},
	}
}

func getDefaultPeersConfig() map[string]networkconfig.PeerConfiguration {
	return map[string]networkconfig.PeerConfiguration{
		"peer0.tictactoe.sivatech.com": {
			URL:           fmt.Sprintf("%s:7051", defaultHost),
			EventURL:      fmt.Sprintf("%s:7053", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("peer0.tictactoe.sivatech.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem"),
		},
		"peer1.tictactoe.sivatech.com": {
			URL:           fmt.Sprintf("%s:8051", defaultHost),
			EventURL:      fmt.Sprintf("%s:8053", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("peer1.tictactoe.sivatech.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem"),
		},
	}
}

func getDefaultOrganizationConfig() map[string]networkconfig.OrganizationConfiguration {
	return map[string]networkconfig.OrganizationConfiguration{
		"sivatech": {
			MSPID:      "TicTacToeMSP",
			CryptoPath: "peerOrganizations/tictactoe.sivatech.com/users/{username}@tictactoe.sivatech.com/msp",
			Peers:      []string{"peer0.tictactoe.sivatech.com", "peer1.tictactoe.sivatech.com"},
			CertificateAuthorities: []string{"ca.sivatech.com"},
		},
		"sivatechordererorg": {
			MSPID:      "SivaTechOrdererOrg",
			CryptoPath: "ordererOrganizations/sivatech.com/users/{username}@sivatech.com/msp",
		},
	}
}

func getDefaultOrderersConfig() map[string]networkconfig.OrdererConfiguration {
	return map[string]networkconfig.OrdererConfiguration{
		"orderer.sivatech.com": {
			URL:           fmt.Sprintf("%s:7050", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("orderer.sivatech.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "ordererOrganizations/sivatech.com/tlsca/tlsca.sivatech.com-cert.pem"),
		},
	}
}

func getDefaultGRPCOption(sslTargetNameOverride string) networkconfig.GRPCOptions {
	return networkconfig.GRPCOptions{
		SSLTargetNameOveride: sslTargetNameOverride,
		KeepAliveTime:        0 * time.Second,
		KeepAliveTimeOut:     20 * time.Second,
		KeepAlivePermit:      false,
		FailFast:             false,
		AllowInsecure:        false,
	}
}

func getDefaultChannelConfig() map[string]networkconfig.ChannelConfiguration {
	return map[string]networkconfig.ChannelConfiguration{
		"t3-sivatech": {
			Orderers:           []string{"orderer.sivatech.com"},
			Peers:              getDefaultChannelPeerConfig(),
			QueryChannelPolicy: getDefaultQueryChannelPolicy(),
		},
	}
}

func getDefaultQueryChannelPolicy() networkconfig.ConfigurationPolicy {
	return networkconfig.ConfigurationPolicy{
		MinResponses:        1,
		MaxTargets:          1,
		RetryAttempts:       5,
		RetryInitialBackoff: 500 * time.Millisecond,
		RetryMaxBackoff:     5 * time.Second,
		RetryBackoffFactor:  2.0,
	}
}

func getDefaultChannelPeerConfig() map[string]networkconfig.ChannelPeerConfiguration {
	return map[string]networkconfig.ChannelPeerConfiguration{
		"peer0.tictactoe.sivatech.com": {
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
		},
		"peer1.tictactoe.sivatech.com": {
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
		},
	}
}

func getDefaultClientConfig() networkconfig.ClientConfiguration {
	return networkconfig.ClientConfiguration{
		Organization:        "sivatech",
		Logging:             networkconfig.DEBUG,
		CryptoConfigPath:    "${GOPATH}/src/hyperledger/hyperledger-tictactoe/network/crypto-config",
		CredentialStorePath: "/tmp/store",
		TLSKeyPair: networkconfig.TLSKeyPathPair{
			KeyPath:  "${GOPATH}/src/hyperledger/hyperledger-tictactoe/network/client-crypto/client-key.pem",
			CertPath: "${GOPATH}/src/hyperledger/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
		},
	}
}
