package fabnetwork

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"os"
)

type ChainCode struct {
	ChannelName string
	Name        string
	Version     string
}

func (m *FabricNetwork) DoesPeerHasChainCode(peerEndpoint, ccName, ccVersion string, resMgmtClient *resmgmt.Client) bool {
	ccResp,err:= resMgmtClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peerEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err!=nil {
		return false
	}

	for _, cc := range ccResp.Chaincodes {
		if cc.Name==ccName && cc.Version==ccVersion {
			return true
		}
	}

	return false
}

func (m *FabricNetwork) GetChainCodeList(orgName string) ([]ChainCode, error) {
	result := []ChainCode{}
	org, ok := m.orgsByName[orgName]
	if !ok {
		return result, errors.Errorf("org \"%s\" is not part of the network", orgName)
	}
	resMgmtClient, err := m.getResourceManagementClient(*org)
	if err != nil {
		return result, errors.WithMessage(err, "unable to get resource management client")
	}

	chainCodePeers, err := m.providers.peerProvider.GetChainCodePeersForOrgId(orgName)
	if err != nil || len(chainCodePeers) == 0 {
		return result, errors.WithMessage(err, "no chaincode peers for org")
	}
	targetPeer := chainCodePeers[0].EndPoint

	ccResp, err := resMgmtClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(targetPeer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return result, errors.WithMessage(err, "error while querying installed chain code")
	}

	for _, cc := range ccResp.Chaincodes {
		result = append(result, ChainCode{
			Name:        cc.Name,
			ChannelName: "",
			Version:     cc.Version,
		})
	}

	return result, nil
}

