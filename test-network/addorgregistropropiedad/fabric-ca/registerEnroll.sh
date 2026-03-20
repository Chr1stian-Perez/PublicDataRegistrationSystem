#!/usr/bin/env bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

function createOrgregistropropiedad {
	infoln "Enrolling the CA admin"
	mkdir -p ../organizations/peerOrganizations/orgregistropropiedad.example.com/

	export FABRIC_CA_CLIENT_HOME=${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname ca-orgregistropropiedad --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgregistropropiedad.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgregistropropiedad.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgregistropropiedad.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-orgregistropropiedad.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/config.yaml"

	infoln "Registering peer0"
  set -x
	fabric-ca-client register --caname ca-orgregistropropiedad --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-orgregistropropiedad --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-orgregistropropiedad --id.name orgregistropropiedadadmin --id.secret orgregistropropiedadadminpw --id.type admin --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
	fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-orgregistropropiedad -M "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-orgregistropropiedad -M "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls" --enrollment.profile tls --csr.hosts peer0.orgregistropropiedad.example.com --csr.hosts localhost --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null


  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/ca.crt"
  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/signcerts/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/server.crt"
  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/keystore/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/server.key"

  mkdir "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/tlscacerts"
  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/tlscacerts/ca.crt"

  mkdir "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/tlsca"
  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls/tlscacerts/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/tlsca/tlsca.orgregistropropiedad.example.com-cert.pem"

  mkdir "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/ca"
  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/msp/cacerts/"* "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/ca/ca.orgregistropropiedad.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
	fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-orgregistropropiedad -M "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/users/User1@orgregistropropiedad.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/users/User1@orgregistropropiedad.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
	fabric-ca-client enroll -u https://orgregistropropiedadadmin:orgregistropropiedadadminpw@localhost:11054 --caname ca-orgregistropropiedad -M "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/users/Admin@orgregistropropiedad.example.com/msp" --tls.certfiles "${PWD}/fabric-ca/orgregistropropiedad/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/msp/config.yaml" "${PWD}/../organizations/peerOrganizations/orgregistropropiedad.example.com/users/Admin@orgregistropropiedad.example.com/msp/config.yaml"
}
