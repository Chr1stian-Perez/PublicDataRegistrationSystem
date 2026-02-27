#!/bin/bash

echo "Looking for dynamic certificates from the Civil Registry clerk..."


BASE_DIR="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp"
CONN_PROFILE="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/connection-orgregistrocivil.json"

PRIV_KEY=$(find $BASE_DIR/keystore -name "*_sk" -print -quit)
CERT=$(find $BASE_DIR/signcerts -name "*.pem" -print -quit)

if [ -z "$PRIV_KEY" ]; then
    echo " ERROR: The office worker's private key was not found. Did you enroll?"
    exit 1
fi

cat << EOF > network.yaml
name: Red Gubernamental Caliper ABAC
version: "2.0.0"
caliper:
  blockchain: fabric
channels:
  - channelName: mychannel
    contracts:
    - id: dtic
organizations:
  - mspid: OrgregistrocivilMSP
    identities:
      certificates:
      - name: 'oficinista_abac'
        clientPrivateKey:
          path: $PRIV_KEY
        clientSignedCert:
          path: $CERT
    connectionProfile:
      path: $CONN_PROFILE
      discover: true
EOF

echo " "Network.yaml file generated!"
