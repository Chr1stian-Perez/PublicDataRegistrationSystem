package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating the DTIC chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting the DTIC chaincode: %s", err.Error())
	}
}
