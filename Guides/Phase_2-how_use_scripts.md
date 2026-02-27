
### Script to setup fabric
```
#!/bin/bash

echo "=========================================================="
echo "Preparing the environment for Hyperledger Fabric"
echo "=========================================================="

# 1. Update package list
echo "[1/6] Updating repositories..."
sudo apt-get update

# 2. Install base tools and dependencies
echo "[2/6] Installing Git, Curl, JQ, and Go..."
sudo apt-get install -y git curl jq golang-go

# 3. Install Docker and Docker Compose
echo "[3/6] Installing Docker Compose..."
sudo apt-get -y install docker-compose

# 4. Configure Docker daemon
echo "[4/6] Starting and enabling Docker service..."
sudo systemctl start docker
sudo systemctl enable docker

# 5. Configure user permissions for Docker
# The $USER variable automatically captures your current username
echo "[5/6] Adding user '$USER' to the docker group..."
sudo usermod -a -G docker $USER

# 6. Download and install Hyperledger Fabric
echo "[6/6] Downloading Fabric script and installing version 2.5.14..."
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh

# Download all (docker images, samples, binaries) specifically for v2.5.14
./install-fabric.sh --fabric-version 2.5.14 docker samples binary

echo "=========================================================="
echo "Installation completed successfully!"
echo "=========================================================="
echo "IMPORTANT: For Docker permissions to take effect without rebooting, run this command right now:"
echo "su - $USER"
```

### Script to change directory names
inside the fabric-samples folder
script name: rename_directories.sh
```
#!/bin/bash

echo "=========================================================="
echo "Running directory and file renaming"
echo "=========================================================="


find . -iname "*org*" | sort -r | while read -r FULL_PATH; do

    DIR=$(dirname "$FULL_PATH")
    OLD_NAME=$(basename "$FULL_PATH")


    NEW_NAME="$OLD_NAME"
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgregistrocivil/orgregistrocivil/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgcne/orgcne/gI')
    NEW_NAME=$(echo "$NEW_NAME" | sed 's/orgcontraloria/orgcontraloria/gI')

    if [ "$OLD_NAME" != "$NEW_NAME" ]; then
        # Execute the move
        mv "$FULL_PATH" "$DIR/$NEW_NAME"
        echo "LOG: $OLD_NAME changed to $NEW_NAME"
    fi
done

echo "=========================================================="
echo "Renaming completed successfully!"
echo "=========================================================="
```


### Script to rename the names of organizations in the files
inside the fabric-samples folder
script name: change_namesorgs.sh
```
#!/bin/bash

echo "=========================================================="
echo "Updating text in files"
echo "=========================================================="


find . -type f \( -name "*.yaml" -o -name "*.sh" -o -name "*.json" \) | while read -r FILE; do

    
    if grep -qiE "orgregistrocivil|orgcne|orgcontraloria" "$FILE"; then
        echo "Modifying: $FILE"

        # Apply changes directly (-i) preserving capitalization
        sed -i 's/Orgregistrocivil/Orgregistrocivil/g; s/orgregistrocivil/orgregistrocivil/g' "$FILE"
        sed -i 's/Orgcne/Orgcne/g; s/orgcne/orgcne/g' "$FILE"
        sed -i 's/Orgcontraloria/Orgcontraloria/g; s/orgcontraloria/orgcontraloria/g' "$FILE"
    fi
done

echo "=========================================================="
echo "Content updated"
echo "=========================================================="
```


### Script to configure the .yaml files of the organization being added


script name: configuration_yaml_files.sh
```
#!/bin/bash

# 1. File: docker-compose-ca-orgcontraloria.yaml
echo "Updating CA YAML..."
cat << 'EOF' > ${PWD}/fabric-samples/test-network/addorgcontraloria/compose/docker/docker-compose-ca-orgcontraloria.yaml
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '3.7'

networks:
  test:
    name: fabric_test

services:
  ca_orgcontraloria:
    image: hyperledger/fabric-ca:latest
    labels:
      service: hyperledger-fabric
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca-orgcontraloria
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_PORT=11054
    ports:
      - "11054:11054"
    command: sh -c 'fabric-ca-server start -b admin:adminpw -d'
    volumes:
      - ../fabric-ca/orgcontraloria:/etc/hyperledger/fabric-ca-server
EOF

# 2.File: docker-compose-couch-orgcontraloria.yaml
echo "Updating CouchDB YAML..."
cat << 'EOF' > ${PWD}/fabric-samples/test-network/addorgcontraloria/compose/docker/docker-compose-couch-orgcontraloria.yaml
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

  peer0.orgcontraloria.example.com:
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

# 3. File: docker-compose-orgcontraloria.yaml
echo "Updating Peer YAML..."
cat << 'EOF' > ${PWD}/fabric-samples/test-network/addorgcontraloria/compose/docker/docker-compose-orgcontraloria.yaml
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '3.7'

volumes:
  peer0.orgcontraloria.example.com:

networks:
  test:
    name: fabric_test

services:

  peer0.orgcontraloria.example.com:
    container_name: peer0.orgcontraloria.example.com
    privileged: true
    image: hyperledger/fabric-peer:2.5
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
      - CORE_PEER_ID=peer0.orgcontraloria.example.com
      - CORE_PEER_ADDRESS=peer0.orgcontraloria.example.com:11051
      - CORE_PEER_LISTENADDRESS=0.0.0.0:11051
      - CORE_PEER_CHAINCODEADDRESS=peer0.orgcontraloria.example.com:11052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:11052
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.orgcontraloria.example.com:11051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.orgcontraloria.example.com:11051
      - CORE_PEER_LOCALMSPID=OrgcontraloriaMSP
      - CORE_CHAINCODE_EXECUTETIMEOUT=300s
      - CORE_CHAINCODE_DEPLOYTIMEOUT=300s
      - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ${PWD}/fabric-samples/config:/etc/hyperledger/fabric
      - ${PWD}/fabric-samples/test-network/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/msp:/etc/hyperledger/fabric/msp
      - ${PWD}/fabric-samples/test-network/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls:/etc/hyperledger/fabric/tls
      - peer0.orgcontraloria.example.com:/var/hyperledger/production
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
```




### Script to configure the ccp-generate.sh files

script name: configuration_ccpgenerates_files.sh
```
#!/bin/bash

echo "=========================================================="
echo "Updating ORG assignments in ccp-generate.sh"
echo "=========================================================="

# 1. addorgcontraloria file (Change ORG=3 to ORG=contraloria)
# This path corresponds to the Comptroller's infrastructure you are designing.
FILE1="${PWD}/fabric-samples/test-network/addorgcontraloria/ccp-generate.sh"
if [ -f "$FILE1" ]; then
    sed -i 's/ORG=3/ORG=contraloria/g' "$FILE1"
    echo "✔ Updated ORG=contraloria in: $FILE1"
else
    echo "✘ Error: Not found $FILE1"
fi

# 2. organizations file (Change ORG=1 to registrocivil and ORG=2 to cne)
# These changes reflect the entities of the public records system.
FILE2="${PWD}/fabric-samples/test-network/organizations/ccp-generate.sh"
if [ -f "$FILE2" ]; then
    sed -i 's/ORG=1/ORG=registrocivil/g' "$FILE2"
    sed -i 's/ORG=2/ORG=cne/g' "$FILE2"
    echo "✔ Updated ORG=registrocivil and ORG=cne in: $FILE2"
else
    echo "✘ Error: Not found $FILE2"
fi

echo "=========================================================="
echo "Changes successfully implemented!"
echo "=========================================================="
```
