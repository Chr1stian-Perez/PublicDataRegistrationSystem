# A Permissioned Blockchain Architecture for Ecuador’s National Public Data Registration System
> **Prototype of decentralized public registries infrastructure for the Ecuadorian context**.

## Project description
This project proposes an architectural model designed to resolve the current fragmentation and lack of interoperability within Ecuador's National Public Data Registration System. To avoid the risks of a centralized model while preserving institutional autonomy, the system implements a permissioned blockchain network that distributes trust evenly among public entities. 
The functional prototype registers citizens' transactions using smart contracts and securely stores supporting documents off-chain.

The system models four registries: Civil Registry (`orgregistrocivil`), Police Registry (`orgregistropolicial`), Property Registry (`orgregistropropiedad`), and Academic Registry (`orgregistroacademico`). Each institution operates an independent peer node, maintains its own Certificate Authority, and participates in transaction endorsement under a `MAJORITY` policy. Supporting documents are stored off-chain in a private IPFS node, and their cryptographic hashes (CIDs) are anchored to the blockchain to guarantee immutability.

## Technologies

The prototype integrates three layers. The blockchain layer runs on **Hyperledger Fabric v2.5.14** with **CouchDB v3.1.1** as the world state database, enabling rich JSON queries over citizen records; the ordering service uses the Raft consensus protocol with a single orderer node. The storage layer uses **IPFS Kubo v0.33** to host supporting documents — birth certificates, property deeds, academic diplomas — whose CIDs are registered on-chain alongside the structured data. The application layer is a **Go v1.26** CLI that connects to peers through the **Fabric Gateway SDK v1.7.1**, authenticates users through X.509 certificates issued by **Fabric CA v1.5.12**, and enforces Role-Based Access Control (RBAC) with four roles: `OPERATOR`, `DIRECTOR`, `REGISTRAR`, and `AUDITOR`.

## Repository structure

```
PublicDataRegistrationSystem/
├── app/              ← Go client application (main.go, go.mod, go.sum)
├── chaincode/        ← Smart contract (contract.go, models.go, utils.go)
├── test-network/     ← Full network infrastructure (see below)
├── guides/           ← Step-by-step deployment documentation
├── README.md
└── LICENSE
```

The `test-network/` directory contains everything needed to bring up the four-organization network:

```
test-network/
├── network.sh                   ← Master Fabric bootstrap script
├── deploy_scalable.sh           ← Full 4-org orchestration script
├── gen_identities.sh            ← RBAC identity generator for all orgs
├── addorgregistropropiedad/     ← Scripts to dynamically add Property Registry
├── addorgregistroacademico/     ← Scripts to dynamically add Academic Registry
├── orgregistropropiedad-scripts/← joinChannel + anchor peer config for Org3
├── orgregistroacademico-scripts/← joinChannel + anchor peer config for Org4
├── scripts/                     ← Core Fabric lifecycle scripts
├── compose/                     ← Docker Compose files for all nodes
├── configtx/                    ← Genesis block and channel policy config
└── organizations/               ← CA crypto material (generated at runtime)
```

---

## Prerequisites

The environment must run Linux or WSL (Ubuntu 20.04 tested). Go 1.26 is required for both the chaincode and the application. Docker and Docker Compose must be available with enough resources to run eight containers simultaneously — approximately 8 vCPUs and 8 GB of RAM are recommended for stable operation under concurrent load.

