# Steps to set up the network
* Note: This guide helps you get started with the fabric network and IPFS.

The following command must be executed within the test-network folder

```
./network.sh up createChannel -ca -c mychannel -s couchdb

```
This will establish a network with the two organizations, "registrocivil" and "cne".

Next, we navigate to the addorgcontraloria folder and run:

```
./addOrgcontraloria.sh up -c mychannel -ca -s couchd
```

This command adds "contraloria" organization to the network.

## Create user with role
Now we return to the test-network folder and create a user called office_rc. To do this, we run:

```
export PATH=${PWD}/../bin:$PATH
```

```
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/
#Registrar al usuario
fabric-ca-client register \
  --caname ca-orgregistrocivil \
  --id.name oficinista_rc \
  --id.secret oficinistaPW \
  --id.type client \
  --id.attrs 'role=OPERATOR:ecert' \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.pem
```

```
#Generar los certificados (Enrolar)
fabric-ca-client enroll \
  -u https://oficinista_rc:oficinistaPW@localhost:7054 \
  --caname ca-orgregistrocivil \
  -M ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.p
```

## Chaincode

### Create chaincode


The chaincode will be created using the files from the dtic_chaincode folder of the repository. To do this, we execute the following commands

```
go mod init dtic_chaincode
go mod tidy
go mod vendor 
```
### Install chaincode

Now we return to the test-network folder

```
cd ~/fabric-samples/test-network
```

And first we add the binaries to the path

```
export PATH=${PWD}/../bin:$PATH
```

We also added the path to the base configuration

```
export FABRIC_CFG_PATH=$PWD/../config/
```

Now we package the chaincode

```
peer lifecycle chaincode package dtic.tar.gz --path ../dtic_chaincode --lang golang --label dtic_1.0
```

We load the variables from "registrocivil"

```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="OrgregistrocivilMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/Admin@orgregistrocivil.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

```

We installed the chaincode

```
peer lifecycle chaincode install dtic.tar.gz
```

At the end of this command, it provides us with a package ID, which we then load with:

```
export PAQUETE_ID=<ID>
```

Now we load the variables from the cne

```
export CORE_PEER_LOCALMSPID="OrgcneMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcne.example.com/users/Admin@orgcne.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

We installed the chaincode:

```
peer lifecycle chaincode install dtic.tar.gz
```

Now we load the variables from "contraloria"

```
export CORE_PEER_LOCALMSPID="OrgcontraloriaMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/users/Admin@orgcontraloria.example.com/msp
export CORE_PEER_ADDRESS=localhost:11051
```

We installed the chaincode:

```
peer lifecycle chaincode install dtic.tar.gz
```

The next step is to define the variable for the Orderer certificate.

```
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

Now we load the cne variables again:

```
export CORE_PEER_LOCALMSPID="OrgcneMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgcne.example.com/users/Admin@orgcne.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

We approve:

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --package-id $PAQUETE_ID --sequence 1 --tls --cafile $ORDERER_CA
```

We now switch to the "registrocivil"

```
export CORE_PEER_LOCALMSPID="OrgregistrocivilMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/Admin@orgregistrocivil.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

```

We approve:

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --package-id $PAQUETE_ID --sequence 1 --tls --cafile $ORDERER_CA
```

The next step is to use the current identity to make the commit

```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name dtic --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcontraloria.example.com/peers/peer0.orgcontraloria.example.com/tls/ca.crt
```

We load the orderer variable

```
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

With this, we now have the chaincode installed and can run tests, for example, creating a citizen:

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n dtic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt -c '{"Args":["Tx_RegisterIdentity", "1712345678", "Juan", "Perez", "1990-01-01", "MAN", "MALE", "QmHashDePruebaIPFS"]}'

```


### IPFS node

For IPFS, we first generate the swarm key with:

```
echo -e "/key/swarm/psk/1.0.0/\n/base16/\n$(tr -dc 'a-f0-9' < /dev/urandom | head -c64)" > swarm.key
```

And now we start the container:

```
docker run -d   --name ipfs_node   -p 4001:4001   -p 5001:5001   -p 8080:8080   -v $(pwd)/swarm.key:/data/ipfs/swarm.key   --entrypoint /bin/sh   ipfs/kubo:latest   -c 'ipfs init --profile=server; \
      ipfs config --json AutoConf.Enabled false; \
      ipfs config --json Bootstrap "[]"; \
      ipfs config --json DNS.Resolvers "{}"; \
      ipfs config --json Routing.DelegatedRouters "[]"; \
      ipfs config --json Ipns.DelegatedPublishers "[]"; \
      ipfs daemon'


```