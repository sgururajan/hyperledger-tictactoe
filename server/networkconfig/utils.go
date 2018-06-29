package networkconfig

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/server/pathUtil"
	"crypto/x509"
	"io/ioutil"
	"encoding/pem"
)

func newTLSConfig(path string) endpoint.TLSConfig {
	config := endpoint.TLSConfig{Path: pathUtil.Substitute(path)}
	if err := config.LoadBytes(); err != nil {
		panic(errors.Errorf("error loading bytes: %s", err))
	}

	return config
}

func tlsCertBytes(path string) *x509.Certificate {
	bytes, err:= ioutil.ReadFile(pathUtil.Substitute(path))
	if err!=nil {
		return nil
	}

	block,_:= pem.Decode(bytes)
	if block!=nil {
		pub, err:= x509.ParseCertificate(block.Bytes)
		if err!=nil {
			return nil
		}

		return pub
	}

	return nil
}