Hyperledger Fabric binaries (`peer`, `orderer`, `fabric-ca-client`) must be installed under `$HOME/fabric-samples/bin`. The official Fabric bootstrap script handles this:

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.14 1.5.12
```

After installation, add the binaries to your `PATH`:

```bash
export PATH=$HOME/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$HOME/fabric-samples/config
```
---
## Reproduction guide

### Step 1 — Clone and prepare the repository

Clone this repository into `$HOME/fabric-samples/` so that the relative paths used by the network scripts resolve correctly:

```bash
cd $HOME/fabric-samples
git clone https://github.com/Chr1stian-Perez/PublicDataRegistrationSystem.git
```

Copy the `chaincode/` and `app/` directories alongside the Fabric binaries so Go can resolve the module paths at build time:

```bash
cp -r PublicDataRegistrationSystem/chaincode $HOME/fabric-samples/dtic_chaincode
cp -r PublicDataRegistrationSystem/app $HOME/fabric-samples/app
```

### Step 2 — Configure file line endings

Scripts authored on Windows may carry CRLF line endings that cause `bash` to fail silently. Install `dos2unix` and normalize all shell scripts before execution:

```bash
sudo apt-get install -y dos2unix
find $HOME/fabric-samples/test-network -name "*.sh" -exec dos2unix {} \;
```

### Step 3 — Bootstrap the base network (Org1 + Org2)

The base network starts Civil Registry and Police Registry, creates `mychannel`, and brings up CouchDB-backed peers with Certificate Authorities enabled:

```bash
cd $HOME/fabric-samples/test-network
./network.sh up createChannel -ca -s couchdb
```

Verify all containers are running before proceeding:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

### Step 4 — Add property registry (Org3)

The `addorgregistropropiedad` scripts extend the running channel to include a third organization without interrupting operations. Execute from within the `test-network` directory:

```bash
cp -r PublicDataRegistrationSystem/test-network/addorgregistropropiedad .
cp -r PublicDataRegistrationSystem/test-network/orgregistropropiedad-scripts scripts/
dos2unix addorgregistropropiedad/*.sh scripts/orgregistropropiedad-scripts/*.sh

cd addorgregistropropiedad
./addorgregistropropiedad.sh up -ca -s couchdb
cd ..
```

### Step 5 — Add academic registry (Org4)

The same pattern applies for the fourth organization. The `addorgregistroacademico` scripts generate crypto material for the CA, launch the peer container, and submit a channel update transaction that was signed by the existing majority:

```bash
cp -r PublicDataRegistrationSystem/test-network/addorgregistroacademico .
cp -r PublicDataRegistrationSystem/test-network/orgregistroacademico-scripts scripts/
dos2unix addorgregistroacademico/*.sh scripts/orgregistroacademico-scripts/*.sh

cd addorgregistroacademico
./addorgregistroacademico.sh up -ca -s couchdb
cd ..
```

At this point `docker ps` should show eight peer/CA/CouchDB containers active.

### Step 6 — Generate RBAC identities

The chaincode enforces role-based access through X.509 certificate attributes. The `gen_identities.sh` script registers and enrolls four users per organization — `operator1`, `director1`, `registrar1`, and `auditor1` — each carrying the corresponding `role` attribute issued by their organization's CA:

```bash
cd $HOME/fabric-samples/test-network
chmod +x gen_identities.sh
./gen_identities.sh
```

### Step 7 — Deploy the chaincode

The smart contract must be packaged, installed on all four peers, approved by a majority, and committed to the channel.

To do this, within the app folder, we run the following commands:

1. First, we load the variables:
   
```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pe
```

2. We create the .tar.gz file:
```bash
peer lifecycle chaincode package dtic.tar.gz --path . --lang golang --label dtic_1.0
```

The following is from the test-network folder:

1. Prepare the Environment
```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pe
```
2. Install the chaincode in each institution:

```bash
for org in Orgregistrocivil Orgregistropolicial Orgregistropropiedad Orgregistroacademico; do
    export $(./setOrgEnv.sh $org | xargs)
    peer lifecycle chaincode install ../dtic_chaincode/dtic.tar.gz
done
```
To do this, we use the setOrgEnv.sh script to load the variables for each organization.

3. Next, we extract the package_id from the chaincode.
```bash
peer lifecycle chaincode queryinstalled
```
And load that into a variable:
```bash
export PACKAGE_ID=<Insert the ID from the previous command>
```
4. Next, we approve it for the four institutions. For this, we use the same setOrgEnv.sh script.

```bash
for org in Orgregistrocivil Orgregistropolicial Orgregistropropiedad Orgregistroacademico; do
    export $(./setOrgEnv.sh $org | xargs)
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" --channelID mychannel --name dtic --version 1.0 --package-id $PACKAGE_ID --sequence 1
done
```

5. Now we load the certificate variables of each organization and the orderer:

```bash
export TEST_NETWORK_HOME=$PWD
export PEER0_ORG1_CA=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistrocivil.example.com/tlsca/tlsca.orgregistrocivil.example.com-cert.pem
export PEER0_ORG2_CA=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistropolicial.example.com/tlsca/tlsca.orgregistropolicial.example.com-cert.pem
export PEER0_ORG3_CA=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistropropiedad.example.com/tlsca/tlsca.orgregistropropiedad.example.com-cert.pem
export PEER0_ORG4_CA=${TEST_NETWORK_HOME}/organizations/peerOrganizations/orgregistroacademico.example.com/tlsca/tlsca.orgregistroacademico.example.com-cert.pem
export ORDERER_CA=${TEST_NETWORK_HOME}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
```
6. Finally we do the commit:
```bash
peer lifecycle chaincode commit -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls --cafile "$ORDERER_CA" \
--channelID mychannel --name dtic --version 1.0 --sequence 1 \
--peerAddresses localhost:7051 --tlsRootCertFiles "$PEER0_ORG1_CA" \
--peerAddresses localhost:9051 --tlsRootCertFiles "$PEER0_ORG2_CA" \
--peerAddresses localhost:11051 --tlsRootCertFiles "$PEER0_ORG3_CA" \
--peerAddresses localhost:12051 --tlsRootCertFiles "$PEER0_ORG4_CA"

```

The deployment uses the `MAJORITY` endorsement policy by default.

### Step 8 — Start the IPFS node

Supporting documents are uploaded to a private IPFS node that runs alongside the Fabric network. The swarm key inside `test-network/` isolates this node from the public IPFS network:

```bash
docker run -d \
  --name ipfs_node \
  -p 4001:4001 -p 5001:5001 -p 8080:8080 \
  -v $HOME/fabric-samples/test-network/swarm.key:/data/ipfs/swarm.key \
  --entrypoint /bin/sh \
  ipfs/kubo:latest \
  -c 'ipfs init --profile=server && \
      ipfs config --json Bootstrap "[]" && \
      ipfs config --json Addresses.API "\"/ip4/0.0.0.0/tcp/5001\"" && \
      ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin "[\"*\"]" && \
      ipfs daemon'
```

Wait approximately 10 seconds for the daemon to initialize, then verify it is reachable:

```bash
curl -s http://localhost:5001/api/v0/id | jq .ID
```

### Step 9 — Build and run the application

The Go application connects to all four organizations through the Fabric Gateway SDK. Build and run it from its source directory:

```bash
cd $HOME/fabric-samples/app
go mod tidy
go run .
```

The interactive CLI presents an organization selector. Selecting an organization triggers a login flow that reads the available X.509 identities from the corresponding `users/` directory and connects to that organization's peer. All field input is normalized to uppercase before submission, ensuring consistent world state values regardless of how operators enter data.

---

## Tear down

To stop all containers and clean up the generated crypto material:

```bash
cd $HOME/fabric-samples/test-network
./network.sh down
docker stop ipfs_node && docker rm ipfs_node
```

This removes volumes, channels, and all generated certificates. Re-running from Step 3 restores a clean environment.

---

## Authors

Developed at the **Departamento de Ciencias de la Computación, Escuela Politécnica Nacional** (Quito, Ecuador):

- **Enrique Mafla** — `enrique.mafla@epn.edu.ec`
- **Christian Pérez** — `christian.perez01@epn.edu.ec`
- **Jeremmy Perugachi** — `jeremmy.perugachi@epn.edu.ec`
