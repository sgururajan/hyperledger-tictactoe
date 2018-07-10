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
		Name: "testnetwork",
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
		"ca.org1.tictactoe.com": {
			URL: fmt.Sprintf("%s:7054", defaultHost),
			TLSCertClientPaths: database.TLSKeyPathPair{
				KeyPath:  "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-key.pem",
				CertPath: "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
			},
			TLSCertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/org1.tictactoe.com/ca/ca.org1.tictactoe.com-cert.pem"),
			RegistrarCredential: database.Credential{
				ID:     "admin",
				Secret: "adminpw",
			},
			CAName: "ca.org1.tictactoe.com",
		},
		"ca.org2.tictactoe.com": {
			URL: fmt.Sprintf("%s:8054", defaultHost),
			TLSCertClientPaths: database.TLSKeyPathPair{
				KeyPath:  "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-key.pem",
				CertPath: "${GOPATH}/src/github.com/sgururajan/hyperledger-tictactoe/network/client-crypto/client-cert.pem",
			},
			TLSCertPath: filepath.Join(cryptoConfigPath, "peerOrganizations/org2.tictactoe.com/ca/ca.org2.tictactoe.com-cert.pem"),
			RegistrarCredential: database.Credential{
				ID:     "admin",
				Secret: "adminpw",
			},
			CAName: "ca.org2.tictactoe.com",
		},
	}
}

func getDefaultPeersConfig() map[string]database.Peer {
	return map[string]database.Peer{
		"peer0.org1.tictactoe.com": {
			URL:             fmt.Sprintf("%s:7051", defaultHost),
			EventURL:        fmt.Sprintf("%s:7053", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer0.org1.tictactoe.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/org1.tictactoe.com/tlsca/tlsca.org1.tictactoe.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer0.org1.tictactoe.com",
			Organization:    "org1",
		},
		"peer1.org1.tictactoe.com": {
			URL:             fmt.Sprintf("%s:8051", defaultHost),
			EventURL:        fmt.Sprintf("%s:8053", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer1.org1.tictactoe.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/org1.tictactoe.com/tlsca/tlsca.org1.tictactoe.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer1.org1.tictactoe.com",
			Organization:    "org1",
		},
		"peer0.org2.tictactoe.com": {
			URL:             fmt.Sprintf("%s:7055", defaultHost),
			EventURL:        fmt.Sprintf("%s:7057", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer0.org2.tictactoe.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/org2.tictactoe.com/tlsca/tlsca.org2.tictactoe.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer0.org2.tictactoe.com",
			Organization:    "org2",
		},
		"peer1.org2.tictactoe.com": {
			URL:             fmt.Sprintf("%s:8055", defaultHost),
			EventURL:        fmt.Sprintf("%s:8057", defaultHost),
			GrpcOptions:     getDefaultGRPCOption("peer1.org2.tictactoe.com"),
			TLSCertPath:     filepath.Join(cryptoConfigPath, "peerOrganizations/org2.tictactoe.com/tlsca/tlsca.org2.tictactoe.com-cert.pem"),
			IsEndrosingPeer: true,
			IsChainCodePeer: true,
			CanQueryLedger:  true,
			EventSource:     true,
			EndPoint:        "peer1.org2.tictactoe.com",
			Organization:    "org2",
		},
	}
}

func getDefaultOrganizationInfo() map[string]database.Organization {
	return map[string]database.Organization{
		"org1": {
			CertificateAuthorities: []string{"ca.org1.tictactoe.com"},
			Name:       "org1",
			MSPID:      "Org1TicTacToeMSP",
			CryptoPath: "peerOrganizations/org1.tictactoe.com/users/{username}@org1.tictactoe.com/msp",
			Peers:      []string{"peer0.org1.tictactoe.com", "peer1.org1.tictactoe.com"},
			MSPDir:     filepath.Join(cryptoConfigPath, "peerOrganizations/org1.tictactoe.com/msp"),
			ID:         "org1",
			AdminUser:  "Admin",
			IsOrderer:  false,
		},
		"org2": {
			CertificateAuthorities: []string{"ca.org2.tictactoe.com"},
			Name:       "org2",
			MSPID:      "Org2TicTacToeMSP",
			CryptoPath: "peerOrganizations/org2.tictactoe.com/users/{username}@org2.tictactoe.com/msp",
			Peers:      []string{"peer0.org2.tictactoe.com", "peer1.org2.tictactoe.com"},
			MSPDir:     filepath.Join(cryptoConfigPath, "peerOrganizations/org2.tictactoe.com/msp"),
			ID:         "org2",
			AdminUser:  "Admin",
			IsOrderer:  false,
		},
		"tictactoeordererorg": {
			CertificateAuthorities: []string{"ca.org2.tictactoe.com"},
			Name:       "tictactoeordererorg",
			MSPID:      "tictactoeorderermsp",
			CryptoPath: "ordererOrganizations/tictactoe.com/users/{username}@tictactoe.com/msp",
			ID:         "tictactoeordererorg",
			AdminUser:  "Admin",
			IsOrderer:  true,
		},
	}
}

func getDefaultOrderersConfig() map[string]database.Orderer {
	return map[string]database.Orderer{
		"ordererorg": {
			Name:          "tictactoeordererorg",
			Organization:  "ordererorg",
			URL:           fmt.Sprintf("%s:7050", defaultHost),
			GRPCOptions:   getDefaultGRPCOption("orderer.tictactoe.com"),
			TLSCACertPath: filepath.Join(cryptoConfigPath, "ordererOrganizations/tictactoe.com/tlsca/tlsca.tictactoe.com-cert.pem"),
			MSPID:         "tictactoeorderermsp",
			CryptoPath:    "ordererOrganizations/tictactoe.com/users/{username}@tictactoe.com/msp",
			MSPDir:        filepath.Join(cryptoConfigPath, "ordererOrganizations/tictactoe.com/msp"),
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
