#!/bin/bash
export FABRIC_CFG_PATH=${PWD}
# GIVE A CHANNEL NAME OF YOUR OWN CHOICE
export CHANNEL_NAME=mychannel
export HLF_VERSION=1.4.6

# Generating crytographic material :- Certs
echo "=========================================="
echo "        Generating Certificates           "
echo "=========================================="
./cryptogen generate --config=./crypto-config.yaml

# Generating Orderer and Genesis block
echo "=========================================="
echo "     Creating Orderer Genesis Block       "
echo "=========================================="
mkdir channel-artifacts
chmod 777 *
./configtxgen -profile SampleMultiNodeEtcdRaft -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

# Creating a Channel Config Transaction
echo "=========================================="
echo "         Creating Channel Config          "
echo "=========================================="
./configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

# Defining Anchor peers from all organisations

echo "=========================================="
echo "          Defining Anchor Peers           "
echo "=========================================="
echo ""
echo ""
echo "Org1......."
./configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
echo "Org2......."
./configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

# Setting up the network using docker compose
for i in {1..2}; do
     export BYFN_CA"${i}"_PRIVATE_KEY=$(ls crypto-config/peerOrganizations/org${i}.example.com/ca/ | grep _sk)
done
echo "=========================================="
echo "           Setting up Network             "
echo "=========================================="
export IMAGE_TAG=1.4.6
docker-compose -f docker-compose-etcdraft2.yaml -f docker-compose-cli.yaml up -d
# docker exec cli scripts/functions.sh
docker ps -a