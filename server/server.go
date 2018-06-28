package main

import (
	"github.com/sgururajan/hyperledger-tictactoe/server/testNetwork"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/sgururajan/hyperledger-tictactoe/server/networkconfig"
	"github.com/sgururajan/hyperledger-tictactoe/server/blockchain"
)

func main() {
	/*testShell()
	os.Exit(0)*/

	defaultNetwork := testNetwork.DefaultNetworkConfiguration()
	networks:= []networkconfig.NetworkConfiguration{defaultNetwork}
	endpointJson, err:= json.MarshalIndent(networks, "", "\t")
	if err!=nil {
		fmt.Printf("unable to serialize default network: %v\n", err)
		os.Exit(100)
	}

	err=ioutil.WriteFile("defaultNetworkConfig.json", endpointJson, os.ModePerm)

	networkConfiguration:= networkconfig.NewFabNetworkConfigurationFromConfig(defaultNetwork)

	fabNetwork:= blockchain.NewFabricNetwork(networkConfiguration)
	err= fabNetwork.Initialize()

	if err!=nil {
		//fmt.Printf("%v", err)
		//os.Exit(101)
		panic(err)
	}

	channelId:= "testchannel"
	//chExists, err:= fabNetwork.IsChannelExists(channelId)
	fabNetwork.CreateChannel(channelId, []string{"sivatech"})

	if err!=nil {
		panic(err)
	}

	//if chExists {
	//	fmt.Printf("channel %s exists", channelId)
	//} else {
	//	fmt.Printf("channel %s does not exists", channelId)
	//}
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

	exePath, _:= os.Executable()
	fmt.Println(exePath)

}


