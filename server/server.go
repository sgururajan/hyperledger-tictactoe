package main

import (
	"encoding/json"
	"fmt"
	"github.com/sgururajan/hyperledger-tictactoe/server/blockchain"
	"github.com/sgururajan/hyperledger-tictactoe/server/networkconfig"
	"github.com/sgururajan/hyperledger-tictactoe/server/testNetwork"
	"io/ioutil"
	"os"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"bytes"
	"github.com/sgururajan/hyperledger-tictactoe/server/serverlog"
	"github.com/sgururajan/hyperledger-tictactoe/server/common"
)

var logBuf bytes.Buffer

func main() {
	/*testShell()
	os.Exit(0)*/

	defaultNetwork := testNetwork.DefaultNetworkConfiguration()
	networks := []networkconfig.NetworkConfiguration{defaultNetwork}
	endpointJson, err := json.MarshalIndent(networks, "", "\t")
	if err != nil {
		fmt.Printf("unable to serialize default network: %v\n", err)
		os.Exit(100)
	}

	err = ioutil.WriteFile("defaultNetworkConfig.json", endpointJson, os.ModePerm)

	networkConfiguration := networkconfig.NewFabNetworkConfigurationFromConfig(defaultNetwork)

	setupLogger()
	logger:= logging.NewLogger("Server")


	fabNetwork := blockchain.NewFabricNetwork(networkConfiguration)
	err = fabNetwork.Initialize()

	logger.Println("successfully intialized network")

	if err != nil {
		//fmt.Printf("%v", err)
		//os.Exit(101)
		panic(err)
	}

	var orgId string
	for _, v := range defaultNetwork.Organizations {
		if v.IsOwner {
			orgId = v.OrgID
		}
	}
	chRequest:= common.CreateChannelRequest{
		ChannelName: "testchannel",
		ConsortiumName: "TicTacToeConsortium",
		OrganizationNames: []string{"sivatech"},
		AnchorPeers: map[string][]string{
			"sivatech": {
				"peer0.tictactoe.sivatech.com",
				"peer1.tictactoe.sivatech.com",
			},
		},
	}

	chExists, err:= fabNetwork.IsChannelExists(orgId, chRequest)
	if err!=nil {
		logger.Errorf("error while checking if channel exists: %v", err)
		os.Exit(101)
	}

	if !chExists {
		err = fabNetwork.CreateChannel(orgId, chRequest)
		if err != nil {
			panic(err)
		}
	} else {
		logger.Infof("channel with name %s already exists", chRequest.ChannelName)
	}



	//if chExists {
	//	fmt.Printf("channel %s exists", channelId)
	//} else {
	//	fmt.Printf("channel %s does not exists", channelId)
	//}
}

func setupLogger()  {
	logging.Initialize(serverlog.GetConsoleLogProvider())
	logging.SetLevel("fabsdk/fab", logging.DEBUG)
	logging.SetLevel("fabricNetwork", logging.DEBUG)
}

func testShell() {
	/*args:= []string{"-c", "cd /home/siva && pwd"}
	cmdText:= "/bin/sh"

	cmd:= exec.Command(cmdText, args...)
	res,err:= cmd.CombinedOutput()
	fmt.Println(string(res))
	if err !=nil {
		fmt.Errorf("error: %v", err.Error())
		return
	}*/

	exePath, _ := os.Executable()
	fmt.Println(exePath)

}
