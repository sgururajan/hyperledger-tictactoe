package fabnetwork

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	msp2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/providers"
)

type FabricNetwork struct {
	Name         string
	providers    entityProviders
	clientConfig entities.ClientConfiguration
	sdkImpl      *fabSdkImpl
	sdkCAImpl    *fabSdkCAImpl
	sdk          *fabsdk.FabricSDK
	orgsByName   map[string]*entities.Organization
}

type entityProviders struct {
	caProvider       providers.CertificateAuthorityProivder
	ordererProvider  providers.OrdererProvider
	orgProvider      providers.OrganizationProvider
	peerProvider     providers.PeerProvider
	securityProvider providers.SecurityProvider
}

type Option func(opt *entityProviders) error

func NewFabricNetwork(clientConfig entities.ClientConfiguration, networkName string, opts ...Option) (*FabricNetwork, error) {
	if len(opts) == 0 {
		return nil, errors.Errorf("options must be provided")
	}

	logging.Initialize(providers.GetAppLogProvider(networkName, "fabNetwork"))

	network := &FabricNetwork{
		Name: networkName,
	}

	err := initializeNetwork(network, clientConfig, opts)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (m *FabricNetwork) Initialize() {
	err := buildOrgsByNameMap(m)
	if err != nil {
	}

	m.updateChannelConfigs()
}

func (m *FabricNetwork) Close() {
	m.sdk.Close()
}

func initializeNetwork(network *FabricNetwork, clientConfig entities.ClientConfiguration, opts []Option) error {
	for _, o := range opts {
		err := o(&network.providers)
		if err != nil {
			return errors.WithMessage(err, "error in options passed to New method")
		}
	}

	var err error
	network.sdkCAImpl, err = newFabSdkCAImpl(network)

	if err != nil {
		return err
	}

	network.sdkImpl, err = newFabSdkImpl(network, clientConfig)

	if err != nil {
		return err
	}

	sdkOptions := []fabsdk.Option{
		fabsdk.WithCryptoSuiteConfig(network.providers.securityProvider),
		fabsdk.WithEndpointConfig(network.sdkImpl),
		fabsdk.WithIdentityConfig(network.sdkCAImpl, network.sdkImpl),
	}

	sdk, err := fabsdk.New(nil, sdkOptions...)
	if err != nil {
		return err
	}

	network.sdk = sdk

	return nil
}

func (m *FabricNetwork) getResourceManagementClient(org entities.Organization) (*resmgmt.Client, error) {
	clientContext := m.sdk.Context(fabsdk.WithOrg(org.Name), fabsdk.WithUser(org.AdminUser))
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return nil, err
	}

	return resMgmtClient, nil
}

func (m *FabricNetwork) getMembershipClient(orgName string) (*msp.Client, error) {
	mspClient, err := msp.New(m.sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}

	return mspClient, nil
}

func (m *FabricNetwork) getSigningIdentity(mspClient *msp.Client, userName string) (msp2.SigningIdentity, error) {
	signingIdentity, err := mspClient.GetSigningIdentity(userName)
	return signingIdentity, err
}

func buildOrgsByNameMap(network *FabricNetwork) error {
	network.orgsByName = make(map[string]*entities.Organization)
	orgs, err := network.providers.orgProvider.GetOrganizations()
	if err != nil {
		return err
	}

	for _, v := range orgs {
		org := v
		network.orgsByName[org.Name] = &org
	}

	return nil
}

func WithOrdererProvider(opts ...interface{}) Option {
	return func(p *entityProviders) error {
		c, err := providers.BuildOrdererProviderFromOptions(opts...)
		if err != nil {
			return err
		}
		p.ordererProvider = c
		return nil
	}
}

func WithPeerProvider(opts ...interface{}) Option {
	return func(p *entityProviders) error {
		c, err := providers.BuildPeerProviderWithOptions(opts...)
		if err != nil {
			return err
		}
		p.peerProvider = c
		return nil
	}
}

func WithOrganizationProvider(opts ...interface{}) Option {
	return func(p *entityProviders) error {
		c, err := providers.BuildOrganizationProviderFromOptions(opts...)
		if err != nil {
			return err
		}
		p.orgProvider = c
		return nil
	}
}

func WithCAProvider(opts ...interface{}) Option {
	return func(p *entityProviders) error {
		c, err := providers.BuildCAProviderFromOptions(opts...)
		if err != nil {
			return err
		}
		p.caProvider = c
		return nil
	}
}

func WithSecurityConfig(opts ...interface{}) Option {
	return func(p *entityProviders) error {
		c, err := providers.BuildSecurityProviderFromOptions(opts...)
		if err != nil {
			return err
		}
		p.securityProvider = c
		return nil
	}
}
