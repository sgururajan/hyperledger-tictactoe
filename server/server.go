package main

import (
	"fmt"
	"hyperledger/hyperledger-tictactoe/server/network"
	"path/filepath"
	"time"
)

var defaultHost = "localhost"
var cryptoConfigPath = "${GOPATH}/src/hyperledger/hyperledger-tictactoe/crypto-config"

func main() {

}

func loadDefaultNetwork() *network.FabEndpointConfiguration {
	clientConfig := getDefaultClientConfig()
	channelsConfig := getDefaultChannelConfig()
	organizationsConfig := getDefaultOrganizationConfig()
	orderersConfig := getDefaultOrderersConfig()
	peersConfig := getDefaultPeersConfig()
	caConfig := getDefaultCAConfig()

	return network.NewFabEndpointConfiguration(caConfig,
		orderersConfig,
		peersConfig,
		organizationsConfig,
		channelsConfig,
		clientConfig)
}

func getDefaultCAConfig() map[string]network.CAConfiguration {
	return map[string]network.CAConfiguration{
		"ca.sivatech.com": {
			URL: fmt.Sprintf("%s:7054", defaultHost),
			TLSCertClientPaths: network.TLSKeyPathPair{
				KeyPath:  "${GOPATH}/src/hyperledger/hyperledger-tictactoe/client-crypto/client-key.pem",
				CertPath: "${GOPATH}/src/hyperledger/hyperledger-tictactoe/client-crypto/client-cert.pem",
			},
			TLSCertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/ca/ca.tictactoe.sivatech.com-cert.pem"),
			RegistrarCredential: network.Credential{
				ID:     "admin",
				Secret: "adminpw",
			},
			CAName: "ca.sivatech.com",
		},
	}
}

func getDefaultPeersConfig() map[string]network.PeerConfiguration {
	return map[string]network.PeerConfiguration{
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

func getDefaultOrganizationConfig() map[string]network.OrganizationConfiguration {
	return map[string]network.OrganizationConfiguration{
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

func getDefaultOrderersConfig() map[string]network.OrdererConfiguration {
	return map[string]network.OrdererConfiguration{
		"orderer.sivatech.com": {
			URL:           fmt.Sprintf("%s:7050", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("orderer.sivatech.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "ordererOrganizations/sivatech.com/tlsca/tlsca.sivatech.com-cert.pem"),
		},
	}
}

func getDefaultGRPCOption(sslTargetNameOverride string) network.GRPCOptions {
	return network.GRPCOptions{
		SSLTargetNameOveride: sslTargetNameOverride,
		KeepAliveTime:        0 * time.Second,
		KeepAliveTimeOut:     20 * time.Second,
		KeepAlivePermit:      false,
		FailFast:             false,
		AllowInsecure:        false,
	}
}

func getDefaultChannelConfig() map[string]network.ChannelConfiguration {
	return map[string]network.ChannelConfiguration{
		"t3-sivatech": {
			Orderers:           []string{"orderer.sivatech.com"},
			Peers:              getDefaultChannelPeerConfig(),
			QueryChannelPolicy: getDefaultQueryChannelPolicy(),
		},
	}
}

func getDefaultQueryChannelPolicy() network.ConfigurationPolicy {
	return network.ConfigurationPolicy{
		MinResponses:        1,
		MaxTargets:          1,
		RetryAttempts:       5,
		RetryInitialBackoff: 500 * time.Millisecond,
		RetryMaxBackoff:     5 * time.Second,
		RetryBackoffFactor:  2.0,
	}
}

func getDefaultChannelPeerConfig() map[string]network.ChannelPeerConfiguration {
	return map[string]network.ChannelPeerConfiguration{
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

func getDefaultClientConfig() network.ClientConfiguration {
	return network.ClientConfiguration{
		Organization:        "sivatech",
		Logging:             network.DEBUG,
		CryptoConfigPath:    "${GOPATH}/src/hyperledger/hyperledger-tictactoe/crypto-config",
		CredentialStorePath: "/tmp/store",
		TLSKeyPair: network.TLSKeyPathPair{
			KeyPath:  "${GOPATH}/src/hyperledger/hyperledger-tictactoe/client-crypto/client-key.pem",
			CertPath: "${GOPATH}/src/hyperledger/hyperledger-tictactoe/client-crypto/client-cert.pem",
		},
	}
}
