package appServer

import (
	"fmt"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"path/filepath"
	"time"
)

var defaultHost = "localhost"
var cryptoConfigPath = "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/crypto-config"

func DefaultNetworkConfiguration() database.Network {
	//organizationsConfig := getDefaultOrganizationConfig()
	organiations := getDefaultOrganizationInfo()
	orderersConfig := getDefaultOrderersConfig()
	peersConfig := getDefaultPeersConfig()
	caConfig := getDefaultCAConfig()
	securityConfig := getSecurityConfiguration()

	return database.Network{
		Name: "testNetwork",
		SecurityConfiguration: securityConfig,
		CertificateAuthority:  caConfig,
		IsSystemCertPool:      false,
		Orderers:              orderersConfig,
		Organizations:         organiations,
		Peers:                 peersConfig,
		//OrganizationsConfiguration: organizationsConfig,
	}
}

func getSecurityConfiguration() database.SecurityConfiguration {
	return database.SecurityConfiguration{
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

func getDefaultCAConfig() map[string]database.CertificateAuthority {
	return map[string]database.CertificateAuthority{
		"ca.sivatech.com": {
			URL: fmt.Sprintf("%s:7054", defaultHost),
			TLSCertClientPaths: database.TLSKeyPathPair{
				KeyPath:  "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-key.pem",
				CertPath: "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
			},
			TLSCertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/ca/ca.tictactoe.sivatech.com-cert.pem"),
			RegistrarCredential: database.Credential{
				ID:     "admin",
				Secret: "adminpw",
			},
			CAName: "ca.sivatech.com",
		},
	}
}

func getDefaultPeersConfig() map[string]database.Peer {
	return map[string]database.Peer{
		"peer0.tictactoe.sivatech.com": {
			URL:             fmt.Sprintf("%s:7051", defaultHost),
			EventURL:        fmt.Sprintf("%s:7053", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer0.tictactoe.sivatech.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer0.tictactoe.sivatech.com",
			Organization:    "sivatech",
		},
		"peer1.tictactoe.sivatech.com": {
			URL:             fmt.Sprintf("%s:8051", defaultHost),
			EventURL:        fmt.Sprintf("%s:8053", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer1.tictactoe.sivatech.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer1.tictactoe.sivatech.com",
			Organization:    "sivatech",
		},
	}
}

func getDefaultOrganizationInfo() map[string]database.Organization {
	return map[string]database.Organization{
		"sivatech": {
			CertificateAuthorities: []string{"ca.sivatech.com"},
			Name:       "sivatech",
			MSPID:      "TicTacToeMSP",
			CryptoPath: "peerOrganizations/tictactoe.sivatech.com/users/{username}@tictactoe.sivatech.com/msp",
			Peers:      []string{"peer0.tictactoe.sivatech.com", "peer1.tictactoe.sivatech.com"},
			MSPDir:     filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/msp"),
			ID:         "SivAtEcHoRgId",
			AdminUser:  "Admin",
		},
	}
}

func getDefaultOrderersConfig() map[string]database.Orderer {
	return map[string]database.Orderer{
		"sivatechordererorg": {
			Name:          "sivatechordererorg",
			Organization:  "sivatech",
			URL:           fmt.Sprintf("%s:7050", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("orderer.sivatech.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "ordererOrganizations/sivatech.com/tlsca/tlsca.sivatech.com-cert.pem"),
			MSPID:         "SivaTechOrdererMSP",
			CryptoPath:    "ordererOrganizations/sivatech.com/users/{username}@sivatech.com/msp",
			MSPDir:        filepath.Join(cryptoConfigPath, "ordererOrganizations/sivatech.com/msp"),
		},
	}
}

func getDefaultGRPCOption(sslTargetNameOverride string) map[string]interface{} {
	return map[string]interface{}{
		"ssl-target-name-override": sslTargetNameOverride,
		"keep-alive-time":          0 * time.Second,
		"keep-alive-timeout":       20 * time.Second,
		"keep-alive-permit":        false,
		"fail-fast":                false,
		"allow-insecure":           false,
	}
}

/*func getDefaultChannelConfig() map[string]networkconfig.ChannelConfiguration {
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
}*/

/*func GetDefaultClientConfig() entities.ClientConfiguration {
	return entities.ClientConfiguration{
		UserName:            "Admin",
		Organization:        "sivatech",
		Logging:             entities.DEBUG,
		CryptoConfigPath:    cryptoConfigPath,
		CredentialStorePath: "/tmp/store",
		TLSKeyPair: entities.TLSKeyPathPair{
			KeyPath:  "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-key.pem",
			CertPath: "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
		},
	}
}*/
