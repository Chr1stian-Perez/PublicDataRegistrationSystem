package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ORGANIZATION CONSTANTS (MSP IDs)

const (
	MSP_CIVIL_REGISTRY    = "OrgregistrocivilMSP" // Civil Registry
	MSP_ELECTORAL_COUNCIL = "OrgcneMSP" // Electoral Council
	MSP_COMPTROLLER       = "OrgcontraloriaMSP" // Comptroller
)
// GetMSPID is an auxiliary function to verify the invoker's identity (Basic RBAC)
func GetMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientIdentity := ctx.GetClientIdentity()
	mspId, err := clientIdentity.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to obtain MSP ID: %v", err)
	}
	return mspId, nil
}


// ATTRIBUTE-BASED ACCESS CONTROL (ABAC)
// CheckABAC verifies that the user's X.509 certificate has a specific attribute.
// E.g.: CheckABAC(ctx, "role", "DIRECTOR")
func CheckABAC(ctx contractapi.TransactionContextInterface, attributeName string, expectedValue string) error {
	clientIdentity := ctx.GetClientIdentity()
	
	// 1. Get the attribute value from the certificate
	value, found, err := clientIdentity.GetAttributeValue(attributeName)
	if err != nil {
		return fmt.Errorf("error reading attribute: %v", err)
	}
	if !found {
		return fmt.Errorf("ABAC Error: the certificate does not possess the required '%s' attribute for this operation", attributeName)
	}
	// 2. Compare with the expected value
	if value != expectedValue {
		return fmt.Errorf("ABAC Denied: the attribute '%s' is '%s', but '%s' is required", attributeName, value, expectedValue)
	}

	return nil
}
