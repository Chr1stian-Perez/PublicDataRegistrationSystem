#!/bin/bash

# 1. File: docker-compose-ca-orgregistropropiedad.yaml
echo "Updating CA YAML..."
cat << 'EOF' > $HOME/fabric-samples/test-network/addorgregistropropiedad/compose/docker/docker-compose-ca-orgregistropropiedad.yaml
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '3.7'

networks:
  test:
    name: fabric_test

services:
  ca_orgregistropropiedad:
    image: hyperledger/fabric-ca:latest
    labels:
      service: hyperledger-fabric
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca-orgregistropropiedad
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_PORT=11054
    ports:
      - "11054:11054"
    command: sh -c 'fabric-ca-server start -b admin:adminpw -d'
    volumes:
      - ../fabric-ca/orgregistropropiedad:/etc/hyperledger/fabric-ca-server
EOF

# 2.File: docker-compose-couch-orgregistropropiedad.yaml
echo "Updating CouchDB YAML..."
cat << 'EOF' > $HOME/fabric-samples/test-network/addorgregistropropiedad/compose/docker/docker-compose-couch-orgregistropropiedad.yaml
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '3.7'

networks:
  test:
    name: fabric_test

services:
  couchdb4:
    container_name: couchdb4
    image: couchdb:3.1.1
    labels:
      service: hyperledger-fabric
    # Populate the COUCHDB_USER and COUCHDB_PASSWORD to set an admin user and password
    # for CouchDB.  This will prevent CouchDB from operating in an "Admin Party" mode.
    environment:
      - COUCHDB_USER=admin
      - COUCHDB_PASSWORD=adminpw
    # Comment/Uncomment the port mapping if you want to hide/expose the CouchDB service,
    # for example map it to utilize Fauxton User Interface in dev environments.
    ports:
      - "9984:5984"
    networks:
      - test

  peer0.orgregistropropiedad.example.com:
    environment:
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb4:5984
      # The CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME and CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD
      # provide the credentials for ledger to connect to CouchDB.  The username and password must
      # match the username and password set for the associated CouchDB.
      - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
      - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw
    depends_on:
      - couchdb4
    networks:
      - test
EOF

# 3. File: docker-compose-orgregistropropiedad.yaml
echo "Updating Peer YAML..."
cat << 'EOF' > $HOME/fabric-samples/test-network/addorgregistropropiedad/compose/docker/docker-compose-orgregistropropiedad.yaml
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '3.7'

volumes:
  peer0.orgregistropropiedad.example.com:

networks:
  test:
    name: fabric_test

services:

  peer0.orgregistropropiedad.example.com:
    container_name: peer0.orgregistropropiedad.example.com
    privileged: true
    image: hyperledger/fabric-peer:latest
    labels:
      service: hyperledger-fabric
    environment:
      #Generic peer variables
      - FABRIC_CFG_PATH=/etc/hyperledger/fabric
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_test
      - FABRIC_LOGGING_SPEC=INFO
      #- FABRIC_LOGGING_SPEC=DEBUG
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      # Peer specific variables
      - CORE_PEER_ID=peer0.orgregistropropiedad.example.com
      - CORE_PEER_ADDRESS=peer0.orgregistropropiedad.example.com:11051
      - CORE_PEER_LISTENADDRESS=0.0.0.0:11051
      - CORE_PEER_CHAINCODEADDRESS=peer0.orgregistropropiedad.example.com:11052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:11052
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.orgregistropropiedad.example.com:11051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.orgregistropropiedad.example.com:11051
      - CORE_PEER_LOCALMSPID=OrgregistropropiedadMSP
      - CORE_CHAINCODE_EXECUTETIMEOUT=300s
      - CORE_CHAINCODE_DEPLOYTIMEOUT=300s
      - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - $HOME/fabric-samples/config:/etc/hyperledger/fabric
      - $HOME/fabric-samples/test-network/organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/msp:/etc/hyperledger/fabric/msp
      - $HOME/fabric-samples/test-network/organizations/peerOrganizations/orgregistropropiedad.example.com/peers/peer0.orgregistropropiedad.example.com/tls:/etc/hyperledger/fabric/tls
      - peer0.orgregistropropiedad.example.com:/var/hyperledger/production
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    ports:
      - 11051:11051
    networks:
      - test
EOF

echo "=========================================================="
echo "YAML files updated successfully."
echo "=========================================================="
