package fabnetwork

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"strings"
	"io/ioutil"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
)

type fabCAconfig struct {
	url       string
	tlsCACert endpoint.MutualTLSConfig
	registrar msp.EnrollCredentials
	caName    string
}

type fabSdkCAImpl struct {
	caConfigs map[string]fabCAconfig
	organizations map[string]entities.Organization
}

func newFabSdkCAImpl(network *FabricNetwork) (*fabSdkCAImpl, error) {
	caImpl:= &fabSdkCAImpl{}

	err:= caImpl.createCAConfigs(network)
	if err!=nil {
		return nil, err
	}

	caImpl.organizations = make(map[string]entities.Organization)
	orgs, err:= network.providers.orgProvider.GetOrganizations()
	if err!=nil {
		return nil, err
	}
	for _, v:= range orgs {
		caImpl.organizations[v.Name]=v
	}

	return caImpl, nil
}

func (m *fabSdkCAImpl) createCAConfigs(network *FabricNetwork) error {
	m.caConfigs = make(map[string]fabCAconfig)

	configs, err:= network.providers.caProvider.GetCertificateAuthorities()
	if err!=nil {
		return err
	}

	for _, v:= range configs {
		m.caConfigs[v.CAName] = fabCAconfig{
			url: v.URL,
			tlsCACert: endpoint.MutualTLSConfig{
				Path: v.TLSCertPath,
				Client: endpoint.TLSKeyPair{
					Key:newTLSConfig(v.TLSCertClientPaths.KeyPath),
					Cert:newTLSConfig(v.TLSCertClientPaths.CertPath),
				},
			},
			registrar: msp.EnrollCredentials{
				EnrollID: v.RegistrarCredential.ID,
				EnrollSecret:v.RegistrarCredential.Secret,
			},
			caName: v.CAName,
		}
	}

	return nil
}



func (m *fabSdkCAImpl) CAConfig(orgName string) (*msp.CAConfig, bool) {
	return m.getCAConfig(orgName)
}

func (m *fabSdkCAImpl) CAServerCerts(orgName string) ([][]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAServerCerts, true
}

func (m *fabSdkCAImpl) CAClientKey(orgName string) ([]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientKey, true
}

func (m *fabSdkCAImpl) CAClientCert(orgName string) ([]byte, bool) {
	caConfig, ok := m.getCAConfig(orgName)
	if !ok {
		return nil, false
	}

	return caConfig.TLSCAClientCert, true
}

func (m *fabSdkCAImpl) CAKeyStorePath() string {
	return "/tmp/msp"
}

func (m *fabSdkCAImpl) CredentialStorePath() string {
	return "/tmp/state-store"
}

func (m *fabSdkCAImpl) getCAConfig(orgName string) (*msp.CAConfig, bool) {
	if len(m.organizations[strings.ToLower(orgName)].CertificateAuthorities) == 0 {
		return nil, false
	}

	org := m.organizations[strings.ToLower(orgName)]
	certAuthName := org.CertificateAuthorities[0]
	if certAuthName == "" {
		return nil, false
	}

	caConfig, ok := m.caConfigs[strings.ToLower(certAuthName)]
	if !ok {
		return nil, false
	}

	mspCAConfig, err := caConfig.getMSPCAConfig()
	if err != nil {
		return nil, false
	}

	return mspCAConfig, true
}

func (m *fabCAconfig) getMSPCAConfig() (*msp.CAConfig, error) {
	serverCerts, err := m.getServerCerts()
	if err != nil {
		return nil, err
	}

	return &msp.CAConfig{
		URL:              m.url,
		Registrar:        m.registrar,
		CAName:           m.caName,
		TLSCAClientCert:  m.tlsCACert.Client.Cert.Bytes(),
		TLSCAClientKey:   m.tlsCACert.Client.Key.Bytes(),
		TLSCAServerCerts: serverCerts,
	}, nil
}

func (m *fabCAconfig) getServerCerts() ([][]byte, error) {
	var serverCerts [][]byte

	pems := m.tlsCACert.Pem
	if len(pems) > 0 {
		serverCerts := make([][]byte, len(pems))
		for i, pem := range pems {
			serverCerts[i] = []byte(pem)
		}

		return serverCerts, nil
	}

	certFiles := strings.Split(m.tlsCACert.Path, ",")
	serverCerts = make([][]byte, len(certFiles))
	for i, certPath := range certFiles {
		bytes, err := ioutil.ReadFile(utils.Substitute(certPath))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to load server certificates")
		}
		serverCerts[i] = bytes
	}

	return serverCerts, nil
}
