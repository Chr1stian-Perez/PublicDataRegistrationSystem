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

## Authors and research
This project is the result of research carried out at the **Department of Computer Science of the National Polytechnic School** (Quito, Ecuador):
* **Enrique Mafla** - `enrique.mafla@epn.edu.ec` 
* **Christian Pérez** - `christian.perez01@epn.edu.ec`
* **Jeremmy Perugachi** `jeremmy.perugachi@epn.edu.ec` 
