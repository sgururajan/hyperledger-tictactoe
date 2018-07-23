package fabnetwork

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	fab2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	msp2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

func (m *FabricNetwork) JoinChannel(orgName string, chRequest entities.CreateChannelRequest, orderer entities.Orderer) error {
	creatorOrg, ok := m.orgsByName[orgName]
	if !ok {
		return errors.Errorf("org %s does not belong to this network", orgName)
	}

	resMgmtClient, err := m.getResourceManagementClient(*creatorOrg)
	if err != nil {
		return err
	}

	peersToJoin := []string{}
	for _, v := range chRequest.AnchorPeers[orgName] {
		hasJoined, err := m.hasPeerJoinedChannel(v, chRequest.ChannelName, resMgmtClient)
		if err != nil {
			return err
		}

		if !hasJoined {
			peersToJoin = append(peersToJoin, v)
		}
	}

	err = resMgmtClient.JoinChannel(chRequest.ChannelName,
		resmgmt.WithTargetEndpoints(peersToJoin...),
		resmgmt.WithOrdererEndpoint(orderer.Name),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))

	if err != nil {
		return err
	}

	log.Infof("org %s joined channel successfully", orgName)

	return nil

}

func (m *FabricNetwork) CreateChannel(creatorOrgName string, chRequest entities.CreateChannelRequest) error {

	hasChannel, err := m.HasChannel(chRequest.ChannelName)
	if err != nil {
		return err
	}
	if hasChannel {
		return nil
	}

	creatorOrg, ok := m.orgsByName[creatorOrgName]
	if !ok {
		return errors.Errorf("org %s does not belong to this network", creatorOrgName)
	}

	txFileName := filepath.Join(utils.Substitute(viper.GetString("channelArtifactsPath")), chRequest.ChannelName+".tx")
	txHelper := newFabConfixTxHelper()
	err = txHelper.createChannelTransaction(m, chRequest, txFileName)
	if err != nil {
		return err
	}

	signatures := []msp2.SigningIdentity{}
	anchorPeers := []string{}
	for _, orgName := range chRequest.OrganizationNames {
		mspClient, err := m.getMembershipClient(orgName)
		if err != nil {
			return err
		}
		org := m.orgsByName[orgName]

		adminIdentity, err := m.getSigningIdentity(mspClient, org.AdminUser)
		if err != nil {
			return err
		}
		signatures = append(signatures, adminIdentity)
		for _, ap := range chRequest.AnchorPeers[orgName] {
			anchorPeers = append(anchorPeers, ap)
		}
	}
	//
	//mspClient, err := m.getMembershipClient(creatorOrg.Name)
	//if err != nil {
	//	return err
	//}
	//
	//adminIdentity, err := m.getSigningIdentity(mspClient, creatorOrg.AdminUser)
	//if err != nil {
	//	return err
	//}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         chRequest.ChannelName,
		ChannelConfigPath: txFileName,
		SigningIdentities: signatures,
		//SigningIdentities: []msp2.SigningIdentity{adminIdentity},
	}

	orderers, err := m.providers.ordererProvider.GetOrderers()
	if err != nil {
		return err
	}

	if len(orderers) == 0 {
		return errors.Errorf("no orderer for organization %s", creatorOrgName)
	}

	context := m.sdk.Context(fabsdk.WithUser(creatorOrg.AdminUser), fabsdk.WithOrg(orderers[0].Name))
	resMgmtClient, err := resmgmt.New(context)
	//resMgmtClient, err := m.getResourceManagementClient(*creatorOrg)
	if err != nil {
		return err
	}

	//saveResp, err := resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(orderers[0].Name))
	saveResp, err := resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(orderers[0].Name), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return err
	}

	log.Infof("channel with name %s created successfully. TxID: %v", chRequest.ChannelName, saveResp.TransactionID)

	for _, v := range chRequest.OrganizationNames {
		err = m.JoinChannel(v, chRequest, orderers[0])
	}

	/*err = resMgmtClient.JoinChannel(chRequest.ChannelName,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(orderers[0].Name))

	if err != nil {
		return err
	}

	log.Infof("channel joined successfully")*/

	err = m.updateChannelConfigs()
	if err != nil {
		return err
	}

	return nil
}

func (m *FabricNetwork) GetChannelList() ([]string, error) {
	result := []string{}
	/*orderers, err:= m.providers.ordererProvider.GetOrderers()
	if err != nil {
		return result, err
	}

	orderer:= m.orgsByName[orderers[0].Name]
	resMgmtClient, err:= m.getResourceManagementClient(*orderer)
	if err != nil {
		return result, err
	}

	chResponse, err:= resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(orderer.Name), resmgmt.WithOrdererEndpoint(orderer.Name))
	if err != nil {
		return result, err
	}

	for _, chInfo:= range chResponse.Channels {
		result = append(result, chInfo.ChannelId)
	}*/

	for k, _ := range m.sdkImpl.channels {
		result = append(result, k)
	}

	return result, nil
}

