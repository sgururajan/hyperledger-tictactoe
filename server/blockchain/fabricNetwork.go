package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/server/blockchainconfig"
	"github.com/sgururajan/hyperledger-tictactoe/server/common"
	"github.com/sgururajan/hyperledger-tictactoe/server/networkconfig"
	"github.com/sgururajan/hyperledger-tictactoe/server/pathUtil"
	"path/filepath"
	"os"
)

var logger = logging.NewLogger("fabricNetwork")

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

	sdkOptions := []fabsdk.Option{
		fabsdk.WithEndpointConfig(m.networkConfig),
		fabsdk.WithCryptoSuiteConfig(m.networkConfig.SecurityConfiguration),
		fabsdk.WithIdentityConfig(m.networkConfig),
	}

	sdk, err := fabsdk.New(nil, sdkOptions...)

	if err != nil {
		return errors.WithMessage(err, "unable to intialize sdk")
	}

	m.sdk = sdk

	return nil
}

func (m *FabricNetwork) IsChannelExists(orgId string, chRequest common.CreateChannelRequest) (bool, error) {

	creatorOrg, ok := m.networkConfig.OrgByOrgId[orgId]
	if !ok {
		return false, errors.Errorf(`orgId %s is not part of the network`, orgId)
	}

	clientContext := m.sdk.Context(fabsdk.WithOrg(creatorOrg.Name), fabsdk.WithUser(creatorOrg.AdminUserName))
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Errorf("unable to obtain resource management client: %v", err)
		return false, err
	}

	peers, ok := m.networkConfig.PeerByOrg[creatorOrg.Name]
	if !ok {
		logger.Errorf("org %s does not have any peers", orgId)
		return false, errors.Errorf("org %s does not have any peers", orgId)
	}

	var targetPeer string
	for _, v := range peers {
		if v.IsEndrosingPeer {
			targetPeer = v.Endpoint
		}
	}

	chResp, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(targetPeer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Errorf("error while querying channels: %v", err)
		return false, err
	}

	for _, cInfo := range chResp.Channels {
		if cInfo.ChannelId == chRequest.ChannelName {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) CreateChannel(orgId string, chRequest common.CreateChannelRequest) error {

	creatorOrg, ok := m.networkConfig.OrgByOrgId[orgId]
	if !ok {
		return errors.Errorf("orgID %s does not belong to this network")
	}
	txFileName, err := m.createChannelTx(orgId, chRequest)
	if err != nil {
		return err
	}

	mspClient, err := msp.New(m.sdk.Context(), msp.WithOrg(creatorOrg.Name))
	if err != nil {
		return errors.WithMessage(err, "unable to obtain msp client")
	}

	adminIdentity, err := mspClient.GetSigningIdentity(creatorOrg.AdminUserName)
	if err != nil {
		return errors.WithMessage(err, "unable to obtain admin identity")
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         chRequest.ChannelName,
		ChannelConfigPath: txFileName,
		SigningIdentities: []mspctx.SigningIdentity{adminIdentity},
	}

	clientContext := m.sdk.Context(fabsdk.WithUser(creatorOrg.AdminUserName), fabsdk.WithOrg(creatorOrg.Name))
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "unable to obtain resource nangement client")
	}

	saveRes, err := resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(m.networkConfig.OrderersInfo[0].Endpoint))

	if err != nil {
		logger.Error("save channel failed")
		return errors.WithMessage(err, "error while saving channel")
	}

	logger.Infof("channel with channel id %v created successfully. TxID: %v", chRequest.ChannelName, saveRes.TransactionID)

	// now join the channel
	err = resMgmtClient.JoinChannel(chRequest.ChannelName, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(m.networkConfig.OrderersInfo[0].Endpoint))

	logger.Infof("successfully added orgs to channel %s", chRequest.ChannelName)

	os.Remove(txFileName)

	return nil
}

func (m *FabricNetwork) createChannelTx(orgId string, chRequest common.CreateChannelRequest) (string, error) {

	channelTxHelper := blockchainconfig.NewChannelTxConfigHelper()
	err := channelTxHelper.CreateChannelTxObject(m.networkConfig, chRequest)
	if err != nil {
		return "", err
	}

	exePath := pathUtil.GetExecutablePath()
	txFileName := filepath.Join(exePath, "channel-artifacts", chRequest.ChannelName+".pb")
	//gbFileName:= filepath.Join(exePath, "channel-artifacts", channelName+"_genesis.block")

	pathUtil.EnsureDirectory(exePath, "channel-artifacts")

	err = channelTxHelper.CreateConfigurationBlocks(txFileName)
	if err != nil {
		logger.Error("error while trying to create channel configuration block")
		return "", err
	}

	return txFileName, nil

}
