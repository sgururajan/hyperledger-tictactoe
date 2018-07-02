package providers

import (
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
)

type OrdererProvider interface {
	GetOrderers() ([]entities.Orderer, error)
	GetOrderersForOrgId(orgId string) ([]entities.Orderer, error)
}

type ordererProviderOptions struct {
	getOrdereres
	getOrderersForOrgId
}

type getOrdereres interface {
	GetOrderers() ([]entities.Orderer, error)
}

type getOrderersForOrgId interface {
	GetOrderersForOrgId(orgId string) ([]entities.Orderer, error)
}



func BuildOrdererProviderFromOptions(opts ...interface{}) (OrdererProvider, error) {
	prov := &ordererProviderOptions{}
	for _, o := range opts {
		err:= setInterfaceFromOption(prov, o)
		if err!=nil {
			return nil, err
		}
	}

	return prov, nil
}

func setInterfaceFromOption(p *ordererProviderOptions, i interface{}) error {
	s := &setter{}

	s.set(p.getOrdereres, func() bool { _, ok := i.(getOrdereres); return ok }, func() { p.getOrdereres = i.(getOrdereres) })
	s.set(p.getOrderersForOrgId, func() bool {
		_, ok := i.(getOrderersForOrgId)
		return ok
	}, func() {
		p.getOrderersForOrgId = i.(getOrderersForOrgId)
	})

	if !s.isSet {
		return errors.Errorf("option %#v is not a sub interface of OrdererProvider, at least one of its functions must be implemented.", i)
	}

	return nil
}
