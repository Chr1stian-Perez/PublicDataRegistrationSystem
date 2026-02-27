package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 1. IDENTITY AND BASE EVIDENCE MANAGEMENT

// Tx_RegisterIdentity: Creates the asset with initial IPFS structure
func (s *SmartContract) Tx_RegisterIdentity(ctx contractapi.TransactionContextInterface, nationalID string, firstNames string, lastNames string, birthDate string, birthPlace string, sex string, gender string, initialCivilRegistryDirCID string, initialRootCID string) error {
	
	msp, _ := GetMSPID(ctx)
	if msp != MSP_CIVIL_REGISTRY {
		return fmt.Errorf("permission denied: only Civil Registry can create identities")
	}

	// ABAC Validation: Only OPERATOR can create identities
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	exists, _ := ctx.GetStub().GetState(nationalID)
	if exists != nil {
		return fmt.Errorf("citizen %s already exists in the world state", nationalID)
	}

	newCitizen := Citizen{
		GlobalID:     nationalID,
		GlobalStatus: "ACTIVE",
		CivilRegistryData: CivilRegistryInfo{
			FirstNames:    firstNames,
			LastNames:     lastNames,
			Sex:           sex,
			Gender:        gender,
			BirthDate:     birthDate,
			BirthPlace:    birthPlace,
			Nationality:   "ECUADORIAN",
			MaritalStatus: "SINGLE",
		},
		ComptrollerData: ComptrollerInfo{
			ComptrollerStatus:  "CLEAN",
			LastNotification:   "",
		},
		// Initialize the hierarchical structure
		DigitalEvidence: IPFSStructure{
			RootCID:             initialRootCID,             // Calculated by the client
			CivilRegistryDirCID: initialCivilRegistryDirCID, // Contains the Birth Certificate
			ElectoralCouncilDirCID:  "",                         // Empty at birth
			ComptrollerDirCID:   "",                         // Empty at birth
		},
		AuditProcesses: []AuditProcess{},
	}

	assetJSON, _ := json.Marshal(newCitizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// Tx_UpdateEvidence: Updates pointers to IPFS directories
func (s *SmartContract) Tx_UpdateEvidence(ctx contractapi.TransactionContextInterface, nationalID string, newRootCID string, newOrgDirCID string) error {
	msp, _ := GetMSPID(ctx)

	// ABAC Validation: Only OPERATOR can upload new evidence
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil { return fmt.Errorf("citizen %s not found in the world state", nationalID) }
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)
	// Validate competency: Each Org only touches ITS directory
	if msp == MSP_CIVIL_REGISTRY {
		citizen.DigitalEvidence.CivilRegistryDirCID = newOrgDirCID
	} else if msp == MSP_ELECTORAL_COUNCIL {
		citizen.DigitalEvidence.ElectoralCouncilDirCID = newOrgDirCID
	} else if msp == MSP_COMPTROLLER {
		citizen.DigitalEvidence.ComptrollerDirCID = newOrgDirCID
	} else {
		return fmt.Errorf("organization not authorized to update evidence")
	}
	// The Root is always updated because it changes if any child changes
	citizen.DigitalEvidence.RootCID = newRootCID

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 2. ELECTORAL MANAGEMENT (CNE) - ELECTORAL ROLL LOGIC

func (s *SmartContract) Tx_BatchElectoralManagement(ctx contractapi.TransactionContextInterface, province string, canton string, parish string, votingPrecinct string, votingTable string, sharedElectoralCouncilDirCID string, updatesJSON string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_ELECTORAL_COUNCIL { return fmt.Errorf("permission denied: only Electoral Council") }

	// ABAC Validation: Only OPERATOR can update electoral roll
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	var updates []ElectoralUpdateInput
	err = json.Unmarshal([]byte(updatesJSON), &updates)
	if err != nil {
		return fmt.Errorf("failed to unmarshal updates JSON: %v", err)
	}

	for _, update := range updates {
		assetBytes, err := ctx.GetStub().GetState(update.NationalID)
		if err != nil { return err }
		if assetBytes == nil {
			continue // Skip if the citizen is not in the world state
		}

		var citizen Citizen
		json.Unmarshal(assetBytes, &citizen)

		if citizen.GlobalStatus == "DECEASED" {
			continue // Skip deceased citizens
		}

		// Update data
		citizen.ElectoralCouncilData.Province = province
		citizen.ElectoralCouncilData.Canton = canton
		citizen.ElectoralCouncilData.Parish = parish
		citizen.ElectoralCouncilData.VotingPrecinct = votingPrecinct
		citizen.ElectoralCouncilData.VotingTable = votingTable
		citizen.ElectoralCouncilData.HasVoted = update.HasVoted
		citizen.ElectoralCouncilData.HasFines = update.HasFines

		// Update pointer to Electoral Council Directory (Electoral Rolls shared CSV)
		if update.NewRootCID != "" && sharedElectoralCouncilDirCID != "" {
			citizen.DigitalEvidence.ElectoralCouncilDirCID = sharedElectoralCouncilDirCID
			citizen.DigitalEvidence.RootCID = update.NewRootCID
		}

		assetJSON, _ := json.Marshal(citizen)
		err = ctx.GetStub().PutState(update.NationalID, assetJSON)
		if err != nil { return err }
	}

	return nil
}

// Tx_PayElectoralFine: Processes payment of an electoral fine and updates the citizen's status
func (s *SmartContract) Tx_PayElectoralFine(ctx contractapi.TransactionContextInterface, nationalID string, newRootCID string, paymentReceiptCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_ELECTORAL_COUNCIL { return fmt.Errorf("permission denied: only Electoral Council") }

	// ABAC Validation: Only OPERATOR can process fine payments
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil { return err }
	if assetBytes == nil { return fmt.Errorf("citizen %s not found in the world state", nationalID) }

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if !citizen.ElectoralCouncilData.HasFines {
		return fmt.Errorf("citizen does not have any pending electoral fines")
	}

	// Remove fine
	citizen.ElectoralCouncilData.HasFines = false

	// Update pointer to Electoral Council Directory (Fine Payment Receipt)
	if newRootCID != "" && paymentReceiptCID != "" {
		citizen.DigitalEvidence.ElectoralCouncilDirCID = paymentReceiptCID
		citizen.DigitalEvidence.RootCID = newRootCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 3. AUDIT (COMPTROLLER) - STATE MACHINE + ABAC

func (s *SmartContract) Tx_EvolveAudit(ctx contractapi.TransactionContextInterface, nationalID string, caseID string, newPhase string, defenseDescription string, defenseEvidenceCID string, newRootCID string, newComptrollerDirCID string) error {
	msp, _ := GetMSPID(ctx)
	// 1. Organizational check
	if msp != MSP_COMPTROLLER {
		return fmt.Errorf("permission denied: only Comptroller manages audits")
	}
	// 2. ABAC CHECK 
	// Only auditors with role 'SENIOR_AUDITOR' can sign these transactions
	err := CheckABAC(ctx, "role", "SENIOR_AUDITOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil { return fmt.Errorf("citizen %s not found in the world state", nationalID) }
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)
	// Find or Create case
	var caseIndex = -1
	for i, process := range citizen.AuditProcesses {
		if process.CaseID == caseID {
			caseIndex = i
			break
		}
	}
	if caseIndex == -1 {
		newCase := AuditProcess{
			CaseID:         caseID,
			CurrentStatus:  "START",
			DefenseHistory: []DefenseEvidence{},
			Conclusion:     "PENDING",
		}
		citizen.AuditProcesses = append(citizen.AuditProcesses, newCase)
		caseIndex = len(citizen.AuditProcesses) - 1
	}
	// Create Defense evidence
	newDefenseEvidence := DefenseEvidence{
		Phase:        newPhase,
		Description:  defenseDescription,
		EvidenceCID:  defenseEvidenceCID,
		Date:         time.Now().Format(time.RFC3339),
		Author:       msp,
	}
	// Update Process status
	citizen.AuditProcesses[caseIndex].CurrentStatus = newPhase
	citizen.AuditProcesses[caseIndex].DefenseHistory = append(citizen.AuditProcesses[caseIndex].DefenseHistory, newDefenseEvidence)
	// AUTOMATIC UPDATE OF CITIZEN STATUS (DIAGRAM LOGIC)
	switch newPhase {
	case "DRAFT":
		citizen.ComptrollerData.ComptrollerStatus = "NOTIFIED"
	case "FINAL_CLEAN_REPORT":
		citizen.ComptrollerData.ComptrollerStatus = "CLEAN"
	case "PREDETERMINATION":
		citizen.ComptrollerData.ComptrollerStatus = "PREDETERMINED"
	case "CIVIL_DETERMINATION":
		citizen.ComptrollerData.ComptrollerStatus = "CIVIL_LIABILITY"
	case "ADMINISTRATIVE_DETERMINATION":
		citizen.ComptrollerData.ComptrollerStatus = "ADMINISTRATIVE_LIABILITY"
	case "CRIMINAL_DETERMINATION":
		citizen.ComptrollerData.ComptrollerStatus = "CRIMINAL_LIABILITY"
	case "APPEAL":
		citizen.ComptrollerData.ComptrollerStatus = "UNDER_APPEAL"
	case "SANCTIONING_SENTENCE":
		citizen.ComptrollerData.ComptrollerStatus = "SANCTIONED"
	case "ACQUITTAL_SENTENCE":
		citizen.ComptrollerData.ComptrollerStatus = "CLEAN"
	}
	citizen.ComptrollerData.LastNotification = time.Now().Format(time.RFC3339)
	// Update IPFS pointers
	if newRootCID != "" && newComptrollerDirCID != "" {
		citizen.DigitalEvidence.ComptrollerDirCID = newComptrollerDirCID
		citizen.DigitalEvidence.RootCID = newRootCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// Tx_RegistryRectification: Now protected with ABAC
func (s *SmartContract) Tx_RegistryRectification(ctx contractapi.TransactionContextInterface, nationalID string, field string, correctedValue string, resolutionCID string, newRootCID string, newCivilRegistryDirCID string) error {
	msp, _ := GetMSPID(ctx)
	
	// Validate ABAC: Only Directors can rectify
	err := CheckABAC(ctx, "role", "DIRECTOR")
	if err != nil { return fmt.Errorf("ABAC Failure: %v", err) }

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil { return fmt.Errorf("citizen %s not found in the world state", nationalID) }
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	// Rectification logic
	if msp == MSP_CIVIL_REGISTRY {
		switch field {
		case "FirstNames":
			citizen.CivilRegistryData.FirstNames = correctedValue
		case "LastNames":
			citizen.CivilRegistryData.LastNames = correctedValue
		case "Sex":
			citizen.CivilRegistryData.Sex = correctedValue
		case "Gender":
			citizen.CivilRegistryData.Gender = correctedValue
		case "Nationality":
			citizen.CivilRegistryData.Nationality = correctedValue
		case "BirthDate":
			citizen.CivilRegistryData.BirthDate = correctedValue
		case "BirthPlace":
			citizen.CivilRegistryData.BirthPlace = correctedValue
		case "MaritalStatus":
			citizen.CivilRegistryData.MaritalStatus = correctedValue
		default:
			return fmt.Errorf("field %s is not valid for Civil Registry rectification", field)
		}
	} else if msp == MSP_ELECTORAL_COUNCIL {
		switch field {
		case "Province":
			citizen.ElectoralCouncilData.Province = correctedValue
		case "Canton":
			citizen.ElectoralCouncilData.Canton = correctedValue
		case "Parish":
			citizen.ElectoralCouncilData.Parish = correctedValue
		case "VotingPrecinct":
			citizen.ElectoralCouncilData.VotingPrecinct = correctedValue
		case "VotingTable":
			citizen.ElectoralCouncilData.VotingTable = correctedValue
		default:
			return fmt.Errorf("field %s is not valid for Electoral Council rectification", field)
		}
	} else {
		return fmt.Errorf("organization without competency")
	}

	// Update IPFS (Rectification resolution)
	if newRootCID != "" {
		citizen.DigitalEvidence.RootCID = newRootCID
		if msp == MSP_CIVIL_REGISTRY { citizen.DigitalEvidence.CivilRegistryDirCID = newCivilRegistryDirCID }
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 4. QUERIES

func (s *SmartContract) Query_PublicProfile(ctx contractapi.TransactionContextInterface, nationalID string) (*Citizen, error) {
	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil { return nil, err }
	if assetBytes == nil { return nil, fmt.Errorf("citizen %s not found in the world state", nationalID) }

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	// LOPDP PROTECTION: Hide sensitive details
	citizen.DigitalEvidence = IPFSStructure{} 
	citizen.AuditProcesses = nil 

	return &citizen, nil
}

func (s *SmartContract) Query_ComprehensiveAudit(ctx contractapi.TransactionContextInterface, nationalID string) (*Citizen, error) {
	msp, _ := GetMSPID(ctx)
	
	// Both Comptroller and Civil Registry have access, but check ABAC
	if msp == MSP_COMPTROLLER {
		err := CheckABAC(ctx, "role", "AUDITOR")
		if err != nil { 
			// Check if it's SENIOR_AUDITOR instead
			errSenior := CheckABAC(ctx, "role", "SENIOR_AUDITOR")
			if errSenior != nil { return nil, fmt.Errorf("ABAC Failure: only AUDITOR or SENIOR_AUDITOR can query") }
		}
	} else if msp != MSP_CIVIL_REGISTRY {
		return nil, fmt.Errorf("access denied")
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil { return nil, err }
	if assetBytes == nil { return nil, fmt.Errorf("citizen %s not found in the world state", nationalID) }
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)
	return &citizen, nil
}

// 5. TRACEABILITY (AUDIT)

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	TxId      string    `json:"tx_id"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"is_delete"`
	Value     *Citizen  `json:"value"`
}

// Query_History: Returns the chain of historical updates for a citizen 
func (s *SmartContract) Query_History(ctx contractapi.TransactionContextInterface, nationalID string) ([]HistoryQueryResult, error) {
	msp, _ := GetMSPID(ctx)
	
	// ABAC Validation: Only AUDITOR or SENIOR_AUDITOR can trace history
	if msp != MSP_COMPTROLLER {
		return nil, fmt.Errorf("permission denied: only Comptroller can query history")
	}

	err := CheckABAC(ctx, "role", "AUDITOR")
	if err != nil { 
		// If not AUDITOR, check if SENIOR_AUDITOR
		errSenior := CheckABAC(ctx, "role", "SENIOR_AUDITOR")
		if errSenior != nil { return nil, fmt.Errorf("ABAC Failure: %v", errSenior) }
	}

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(nationalID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var citizen Citizen
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &citizen)
			if err != nil {
				return nil, err
			}
		}

		// FECHA CORREGIDA
		timestamp := response.Timestamp.AsTime()

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			IsDelete:  response.IsDelete,
			Value:     &citizen,
		}
		records = append(records, record)
	}

	return records, nil
}
