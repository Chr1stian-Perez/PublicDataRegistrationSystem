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
