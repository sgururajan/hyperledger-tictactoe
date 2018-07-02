package fabnetwork

import (
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/providers"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"strings"
	fab2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"path/filepath"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"github.com/spf13/viper"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	msp2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/cloudflare/cfssl/log"
)

type FabricNetwork struct {
	providers entityProviders
	clientConfig entities.ClientConfiguration
	sdkImpl *fabSdkImpl
	sdkCAImpl *fabSdkCAImpl
	sdk *fabsdk.FabricSDK
	orgsByName map[string]*entities.Organization
}

type entityProviders struct {
	caProvider       providers.CertificateAuthorityProivder
	ordererProvider  providers.OrdererProvider
	orgProvider      providers.OrganizationProvider
	peerProvider     providers.PeerProvider
	securityProvider providers.SecurityProvider
}

type Option func(opt *entityProviders) error

func NewFabricNetwork(clientConfig entities.ClientConfiguration, opts ...Option) (*FabricNetwork, error) {
	if len(opts) == 0 {
		return nil, errors.Errorf("options must be provided")
	}

	network := &FabricNetwork{}

	err:= initializeNetwork(network, clientConfig, opts)
	if err != nil {
		return nil, err
	}

	err = buildOrgsByNameMap(network)
	if err != nil {
		return nil, err
	}

	network.updateChannelConfigs()

	return network, nil
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

	if err!=nil {
		return err
	}

	network.sdkImpl, err = newFabSdkImpl(network, clientConfig)

	if err!=nil {
		return err
	}

	sdkOptions:= []fabsdk.Option{
		fabsdk.WithCryptoSuiteConfig(network.providers.securityProvider),
		fabsdk.WithEndpointConfig(network.sdkImpl),
		fabsdk.WithIdentityConfig(network.sdkCAImpl, network.sdkImpl),
	}

	sdk,err:= fabsdk.New(nil, sdkOptions...)
	if err != nil {
		return err
	}

	network.sdk=sdk;

	return nil
}

func (m *FabricNetwork) updateChannelConfigs() error {
	fabChannels:= make(map[string][]entities.Peer)
	orgs, err:= m.providers.orgProvider.GetOrganizations()
	if err != nil {
		return err
	}

	for _, o:= range orgs {
		chResps,err:= m.getChannelsForOrg(o)
		if err != nil {
			return err
		}

		for k,v:= range chResps {
			if _,exists:= fabChannels[k];!exists{
				fabChannels[k] = []entities.Peer{}
			}

			fabChannels[k] = append(fabChannels[k], v...)
		}
	}

	orderers,err:= m.providers.ordererProvider.GetOrderers()
	if err != nil {
		return err
	}
	chOrderers:= []string{}
	for _,o:= range orderers {
		chOrderers = append(chOrderers, o.Name)
	}

	chConfigs:= m.sdkImpl.channels

	for k,v:=range fabChannels {
		chConfig:= fab2.ChannelEndpointConfig{
			Orderers: chOrderers,
			Policies: fab2.ChannelPolicies{
				QueryChannelConfig: defaultQueryChannelPolicy(),
			},
		}

		//m.sdkImpl.channels[strings.ToLower(k)].Peers = make(map[string]fab2.PeerChannelConfig)
		chPeers:= make(map[string]fab2.PeerChannelConfig)
		for _,pv:= range v{
			chPeers[pv.EndPoint] = fab2.PeerChannelConfig{
				EventSource: pv.EventSource,
				LedgerQuery: pv.CanQueryLedger,
				ChaincodeQuery: pv.IsChainCodePeer,
				EndorsingPeer: pv.IsEndrosingPeer,
			}
		}

		chConfig.Peers=make(map[string]fab2.PeerChannelConfig)
		for pk,pv:= range chPeers {
			chConfig.Peers[pk]=pv
		}
		//m.sdkImpl.channels[strings.ToLower(k)]=chConfig
		chConfigs[strings.ToLower(k)] = chConfig
	}

	m.sdkImpl.channels = chConfigs

	return nil
}

