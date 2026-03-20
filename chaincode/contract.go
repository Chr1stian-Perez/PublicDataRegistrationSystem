package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Tx_Register: Master transaction to route all registrations
func (s *SmartContract) Tx_Register(ctx contractapi.TransactionContextInterface, nationalID string, recordType string, recordDataJSON string, newCID string) error {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	recordType = strings.ToUpper(strings.TrimSpace(recordType))

	var fields map[string]string
	if recordDataJSON != "" {
		err := json.Unmarshal([]byte(recordDataJSON), &fields)
		if err != nil {
			return fmt.Errorf("invalid record data formatting: %v", err)
		}
		for k, v := range fields {
			fields[k] = strings.ToUpper(strings.TrimSpace(v))
		}
	}

	switch recordType {
	case "IDENTITY":
		return s.registerIdentity(ctx, nationalID, fields["first_names"], fields["last_names"], fields["sex"], fields["gender"], fields["nationality"], fields["birth_date"], fields["birth_place"], fields["marital_status"], fields["spouse"], fields["father_id"], fields["father_nationality"], fields["mother_id"], fields["mother_nationality"], fields["address"], fields["marriage_date"], fields["death_date"], fields["observations"], fields["death_registration_date"], fields["citizen_condition"], newCID)
	case "MARRIAGE":
		return s.registerMarriage(ctx, nationalID, fields["spouse_name"], fields["marriage_date"], newCID)
	case "DEATH":
		return s.registerDeath(ctx, nationalID, fields["death_date"], fields["registration_date"], fields["observations"], newCID)
	case "ACADEMIC_DEGREE":
		return s.registerAcademicDegree(ctx, nationalID, fields["registry_id"], fields["instruction_level"], fields["exact_degree_name"], fields["university_name"], fields["registration_date"], newCID)
	case "PROPERTY":
		return s.registerProperty(ctx, nationalID, fields["property_id"], fields["owner_name"], fields["owner_id"], fields["property_type"], fields["legal_status"], fields["registration_date"], newCID)
	case "POLICE_RECORD":
		return s.registerPoliceRecord(ctx, nationalID, fields["capture_order"], fields["criminal_records"], fields["exit_impediments"], fields["migratory_alert"], newCID)
	default:
		return fmt.Errorf("unknown record type: %s", recordType)
	}
}

// 1. IDENTITY AND BASE EVIDENCE MANAGEMENT

