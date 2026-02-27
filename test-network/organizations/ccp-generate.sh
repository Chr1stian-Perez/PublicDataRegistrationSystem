#!/usr/bin/env bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG=registrocivil
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/orgregistrocivil.example.com/tlsca/tlsca.orgregistrocivil.example.com-cert.pem
CAPEM=organizations/peerOrganizations/orgregistrocivil.example.com/ca/ca.orgregistrocivil.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/orgregistrocivil.example.com/connection-orgregistrocivil.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/orgregistrocivil.example.com/connection-orgregistrocivil.yaml

ORG=cne
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/orgcne.example.com/tlsca/tlsca.orgcne.example.com-cert.pem
CAPEM=organizations/peerOrganizations/orgcne.example.com/ca/ca.orgcne.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/orgcne.example.com/connection-orgcne.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/orgcne.example.com/connection-orgcne.yaml
