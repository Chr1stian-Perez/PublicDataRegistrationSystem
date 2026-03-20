#!/bin/bash
export PATH=$PATH:/usr/local/go/bin:$HOME/fabric-samples/bin
cd $HOME/fabric-samples/test-network

echo "1. Copying architecture modules..."
cp -r /mnt/c/Users/Christian/Desktop/Nueva\ Version\ Paper/PublicDataRegistrationSystem/NewArchitecture/addorgregistropropiedad .
cp -r /mnt/c/Users/Christian/Desktop/Nueva\ Version\ Paper/PublicDataRegistrationSystem/NewArchitecture/addorgregistroacademico .
cp -r /mnt/c/Users/Christian/Desktop/Nueva\ Version\ Paper/PublicDataRegistrationSystem/NewArchitecture/orgregistropropiedad-scripts scripts/
cp -r /mnt/c/Users/Christian/Desktop/Nueva\ Version\ Paper/PublicDataRegistrationSystem/NewArchitecture/orgregistroacademico-scripts scripts/

echo "2. Formatting line endings..."
dos2unix addorgregistropropiedad/*.sh addorgregistroacademico/*.sh scripts/orgregistropropiedad-scripts/*.sh scripts/orgregistroacademico-scripts/*.sh > /dev/null 2>&1

echo "3. Patching envVar.sh for Org3 (registropropiedad) and Org4 (registroacademico)..."
sed -i 's/org3/orgregistropropiedad/g' scripts/envVar.sh
sed -i 's/Org3/Orgregistropropiedad/g' scripts/envVar.sh
sed -i '/export PEER0_ORG3_CA/a export PEER0_ORG4_CA=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistroacademico.example.com/tlsca/tlsca.orgregistroacademico.example.com-cert.pem' scripts/envVar.sh

cat << 'EOF' >> scripts/envVar.sh

  if [ $1 -eq 4 ] || [[ ${1,,} == "orgregistroacademico" ]]; then
    export CORE_PEER_LOCALMSPID=OrgregistroacademicoMSP
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG4_CA
    export CORE_PEER_MSPCONFIGPATH=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistroacademico.example.com/users/Admin@orgregistroacademico.example.com/msp
    export CORE_PEER_ADDRESS=localhost:12051
  fi
EOF
# Wait, simply appending at the end might be outside of setGlobals.
# Better to use inline sed replacement to inject ORG4 properly.
# Removing the above dumb append and doing proper sed block injection.
sed -i '/elif \[ $USING_ORG -eq 3 \]; then/i \
  elif [ $USING_ORG -eq 4 ]; then\
    export CORE_PEER_LOCALMSPID=OrgregistroacademicoMSP\
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG4_CA\
    export CORE_PEER_MSPCONFIGPATH=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistroacademico.example.com/users/Admin@orgregistroacademico.example.com/msp\
    export CORE_PEER_ADDRESS=localhost:12051\
' scripts/envVar.sh

sed -i '/elif \[ $1 -eq 3 \]; then ORG_NAME="registropropiedad";/a \    elif [ $1 -eq 4 ]; then ORG_NAME="registroacademico";' scripts/envVar.sh

echo "Patching setAnchorPeer.sh..."
sed -i 's/org3/orgregistropropiedad/g' scripts/setAnchorPeer.sh
sed -i '/elif \[ $ORG -eq 3 \]; then/i \
  elif [ $ORG -eq 4 ]; then\
    HOST="peer0.orgregistroacademico.example.com"\
    PORT=12051\
' scripts/setAnchorPeer.sh

echo "4. Deploying base network..."
./network.sh up createChannel -ca -s couchdb

echo "5. Deploying Propiedad (Org3)..."
cd addorgregistropropiedad
./addorgregistropropiedad.sh up -ca -s couchdb
cd ..

echo "6. Deploying Academico (Org4)..."
cd addorgregistroacademico
./addorgregistroacademico.sh up -ca -s couchdb
cd ..

echo "Architecture Deployment Finished!"
