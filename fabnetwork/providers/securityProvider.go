package providers

import "github.com/pkg/errors"

type SecurityProvider interface {
	IsSecurityEnabled() bool
	SecurityAlgorithm() string
	SecurityLevel() int
	SecurityProvider() string
	SoftVerify() bool
	SecurityProviderLibPath() string
	SecurityProviderPin() string
	SecurityProviderLabel() string
	KeyStorePath() string
}

type securityProviderOptions struct {
	isSecurityEnabled
	securityAlgorithm
	securityLevel
	securityProvider
	softVerify
	securityProviderLibPath
	securityProviderPin
	securityProviderLabel
	keyStorePath
}

type isSecurityEnabled interface {
	IsSecurityEnabled() bool
}

type securityAlgorithm interface {
	SecurityAlgorithm() string
}

type securityLevel interface {
	SecurityLevel() int
}

type securityProvider interface {
	SecurityProvider() string
}

type softVerify interface {
	SoftVerify() bool
}

type securityProviderLibPath interface {
	SecurityProviderLibPath() string
}

type securityProviderPin interface {
	SecurityProviderPin() string
}

type securityProviderLabel interface {
	SecurityProviderLabel() string
}

type keyStorePath interface {
	KeyStorePath() string
}

func BuildSecurityProviderFromOptions(opts ...interface{}) (SecurityProvider, error) {
	p := &securityProviderOptions{}

	for _, o := range opts {
		err := setSecurityProviderOptionFromInterface(p, o)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func setSecurityProviderOptionFromInterface(p *securityProviderOptions, opt interface{}) error {
	s := &setter{}

	s.set(p.isSecurityEnabled, func() bool { _, ok := opt.(isSecurityEnabled); return ok }, func() { p.isSecurityEnabled = opt.(isSecurityEnabled) })
	s.set(p.securityAlgorithm, func() bool { _, ok := opt.(securityAlgorithm); return ok }, func() { p.securityAlgorithm = opt.(securityAlgorithm) })
	s.set(p.securityLevel, func() bool { _, ok := opt.(securityLevel); return ok }, func() { p.securityLevel = opt.(securityLevel) })
	s.set(p.securityProvider, func() bool { _, ok := opt.(securityProvider); return ok }, func() { p.securityProvider = opt.(securityProvider) })
	s.set(p.softVerify, func() bool { _, ok := opt.(softVerify); return ok }, func() { p.softVerify = opt.(softVerify) })
	s.set(p.securityProviderLibPath, func() bool { _, ok := opt.(securityProviderLibPath); return ok }, func() { p.securityProviderLibPath = opt.(securityProviderLibPath) })
	s.set(p.securityProviderPin, func() bool { _, ok := opt.(securityProviderPin); return ok }, func() { p.securityProviderPin = opt.(securityProviderPin) })
	s.set(p.securityProviderLabel, func() bool { _, ok := opt.(securityProviderLabel); return ok }, func() { p.securityProviderLabel = opt.(securityProviderLabel) })
	s.set(p.keyStorePath, func() bool { _, ok := opt.(keyStorePath); return ok }, func() { p.keyStorePath = opt.(keyStorePath) })

	if !s.isSet {
		return errors.Errorf("option %#v is not a sub interface of PeerProvider, at least one of its functions must be implemented.", opt)
	}

	return nil
}
