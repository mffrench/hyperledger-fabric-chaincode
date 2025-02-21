package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func history(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	resultIterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to get asset: %s, error: %s", args[0], err.Error())
	}

	if resultIterator == nil {
		return nil, fmt.Errorf("Asset not found: %s", args[0])
	}

	defer resultIterator.Close()

	batch := make([]string, 0, 8)
	for resultIterator.HasNext() {
		modif, err := resultIterator.Next()
		if err != nil {
			continue
		}
		batch = append(batch, modif.String())
	}

	return []byte(strings.Join(batch, "\n")), nil
}

func get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err.Error())
	}

	// asset not found
	if value == nil {
		return nil, nil
	}

	return value, nil
}

func set(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var asset SimpleAsset
	if err := json.Unmarshal([]byte(args[1]), &asset); err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %s", err.Error())
	}

	asset.TxID = stub.GetTxID()

	b, err := json.Marshal(&asset)
	if err != nil {
		return nil, fmt.Errorf("Marshal failed: %s", err.Error())
	}

	if err := stub.PutState(args[0], b); err != nil {
		return nil, fmt.Errorf("Failed to set asset: %s", args[0])
	}

	return nil, nil
}

func setEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("Incorrect arguments. Expecting an event and its message")
    }

    eventID := args[0]
    message := args[1]

    if err := stub.SetEvent(eventID, []byte(message)); err != nil {
        return nil, fmt.Errorf("Failed to set event: %s", args[0])
    }

    return nil, nil
}
