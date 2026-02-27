#!/bin/bash

echo "=========================================================="
echo "Preparing the environment for Hyperledger Fabric"
echo "=========================================================="

# 1. Update package list
echo "[1/6] Updating repositories..."
sudo apt-get update

# 2. Install base tools and dependencies
echo "[2/6] Installing Git, Curl, JQ, and Go..."
sudo apt-get install -y git curl jq golang-go

# 3. Install Docker and Docker Compose
echo "[3/6] Installing Docker Compose..."
sudo apt-get -y install docker-compose

# 4. Configure Docker daemon
echo "[4/6] Starting and enabling Docker service..."
sudo systemctl start docker
sudo systemctl enable docker

# 5. Configure user permissions for Docker
echo "[5/6] Adding user '$USER' to the docker group..."
sudo usermod -a -G docker $USER

# 6. Download and install Hyperledger Fabric
echo "[6/6] Downloading Fabric script and installing version 2.5.14..."
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh

# Download all (docker images, samples, binaries) specifically for v2.5.14
./install-fabric.sh --fabric-version 2.5.14 docker samples binary

# 7. Pull additional Chaincode Environment images
echo "[7/6] Pulling additional chaincode environment images (2.5 & 3.0)..."
sudo docker pull hyperledger/fabric-ccenv:2.5
sudo docker pull hyperledger/fabric-ccenv:3.0

echo "=========================================================="
echo "Installation completed successfully!"
echo "=========================================================="
echo "IMPORTANT: For Docker permissions to take effect without rebooting, run this command right now:"
echo "su - $USER"