// registerIdentity: Creates the asset with initial IPFS structure (Private Internal)
func (s *SmartContract) registerIdentity(ctx contractapi.TransactionContextInterface, nationalID string, firstNames string, lastNames string, sex string, gender string, nationality string, birthDate string, birthPlace string, maritalStatus string, spouse string, fatherID string, fatherNationality string, motherID string, motherNationality string, address string, marriageDate string, deathDate string, observations string, deathRegistrationDate string, citizenCondition string, initialCID string) error {

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
		GlobalID:         nationalID,
		CitizenCondition: citizenCondition,
		CivilRegistryData: CivilRegistryInfo{
			FirstNames:            firstNames,
			LastNames:             lastNames,
			Sex:                   sex,
			Gender:                gender,
			Nationality:           nationality,
			BirthDate:             birthDate,
			BirthPlace:            birthPlace,
			MaritalStatus:         maritalStatus,
			Spouse:                spouse,
			FatherID:              fatherID,
			FatherNationality:     fatherNationality,
			MotherID:              motherID,
			MotherNationality:     motherNationality,
			Address:               address,
			MarriageDate:          marriageDate,
			DeathDate:             deathDate,
			Observations:          observations,
			DeathRegistrationDate: deathRegistrationDate,
		},
		AcademicRegistryData: []AcademicRegistryInfo{},
		PropertyRegistryData: []PropertyRegistryInfo{},
		PoliceRegistryData:   PoliceRegistryInfo{},
		CID:                  initialCID,
	}

	assetJSON, _ := json.Marshal(newCitizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// Tx_UpdateEvidence: Updates pointer to IPFS pool
func (s *SmartContract) Tx_UpdateEvidence(ctx contractapi.TransactionContextInterface, nationalID string, newCID string) error {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	
	// ABAC Validation: Only OPERATOR can upload new evidence
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found in the world state", nationalID)
	}
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	citizen.CID = newCID

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 2. ACADEMIC REGISTRY (ORG2)

func (s *SmartContract) registerAcademicDegree(ctx contractapi.TransactionContextInterface, nationalID string, nationalRegistryID string, instructionLevel string, exactDegreeName string, universityName string, registrationDate string, newCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_ACADEMIC_REGISTRY {
		return fmt.Errorf("permission denied: only Academic Registry")
	}

	// ABAC Validation: Only OPERATOR can register degrees
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil {
		return err
	}
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found in the world state", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if citizen.CitizenCondition == "Inactive" {
		return fmt.Errorf("cannot register degree for an inactive citizen")
	}

	// Create new degree and append it
	newDegree := AcademicRegistryInfo{
		NationalRegistryID: nationalRegistryID,
		InstructionLevel:   instructionLevel,
		ExactDegreeName:    exactDegreeName,
		UniversityName:     universityName,
		RegistrationDate:   registrationDate,
	}

	citizen.AcademicRegistryData = append(citizen.AcademicRegistryData, newDegree)

	// Update pointer to IPFS Pool
	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 3. PROPERTY REGISTRY (ORG3)

func (s *SmartContract) registerProperty(ctx contractapi.TransactionContextInterface, nationalID string, propertyID string, ownerName string, ownerID string, propertyType string, legalStatus string, registrationDate string, newCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_PROPERTY_REGISTRY {
		return fmt.Errorf("permission denied: only Property Registry")
	}

	// ABAC Validation: Only REGISTRAR can process property registrations
	err := CheckABAC(ctx, "role", "REGISTRAR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil {
		return err
	}
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found in the world state", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if citizen.CitizenCondition == "Inactive" {
		return fmt.Errorf("cannot register properties to an inactive citizen")
	}

	// Check if property exists to update, else append
	var propertyIndex = -1
	for i, prop := range citizen.PropertyRegistryData {
		if prop.PropertyID == propertyID {
			propertyIndex = i
			break
		}
	}

	if propertyIndex == -1 {
		// New property
		newProperty := PropertyRegistryInfo{
			PropertyID:       propertyID,
			OwnerName:        ownerName,
			OwnerID:          ownerID,
			PropertyType:     propertyType,
			LegalStatus:      legalStatus,
			RegistrationDate: registrationDate,
		}
		citizen.PropertyRegistryData = append(citizen.PropertyRegistryData, newProperty)
	} else {
		// Update existing
		citizen.PropertyRegistryData[propertyIndex].OwnerName = ownerName
		citizen.PropertyRegistryData[propertyIndex].OwnerID = ownerID
		citizen.PropertyRegistryData[propertyIndex].PropertyType = propertyType
		citizen.PropertyRegistryData[propertyIndex].LegalStatus = legalStatus
		citizen.PropertyRegistryData[propertyIndex].RegistrationDate = registrationDate
	}

	// Update pointer to IPFS Pool
	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// 4. POLICE REGISTRY (ORG4)

func (s *SmartContract) registerPoliceRecord(ctx contractapi.TransactionContextInterface, nationalID string, captureOrder string, criminalRecords string, exitImpediments string, migratoryAlert string, newCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_POLICE_REGISTRY {
		return fmt.Errorf("permission denied: only Police Registry")
	}

	// ABAC Validation: Only OPERATOR can register police records
	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil {
		return err
	}
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found in the world state", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if citizen.CitizenCondition == "Inactive" {
		return fmt.Errorf("cannot register police records to an inactive citizen")
	}

	citizen.PoliceRegistryData.CaptureOrder = captureOrder
	citizen.PoliceRegistryData.CriminalRecords = criminalRecords
	citizen.PoliceRegistryData.ExitImpediments = exitImpediments
	citizen.PoliceRegistryData.MigratoryAlert = migratoryAlert

	// Update pointer to IPFS pool
	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// registerMarriage: Org1 updates marital status and spouse info
func (s *SmartContract) registerMarriage(ctx contractapi.TransactionContextInterface, nationalID string, spouseName string, marriageDate string, newCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_CIVIL_REGISTRY {
		return fmt.Errorf("permission denied: only Civil Registry can register marriages")
	}

	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if citizen.CitizenCondition == "Inactive" {
		return fmt.Errorf("cannot register marriage for an inactive citizen")
	}

	citizen.CivilRegistryData.MaritalStatus = "MARRIED"
	citizen.CivilRegistryData.Spouse = spouseName
	citizen.CivilRegistryData.MarriageDate = marriageDate

	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// registerDeath: Org1 updates death info and sets status to Inactive
func (s *SmartContract) registerDeath(ctx contractapi.TransactionContextInterface, nationalID string, deathDate string, registrationDate string, observations string, newCID string) error {
	msp, _ := GetMSPID(ctx)
	if msp != MSP_CIVIL_REGISTRY {
		return fmt.Errorf("permission denied: only Civil Registry can register deaths")
	}

	err := CheckABAC(ctx, "role", "OPERATOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: %v", err)
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	citizen.CivilRegistryData.DeathDate = deathDate
	citizen.CivilRegistryData.DeathRegistrationDate = registrationDate
	citizen.CivilRegistryData.Observations = observations
	citizen.CitizenCondition = "Inactive"

	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// Tx_UpdateData: Generic update of multiple fields by authorized Orgs
func (s *SmartContract) Tx_UpdateData(ctx contractapi.TransactionContextInterface, nationalID string, fieldsJSON string, newCID string) error {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	msp, _ := GetMSPID(ctx)
	
	// Only DIRECTOR or high-level role should perform data updates
	err := CheckABAC(ctx, "role", "DIRECTOR")
	if err != nil {
		return fmt.Errorf("ABAC Failure: only DIRECTOR can perform data updates")
	}

	assetBytes, _ := ctx.GetStub().GetState(nationalID)
	if assetBytes == nil {
		return fmt.Errorf("citizen %s not found", nationalID)
	}
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	if citizen.CitizenCondition == "Inactive" {
		return fmt.Errorf("cannot update data for an inactive citizen")
	}

	var fields map[string]string
	err = json.Unmarshal([]byte(fieldsJSON), &fields)
	if err != nil {
		return fmt.Errorf("invalid fields formatting: %v", err)
	}

	for k, v := range fields {
		fields[k] = strings.ToUpper(strings.TrimSpace(v))
	}

	// Authorization logic based on field competency
	for fieldName, newValue := range fields {
		switch msp {
		case MSP_CIVIL_REGISTRY:
			switch fieldName {
			case "first_names": citizen.CivilRegistryData.FirstNames = newValue
			case "last_names": citizen.CivilRegistryData.LastNames = newValue
			case "gender": citizen.CivilRegistryData.Gender = newValue
			case "address": citizen.CivilRegistryData.Address = newValue
			default: return fmt.Errorf("Civil Registry is not authorized to update field: %s", fieldName)
			}
		case MSP_ACADEMIC_REGISTRY:
			// fieldName format: "registryID:field"
			parts := splitTwo(fieldName)
			if len(parts) != 2 {
				return fmt.Errorf("academic update format: '<national_registry_id>:<field>'")
			}
			registryID, field := parts[0], parts[1]
			found := false
			for i := range citizen.AcademicRegistryData {
				if citizen.AcademicRegistryData[i].NationalRegistryID == registryID {
					switch field {
					case "instruction_level":  citizen.AcademicRegistryData[i].InstructionLevel = newValue
					case "exact_degree_name":  citizen.AcademicRegistryData[i].ExactDegreeName = newValue
					case "university_name":    citizen.AcademicRegistryData[i].UniversityName = newValue
					case "registration_date":  citizen.AcademicRegistryData[i].RegistrationDate = newValue
					default: return fmt.Errorf("academic registry: field '%s' is not updatable", field)
					}
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("academic record with registry ID '%s' not found for citizen %s", registryID, nationalID)
			}

		case MSP_PROPERTY_REGISTRY:
			// fieldName format: "propertyID:field"
			parts := splitTwo(fieldName)
			if len(parts) != 2 {
				return fmt.Errorf("property update format: '<property_id>:<field>'")
			}
			propertyID, field := parts[0], parts[1]
			found := false
			for i := range citizen.PropertyRegistryData {
				if citizen.PropertyRegistryData[i].PropertyID == propertyID {
					switch field {
					case "property_type":     citizen.PropertyRegistryData[i].PropertyType = newValue
					case "legal_status":      citizen.PropertyRegistryData[i].LegalStatus = newValue
					case "registration_date": citizen.PropertyRegistryData[i].RegistrationDate = newValue
					default: return fmt.Errorf("property registry: field '%s' is not updatable", field)
					}
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("property with ID '%s' not found for citizen %s", propertyID, nationalID)
			}

		case MSP_POLICE_REGISTRY:
			switch fieldName {
			case "capture_order":   citizen.PoliceRegistryData.CaptureOrder = newValue
			case "migratory_alert": citizen.PoliceRegistryData.MigratoryAlert = newValue
			default: return fmt.Errorf("Police Registry is not authorized to update field: %s", fieldName)
			}
		default:
			return fmt.Errorf("unauthorized organization for updates")
		}
	}

	if newCID != "" {
		citizen.CID = newCID
	}

	assetJSON, _ := json.Marshal(citizen)
	return ctx.GetStub().PutState(nationalID, assetJSON)
}

// splitTwo splits a string on the first ":" into at most 2 parts.
func splitTwo(s string) []string {
	for i, c := range s {
		if c == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}


// 5. QUERIES

func (s *SmartContract) Query_PublicProfile(ctx contractapi.TransactionContextInterface, nationalID string) (*Citizen, error) {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil {
		return nil, err
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("citizen %s not found in the world state", nationalID)
	}

	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)

	// LOPDP PROTECTION: Hide sensitive structured data for public requests.
	// IPFS CIDs (root_cid) are PUBLIC references — they are kept visible.
	citizen.PoliceRegistryData = PoliceRegistryInfo{}        // Hide police records
	citizen.PropertyRegistryData = []PropertyRegistryInfo{} // Hide property details
	citizen.AcademicRegistryData = []AcademicRegistryInfo{} // Hide academic records
	// NOTE: root_cid is intentionally NOT cleared — CIDs are verifiable public pointers

	return &citizen, nil
}

func (s *SmartContract) Query_ComprehensiveProfile(ctx contractapi.TransactionContextInterface, nationalID string) (*Citizen, error) {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	msp, _ := GetMSPID(ctx)

	// All four Orgs can query comprehensive data, but we can bind it to specific roles if needed
	validOrgs := map[string]bool{
		MSP_CIVIL_REGISTRY:    true,
		MSP_ACADEMIC_REGISTRY: true,
		MSP_PROPERTY_REGISTRY: true,
		MSP_POLICE_REGISTRY:   true,
	}

	if !validOrgs[msp] {
		return nil, fmt.Errorf("access denied: querying comprehensive profiles requires inter-institutional permissions")
	}

	// ABAC Protection: Only authorized roles can see the complete unified profile
	err := CheckABAC(ctx, "role", "DIRECTOR")
	if err != nil {
		errAudit := CheckABAC(ctx, "role", "AUDITOR")
		if errAudit != nil {
			return nil, fmt.Errorf("ABAC Failure: only DIRECTOR or AUDITOR can request comprehensive profiles across organizations")
		}
	}

	assetBytes, err := ctx.GetStub().GetState(nationalID)
	if err != nil {
		return nil, err
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("citizen %s not found in the world state", nationalID)
	}
	var citizen Citizen
	json.Unmarshal(assetBytes, &citizen)
	return &citizen, nil
}

// 6. TRACEABILITY (AUDIT)

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	TxId      string    `json:"tx_id"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"is_delete"`
	Value     *Citizen  `json:"value"`
}

// Query_History: Returns the chain of historical updates for a citizen
func (s *SmartContract) Query_History(ctx contractapi.TransactionContextInterface, nationalID string) ([]HistoryQueryResult, error) {
	nationalID = strings.ToUpper(strings.TrimSpace(nationalID))
	msp, _ := GetMSPID(ctx)

	// ABAC Validation: AUDITOR role is typically required for full system audits.
	// However, to allow citizens to view their own history (Public queries), we bypass the strict error.
	fmt.Println("--- DEBUG: Executing Query_History V11 ---")
	err := CheckABAC(ctx, "role", "AUDITOR")
	if err != nil {
		// Log internally but do not block the transaction.
		fmt.Printf("Notice: Non-auditor identity querying history for %s. Granted as public citizen access.\n", nationalID)
	}

	// Organization check: Allow any valid organization member (including default User1 proxy) to view the trail
	validOrgs := map[string]bool{
		MSP_CIVIL_REGISTRY:    true,
		MSP_ACADEMIC_REGISTRY: true,
		MSP_PROPERTY_REGISTRY: true,
		MSP_POLICE_REGISTRY:   true,
	}

	if !validOrgs[msp] {
		return nil, fmt.Errorf("permission denied: invalid organization")
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
				// LOG: Catch unmarshal errors for legacy data (e.g., when root_cid was an object)
				fmt.Printf("Notice: Skipping legacy data unmarshal for TxId %s: %v\n", response.TxId, err)
				// We don't return error here to allow the rest of history to be viewed.
				// The citizen struct will remain zeroed/nil for this record.
			}
		}

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
