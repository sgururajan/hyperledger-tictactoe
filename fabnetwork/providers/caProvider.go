package providers

import (
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/pkg/errors"
)

type CertificateAuthorityProivder interface {
	GetCertificateAuthorities() ([]entities.CertificateAuthority, error)
}

type caProviderOptions struct {
	getCAAuthorities
}

type getCAAuthorities interface {
	GetCertificateAuthorities() ([]entities.CertificateAuthority, error)
}

func BuildCAProviderFromOptions(opts ...interface{}) (CertificateAuthorityProivder, error) {
	p:= &caProviderOptions{}
	for _, o:= range opts {
		err:= setCAInterfaceFromOption(p, o)
		if err!=nil {
			return nil, err
		}
	}

	return p, nil
}

func setCAInterfaceFromOption(p *caProviderOptions, o interface{}) error {
	s:= setter{}

	s.set(p.getCAAuthorities, func() bool {
		_,ok:= o.(getCAAuthorities)
		return ok
	}, func() {
		p.getCAAuthorities = o.(getCAAuthorities)
	})

	if !s.isSet {
		return errors.Errorf("option %#v is not a sub interface of CertificateAuthorityProivder, at least one of its functions must be implemented.", o)
	}

	return nil
}

