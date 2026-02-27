# Hyperledger caliper


* Note: All the scripts used in this documentation are located in the caliper_metrics folder of the repository.

To measure latency and throughput, the Hyperledger Caliper tool will be utilized. For this purpose, a dedicated directory was created within the Fabric folder:

```
mkdir caliper_metrics
cd caliper_metrics
```
The setup process begins with installing the essential dependencies. Specifically, the following environment requirements must be met before running the performance test:

```
sudo apt update
sudo apt install nodejs npm -y

npm install @hyperledger/caliper-cli@0.5.0
npx caliper bind --caliper-bind-sut fabric:2.4
```

To simulate network traffic, a custom workload module named workload.js is implemented. This script defines the transaction logic and submission rate for the benchmarking process

## workload.js:
```
'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegistroWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async submitTransaction() {
        const randomCI = Math.floor(Math.random() * 10000000000).toString();
        
        
        const txArgs = [
            randomCI,            // nationalID
            "Juan",                     // firstNames
            "Perez",                    // lastNames
            "1990-01-01",               // birthDate
            "QUITO",                    // birthPlace (NUEVO)
            "MALE",                     // sex (NUEVO - INGLÉS)
            "MASCULINE",                // gender (NUEVO - INGLÉS)
            "QmHashCertificado",        // initialCivilRegistryDirCID
            "QmHashRaizIPFS"            // initialRootCID
        ];

        const request = {
            contractId: "dtic",
            contractFunction: "Tx_RegisterIdentity", // NOMBRE EN INGLÉS
            invokerIdentity: "oficinista_abac",    // COINCIDE CON network.yaml
            contractArguments: txArgs,
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new RegistroWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
```

To orchestrate the execution of the tests, the benchconfig.yaml file is established. This configuration file specifies the benchmark rounds, the submission rate, and the duration of the performance evaluation.
## benchconfig.yaml:
```
test:
  name: "Network saturation test"
  description: "Measuring the impact up to 200 TPS"
  workers:
    number: 1 
  rounds:
    - label: "Round 1 - 10 TPS"
      txDuration: 30
      rateControl: { type: fixed-rate, opts: { tps: 10 } }
      workload: { module: workload.js }
    - label: "Round 2 - 50 TPS"
      txDuration: 30
      rateControl: { type: fixed-rate, opts: { tps: 50 } }
      workload: { module: workload.js }
    - label: "Round 3 - 100 TPS"
      txDuration: 30
      rateControl: { type: fixed-rate, opts: { tps: 100 } }
      workload: { module: workload.js }
    - label: "Round 4 - Estrés 150 TPS"
      txDuration: 30
      startDelay: 15
      rateControl: { type: fixed-rate, opts: { tps: 150 } }
      workload: { module: workload.js }
    - label: "Round 5 - Saturación 200 TPS"
      txDuration: 30
      startDelay: 20 
      rateControl: { type: fixed-rate, opts: { tps: 200 } }
      workload: { module: workload.js }
```

To automate the integration of cryptographic material, the generate_network.sh script is implemented. This script dynamically retrieves the Private Keys and Certificates and populates the network.yaml configuration file.

## generate_network.sh:
```
#!/bin/bash

echo "Looking for dynamic certificates from the Civil Registry clerk..."


BASE_DIR="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp"
CONN_PROFILE="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/connection-orgregistrocivil.json"

PRIV_KEY=$(find $BASE_DIR/keystore -name "*_sk" -print -quit)
CERT=$(find $BASE_DIR/signcerts -name "*.pem" -print -quit)

if [ -z "$PRIV_KEY" ]; then
    echo " ERROR: The office worker's private key was not found. Did you enroll?"
    exit 1
fi

cat << EOF > network.yaml
name: Red Gubernamental Caliper ABAC
version: "2.0.0"
caliper:
  blockchain: fabric
channels:
  - channelName: mychannel
    contracts:
    - id: dtic
organizations:
  - mspid: OrgregistrocivilMSP
    identities:
      certificates:
      - name: 'oficinista_abac'
        clientPrivateKey:
          path: $PRIV_KEY
        clientSignedCert:
          path: $CERT
    connectionProfile:
      path: $CONN_PROFILE
      discover: true
EOF

echo " "Network.yaml file generated!"

echo " ¡Archivo network.yaml generado con identidad ABAC!"
```