func (m *FabricNetwork) HasChainCode(orgName string, ccRequet entities.InstallChainCodeRequest) (bool, error) {
	ccList, err := m.GetChainCodeList(orgName)
	if err != nil {
		return false, err
	}

	for _, v := range ccList {
		if v.Name == ccRequet.ChainCodeName && v.Version == ccRequet.ChainCodeVersion {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) installChainCodeForOrg(orgName string, ccRequest entities.InstallChainCodeRequest, installReq resmgmt.InstallCCRequest) error {
	/*hasChainCode, err:= m.HasChainCode(orgName, ccRequest)
	if err != nil {
		return err
	}

	if hasChainCode {
		return nil
	}*/

	creatorOrg, ok := m.orgsByName[orgName]
	if !ok {
		return errors.Errorf("org \"%s\" is not part of the network", orgName)
	}

	resMgmtClient, err := m.getResourceManagementClient(*creatorOrg)
	if err != nil {
		return errors.WithMessage(err, "unable to obtain resource mamangement client")
	}

	ePeers, err := m.providers.peerProvider.GetChainCodePeersForOrgId(orgName)
	if err != nil {
		return err
	}

	pEndPoints := []string{}
	for _, v := range ePeers {
		if !m.DoesPeerHasChainCode(v.EndPoint, ccRequest.ChainCodeName, ccRequest.ChainCodeVersion, resMgmtClient) {
			pEndPoints = append(pEndPoints, v.EndPoint)
		}
	}

	installResp, err := resMgmtClient.InstallCC(installReq, resmgmt.WithTargetEndpoints(pEndPoints...), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "error when installing chaincode")
	}

	log.Infof("chaincode \"%s\" installed succssfully")
	for _, v := range installResp {
		log.Infof("status: %v, target: %s, status: %s", v.Status, v.Target, v.Info)
	}


	return nil
}

func (m *FabricNetwork) InstallChainCode(orgNames []string, ccRequest entities.InstallChainCodeRequest) error {
	if len(orgNames)==0 {
		return nil
	}

	hasChainCode, err:= m.HasChainCode(orgNames[0], ccRequest)
	if err != nil {
		return err
	}

	if hasChainCode {
		return nil
	}

	ccPackage, err := gopackager.NewCCPackage(ccRequest.ChainCodePath, os.Getenv("GOPATH"))
	if err != nil {
		return errors.WithMessage(err, "error while packaging chaincode")
	}

	installReq := resmgmt.InstallCCRequest{
		Name:    ccRequest.ChainCodeName,
		Path:    ccRequest.ChainCodePath,
		Version: ccRequest.ChainCodeVersion,
		Package: ccPackage,
	}

	for _, o := range orgNames {

		err = m.installChainCodeForOrg(o, ccRequest,installReq)
		if err != nil {
			return err
		}
	}

	creatorOrg, ok := m.orgsByName[orgNames[0]]
	if !ok {
		return errors.Errorf("org \"%s\" is not part of the network", orgNames[0])
	}
	resMgmtClient, err := m.getResourceManagementClient(*creatorOrg)
	ccSignature := []string{creatorOrg.MSPID}
	/*for _, o:= range orgs {
		ccSignature = append(ccSignature, o.MSPID)
	}*/

	//ccPolicy:= cauthdsl.SignedByAnyAdmin([]string{creatorOrg.MSPID})

	ccPolicy := cauthdsl.SignedByAnyAdmin(ccSignature)
	initArg := [][]byte{[]byte("init")}

	initResp, err := resMgmtClient.InstantiateCC(
		ccRequest.ChannelName,
		resmgmt.InstantiateCCRequest{
			Version: ccRequest.ChainCodeVersion,
			Name:    ccRequest.ChainCodeName,
			Path:    ccRequest.ChainCodePath,
			Policy:  ccPolicy,
			Args:    initArg,
		})

	if err != nil {
		return errors.WithMessage(err, "error when instantiating chain code")
	}

	log.Infof("chaincode instantiated and installed successfully. TxID: %s", initResp.TransactionID)

	return nil
}

func (m *FabricNetwork) ExecuteChainCode(orgName, channelName, chainCodeName, cmd string, endorsingPeers []string, args []string) error {
	creatorOrg, ok := m.orgsByName[orgName]
	if !ok {
		return errors.Errorf("org \"%s\" is not part of the network", orgName)
	}

	channelClient, err := m.getChannelClient(creatorOrg, channelName)
	if err != nil {
		return err
	}

	chArgs := [][]byte{[]byte(cmd)}
	chArgs = append(chArgs, convertToChannelArgs(args)...)
	resp, err := channelClient.Execute(channel.Request{
		ChaincodeID: chainCodeName,
		Args:        chArgs,
		Fcn:         "invoke",
	}, channel.WithTargetEndpoints(endorsingPeers...), channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return err
	}

	allGood := areChannelResponsesGood(resp.Responses)
	if !allGood {
		return errors.Errorf("warning: did not get positive response from all the peers")
	}

	return nil
}

func (m *FabricNetwork) QueryChainCode(orgName, channelName, chaincodeName, cmd string, targetPeers []string, args []string) ([]byte, error) {
	creatorOrg, ok := m.orgsByName[orgName]
	if !ok {
		return nil, errors.Errorf("org \"%s\" is not part of the network", orgName)
	}

	channelClient, err := m.getChannelClient(creatorOrg, channelName)
	if err != nil {
		return nil, err
	}

	chArgs := [][]byte{[]byte(cmd)}
	chArgs = append(chArgs, convertToChannelArgs(args)...)
	resp, err := channelClient.Query(channel.Request{
		Fcn:         "invoke",
		Args:        chArgs,
		ChaincodeID: chaincodeName,
	}, channel.WithTargetEndpoints(targetPeers...), channel.WithRetry(retry.DefaultChannelOpts))

	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func areChannelResponsesGood(responses []*fab.TransactionProposalResponse) bool {
	allGood := true
	for _, p := range responses {
		oneGood := false
		if p.Response.Status == 200 {
			oneGood = true
		}
		allGood = allGood && oneGood
	}

	return allGood
}

func convertToChannelArgs(args []string) [][]byte {
	result := [][]byte{}
	for _, v := range args {
		result = append(result, []byte(v))
	}
	return result
}

func (m *FabricNetwork) getChannelClient(creatorOrg *entities.Organization, channelName string) (*channel.Client, error) {
	context := m.sdk.ChannelContext(channelName, fabsdk.WithUser(creatorOrg.AdminUser), fabsdk.WithOrg(creatorOrg.Name))
	channelClient, err := channel.New(context)
	if err != nil {
		return nil, errors.WithMessage(err, "error when creating channel client")
	}

	//ledgerClient,_:= ledger.New(context)
	//event,_:=event.New(context, event.WithSeekType(seek.Newest))

	return channelClient, nil

}
