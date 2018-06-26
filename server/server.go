package main

import (
	"hyperledger/hyperledger-tictactoe/server/testNetwork"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	defaultEndpointConfig:= testNetwork.DefaultNetworkConfiguration()
	endpointJson, err:= json.MarshalIndent(defaultEndpointConfig, "", "\t")
	if err!=nil {
		fmt.Printf("unable to serialize default network: %v\n", err)
		os.Exit(100)
	}

	err=ioutil.WriteFile("defaultNetworkConfig.json", endpointJson, os.ModePerm)

}


