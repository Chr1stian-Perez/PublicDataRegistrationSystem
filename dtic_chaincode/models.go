package main

// DATA STRUCTURE DEFINITIONS (ASSET DEFINITION)

// 1. INFORMATIVE DATA (World State - Fast Queries)

type CivilRegistryInfo struct {
	FirstNames    string `json:"first_names"`
	LastNames     string `json:"last_names"`
	Sex           string `json:"sex"`    // E.g.: "MALE", "FEMALE"
	Gender        string `json:"gender"` // E.g.: "MASCULINE", "FEMININE"
	Nationality   string `json:"nationality"`
	BirthDate     string `json:"birth_date"`
	BirthPlace    string `json:"birth_place"`
	MaritalStatus string `json:"marital_status"`
}

type ElectoralCouncilInfo struct {
	Province       string `json:"province"`
	Canton         string `json:"canton"`
	Parish         string `json:"parish"`
	VotingPrecinct string `json:"voting_precinct"`
	HasVoted       bool   `json:"has_voted"`
	HasFines       bool   `json:"has_fines"`
	VotingTable    string `json:"voting_table"` // Links to the Electoral Roll
}

// ComptrollerInfo: New struct to handle the citizen's legal status
type ComptrollerInfo struct {
	ComptrollerStatus string `json:"comptroller_status"` // CLEAN, NOTIFIED, PREDETERMINED, LIABLE...
	LastNotification  string `json:"last_notification"`  // Date or reference
}

// 2. HIERARCHICAL EVIDENCE STRUCTURE (IPFS MERKLE DAG)
type IPFSStructure struct {
	RootCID             string `json:"root_cid"`                  // User's Root Hash (Changes if any child changes)
	CivilRegistryDirCID string `json:"civil_registry_dir_cid"`    // Hash of the "Civil Registry Folder" directory
	ElectoralCouncilDirCID  string `json:"electoral_council_dir_cid"` // Hash of the "Electoral Council Folder" directory (Electoral Rolls)
	ComptrollerDirCID   string `json:"comptroller_dir_cid"`       // Hash of the "Comptroller Folder" directory
}

// 3. AUDIT STRUCTURES

type DefenseEvidence struct {
	Phase        string `json:"phase"`        // E.g.: "EXECUTION", "APPEAL"
	Description  string `json:"description"`  // Summary of the action
	EvidenceCID  string `json:"evidence_cid"` // Supporting document specific to this step
	Date         string `json:"date"`
	Author       string `json:"author"`       // MSP ID of the author
}

type AuditProcess struct {
	CaseID         string            `json:"case_id"`
	CurrentStatus  string            `json:"current_status"` // DRAFT, REPORT, LIABILITY, ETC.
	DefenseHistory []DefenseEvidence `json:"defense_history"`
	Conclusion     string            `json:"conclusion"`     // PENDING, SANCTION, ACQUITTAL
}

// 4. MAIN ASSET: UNIFIED CITIZEN

type Citizen struct {
	GlobalID        string          `json:"national_id"`   // PK
	GlobalStatus    string          `json:"global_status"` // ACTIVE, DECEASED
	
	// State Data
	CivilRegistryData    CivilRegistryInfo    `json:"civil_registry_data"`
	ElectoralCouncilData ElectoralCouncilInfo `json:"electoral_council_data"`
	ComptrollerData      ComptrollerInfo      `json:"comptroller_data"` // New segment
	
	// Evidence (Pointers to IPFS Directories)
	DigitalEvidence IPFSStructure `json:"digital_evidence"`
	
	// Control Processes
	AuditProcesses []AuditProcess `json:"audit_processes"`
}

// 5. BATCH TRANSACTION & ANCILLARY INPUTS


type ElectoralUpdateInput struct {
	NationalID string `json:"national_id"`
	HasVoted   bool   `json:"has_voted"`
	HasFines   bool   `json:"has_fines"`
	NewRootCID string `json:"new_root_cid"` // Computed off-chain for each user
}
