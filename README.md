# A Permissioned Blockchain Architecture for Ecuador’s National Public Data Registration System
> **Prototype of decentralized public registries infrastructure for the Ecuadorian context**.

## Project Description
This project proposes an architectural model designed to resolve the current fragmentation and lack of interoperability within Ecuador's National Public Data Registration System. To avoid the risks of a centralized model while preserving institutional autonomy, the system implements a permissioned blockchain network that distributes trust evenly among public entities. 
The functional prototype registers citizens' transactions using smart contracts and securely stores supporting documents off-chain.

## Technologies and architecture
The prototype was designed to simulate a multi-institutional scenario and is divided into the following main technologies:

* **Blockchain:** [Hyperledger Fabric v2.5.14] for the permissioned network, utilizing CouchDB databases for the world state, institutional peer nodes (Registro Civil, CNE, Contraloría), and an ordering service cluster.
* **Decentralized Storage:** [IPFS] to store supporting documents across a distributed cluster, generating cryptographic hashes (CIDs) to ensure immutability.
* **Application Layer:** Role-based Identity and Access Management (IAM) system using Attribute-Based Access Control (ABAC), alongside integration with Public Key Infrastructure (PKI).
* **Performance Evaluation:** [Hyperledger Caliper] used to measure latency and throughput under stress scenarios.

## Prerequisites and Environment
To reproduce the multi-institutional public registry scenario, the system must be evaluated in a containerized environment.
* **OS:** Ubuntu (WSL).
* **Requirements:** Docker, Docker Compose, Git, Curl, JQ, and Go.
* **Hardware:** Approximate resources of 8 vCPUs and 8 GB of RAM are recommended for stability during stress testing.

---

## Experiment Reproduction Guide

The setup and testing of the prototype are divided into 4 sequential phases. For detailed execution commands, please refer to the specific `.md` guides linked in each phase.

### Phase 1: Fabric Infrastructure Provisioning
The automated deployment of institutional nodes and the permissioned network is performed through a series of bash scripts. Run them in this exact order:

1. [`setup-fabric.sh`](./test-network/scripts_setup/): Installs Fabric binaries, Docker images, and dependencies.
2. [`rename_directories.sh`](./test-network/scripts_setup/): Customizes the default Fabric directories to our specific institutions (Registro Civil, CNE, Contraloría).
3. [`change_namesorgs.sh`](./test-network/scripts_setup/): Replaces organization names internally across all configuration files.
4. [`configuration_yaml_files.sh`](./test-network/scripts_setup/): Generates the `.yaml` Docker Compose files for the new organizations (e.g., Contraloría).
5. [`configuration_ccpgenerates_files.sh`](./test-network/scripts_setup/): Configures the connection profiles (`ccp-generate.sh`) for the application layer.
> **Detailed instructions:** See [`Phase_1-setup_fabric.md`](./Guides/Phase_1-setup_fabric.md)

### Phase 2: Network Bootstrapping & Identity Management (IAM)
Once the files are provisioned, the network must be started and the ABAC users created:

1. **Start the network:** Bring up the `mychannel` with CouchDB and the base organizations (Registro Civil and CNE).
```
cd test-network
./network.sh up createChannel -ca -c mychannel -s couchdb
```
2. **Add Org3:** Execute the script to dynamically add the `Contraloría` organization to the channel. Note: [`addOrgcontraloria.sh`](./test-network/addorgcontraloria/)
```
cd addorgcontraloria
./addOrgcontraloria.sh up -c mychannel -ca -s couchdb
```
3. **Register & Enroll ABAC User:** Use the Fabric CA Client to register a new user (`oficinista_rc`) with specific ABAC roles (`role=OPERATOR:ecert`) and enroll their cryptographic certificates to generate the MSP folder.
```
cd ../
export PATH=${PWD}/../bin:$PATH
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/

# Register the user with specific ABAC role
fabric-ca-client register \
  --caname ca-orgregistrocivil \
  --id.name oficinista_rc \
  --id.secret oficinistaPW \
  --id.type client \
  --id.attrs 'role=OPERATOR:ecert' \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.pem

# Enroll to generate certificates
fabric-ca-client enroll \
  -u https://oficinista_rc:oficinistaPW@localhost:7054 \
  --caname ca-orgregistrocivil \
  -M ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp \
  --tls.certfiles ${PWD}/organizations/fabric-ca/orgregistrocivil/tls-cert.pem
```
> **Detailed instructions:** See the first half of [`Phase_2-how_use_scripts.md`](./Guides/Phase_2-how_use_scripts.md)

