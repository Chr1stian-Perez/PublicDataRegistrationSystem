#!/usr/bin/env bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

function createOrgcontraloria {
	infoln "Enrolling the CA admin"
	mkdir -p ../organizations/peerOrganizations/orgcontraloria.example.com/

	export FABRIC_CA_CLIENT_HOME=${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname ca-orgcontraloria --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgcontraloria.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgcontraloria.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgcontraloria.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgcontraloria.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/config.yaml"

	infoln "Registering peer0"
  set -x
	fabric-ca-client register --caname ca-orgcontraloria --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-orgcontraloria --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-orgcontraloria --id.name orgcontraloriaadmin --id.secret orgcontraloriaadminpw --id.type admin --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
	fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-orgcontraloria -M "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-orgcontraloria -M "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls" --enrollment.profile tls --csr.hosts peer0.orgcontraloria.example.com --csr.hosts localhost --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null


  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/ca.crt"
  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/signcerts/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/server.crt"
  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/keystore/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/server.key"

  mkdir "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/tlscacerts"
  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/tlscacerts/ca.crt"

  mkdir "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/tlsca"
  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/tlsca/tlsca.orgcontraloria.example.com-cert.pem"

  mkdir "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/ca"
  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/msp/cacerts/"* "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/ca/ca.orgcontraloria.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
	fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-orgcontraloria -M "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/users/User1@orgcontraloria.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/users/User1@orgcontraloria.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
	fabric-ca-client enroll -u https://orgcontraloriaadmin:orgcontraloriaadminpw@localhost:11054 --caname ca-orgcontraloria -M "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/users/Admin@orgcontraloria.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgcontraloria/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgcontraloria.example.com/users/Admin@orgcontraloria.example.com/msp/config.yaml"
}
