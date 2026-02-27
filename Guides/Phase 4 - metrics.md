# Hyperledger caliper


To measure latency and throughput, the Caliper tool will be used. For this, create a folder in the Fabric directory:

```
mkdir caliper_metrics
```

Inside, create the file workload.js:

```
'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegistroWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async submitTransaction() {
        const cedulaAleatoria = Math.floor(Math.random() * 10000000000).toString();
        
        // THE 9 EXACT ARGUMENTS OF THE NEW SMART CONTRACT
        const txArgs = [
            cedulaAleatoria,            // nationalID
            "Juan",                     // firstNames
            "Perez",                    // lastNames
            "1990-01-01",               // birthDate
            "QUITO",                    // birthPlace (NEW)
            "MALE",                     // sex (NEW - ENGLISH)
            "MASCULINE",                // gender (NEW - ENGLISH)
            "QmHashCertificado",        // initialCivilRegistryDirCID
            "QmHashRaizIPFS"            // initialRootCID
        ];

        const request = {
            contractId: "dtic",
            contractFunction: "Tx_RegisterIdentity", // NAME IN ENGLISH
            invokerIdentity: "oficinista_abac",    // MATCHES network.yaml
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

This file is used to send transactions using the chaincode parameters.

Next, create the benchconfig.yaml file:

```
test:
  name: "Thesis Saturation Test - ABAC"
  description: "Measuring the impact of ABAC up to 200 TPS"
  workers:
    number: 1 # <--- Crucial for stability
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
    - label: "Round 4 - Stress 150 TPS"
      txDuration: 30
      startDelay: 15 # <--- Time for the Peer to recover
      rateControl: { type: fixed-rate, opts: { tps: 150 } }
      workload: { module: workload.js }
    - label: "Round 5 - Saturation 200 TPS"
      txDuration: 30
      startDelay: 20 # <--- More recovery time
      rateControl: { type: fixed-rate, opts: { tps: 200 } }
      workload: { module: workload.js }
```

This is the configuration for the test and its phases.

Then, create the generate_network.sh script:

```
#!/bin/bash

echo "Searching for dynamic certificates of the Clerk (ABAC) of the Civil Registry..."


BASE_DIR="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/users/oficinista_rc@orgregistrocivil.example.com/msp"
CONN_PROFILE="../test-network/organizations/peerOrganizations/orgregistrocivil.example.com/connection-orgregistrocivil.json"

PRIV_KEY=$(find $BASE_DIR/keystore -name "*_sk" -print -quit)
CERT=$(find $BASE_DIR/signcerts -name "*.pem" -print -quit)

if [ -z "$PRIV_KEY" ]; then
    echo " ERROR: Clerk's private key not found. Did you enroll?"
    exit 1
fi

cat << EOF > network.yaml
name: Government Network Caliper ABAC
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

echo " network.yaml file generated with ABAC identity!"
```

This script creates the network.yaml file and fetches the necessary keys.