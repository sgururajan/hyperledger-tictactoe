package providers

import (
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/pkg/errors"
)

type OrganizationProvider interface {
	GetOrganizations() ([]entities.Organization, error)
}

type organizationProviderOption struct {
	getOrganizations
}

type getOrganizations interface {
	GetOrganizations() ([]entities.Organization, error)
}

func BuildOrganizationProviderFromOptions(opts ...interface{}) (OrganizationProvider, error) {
	p:= &organizationProviderOption{}
	for _,o:= range opts {
		err:= setOrgProviderInterfaceFromOption(p, o)
		if err!=nil {
			return nil, err
		}
	}

	return p, nil
}

func setOrgProviderInterfaceFromOption(c *organizationProviderOption, o interface{}) error {
	s:= &setter{}

	s.set(c.getOrganizations, func() bool {
		_, ok:= o.(getOrganizations)
		return ok
	}, func() {
		c.getOrganizations = o.(getOrganizations)
	})

	if !s.isSet {
		return errors.Errorf("option %#v is not a sub interface of OrganizationProvider, at least one of its functions must be implemented.", o)
	}

	return nil
}

