package main

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	channelName   = "mychannel"
	chaincodeName = "dtic_chaincode"
	ipfsAPI       = "http://localhost:5001/api/v0"
)
func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n====================================================")
		fmt.Println("    PUBLIC DATA REGISTRATION SYSTEM (4-ORG)")
		fmt.Println("====================================================")
		fmt.Println("Select organization:")
		fmt.Println("1. Civil registry        (registro_civil)")
		fmt.Println("2. Academic registry     (registro_academico)")
		fmt.Println("3. Property registry     (registro_propiedad)")
		fmt.Println("4. Police registry       (registro_policial)")
		fmt.Println("5. Public queries")
		fmt.Println("0. Exit")
		fmt.Print("Option: ")

		var orgIdx int
		_, err := fmt.Scanln(&orgIdx)
		if err != nil {
			reader.ReadString('\n')
			continue
		}
		if orgIdx == 0 {
			fmt.Println("Goodbye")
			break
		}
		if orgIdx == 5 {
			runPublicQueries(reader)
			continue
		}
		if orgIdx < 1 || orgIdx > 4 {
			fmt.Println("Invalid option.")
			continue
		}
		runOrgInterface(orgIdx, reader)
	}
}

// ─── ORG INTERFACE ────────────────────────────────────────────────────────────

func runOrgInterface(orgIdx int, reader *bufio.Reader) {
	orgName := getOrgName(orgIdx)
	fmt.Printf("\n--- %s Interface ---\n", strings.ToUpper(orgName))

	// LOGIN (Select Identity)
	usersDir := fmt.Sprintf("../test-network/organizations/peerOrganizations/%s.example.com/users", orgName)
	entries, err := ioutil.ReadDir(usersDir)
	if err != nil {
		fmt.Printf("Error reading identities: %v\n", err)
		return
	}

	var identities []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Extract just the username part Before @
			parts := strings.Split(entry.Name(), "@")
			identities = append(identities, parts[0])
		}
	}

	if len(identities) == 0 {
		fmt.Println("No identities found for this organization. Please run CA enrollment scripts.")
		return
	}

	fmt.Println("\n[LOGIN]")
	fmt.Println("Select your identity (Certificate role):")
	for i, id := range identities {
		fmt.Printf("%d. %s\n", i+1, id)
	}
	fmt.Print("Choice: ")

	var idChoice int
	_, err = fmt.Scanln(&idChoice)
	if err != nil || idChoice < 1 || idChoice > len(identities) {
		reader.ReadString('\n')
		fmt.Println("Invalid identity choice.")
		return
	}

	selectedIdentity := identities[idChoice-1]
	fmt.Printf("Logging in as: %s \n", selectedIdentity)

	gw, conn, err := connectToGateway(orgIdx, selectedIdentity)
	if err != nil {
		fmt.Printf("Error connecting to gateway: %v\n", err)
		return
	}
	defer conn.Close()
	defer gw.Close()

	contract := gw.GetNetwork(channelName).GetContract(chaincodeName)

	for {
		fmt.Printf("\n[%s MENU]\n", strings.ToUpper(orgName))
		switch orgIdx {
		case 1:
			fmt.Println("1. Register identity (Birth)")
			fmt.Println("2. Register marriage")
			fmt.Println("3. Register death")
		case 2:
			fmt.Println("1. Register academic degree")
		case 3:
			fmt.Println("1. Register property")
		case 4:
			fmt.Println("1. Register police record")
		}
		fmt.Println("8. View complete citizen profile (this institution only)")
		fmt.Println("9. Data update")
		fmt.Println("0. Back to main menu")
		fmt.Print("Option: ")

		var op int
		_, err := fmt.Scanln(&op)
		if err != nil {
			reader.ReadString('\n')
			continue
		}
		if op == 0 {
			break
		}
		if op == 9 {
			performUpdateData(orgIdx, contract, reader)
			continue
		}
		if op == 8 {
			id := readInput(reader, "National ID to query: ")
			showOrgProfile(orgIdx, contract, id)
			continue
		}
		processTransaction(orgIdx, op, contract, reader)
	}
}