func (m *FabricNetwork) getChannelsForOrg(org entities.Organization) (map[string][]entities.Peer, error) {
	resMgmtClient, err:= m.getResourceManagementClient(org)
	if err != nil {
		return nil, err
	}
	endrosingPeers,err:= m.providers.peerProvider.GetEndrosingPeersForOrgId(org.Name)
	if err != nil {
		return nil, err
	}

	result:= make(map[string][]entities.Peer)

	for _,p:= range endrosingPeers {
		chResp, err:= resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(p.EndPoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return nil, err
		}

		for _,chInfo:= range chResp.Channels {
			if _,ok:= result[chInfo.ChannelId]; !ok {
				result[chInfo.ChannelId] = []entities.Peer{}
			}
			result[chInfo.ChannelId] = append(result[chInfo.ChannelId], p)
		}
	}

	return result, nil
}



func (m *FabricNetwork) getResourceManagementClient(org entities.Organization) (*resmgmt.Client, error)  {
	clientContext:= m.sdk.Context(fabsdk.WithOrg(org.Name), fabsdk.WithUser(org.AdminUser))
	resMgmtClient, err:= resmgmt.New(clientContext)
	if err!=nil {
		return nil, err
	}

	return resMgmtClient, nil
}

func (m *FabricNetwork) getMembershipClient(orgName string) (*msp.Client, error) {
	mspClient, err:= msp.New(m.sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient, nil
}

func (m *FabricNetwork) getSigningIdentity(mspClient *msp.Client, userName string) (msp2.SigningIdentity, error) {
	signingIdentity,err:= mspClient.GetSigningIdentity(userName)
	return signingIdentity,err
}

func buildOrgsByNameMap(network *FabricNetwork) error {
	network.orgsByName = make(map[string]*entities.Organization)
	orgs, err:= network.providers.orgProvider.GetOrganizations()
	if err != nil {
		return err
	}

	for _,v:= range orgs {
		network.orgsByName[v.Name] = &v
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

func (m *FabricNetwork) GetChannelList() ([]string, error) {
	result:= []string{}
	for k,_:= range m.sdkImpl.channels {
		result = append(result, k)
	}

	return result, nil
}

func (m *FabricNetwork) HasChannel(channelName string) (bool, error) {
	chList, err:= m.GetChannelList()
	if err != nil {
		return false, err
	}

	for _,v:= range chList {
		if v==channelName {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) CreateChannel(creatorOrgName string, chRequest entities.CreateChannelRequest) error {

	hasChannel, err:= m.HasChannel(chRequest.ChannelName)
	if err != nil {
		return err
	}
	if hasChannel {
		return nil
	}

	creatorOrg, ok:= m.orgsByName[creatorOrgName]
	if !ok {
		return errors.Errorf("org %s does not belong to this network", creatorOrgName)
	}

	txFileName:= filepath.Join(utils.Substitute(viper.GetString("channelArtifactsPath")), chRequest.ChannelName + ".tx")
	txHelper:= newFabConfixTxHelper()
	err = txHelper.createChannelTransaction(m, chRequest, txFileName)
	if err != nil {
		return err
	}

	mspClient,err:= m.getMembershipClient(creatorOrg.Name)
	if err != nil {
		return err
	}

	adminIdentity,err:= m.getSigningIdentity(mspClient, creatorOrg.AdminUser)
	if err != nil {
		return err
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         chRequest.ChannelName,
		ChannelConfigPath: txFileName,
		SigningIdentities: []msp2.SigningIdentity{adminIdentity},
	}
	resMgmtClient, err:= m.getResourceManagementClient(*creatorOrg)
	if err != nil {
		return err
	}

	orderers, err:= m.providers.ordererProvider.GetOrderers()
	if err != nil {
		return err
	}

	if len(orderers)==0 {
		return errors.Errorf("no orderer for organization %s", creatorOrgName)
	}

	saveResp, err:= resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(orderers[0].Name))
	if err != nil {
		return err
	}

	log.Infof("channel with name %s created successfully. TxID: %v", chRequest.ChannelName, saveResp.TransactionID)

	err = resMgmtClient.JoinChannel(chRequest.ChannelName, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderers[0].Name))

	if err != nil {
		return err
	}

	log.Infof("channel joined successfully")

	err = m.updateChannelConfigs()
	if err != nil {
		return err
	}

	return nil
}

