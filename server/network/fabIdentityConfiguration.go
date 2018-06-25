package network

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

type FabCAConfig struct {
	URL       string
	TLSCACert endpoint.MutualTLSConfig
	Registrar msp.EnrollCredentials
	CAName    string
}