func processTransaction(orgIdx, op int, contract *client.Contract, reader *bufio.Reader) {
	var recordType string
	var id string
	var cid string
	recordData := make(map[string]string)

	switch orgIdx {
	case 1: // ══ REGISTRO CIVIL ══════════════════════════════════════════
		switch op {
		case 1:
			recordType = "IDENTITY"
			fmt.Println("\n[BIRTH REGISTRATION - CIVIL REGISTRY]")
			id                             = readInput(reader, "National ID                         : ")
			recordData["first_names"]      = readInput(reader, "First names                         : ")
			recordData["last_names"]       = readInput(reader, "Last names                          : ")
			recordData["sex"]              = inputOrDefault(reader, "Sex [MALE/FEMALE]                   : ", "MALE")
			recordData["gender"]           = recordData["sex"]
			recordData["nationality"]      = inputOrDefault(reader, "Nationality                         : ", "ECUADORIAN")
			recordData["birth_date"]       = readInput(reader, "Birth date (YYYY-MM-DD)             : ")
			recordData["birth_place"]      = inputOrDefault(reader, "Birth place (State/County/Parish)   : ", "N/A")
			recordData["marital_status"]   = "SINGLE"
			recordData["spouse"]           = "N/A"
			recordData["father_id"]        = inputOrDefault(reader, "Father national ID  (ENTER=N/A)     : ", "N/A")
			recordData["father_nationality"] = inputOrDefault(reader, "Father nationality  (ENTER=ECUADORIAN): ", "ECUADORIAN")
			recordData["mother_id"]        = inputOrDefault(reader, "Mother national ID  (ENTER=N/A)     : ", "N/A")
			recordData["mother_nationality"] = inputOrDefault(reader, "Mother nationality  (ENTER=ECUADORIAN): ", "ECUADORIAN")
			recordData["address"]          = inputOrDefault(reader, "Home address (State/County/Parish)  : ", "N/A")
			recordData["marriage_date"]    = "N/A"
			recordData["death_date"]       = "N/A"
			recordData["observations"]     = "N/A"
			recordData["death_registration_date"] = "N/A"
			recordData["citizen_condition"] = "Active"
			cid = promptDocumentUpload(reader, "Birth Certificate")

		case 2:
			recordType = "MARRIAGE"
			fmt.Println("\n[MARRIAGE REGISTRATION - CIVIL REGISTRY]")
			id                        = readInput(reader, "National ID of citizen              : ")
			recordData["spouse_name"] = readInput(reader, "Spouse full name                    : ")
			recordData["marriage_date"] = readInput(reader, "Marriage date (YYYY-MM-DD)          : ")
			cid = promptDocumentUpload(reader, "Marriage certificate")

		case 3:
			recordType = "DEATH"
			fmt.Println("\n[DEATH REGISTRATION - CIVIL REGISTRY]")
			id                              = readInput(reader, "National ID of citizen              : ")
			recordData["death_date"]        = readInput(reader, "Death date (YYYY-MM-DD)             : ")
			recordData["registration_date"] = readInput(reader, "Registration date (YYYY-MM-DD)      : ")
			recordData["observations"]      = inputOrDefault(reader, "Observations (ENTER=N/A)            : ", "N/A")
			cid = promptDocumentUpload(reader, "Death certificate")
		}

	case 2: // ══ REGISTRO ACADÉMICO ══════════════════════════════════════
		if op == 1 {
			recordType = "ACADEMIC_DEGREE"
			fmt.Println("\n[ACADEMIC DEGREE REGISTRATION - ACADEMIC REGISTRY]")
			id                               = readInput(reader, "National ID                         : ")
			recordData["registry_id"]        = readInput(reader, "Unique registry number              : ")
			recordData["instruction_level"]  = inputOrDefault(reader, "Instruction level (e.g. HIGHER)     : ", "HIGHER")
			recordData["exact_degree_name"]  = readInput(reader, "Exact degree name                   : ")
			recordData["university_name"]    = readInput(reader, "University / Institution            : ")
			recordData["registration_date"]  = inputOrDefault(reader, "Registration date (YYYY-MM-DD)      : ", "N/A")
			cid = promptDocumentUpload(reader, "Academic Degree Diploma")
		}

	case 3: // ══ REGISTRO DE LA PROPIEDAD ════════════════════════════════
		if op == 1 {
			recordType = "PROPERTY"
			fmt.Println("\n[PROPERTY REGISTRATION - PROPERTY REGISTRY]")
			ownerID := readInput(reader, "Owner National ID                   : ")
			id = ownerID
			recordData["owner_id"]          = ownerID
			recordData["property_id"]       = readInput(reader, "Property ID number                  : ")
			recordData["owner_name"]        = readInput(reader, "Owner full name                     : ")
			recordData["property_type"]     = inputOrDefault(reader, "Property type [HOUSE/APARTMENT/COMMERCIAL]: ", "HOUSE")
			recordData["legal_status"]      = inputOrDefault(reader, "Legal status [CLEAR/SEIZED]         : ", "CLEAR")
			recordData["registration_date"] = inputOrDefault(reader, "Registration date (YYYY-MM-DD)      : ", "N/A")
			cid = promptDocumentUpload(reader, "Property Deed")
		}

	case 4: // ══ REGISTRO POLICIAL ════════════════════════════════════════
		if op == 1 {
			recordType = "POLICE_RECORD"
			fmt.Println("\n[POLICE RECORD - POLICE REGISTRY]")
			id                               = readInput(reader, "Document ID                         : ")
			recordData["capture_order"]      = inputOrDefault(reader, "Capture order      [N/A or desc]    : ", "N/A")
			recordData["criminal_records"]   = inputOrDefault(reader, "Criminal records   [N/A or desc]    : ", "N/A")
			recordData["exit_impediments"]   = inputOrDefault(reader, "Exit impediments   [No/Yes]         : ", "No")
			recordData["migratory_alert"]    = inputOrDefault(reader, "Migratory alert    [N/A or desc]    : ", "N/A")
			cid = promptDocumentUpload(reader, "Police Reference Document")
		}
	}

	if recordType == "" {
		fmt.Println("Invalid operation for this organization.")
		return
	}

	recordDataJSON, err := json.Marshal(recordData)
	if err != nil {
		fmt.Printf("Failed to encode data to JSON: %v\n", err)
		return
	}

	fmt.Printf("\nSubmitting Tx_Register [%s]...\n", recordType)
	_, err = contract.SubmitTransaction("Tx_Register", id, recordType, string(recordDataJSON), cid)
	if err != nil {
		handleResult(err, "")
		return
	}
	fmt.Println("✓ Transaction committed to blockchain successfully!")
	// Auto-display the profile immediately to confirm what was saved
	if id != "" {
		fmt.Printf("\n→ Retrieving saved record for citizen %s \n", id)
		showOrgProfile(orgIdx, contract, id)
	}
}

