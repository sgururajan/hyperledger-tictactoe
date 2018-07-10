package database

import (
	"testing"
	"fmt"
	"time"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"path/filepath"
)

var cryptoConfigPath="${GOPATH}/github.com/sgururajan/hyperledter-tictactoe/network/crypto-config"
var host="localhost"
var testnetworkName="testnetwork"

var mockOrgs= map[string]Organization{
	"sivatech": {
		MSPID:      "TicTacToeMSP",
		CryptoPath: "peerOrganizations/tictactoe.sivatech.com/users/{username}@tictactoe.sivatech.com/msp",
		Peers:      []string{"peer0.tictactoe.sivatech.com", "peer1.tictactoe.sivatech.com"},
		CertificateAuthorities: []string{"ca.sivatech.com"},
	},
}

var mockOrderers=map[string]Orderer{
	"orderer.sivatech.com": {
		//URL: "orderer.sivatech.com:7050",
		//URL: "192.168.2.5:7050",
		Name: "sivatechordererorg",
		URL: fmt.Sprintf("%s:7050", host),
		GRPCOptions: map[string]interface{}{
			"ssl-target-name-override": "orderer.sivatech.com",
			"keep-alive-time":          0 * time.Second,
			"keep-alive-timeout":       20 * time.Second,
			"keep-alive-permit":        false,
			"fail-fast":                false,
			"allow-insecure":           false,
		},
		TLSCACertPath: utils.Substitute(filepath.Join(cryptoConfigPath,"ordererOrganizations/sivatech.com/tlsca/tlsca.sivatech.com-cert.pem")),
		//TLSCACerts: newTLSConfig(path.Join(CryptoConfigPath, "/ordererOrganizations/sivatech.com/tlsca/tlsca.sivatech.com-cert.pem")),
	},
}

var mockPeers = map[string]Peer{
	"peer0.tictactoe.sivatech.com": {
		// URL:      "peer0.tictactoe.sivatech.com:7051",
		// EventURL: "peer0.tictactoe.sivatech.com:7053",
		//URL:      "192.168.2.5:7051",
		//EventURL: "192.168.2.5:7053",
		URL:      fmt.Sprintf("%s:7051", host),
		EventURL: fmt.Sprintf("%s:7053", host),
		GrpcOptions: map[string]interface{}{
			"ssl-target-name-override": "peer0.tictactoe.sivatech.com",
			"keep-alive-time":          0 * time.Second,
			"keep-alive-timeout":       20 * time.Second,
			"keep-alive-permit":        false,
			"fail-fast":                false,
			"allow-insecure":           false,
		},
		TLSCertPath: utils.Substitute(filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem")),
		//TLSCACerts: newTLSConfig(path.Join(CryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem")),
	},
	"peer1.tictactoe.sivatech.com": {
		// URL:      "peer1.tictactoe.sivatech.com:8051",
		// EventURL: "peer1.tictactoe.sivatech.com:8053",
		//URL:      "192.168.2.5:8051",
		//EventURL: "192.168.2.5:8053",
		URL:      fmt.Sprintf("%s:8051", host),
		EventURL: fmt.Sprintf("%s:8053", host),
		GrpcOptions: map[string]interface{}{
			"ssl-target-name-override": "peer1.tictactoe.sivatech.com",
			"keep-alive-time":          0 * time.Second,
			"keep-alive-timeout":       20 * time.Second,
			"keep-alive-permit":        false,
			"fail-fast":                false,
			"allow-insecure":           false,
		},
		TLSCertPath: utils.Substitute(filepath.Join(cryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem")),
		//TLSCACerts: newTLSConfig(path.Join(CryptoConfigPath, "peerOrganizations/tictactoe.sivatech.com/tlsca/tlsca.tictactoe.sivatech.com-cert.pem")),
	},
}

var mockCAs = map[string]CertificateAuthority{
	"ca.sivatech.com": {
		// URL: "https://ca.sivatech.com:7054",
		//URL: "https://192.168.2.5:7054",
		URL: fmt.Sprintf("%s:7054", host),
		TLSCertClientPaths: TLSKeyPathPair{
			KeyPath: utils.Substitute("${GOPATH}/src/hyperledger/tictactoe/client-crypto/client-key.pem"),
			CertPath: utils.Substitute("${GOPATH}/src/hyperledger/tictactoe/client-crypto/client-cert.pem"),
		},
		TLSCertPath: utils.Substitute(filepath.Join(cryptoConfigPath,"peerOrganizations/tictactoe.sivatech.com/ca/ca.tictactoe.sivatech.com-cert.pem" )),
		RegistrarCredential: Credential{
			ID:     "admin",
			Secret: "adminpw",
		},
		CAName: "ca.sivatech.com",
	},
}

var mockSecurityConfig = SecurityConfiguration{
	KeyStoragePath: "/tmp/msp/keystore",
	ProviderLabel:"",
	ProviderPin: "",
	ProviderLibPath: "",
	IsSoftVerify: true,
	Provider: "sw",
	Level: 256,
	Algorithm: "SHA2",
	IsEnabled: true,
}

var mockFileRepo = &NetworkFileRepository{
	networks:map[string]Network {
		"testnetwork" :{
			Orderers:mockOrderers,
			CertificateAuthority:mockCAs,
			Name:"testnetwork",
			Peers:mockPeers,
			Organizations:mockOrgs,
			IsSystemCertPool:true,
			SecurityConfiguration:mockSecurityConfig,
			Consortiums:map[string][]string{},
		},
	},
}

func TestNewNetworkFileRepository(t *testing.T) {
	repo:= mockFileRepo
	if repo==nil {
		t.Fatal("repository intialization failed")
	}
}

func TestNetworkFileRepository_GetCertificateAuthorities(t *testing.T) {
	repo:= mockFileRepo
	cas,err:= repo.GetCertificateAuthorities(testnetworkName)
	if err!=nil {
		t.Fatal("get certificate authorities should return atleast one")
	}

	if len(cas)==0 {
		t.Fatal("atleast one ca expected")
	}

	if (cas[0].CAName!="ca.sivatech.com") {
		t.Fatal("unexpected ca name")
	}
}

func TestNetworkFileRepository_GetNetworks(t *testing.T) {
	repo:= mockFileRepo
	networks, err:= repo.GetNetworks()
	if err!=nil {
		t.Fatal("get networks should return valid response")
	}

	if len(networks)==0 {
		t.Fatal("at least one network expected")
	}
}

func TestNetworkFileRepository_GetOrderers(t *testing.T) {
	repo:= mockFileRepo
	ords,err:= repo.GetOrderers(testnetworkName)
	if err!=nil {
		t.Fatal("get orderers should return valid response")
	}

	if len(ords)==0 {
		t.Fatal("network should have atleast one orderer")
	}

	if ords[0].Name != "sivatechordererorg" {
		t.Fatal("unexpected orderer name")
	}
}