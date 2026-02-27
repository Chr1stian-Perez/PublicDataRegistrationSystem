# A Permissioned Blockchain Architecture for Ecuador’s National Public Data Registration System
> **Prototype of decentralized public registries infrastructure for the Ecuadorian context**.

## Project description
This project proposes an architectural model to resolve the current fragmentation and lack of interoperability of Ecuador's National Public Data Registration System. To avoid the risks of a centralized model while preserving institutional autonomy, the system implements a permissioned blockchain network that distributes trust among public entities. 
The functional prototype registries citizens' transactions using smart contracts and securely stores supporting documents off-chain.

## Technologies and architecture
The prototype was designed to simulate a multi-institutional scenario and is divided into the following main technologies:

* **Blockchain:** [Hyperledger Fabric] for the permissioned network, using CouchDB databases for the world state, institutional peer nodes, and an ordering service cluster.
* **Decentralized storage:** [IPFS] (InterPlanetary File System) to store supporting documents in a distributed cluster, generating cryptographic hashes to ensure immutability.
* **Application layer:** REST API Gateway, role-based Identity and Access Management (IAM) system, and integration with Public Key Infrastructure (PKI).
* **Deployment and environment:** The system is containerized using [Docker] and secures communications through TLS channels. It was evaluated on an Ubuntu-based environment (WSL).

## Reproduction of the experiment and installation

To simulate the multi-institutional public registry scenario, the system was evaluated in a containerized environment using **Docker** on **Ubuntu (WSL)**. To reproduce the optimal environment, approximate hardware resources of 8 vCPUs and 8 GB of RAM are recommended.

Steps to set up the **Hyperledger Fabric** infrastructure.

> **Important:** If you want to see the detailed content of the installation scripts, please refer to our complete guide in the [`setup_fabric.md`](setup_fabric.md) file.

### Fabric network configuration

The automated deployment of institutional nodes and the permissioned network is performed through a series of bash scripts included in this repository:

1. **Initial setup:** Run the main Fabric installation script (`setup_fabric.sh`).
2. **Organization structure (within `fabric-samples`):** 
   * Run `rename_directories.sh` to change the directory names to those of the corresponding institutions.
   * Run `change_namesorgs.sh` to replace the names of the organizations internally in the files.
3. **Node configuration:** Run the `configuration_yaml_files.sh` script to adjust the `.yaml` files corresponding to the organization being added to the network.
4. **Connection profiles:** Finally, run `configuration_ccpgenerates_files.sh` to configure the `ccp-generate.sh` files, which enable the application layer to interact and connect with the blockchain network.

## Authors and research
This project is the result of research carried out at the **Department of Computer Science of the National Polytechnic School** (Quito, Ecuador):
* **Enrique Mafla** - `enrique.mafla@epn.edu.ec` 
* **Christian Pérez** - `christian.perez01@epn.edu.ec`
* **Jeremmy Perugachi** `jeremmy.perugachi@epn.edu.ec` 
