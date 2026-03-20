#!/bin/bash
# Uses $HOME to dynamically resolve the current user's fabric-samples path
export PATH=$HOME/fabric-samples/bin:/usr/local/go/bin:/usr/bin:/bin:$PATH
export FABRIC_CFG_PATH=$HOME/fabric-samples/config/

cd $HOME/fabric-samples/test-network

# Organization names with 'org' prefix to match the directory structure
ORGS=("orgregistrocivil" "orgregistroacademico" "orgregistropropiedad" "orgregistropolicial")
ORG_PORTS=(7054 12054 11054 8054)

for i in "${!ORGS[@]}"; do
  ORG_NAME=${ORGS[$i]}
  CA_PORT=${ORG_PORTS[$i]}
  
  echo ">>> Registering and Enrolling users for $ORG_NAME..."

  TLS_CERT_FILE="${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
  
  # Base directory for the organization
  ORG_DIR="${PWD}/organizations/peerOrganizations/${ORG_NAME}.example.com"

  # 1. Register identities in the CA with their role attributes
  # The CA admin is used for registration
  export FABRIC_CA_CLIENT_HOME="${ORG_DIR}"

  # Register the 4 roles required for contract ABAC enforcement
  fabric-ca-client register --caname ca-${ORG_NAME} --id.name "operator1" --id.secret "operator1pw" --id.type client --id.attrs 'role=OPERATOR:ecert' --tls.certfiles "${TLS_CERT_FILE}" || true
  fabric-ca-client register --caname ca-${ORG_NAME} --id.name "director1" --id.secret "director1pw" --id.type client --id.attrs 'role=DIRECTOR:ecert' --tls.certfiles "${TLS_CERT_FILE}" || true
  fabric-ca-client register --caname ca-${ORG_NAME} --id.name "registrar1" --id.secret "registrar1pw" --id.type client --id.attrs 'role=REGISTRAR:ecert' --tls.certfiles "${TLS_CERT_FILE}" || true
  fabric-ca-client register --caname ca-${ORG_NAME} --id.name "auditor1" --id.secret "auditor1pw" --id.type client --id.attrs 'role=AUDITOR:ecert' --tls.certfiles "${TLS_CERT_FILE}" || true

  # 2. Enrollment phase (generate the actual certificates)
  USERS=("operator1" "director1" "registrar1" "auditor1")
  PWS=("operator1pw" "director1pw" "registrar1pw" "auditor1pw")
  
  for j in "${!USERS[@]}"; do
    USER=${USERS[$j]}
    PW=${PWS[$j]}
    
    DEST_DIR="${ORG_DIR}/users/${USER}@${ORG_NAME}.example.com"
    rm -rf "$DEST_DIR"
    mkdir -p "$DEST_DIR/msp"

    export FABRIC_CA_CLIENT_HOME="${DEST_DIR}"
    
    # Enroll the user against their organization's local CA
    fabric-ca-client enroll -u "https://${USER}:${PW}@localhost:${CA_PORT}" \
      --caname ca-${ORG_NAME} \
      -M "${DEST_DIR}/msp" \
      --tls.certfiles "${TLS_CERT_FILE}"
      
    # Copy config.yaml so the SDK can correctly resolve the MSP
    cp "${ORG_DIR}/msp/config.yaml" "${DEST_DIR}/msp/config.yaml"
  done
done

echo "Done generating ABAC accounts for all organizations."
