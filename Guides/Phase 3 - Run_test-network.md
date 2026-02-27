# 

Inside the test-network folder, run:

```
./network.sh up createChannel -ca -c mychannel -s couchdb

```
This will bring up the network with the 2 organizations: registrocivil and cne.

Then, go to the addorgcontraloria folder and run:

```
./addOrgcontraloria.sh up -c mychannel -ca -s couchd
```

This command will add the contraloria organization to the network.

## Create user with role
Now return to the test-network folder and create a user named oficinista_rc by running:

```
export PATH=${PWD}/../bin:$PATH
```

```
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/
# Register the user
fabric-ca-client register \
  --caname ca-orgregistrocivil \
  --id.name oficinista_rc \
  --id.secret oficinistaPW \
  --id.type client \
  --id.attrs 'role=OPERATOR:ecert' \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.pem
```

```
# Generate certificates (Enroll)
fabric-ca-client enroll \
  -u https://oficinista_rc:oficinistaPW@localhost:7054 \
  --caname ca-orgregistrocivil \
  -M ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.p
```

## Chaincode

### Create chaincode

Change directory to fabric-samples and create a directory named dtic_chaincode

There, create the 4 necessary files, then run the following commands:


```
go mod init dtic_chaincode
go mod tidy
go mod vendor 
```
### Install chaincode

Now return to the test-network folder 

```
cd ~/fabric-samples/test-network
```

First, add the binaries to the path:

```
export PATH=${PWD}/../bin:$PATH
```

Also add the path of the base configuration:

```
export FABRIC_CFG_PATH=$PWD/../config/
```

Now package the chaincode:
```
peer lifecycle chaincode package dtic.tar.gz --path ../dtic_chaincode --lang golang --label dtic_1.0
```

Load the variables for registrocivil:

```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="OrgregistrocivilMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/Admin@orgregistrocivil.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

```

Install the chaincode:

```
peer lifecycle chaincode install dtic.tar.gz
```
Upon completion, this command will provide a package ID, load it with: 

```
export PAQUETE_ID=<ID>
```
Now load the variables for cne:

```
export CORE_PEER_LOCALMSPID="OrgcneMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcne.example.com/users/Admin@orgcne.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

Install the chaincode: 

```
peer lifecycle chaincode install dtic.tar.gz
```

Load variables for contraloria:

```
export CORE_PEER_LOCALMSPID="OrgcontraloriaMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/users/Admin@orgcontraloria.example.com/msp
export CORE_PEER_ADDRESS=localhost:11051
```
Install the chaincode:

```
peer lifecycle chaincode install dtic.tar.gz
```

Next, define the variable for the Orderer certificate:

```
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

Now load the variables for cne again:

```
export CORE_PEER_LOCALMSPID="OrgcneMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcne.example.com/users/Admin@orgcne.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

Approve: 
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --package-id $PAQUETE_ID --sequence 1 --tls --cafile $ORDERER_CA
```

Switch now to registrocivil 

```
export CORE_PEER_LOCALMSPID="OrgregistrocivilMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/Admin@orgregistrocivil.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

```

And approve:

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --package-id $PAQUETE_ID --sequence 1 --tls --cafile $ORDERER_CA
```

Next, using the current identity, make the commit:

```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/ca.crt
```
Load the orderer variable:

```
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

With this, you have installed the chaincode and can test it, for example, by creating a citizen: 

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n dtic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt -c '{"Args":["Tx_RegisterIdentity", "1712345678", "Juan", "Perez", "1990-01-01", "HOMBRE", "MASCULINO", "QmHashDePruebaIPFS"]}'

```


### IPFS Node

For IPFS, first generate the swarm key with:

```
echo -e "/key/swarm/psk/1.0.0/\n/base16/\n$(tr -dc 'a-f0-9' < /dev/urandom | head -c64)" > swarm.key
```

Then start the container:

```
docker run -d   --name ipfs_node   -p 4001:4001   -p 5001:5001   -p 8080:8080   -v $(pwd)/swarm.key:/data/ipfs/swarm.key   --entrypoint /bin/sh   ipfs/kubo:latest   -c 'ipfs init --profile=server; \
      ipfs config --json AutoConf.Enabled false; \
      ipfs config --json Bootstrap "[]"; \
      ipfs config --json DNS.Resolvers "{}"; \
      ipfs config --json Routing.DelegatedRouters "[]"; \
      ipfs config --json Ipns.DelegatedPublishers "[]"; \
      ipfs daemon'


```