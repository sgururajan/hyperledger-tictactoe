package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/sgururajan/hyperledger-tictactoe/server/networkconfig"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/server/blockchainconfig"
	"path/filepath"
	"github.com/sgururajan/hyperledger-tictactoe/server/pathUtil"
)

type FabricNetwork struct {
	networkConfig *networkconfig.FabNetworkConfiguration
	sdk           *fabsdk.FabricSDK
	identity      mspctx.Identity
	adminClient   resmgmt.Client
	isInitialized bool
}

func NewFabricNetwork(config *networkconfig.FabNetworkConfiguration) *FabricNetwork {
	return &FabricNetwork{
		networkConfig: config,
		isInitialized: false,
	}
}

func (m *FabricNetwork) Initialize() error {
	if m.isInitialized {
		return fmt.Errorf("network %s already intialized", m.networkConfig.Name)
	}

	sdkOptions:= []fabsdk.Option{
		fabsdk.WithEndpointConfig(m.networkConfig),
		fabsdk.WithCryptoSuiteConfig(m.networkConfig.SecurityConfiguration),
		fabsdk.WithIdentityConfig(m.networkConfig),
	}

	sdk, err:= fabsdk.New(nil, sdkOptions...)

	if err!=nil {
		return errors.WithMessage(err, "unable to intialize sdk")
	}

	m.sdk = sdk

	return nil
}

func (m *FabricNetwork) IsChannelExists(channelId string) (bool, error) {

	//clientContext := network.sdk.Context(fabsdk.WithUser(network.OrgAdmin), fabsdk.WithOrg(network.OrgName))
	//
	//resMgmtClient, err := resmgmt.New(clientContext)
	//if err != nil {
	//	panic(err)
	//}

	//ordererName, ok:= getDefaultOrdererEndpoint(m.networkConfig)
	//if !ok {
	//	return false,errors.Errorf("unable to get default orderer")
	//}

	/*orgInfo, ok:= getDefaultOrgInfo(m.networkConfig)
	if !ok || orgInfo.Endpoint=="" {
		return false, errors.Errorf("no organization found")
	}

	orgKey:= orgInfo.Endpoint
	peerKey, ok:= getDefaultEndrosingPeerEndpointForOrganization(m.networkConfig, orgKey)
	if !ok {
		return false, errors.Errorf("no peers found for org %s", orgKey)
	}

	contextProvider:= m.sdk.Context(fabsdk.WithUser(orgInfo.AdminUserName), fabsdk.WithOrg(orgKey))

	resMgmntClient, err:= resmgmt.New(contextProvider)
	if err!= nil {
		return false, errors.WithMessage(err, "unable to create resource mangement client")
	}

	qResp, err:= resMgmntClient.QueryChannels(resmgmt.WithTargetEndpoints(peerKey), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err!=nil{
		return false,errors.WithMessage(err, "query for channels failed")
	}

	for _,chInfo:= range qResp.Channels {
		if chInfo.ChannelId==channelId {
			return true, nil
		}
	}*/

	return false, nil
}

func (m *FabricNetwork) CreateChannel(channelName string, orgs []string) error {

	txFileName, err:= m.createChannelTx(channelName, orgs)
	if err!=nil {
		return err
	}

	fmt.Println(txFileName)

	return nil
}

func (m *FabricNetwork) createChannelTx(channelName string, orgs []string) (string, error) {

	chOrgs:= []blockchainconfig.ChannelOrg{}
	// add the first orderers by default
	for _, o:= range m.networkConfig.OrgsInfo {
		if o.IsOrderer {
			chOrgs = append(chOrgs, blockchainconfig.ChannelOrg{
				Name: o.Name,
				IsOrderer:true,
				Endpoint:o.Endpoint,
			})
		}
	}

	for _,o:= range orgs {
		org, ok:= m.networkConfig.OrgsByName[o]
		if ok && !org.IsOrderer {
			peers:= []string{}
			for _,p:= range m.networkConfig.PeersInfo {
				if p.OrgName==org.Name {
					peers = append(peers, p.Endpoint)
				}
			}
			chOrgs = append(chOrgs, blockchainconfig.ChannelOrg{
				Name: org.Name,
				IsOrderer: false,
				AnchorPeers: peers,
				Endpoint:"",
			})
		}
	}

	channelTxHelper:= blockchainconfig.NewChannelTxConfigHelper()
	err:= channelTxHelper.CreateChannelTxObject(m.networkConfig,"testchannel", chOrgs, "orderer.sivatech.com" )
	if err!=nil {
		return "","",err
	}

	exePath:= pathUtil.GetExecutablePath()
	txFileName:= filepath.Join(exePath, "channel-artifacts", channelName+".pb")
	gbFileName:= filepath.Join(exePath, "channel-artifacts", channelName+"_genesis.block")

	pathUtil.EnsureDirectory(exePath, "channel-artifacts")

	err = channelTxHelper.CreateConfigurationBlocks(txFileName, gbFileName)
	if err!=nil {
		fmt.Println(err.Error())
		return "","",err
	}

	return txFileName, gbFileName, nil

}
