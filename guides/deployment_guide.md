# Deployment reference guide

This document expands on the reproduction steps outlined in the main `README.md`, providing command-level detail, expected outputs, and common failure points encountered during development. Each section maps directly to a numbered step in the README.

---

## Step 1 — Repository layout and Fabric binaries

Hyperledger Fabric resolves many paths relative to `$HOME/fabric-samples`. The simplest approach is to clone this repository into that directory and work from there. Placing the chaincode and application inside `fabric-samples` ensures that `go mod tidy` and the peer lifecycle commands can locate all dependencies without additional `GOPATH` configuration.

If the Fabric binaries are not yet installed, the bootstrap script downloads them along with the Docker images for the exact versions used in development:

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.14 1.5.12
```

Confirm that `peer version` returns `2.5.x` before proceeding. Mismatched binaries are the most common source of lifecycle errors.

---

## Step 2 — Line ending normalization

Scripts transferred from Windows to Linux carry `\r\n` line endings that `bash` interprets as part of variable names, producing errors like `$'\r': command not found`. Running `dos2unix` recursively on the `test-network/` tree eliminates this before any script is executed:

```bash
find $HOME/fabric-samples/test-network -name "*.sh" -exec dos2unix {} \;
```

---

## Step 3 — Base network bootstrap

The `network.sh up` command starts the ordering service, the Certificate Authorities for Civil Registry and Police Registry, and their respective peers with CouchDB state databases. The `-ca` flag is mandatory — without it, the CA containers do not start and identity enrollment in Step 6 will fail.

```bash
cd $HOME/fabric-samples/test-network
./network.sh up createChannel -ca -s couchdb
```

**Expected:** Eight containers start (`orderer`, `peer0.orgregistrocivil`, `peer0.orgregistropolicial`, two CouchDB instances, two CA instances). Confirm with `docker ps`.

---

## Step 4 & 5 — Dynamic organization addition

Each addorg script internally calls `configtxlator` to decode the current channel configuration, appends the new organization's MSP definition, re-encodes it, and submits the update transaction signed by the existing majority. The operation is non-disruptive — active peers continue serving requests while the update propagates.

**Property registry (Org3):**
```bash
cp -r $REPO/test-network/addorgregistropropiedad $HOME/fabric-samples/test-network/
cp -r $REPO/test-network/orgregistropropiedad-scripts $HOME/fabric-samples/test-network/scripts/
dos2unix $HOME/fabric-samples/test-network/addorgregistropropiedad/*.sh
cd $HOME/fabric-samples/test-network/addorgregistropropiedad
./addorgregistropropiedad.sh up -ca -s couchdb
```

**Academic registry (Org4):**
```bash
cp -r $REPO/test-network/addorgregistroacademico $HOME/fabric-samples/test-network/
cp -r $REPO/test-network/orgregistroacademico-scripts $HOME/fabric-samples/test-network/scripts/
dos2unix $HOME/fabric-samples/test-network/addorgregistroacademico/*.sh
cd $HOME/fabric-samples/test-network/addorgregistroacademico
./addorgregistroacademico.sh up -ca -s couchdb
```

> Replace `$REPO` with the full path to this repository.

**Common failure:** If the script reports `context deadline exceeded`, the Docker network `fabric_test` may not have fully initialized. Wait 10 seconds and retry.

---

## Step 6 — RBAC identity generation

The `gen_identities.sh` script iterates over all four organizations and calls `fabric-ca-client register` followed by `fabric-ca-client enroll` for each of the four ABAC roles. The resulting MSP directories are created directly under each organization's `users/` path, making them immediately available to the Go application.

```bash
cd $HOME/fabric-samples/test-network
chmod +x gen_identities.sh
./gen_identities.sh
```

**Verify a generated identity:** The file `users/operator1@orgregistrocivil.example.com/msp/signcerts/cert.pem` must exist. Its contents can be inspected with `openssl x509 -in <path> -text -noout` — look for the `role=OPERATOR` extension in the subject attributes section.

---

## Step 7 — Chaincode deployment

The `deployCC` wrapper packages the chaincode as a `.tar.gz`, installs it on all four peers, collects approvals, verifies commit readiness, and commits the definition. The sequence number (`-ccs`) must be incremented by one for each upgrade:

```bash
cd $HOME/fabric-samples/test-network
./network.sh deployCC -ccn dtic_chaincode -ccp ../dtic_chaincode -ccl go -ccs 1 -ccv 1.0
```

To upgrade after modifying the smart contract source:
```bash
./network.sh deployCC -ccn dtic_chaincode -ccp ../dtic_chaincode -ccl go -ccs 2 -ccv 2.0
```

**Verify the committed policy:**
```bash
export PATH=$HOME/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$HOME/fabric-samples/config
source scripts/envVar.sh && setGlobals 1
peer lifecycle chaincode querycommitted -C mychannel -n dtic_chaincode
```

The output must show `sequence: 1` and the `MAJORITY` endorsement rule listing all four MSP IDs.

---

## Step 8 — IPFS node

The private swarm key located at `test-network/swarm.key` prevents the node from connecting to the global IPFS network, ensuring documents remain within the controlled environment. After the `docker run` command, CORS headers must be configured to allow the Go application to reach the API:

```bash
docker exec ipfs_node ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
docker exec ipfs_node ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST"]'
docker restart ipfs_node
```

Test the API endpoint: `curl -s -X POST http://localhost:5001/api/v0/id` should return a JSON object containing the node's peer ID.

---

## Step 9 — Application

The Go application reads certificate paths relative to `../test-network/organizations/`. This path is resolved from wherever `go run main.go` is executed, so it must be run from inside the `app/` directory:

```bash
cd $HOME/fabric-samples/app
go run .
```

The CLI prompts for an organization (1–4), a user identity, and then presents the available transactions for that institution. All text input is converted to uppercase before submission to ensure consistent world state keys. CIDs returned by the IPFS upload are displayed immediately and stored permanently on-chain.

---

## Network ports reference

| Service | Host Port |
|---|---|
| `peer0.orgregistrocivil` | 7051 |
| `peer0.orgregistropolicial` | 9051 |
| `peer0.orgregistropropiedad` | 11051 |
| `peer0.orgregistroacademico` | 13051 |
| Orderer | 7050 |
| CA orgregistrocivil | 7054 |
| CA orgregistropolicial | 8054 |
| CA orgregistropropiedad | 11054 |
| CA orgregistroacademico | 12054 |
| IPFS API | 5001 |
| IPFS Gateway | 8080 |
| CouchDB (Civil) | 5984 |
| CouchDB (Policial) | 7984 |
| CouchDB (Propiedad) | 9984 |
| CouchDB (Académico) | 11984 |