func (m *FabricNetwork) DoesPeerHasChannel(peerName, channelName string) bool {
	chConfig, ok := m.sdkImpl.channels[channelName]
	if !ok {
		return false
	}

	_, ok = chConfig.Peers[peerName]
	if !ok {
		return false
	}

	return true
}

func (m *FabricNetwork) hasPeerJoinedChannel(peerEndpoint, channelName string, resMgmtClient *resmgmt.Client) (bool, error) {
	chQueryResponse, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(peerEndpoint))
	if err != nil {
		return false, err
	}

	for _, chInfo := range chQueryResponse.Channels {
		if chInfo.ChannelId == channelName {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) HasChannel(channelName string) (bool, error) {
	chList, err := m.GetChannelList()
	if err != nil {
		return false, err
	}

	for _, v := range chList {
		if v == channelName {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) getChannelsForOrg(org entities.Organization) (map[string][]entities.Peer, error) {
	resMgmtClient, err := m.getResourceManagementClient(org)
	if err != nil {
		return nil, err
	}
	endrosingPeers, err := m.providers.peerProvider.GetEndrosingPeersForOrgId(org.Name)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]entities.Peer)

	for _, p := range endrosingPeers {
		chResp, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(p.EndPoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return nil, err
		}

		for _, chInfo := range chResp.Channels {
			if _, ok := result[chInfo.ChannelId]; !ok {
				result[chInfo.ChannelId] = []entities.Peer{}
			}
			result[chInfo.ChannelId] = append(result[chInfo.ChannelId], p)
		}
	}

	return result, nil
}

func (m *FabricNetwork) getOrgsForChannel(channelName string) ([]entities.Organization, []string, error) {
	orgsResult := []entities.Organization{}
	peerEndpoints := []string{}
	chPeers, ok := m.sdkImpl.ChannelPeers(channelName)
	if !ok {
		return orgsResult, peerEndpoints, nil
	}

	orgMspIds := make(map[string]bool)
	for _, p := range chPeers {
		if !orgMspIds[p.MSPID] {
			orgMspIds[p.MSPID] = true
		}
	}

	orgs, err := m.providers.orgProvider.GetOrganizations()
	if err != nil {
		return orgsResult, peerEndpoints, err
	}

	for _, o := range orgs {
		if orgMspIds[o.MSPID] {
			orgsResult = append(orgsResult, o)
			oPeers, err := m.providers.peerProvider.GetChainCodePeersForOrgId(o.Name)
			if err != nil {
				return []entities.Organization{}, []string{}, err
			}

			for _, p := range oPeers {
				peerEndpoints = append(peerEndpoints, p.EndPoint)
			}
		}
	}

	return orgsResult, peerEndpoints, nil
}

func (m *FabricNetwork) updateChannelConfigs() error {
	fabChannels := make(map[string][]entities.Peer)
	orgs, err := m.providers.orgProvider.GetOrganizations()
	if err != nil {
		return err
	}

	for _, o := range orgs {
		chResps, err := m.getChannelsForOrg(o)
		if err != nil {
			return err
		}

		for k, v := range chResps {
			if _, exists := fabChannels[k]; !exists {
				fabChannels[k] = []entities.Peer{}
			}

			fabChannels[k] = append(fabChannels[k], v...)
		}
	}

	orderers, err := m.providers.ordererProvider.GetOrderers()
	if err != nil {
		return err
	}
	chOrderers := []string{}
	for _, o := range orderers {
		chOrderers = append(chOrderers, o.Name)
	}

	chConfigs := m.sdkImpl.channels

	for k, v := range fabChannels {
		chConfig := fab2.ChannelEndpointConfig{
			Orderers: chOrderers,
			Policies: fab2.ChannelPolicies{
				QueryChannelConfig: defaultQueryChannelPolicy(),
			},
		}

		//m.sdkImpl.channels[strings.ToLower(k)].Peers = make(map[string]fab2.PeerChannelConfig)
		chPeers := make(map[string]fab2.PeerChannelConfig)
		for _, pv := range v {
			chPeers[pv.EndPoint] = fab2.PeerChannelConfig{
				EventSource:    pv.EventSource,
				LedgerQuery:    pv.CanQueryLedger,
				ChaincodeQuery: pv.IsChainCodePeer,
				EndorsingPeer:  pv.IsEndrosingPeer,
			}
		}

		chConfig.Peers = make(map[string]fab2.PeerChannelConfig)
		for pk, pv := range chPeers {
			chConfig.Peers[pk] = pv
		}
		//m.sdkImpl.channels[strings.ToLower(k)]=chConfig
		chConfigs[strings.ToLower(k)] = chConfig
	}

	m.sdkImpl.channels = chConfigs


	return nil
}