// ─── DATA UPDATE ──────────────────────────────────────────────────────────────

func performUpdateData(orgIdx int, contract *client.Contract, reader *bufio.Reader) {
	fmt.Println("\n[DATA UPDATE - Tx_UpdateData]")

	id := readInput(reader, "National ID                          : ")
	fields := make(map[string]string)

	switch orgIdx {
	case 1:
		fmt.Println("Updatable fields (Civil Registry):")
		fmt.Println("  first_names | last_names | gender | address")
	case 2:
		fmt.Println("Updatable fields (Academic Registry — format: <registry_id>:<field>):")
		fmt.Println("  instruction_level | exact_degree_name | university_name | registration_date")
		fmt.Println("  Example field: 'REG-001:exact_degree_name'")
	case 3:
		fmt.Println("Updatable fields (Property Registry — format: <property_id>:<field>):")
		fmt.Println("  property_type | legal_status | registration_date")
		fmt.Println("  Example field: '375736:legal_status'")
	case 4:
		fmt.Println("Updatable fields (Police Registry):")
		fmt.Println("  capture_order | migratory_alert")
	}

	for {
		field := readInput(reader, "Field name (or type 'DONE' to finish) : ")
		if strings.ToUpper(strings.TrimSpace(field)) == "DONE" {
			if len(fields) == 0 {
				fmt.Println("No fields entered to update. Cancelling.")
				return
			}
			break
		}
		newVal := readInput(reader, "New value                            : ")
		fields[field] = newVal
	}

	cid := promptDocumentUpload(reader, "Supporting evidence for update (Doc. Respaldo)")

	fieldsJSONBytes, err := json.Marshal(fields)
	if err != nil {
		fmt.Printf("Failed to encode fields to JSON: %v\n", err)
		return
	}

	fmt.Println("Submitting Tx_UpdateData ")
	_, err = contract.SubmitTransaction("Tx_UpdateData", id, string(fieldsJSONBytes), cid)
	handleResult(err, "✓ Data update applied successfully.")
}

