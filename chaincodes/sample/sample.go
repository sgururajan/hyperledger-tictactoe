package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

func main() {
	err:= shim.Start(new(Sample))
	if err != nil {
		fmt.Println("error staring sample chain code")
		fmt.Printf("%#v", err)
	}
}

type Sample struct {

}

func (m *Sample) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("successfully instantiated sample chaincode"))
}

func (m *Sample) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("successfully invoked sample chaincode"))
}