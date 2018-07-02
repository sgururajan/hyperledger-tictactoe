package fabnetwork

import (
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"os"
	"github.com/cloudflare/cfssl/log"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

type ChainCode struct {
	ChannelName string
	Name string
	Version string
}

func (m *FabricNetwork) GetChainCodeList(orgName string) ([]ChainCode, error) {
	result:= []ChainCode{}
	org, ok:= m.orgsByName[orgName]
	if !ok {
		return result, errors.Errorf("org \"%s\" is not part of the network", orgName)
	}
	resMgmtClient, err:= m.getResourceManagementClient(*org)
	if err != nil {
		return result, errors.WithMessage(err, "unable to get resource management client")
	}

	chainCodePeers, err:= m.providers.peerProvider.GetChainCodePeersForOrgId(orgName)
	if err != nil || len(chainCodePeers)==0 {
		return result, errors.WithMessage(err, "no chaincode peers for org")
	}
	targetPeer:= chainCodePeers[0].EndPoint

	ccResp, err:= resMgmtClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(targetPeer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return result, errors.WithMessage(err, "error while querying installed chain code")
	}

	for _, cc:= range ccResp.Chaincodes {
		result = append(result, ChainCode{
			Name:cc.Name,
			ChannelName: "",
			Version: cc.Version,
		})
	}

	return result, nil
}

func (m *FabricNetwork) HasChainCode(orgName string, ccRequet entities.InstallChainCodeRequest) (bool, error) {
	ccList, err:= m.GetChainCodeList(orgName)
	if err != nil {
		return false, err
	}

	for _, v:= range ccList {
		if v.Name==ccRequet.ChainCodeName && v.Version==ccRequet.ChainCodeVersion {
			return true, nil
		}
	}

	return false, nil
}

func (m *FabricNetwork) InstallChainCode(orgName string, ccRequest entities.InstallChainCodeRequest) error {
	hasChainCode, err:= m.HasChainCode(orgName, ccRequest)
	if err != nil {
		return err
	}

	if hasChainCode {
		return nil
	}

	ccPackage, err:= gopackager.NewCCPackage(ccRequest.ChainCodePath, os.Getenv("GOPATH"))
	if err != nil {
		return errors.WithMessage(err, "error while packaging chaincode")
	}

	installReq:= resmgmt.InstallCCRequest{
		Name:    ccRequest.ChainCodeName,
		Path:    ccRequest.ChainCodePath,
		Version: ccRequest.ChainCodeVersion,
		Package: ccPackage,
	}

	creatorOrg, ok:= m.orgsByName[orgName]
	if !ok {
		return errors.Errorf("org \"%s\" is not part of the network", orgName)
	}

	resMgmtClient, err:= m.getResourceManagementClient(*creatorOrg)
	if err != nil {
		return errors.WithMessage(err, "unable to obtain resource mamangement client")
	}

	installResp, err:= resMgmtClient.InstallCC(installReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "error when installing chaincode")
	}

	log.Infof("chaincode \"%s\" installed succssfully")
	for _, v:= range installResp {
		log.Infof("status: %v, target: %s, status: %s", v.Status, v.Target, v.Info)
	}

	ccPolicy:= cauthdsl.SignedByAnyAdmin([]string{creatorOrg.MSPID})
	initArg:= [][]byte{[]byte("init")}

	initResp, err:= resMgmtClient.InstantiateCC(
		ccRequest.ChannelName,
		resmgmt.InstantiateCCRequest{
			Version: ccRequest.ChainCodeVersion,
			Name: ccRequest.ChainCodeName,
			Path: ccRequest.ChainCodePath,
			Policy: ccPolicy,
			Args: initArg,
		})

	if err != nil {
		return errors.WithMessage(err, "error when instantiating chain code")
	}

	log.Infof("chaincode instantiated and installed successfully. TxID: %s", initResp.TransactionID)

	return nil
}
