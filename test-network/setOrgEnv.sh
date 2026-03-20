#!/usr/bin/env bash
#
# SPDX-License-Identifier: Apache-2.0




# default to using Orgregistrocivil
ORG=${1:-Orgregistrocivil}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

ORDERER_CA=${DIR}/test-network/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
PEER0_ORG1_CA=${DIR}/test-network/organizations/peerOrganizations/orgregistrocivil.example.com/tlsca/tlsca.orgregistrocivil.example.com-cert.pem
PEER0_ORG2_CA=${DIR}/test-network/organizations/peerOrganizations/orgregistropolicial.example.com/tlsca/tlsca.orgregistropolicial.example.com-cert.pem
PEER0_ORG3_CA=${DIR}/test-network/organizations/peerOrganizations/orgregistropropiedad.example.com/tlsca/tlsca.orgregistropropiedad.example.com-cert.pem
PEER0_ORG4_CA=${DIR}/test-network/organizations/peerOrganizations/orgregistroacademico.example.com/tlsca/tlsca.orgregistroacademico.example.com-cert.pem

if [[ ${ORG,,} == "orgregistrocivil" || ${ORG,,} == "digibank" ]]; then

   CORE_PEER_LOCALMSPID=OrgregistrocivilMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/test-network/organizations/peerOrganizations/orgregistrocivil.example.com/users/Admin@orgregistrocivil.example.com/msp
   CORE_PEER_ADDRESS=localhost:7051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/test-network/organizations/peerOrganizations/orgregistrocivil.example.com/tlsca/tlsca.orgregistrocivil.example.com-cert.pem

elif [[ ${ORG,,} == "orgregistropolicial" || ${ORG,,} == "magnetocorp" ]]; then

   CORE_PEER_LOCALMSPID=OrgregistropolicialMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/test-network/organizations/peerOrganizations/orgregistropolicial.example.com/users/Admin@orgregistropolicial.example.com/msp
   CORE_PEER_ADDRESS=localhost:9051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/test-network/organizations/peerOrganizations/orgregistropolicial.example.com/tlsca/tlsca.orgregistropolicial.example.com-cert.pem
elif [[ ${ORG,,} == "orgregistropropiedad" ]]; then
   CORE_PEER_LOCALMSPID=OrgregistropropiedadMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/test-network/organizations/peerOrganizations/orgregistropropiedad.example.com/users/Admin@orgregistropropiedad.example.com/msp
   CORE_PEER_ADDRESS=localhost:11051
   CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG3_CA}

# Nueva sección para Orgregistroacademico (Org4)
elif [[ ${ORG,,} == "orgregistroacademico" ]]; then
   CORE_PEER_LOCALMSPID=OrgregistroacademicoMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/test-network/organizations/peerOrganizations/orgregistroacademico.example.com/users/Admin@orgregistroacademico.example.com/msp
   CORE_PEER_ADDRESS=localhost:12051
   CORE_PEER_TLS_ROOTCERT_FILE=${PEER0_ORG4_CA}
else
   echo "Unknown \"$ORG\", please choose Orgregistrocivil/Digibank or Orgregistropolicial/Magnetocorp"
   echo "For example to get the environment variables to set upa Orgregistropolicial shell environment run:  ./setOrgEnv.sh Orgregistropolicial"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh Orgregistropolicial | xargs)'
   exit 1
fi

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_ORG1_CA=${PEER0_ORG1_CA}"
echo "PEER0_ORG2_CA=${PEER0_ORG2_CA}"
echo "PEER0_ORG3_CA=${PEER0_ORG3_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"
