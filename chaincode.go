package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var (
	logger          = shim.NewLogger("Chaincode")
	unknownFunction = "Unknown function"
)

type SimpleAssetChaincode struct {
}

func (cc *SimpleAssetChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Init chaincode for asset ....")
    args := stub.GetStringArgs()
	return shim.Success([]byte{byte(len(args))})
}

func (cc *SimpleAssetChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
    _, args = args[0], args[1:]
	logger.Infof("received %s(%v)", function, args)

	var err error
	var result []byte

	switch function {
	case "init":
		return cc.Init(stub)
	case "history":
		result, err = history(stub, args)
	case "get":
		result, err = get(stub, args)
	case "set":
		result, err = set(stub, args)
    case "setEvent":
        result, err = setEvent(stub, args)
	default:
		logger.Error(unknownFunction)
		return shim.Error(unknownFunction)
	}

	if err != nil {
		logger.Errorf("[%s] with args (%s) failed: %s", function, strings.Join(args, " | "), err.Error())
		return shim.Error(fmt.Sprintf("[%s] with args (%s) failed: %s", function, strings.Join(args, " | "), err.Error()))
	}

	return shim.Success(result)

}