// showOrgProfile shows only the data of the citizen that belongs to the calling org.
func showOrgProfile(orgIdx int, contract *client.Contract, nationalID string) {
	orgLabel := strings.ToUpper(getOrgName(orgIdx))
	fmt.Printf("\n[%s — CITIZEN PROFILE | ID: %s]\n", orgLabel, nationalID)
	fmt.Println(strings.Repeat("─", 60))

	res, err := contract.EvaluateTransaction("Query_ComprehensiveProfile", nationalID)
	if err != nil {
		// fallback: show public profile
		res2, err2 := contract.EvaluateTransaction("Query_PublicProfile", nationalID)
		if err2 != nil {
			fmt.Printf("Error retrieving profile: %v\n", err2)
			return
		}
		fmt.Println(formatJSON(res2))
		fmt.Println(strings.Repeat("─", 60))
		return
	}

	var full map[string]interface{}
	if err := json.Unmarshal(res, &full); err != nil {
		fmt.Println(formatJSON(res))
		return
	}

	filtered := filterForOrg(orgIdx, full)
	out, _ := json.MarshalIndent(filtered, "", "  ")
	fmt.Println(string(out))
	fmt.Println(strings.Repeat("─", 60))
}

// filterForOrg extracts only the JSON fields the given org is authorized to see.
func filterForOrg(orgIdx int, data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	// Every org sees the citizen's basic identity header
	for _, f := range []string{"national_id", "global_status"} {
		if v, ok := data[f]; ok {
			result[f] = v
		}
	}

	// Helper: get the CID
	getCID := func() string {
		if cid, ok := data["CID"].(string); ok {
			return cid
		}
		return ""
	}
	switch orgIdx {
	case 1: // Civil Registry — sees all civil registration data
		for _, f := range []string{
			"first_names", "last_names", "sex", "gender", "nationality",
			"birth_date", "birth_place", "marital_status", "spouse",
			"father_name", "father_nationality",
			"mother_name", "mother_nationality",
			"address", "marriage_date", "death_date",
			"observations", "death_registration_date",
		} {
			if v, ok := data["civil_registry_data"]; ok {
				if m, ok := v.(map[string]interface{}); ok {
					if val, ok := m[f]; ok {
						result[f] = val
					}
				}
			}
		}
		result["CID"] = getCID()

	case 2: // Academic Registry — sees only academic degrees
		if cr, ok := data["civil_registry_data"].(map[string]interface{}); ok {
			result["first_names"] = cr["first_names"]
			result["last_names"] = cr["last_names"]
		}
		result["academic_registry_data"] = data["academic_registry_data"]
		result["CID"] = getCID()

	case 3: // Property Registry — sees only property holdings
		if cr, ok := data["civil_registry_data"].(map[string]interface{}); ok {
			result["owner_name"] = fmt.Sprintf("%s %s", cr["first_names"], cr["last_names"])
		}
		result["property_registry_data"] = data["property_registry_data"]
		result["CID"] = getCID()

	case 4: // Police Registry — sees only police record
		result["police_registry_data"] = data["police_registry_data"]
		result["CID"] = getCID()
	}

	return result
}

