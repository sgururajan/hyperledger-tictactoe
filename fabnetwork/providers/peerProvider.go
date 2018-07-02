package providers

import (
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/pkg/errors"
)

type PeerProvider interface {
	GetPeersForOrgId(orgId string)([]entities.Peer, error)
	GetEndrosingPeersForOrgId(orgId string) ([]entities.Peer, error)
	GetChainCodePeersForOrgId(orgId string) ([]entities.Peer, error)
	GetPeers() ([]entities.Peer, error)
}

type peerProviderOptions struct {
	getPeersForOrgId
	getEndrosingPeersForOrgId
	getChaincodePeersForOrgId
	getPeers
}

type getPeers interface {
GetPeers() ([]entities.Peer, error)
}

type getPeersForOrgId interface {
	GetPeersForOrgId(orgId string)([]entities.Peer, error)
}

type getEndrosingPeersForOrgId interface {
	GetEndrosingPeersForOrgId(orgId string) ([]entities.Peer, error)
}

type getChaincodePeersForOrgId interface {
	GetChainCodePeersForOrgId(orgId string) ([]entities.Peer, error)
}

func BuildPeerProviderWithOptions(opts ...interface{}) (PeerProvider, error) {
	p:= &peerProviderOptions{}
	for _,o:= range opts {
		err:= setPeerProviderInterfaceFromOptions(p, o)
		if err!=nil {
			return nil, err
		}
	}

	return p, nil
}

func setPeerProviderInterfaceFromOptions(c *peerProviderOptions, o interface{}) error {
	s:= &setter{}

	s.set(c.getChaincodePeersForOrgId, func() bool {
		_,ok:= o.(getChaincodePeersForOrgId)
		return ok
	}, func() {
		c.getChaincodePeersForOrgId = o.(getChaincodePeersForOrgId)
	})

	s.set(c.getEndrosingPeersForOrgId, func() bool {
		_,ok:= o.(getEndrosingPeersForOrgId)
		return ok
	}, func() {
		c.getEndrosingPeersForOrgId = o.(getEndrosingPeersForOrgId)
	})

	s.set(c.getPeersForOrgId, func() bool {
		_,ok:= o.(getPeersForOrgId)
		return ok
	}, func() {
		c.getPeersForOrgId = o.(getPeersForOrgId)
	})

	s.set(c.getPeers, func() bool {
		_,ok:= o.(getPeers)
		return ok
	}, func() {
		c.getPeers = o.(getPeers)
	})

	if !s.isSet {
		return errors.Errorf("option %#v is not a sub interface of PeerProvider, at least one of its functions must be implemented.", o)
	}

	return nil
}
