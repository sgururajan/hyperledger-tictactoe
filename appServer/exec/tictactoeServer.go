package main

import (
	"github.com/sgururajan/hyperledger-tictactoe/appServer"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"github.com/cloudflare/cfssl/log"
	"github.com/spf13/viper"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
)

var networksDbFile="networks.json"

func main() {

	viper.SetConfigFile("appsettings.json")
	verr:= viper.ReadInConfig()

	if verr != nil {
		log.Fatalf("error while reading config")
		panic(verr)
	}

	log.Infof("appsetting - configTxGenToolsPath: %s", viper.GetString("configTxGenToolsPath"))

	err:= setupDefaultNetwork()
	if err!=nil {
		panic(err)
	}
	repo:= setupRepository()
	networkHandler,err:= networkHandlers.NewNetworkHandler(repo)
	if err != nil {
		panic(err)
	}

	defer networkHandler.Close()

	testNetwork(networkHandler)

}

func setupRepository() database.NetworkRepository {
	return database.NewNetworkFileRepository(networksDbFile)
}

func setupDefaultNetwork() error {
	if _,err:= os.Stat(networksDbFile); err!= nil {
		networks:= getDefaultNetwork()

		cbytes, err:= json.MarshalIndent(networks, "", "\t")
		if err!=nil {
			return err
		}

		err = ioutil.WriteFile(networksDbFile, cbytes, os.ModePerm)
		if err !=nil {
			return err
		}
	}

	return nil
}

func getDefaultNetwork() map[string]database.Network {
	dNetwork:= appServer.DefaultNetworkConfiguration()
	networks:= make(map[string]database.Network)
	networks[dNetwork.Name]=dNetwork

	return networks
}

func testNetwork(networkHandler *networkHandlers.NetworkHandler) {
	network, err:= networkHandler.GetNetwork("testNetwork")
	if err != nil {
		log.Errorf("error while getting network. err: %v", err)
		os.Exit(0)
	}

	chReq:=entities.CreateChannelRequest{
		ChannelName: "testchannel1",
		OrganizationNames: []string{
			"sivatech",
		},
		AnchorPeers: map[string][]string {
			"sivatech": []string{
				"peer0.tictactoe.sivatech.com",
			},
		},
		ConsortiumName: "TicTacToeConsortium",
	}

	err = network.CreateChannel("sivatech", chReq)

	if err != nil {
		panic(err)
	}

	ccRequest:= entities.InstallChainCodeRequest{
		ChainCodeName:"sample",
		ChainCodePath:"github.com/sgururajan/hyperledger-tictactoe/chaincodes/sample/",
		ChainCodeVersion: "0.0.4",
		ChannelName: "testchannel",
	}

	err = network.InstallChainCode("sivatech", ccRequest)

	if err != nil {
		panic(err)
	}
}