// ─── PUBLIC QUERIES 
func runPublicQueries(reader *bufio.Reader) {
	fmt.Println("\n--- PUBLIC QUERIES ---")

	// LOGIN (Select identity)
	usersDir := "../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/users"
	entries, err := ioutil.ReadDir(usersDir)
	if err != nil {
		fmt.Printf("Error reading identities: %v\n", err)
		return
	}
	var identities []string
	for _, entry := range entries {
		if entry.IsDir() {
			parts := strings.Split(entry.Name(), "@")
			identities = append(identities, parts[0])
		}
	}

	if len(identities) == 0 {
		fmt.Println("No identities found. Please run CA enrollment scripts.")
		return
	}

	fmt.Println("\n[LOGIN]")
	fmt.Println("Select your identity (Role):")
	for i, id := range identities {
		fmt.Printf("%d. %s\n", i+1, id)
	}
	fmt.Print("Choice: ")

	var idChoice int
	_, err = fmt.Scanln(&idChoice)
	if err != nil || idChoice < 1 || idChoice > len(identities) {
		reader.ReadString('\n')
		fmt.Println("Invalid identity choice.")
		return
	}

	selectedIdentity := identities[idChoice-1]
	fmt.Printf("Logging in as: %s \n", selectedIdentity)

	gw, conn, err := connectToGateway(1, selectedIdentity) // Civil Registry peer
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer conn.Close()
	defer gw.Close()

	contract := gw.GetNetwork(channelName).GetContract(chaincodeName)

	for {
		fmt.Println("\n[PUBLIC QUERIES]")
		fmt.Println("1. View public profile")
		fmt.Println("2. View blockchain history")
		fmt.Println("0. Back")
		fmt.Print("Option: ")

		var op int
		fmt.Scanln(&op)
		if op == 0 {
			break
		}

		id := readInput(reader, "National ID: ")
		switch op {
		case 1:
			res, err := contract.EvaluateTransaction("Query_PublicProfile", id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nPublic Profile:\n%s\n", formatJSON(res))
			}
		case 2:
			res, err := contract.EvaluateTransaction("Query_History", id)
			if err != nil {
				fmt.Printf("Error (History may require AUDITOR role): %v\n", err)
			} else {
				fmt.Printf("\nBlockchain History:\n%s\n", formatJSON(res))
			}
		}
	}
}

// ─── IPFS INTEGRATION
// promptDocumentUpload offers to upload a reference PDF to IPFS.
// Defaults to Certificate.pdf (symbolic document) if the user presses ENTER.
// Returns the CID.
func promptDocumentUpload(reader *bufio.Reader, docLabel string) string {
	const defaultCert = "/home/cperez/fabric-samples/app/Certificate.pdf"
	fmt.Printf("\n[IPFS - %s]\n", docLabel)
	fmt.Printf("Reference PDF (ENTER = use symbolic Certificate.pdf): ")
	filePath := readInput(reader, "")
	if filePath == "" {
		filePath = defaultCert
		fmt.Printf("  → Using: %s\n", filePath)
	}

	var cid string
	var err error

	if filePath == defaultCert {
		// Alter the default certificate to generate a unique CID
		cid, err = alterAndUploadToIPFS(filePath)
	} else {
		cid, err = uploadToIPFS(filePath)
	}
	if err != nil {
		fmt.Printf("IPFS upload failed: %v\n", err)
		fmt.Println("   CID will be empty for this transaction.")
		return ""
	}

	fmt.Printf("✓  Document uploaded to IPFS!\n")
	fmt.Printf("   CID : %s\n", cid)
	fmt.Printf("   URL : http://localhost:8080/ipfs/%s\n", cid)
	return cid
}


// alterAndUploadToIPFS appends a random timestamp/nonce to the file before uploading
// ensuring that the resulting IPFS CID is always unique.
func alterAndUploadToIPFS(filePath string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot read file for alteration: %v", err)
	}

	// Append a unique comment at the end of the PDF
	nonce := fmt.Sprintf("\n%% IPFS-NONCE: %d\n", time.Now().UnixNano())
	modifiedBytes := append(fileBytes, []byte(nonce)...)

	return uploadBytesToIPFS(modifiedBytes, filepath.Base(filePath))
}