### Phase 3: Smart Contract (Chaincode) Deployment
To deploy the business logic across the consortium, you must package, install, and approve the chaincode for all 3 organizations:
**1. Package the chaincode:**
Change the directory to `fabric-samples` and create a folder named [`dtic_chaincode`](./dtic_chaincode). Then, run the following commands:
```
go mod init dtic_chaincode
go mod tidy
go mod vendor

cd ~/fabric-samples/test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

peer lifecycle chaincode package dtic.tar.gz --path ../dtic_chaincode --lang golang --label dtic_1.0

```
**2. Install, Approve, and Commit (Multi-Org):**
Because this is a multi-institutional network, you must export the specific environment variables (`CORE_PEER_LOCALMSPID`, `CORE_PEER_TLS_ROOTCERT_FILE`, etc.) for **each** organization (Registro Civil, CNE, Contraloría) before running the install and approve commands.
**CRITICAL STEP:** Due to the length of the environment variable exports required for the consensus mechanism, please copy and run the exact commands detailed in the **"Install chaincode"** section of the [`Phase_3-Run_test-network.md`](./Guides/Phase_3-Run_test-network.md)  and  [`set_up_the_network.md`](./Guides/set_up_the_network.md) guide.

**3. Test the Network (Invoke Transaction):**
Once the chaincode is committed, test the creation of a citizen identity by executing an invoke command:

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n dtic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgregistrocivil.example.com/peers/peer0.orgregistrocivil.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/orgcne.example.com/peers/peer0.orgcne.example.com/tls/ca.crt -c '{"Args":["Tx_RegisterIdentity", "1712345678", "Juan", "Perez", "1990-01-01", "HOMBRE", "MASCULINO", "QmHashDePruebaIPFS"]}'
```
**4. Decentralized Storage (IPFS) Setup**

Deploy the private IPFS node to store supporting documents off-chain.

```
# Generate private swarm key
echo -e "/key/swarm/psk/1.0.0/\n/base16/\n$(tr -dc 'a-f0-9' < /dev/urandom | head -c64)" > swarm.key

# Run IPFS Docker Container
docker run -d   --name ipfs_node   -p 4001:4001   -p 5001:5001   -p 8080:8080   -v $(pwd)/swarm.key:/data/ipfs/swarm.key   --entrypoint /bin/sh   ipfs/kubo:latest   -c 'ipfs init --profile=server; \
      ipfs config --json AutoConf.Enabled false; \
      ipfs config --json Bootstrap "[]"; \
      ipfs config --json DNS.Resolvers "{}"; \
      ipfs config --json Routing.DelegatedRouters "[]"; \
      ipfs config --json Ipns.DelegatedPublishers "[]"; \
      ipfs daemon'
```
>  **Detailed instructions:** See the IPFS section in [`Phase_3-Run_test-network.md`](./Guides/Phase_3-Run_test-network.md) and  [`set_up_the_network.md`](./Guides/set_up_the_network.md)

### Phase 4: Performance Benchmarking (Hyperledger Caliper)
To measure the transaction throughput and latency under stress:

1. Set up the `caliper_metrics` directory.
2. Configure `workload.js` to simulate random citizen identity registrations matching the exact arguments of the Smart Contract.
3. Define the stress scenarios in `benchconfig.yaml` (from 10 TPS up to 200 TPS).
4. Run [`generate_network.sh`](./caliper_metrics/) to dynamically fetch the `oficinista_rc` ABAC private keys and certificates and inject them into Caliper's `network.yaml`.
5. Execute the Caliper benchmark to generate the performance reports.
> **Detailed instructions:** See [`Phase_4-metrics.md`](./Guides/Phase_4-metrics.md) and [`How_to_use_scripts_metrics_caliper.md`](./Guides/How_to_use_scripts_metrics_caliper.md) guide

---
## test-network Directory Structure

The test-network directory serves as the core for deploying the project's blockchain infrastructure. Its primary purpose is to provide an automated, local environment, based on Docker containers, to simulate the consortium. It is divided into components: The compose folder stores the YAML files that define and launch the network nodes (peers, orderers, and CouchDB databases). The configtx directory contains the configuration rules for generating the genesis block and establishing channel policies. Identity management is handled by the organizations directory, which safeguards the cryptographic material (certificates and keys) generated by the Certification Authorities to authenticate the Civil Registry and the National Electoral Council (CNE). Automation is managed through scripts and scripts_setup, which execute network startup, channel creation, and smart contract packaging. Finally, the addorgcontraloria directory includes the files necessary to scale the system, allowing the integration of a third entity into the blockchain in real time without interrupting operations.
---

## `test-network` Directory Structure
The core Fabric network environment is organized into the following key subdirectories to separate configurations, cryptographic materials, and automations:

```text
test-network/
├── addorgcontraloria/    # Scripts and configs to dynamically add Org3 (Contraloría) to the channel
├── compose/              # Docker compose YAML files for base nodes (CouchDB, CAs, Peers, Orderer)
├── configtx/             # Configuration for the genesis block and channel creation
├── organizations/        # Cryptographic materials (MSP, TLS) generated by Fabric CA
├── scripts/              # Core Fabric scripts to manage the network lifecycle
└── scripts_setup/        # Custom automated bash scripts for Phase 1 infrastructure provisioning
```

## Authors and Research
This project is the result of research carried out at the **Department of Computer Science of the National Polytechnic School** (Quito, Ecuador):
* **Enrique Mafla** - `enrique.mafla@epn.edu.ec` 
* **Christian Pérez** - `christian.perez01@epn.edu.ec`
* **Jeremmy Perugachi** - `jeremmy.perugachi@epn.edu.ec`
