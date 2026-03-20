package main

// DATA STRUCTURE DEFINITIONS (ASSET)
// 1. INFORMATIVE DATA (World state - Fast queries)

type CivilRegistryInfo struct {
	FirstNames            string `json:"first_names"`
	LastNames             string `json:"last_names"`
	Sex                   string `json:"sex"`
	Gender                string `json:"gender"`
	Nationality           string `json:"nationality"`
	BirthDate             string `json:"birth_date"`
	BirthPlace            string `json:"birth_place"`
	MaritalStatus         string `json:"marital_status"`
	Spouse                string `json:"spouse"`
	FatherID              string `json:"father_id"`
	FatherNationality     string `json:"father_nationality"`
	MotherID              string `json:"mother_id"`
	MotherNationality     string `json:"mother_nationality"`
	Address               string `json:"address"`
	MarriageDate          string `json:"marriage_date"`
	DeathDate             string `json:"death_date"`
	Observations          string `json:"observations"`
	DeathRegistrationDate string `json:"death_registration_date"`
}

type AcademicRegistryInfo struct {
	NationalRegistryID string `json:"national_registry_id"`
	InstructionLevel   string `json:"instruction_level"`
	ExactDegreeName    string `json:"exact_degree_name"`
	UniversityName     string `json:"university_name"`
	RegistrationDate   string `json:"registration_date"`
}

type PropertyRegistryInfo struct {
	PropertyID       string `json:"property_id"`
	OwnerName        string `json:"owner_name"`
	OwnerID          string `json:"owner_id"`
	PropertyType     string `json:"property_type"`
	LegalStatus      string `json:"legal_status"`
	RegistrationDate string `json:"registration_date"`
}

type PoliceRegistryInfo struct {
	CaptureOrder    string `json:"capture_order"`
	CriminalRecords string `json:"criminal_records"`
	ExitImpediments string `json:"exit_impediments"`
	MigratoryAlert  string `json:"migratory_alert"`
}

// 2. EVIDENCE STRUCTURE (SINGLE CID POOL)
// 3. MAIN ASSET: CITIZEN

type Citizen struct {
	GlobalID         string `json:"national_id"`   // PK
	CitizenCondition string `json:"global_status"` // Active, Inactive

	// State Data
	CivilRegistryData    CivilRegistryInfo      `json:"civil_registry_data"`
	AcademicRegistryData []AcademicRegistryInfo `json:"academic_registry_data"` // Array for multiple degrees
	PropertyRegistryData []PropertyRegistryInfo `json:"property_registry_data"` // Array for multiple properties
	PoliceRegistryData   PoliceRegistryInfo     `json:"police_registry_data"`

	// Evidence (Pointer to IPFS pool)
	CID string `json:"CID"`
}