// uploadBytesToIPFS uploads in-memory bytes to the IPFS node and returns its CID hash.
func uploadBytesToIPFS(fileBytes []byte, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://localhost:5001/api/v0/add", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IPFS Add failed with status: %s", resp.Status)
	}

	var ipfsResp struct {
		Hash string `json:"Hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ipfsResp); err != nil {
		return "", err
	}

	if ipfsResp.Hash == "" {
		return "", fmt.Errorf("IPFS returned empty CID")
	}

	return ipfsResp.Hash, nil
}

// uploadToIPFS uploads a local file to the IPFS node and returns its CID hash.
func uploadToIPFS(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("form create error: %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("file copy error: %v", err)
	}
	writer.Close()

	resp, err := http.Post(ipfsAPI+"/add?pin=true", writer.FormDataContentType(), body)
	if err != nil {
		return "", fmt.Errorf("IPFS API not reachable at %s (is IPFS running?): %v", ipfsAPI, err)
	}
	defer resp.Body.Close()

	var result struct {
		Hash string `json:"Hash"`
		Name string `json:"Name"`
		Size string `json:"Size"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("IPFS response parse error: %v", err)
	}
	if result.Hash == "" {
		return "", fmt.Errorf("IPFS returned empty CID")
	}
	return result.Hash, nil
}

// ─── HELPERS
func inputOrDefault(reader *bufio.Reader, prompt, defaultVal string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	val := strings.TrimSpace(text)
	if val == "" {
		if defaultVal != "" {
			fmt.Printf("  → using default: %s\n", defaultVal)
		}
		return defaultVal
	}
	return val
}

func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func handleResult(err error, successMsg string) {
	if err != nil {
		fmt.Printf("\n✗ Transaction Failed: %v\n", err)
		
		// Attempt to extract gRPC status from the error to get chaincode details
		if s, ok := status.FromError(err); ok {
			fmt.Printf("  Message: %s\n", s.Message())
			for _, detail := range s.Details() {
				fmt.Printf("  Detail: %v\n", detail)
			}
		}
	} else {
		fmt.Println(successMsg)
	}
}

func formatJSON(data []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return string(data)
	}
	return out.String()
}

func getOrgName(idx int) string {
	return []string{"", "orgregistrocivil", "orgregistroacademico", "orgregistropropiedad", "orgregistropolicial"}[idx]
}

// ─── FABRIC GATEWAY CONNECTION
func connectToGateway(orgIdx int, identityName string) (*client.Gateway, *grpc.ClientConn, error) {
	org := getOrgName(orgIdx)
	cryptoBase := fmt.Sprintf("../test-network/organizations/peerOrganizations/%s.example.com", org)
	
	certPath   := path.Join(cryptoBase, fmt.Sprintf("users/%s@%s.example.com/msp/signcerts/cert.pem", identityName, org))
	keyDir     := path.Join(cryptoBase, fmt.Sprintf("users/%s@%s.example.com/msp/keystore", identityName, org))
	tlsPath    := path.Join(cryptoBase, fmt.Sprintf("peers/peer0.%s.example.com/tls/ca.crt", org))
	ports      := []int{0, 7051, 9051, 11051, 13051}
	endpoint   := fmt.Sprintf("localhost:%d", ports[orgIdx])

	conn, err := newGrpcConnection(tlsPath, endpoint, fmt.Sprintf("peer0.%s.example.com", org))
	if err != nil {
		return nil, nil, err
	}
	id, err := newIdentity(orgIdx, certPath)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	sign, err := newSign(keyDir)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	gw, err := client.Connect(id, client.WithSign(sign), client.WithClientConnection(conn))
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return gw, conn, nil
}

func newGrpcConnection(tlsCertPath, endpoint, gatewayPeer string) (*grpc.ClientConn, error) {
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to append certs from PEM")
	}
	creds := credentials.NewClientTLSFromCert(cp, gatewayPeer)
	return grpc.Dial(endpoint, grpc.WithTransportCredentials(creds))
}

func newIdentity(orgIdx int, certPath string) (*identity.X509Identity, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}
	mspIDs := []string{"", "OrgregistrocivilMSP", "OrgregistroacademicoMSP", "OrgregistropropiedadMSP", "OrgregistropolicialMSP"}
	return identity.NewX509Identity(mspIDs[orgIdx], cert)
}

func newSign(keyDir string) (identity.Sign, error) {
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("keystore is empty in %s", keyDir)
	}
	keyPEM, err := os.ReadFile(path.Join(keyDir, files[0].Name()))
	if err != nil {
		return nil, err
	}
	privateKey, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewPrivateKeySign(privateKey)
}
