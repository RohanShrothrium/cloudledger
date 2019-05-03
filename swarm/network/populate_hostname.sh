#!/bin/bash

set -ev

ENV_LOCATION=$PWD/.env
echo $ENV_LOCATION
source $ENV_LOCATION

FLAG="-i"
ARCH=$(uname)
if [ "$ARCH" == "Linux" ]; then
  FLAG="-i"
elif [ "$ARCH" == "Darwin" ]; then
  FLAG="-it"
fi

# Zookeeper + Kafka
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG1_HOSTNAME/g" $ZK_COMPOSE_PATH
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG1_HOSTNAME/g" $KAFKA_COMPOSE_PATH

# ORG1
ORG1_CA_PATH=$(ls /var/mynetwork/certs/crypto-config/peerOrganizations/org1.example.com/ca/ | grep "_sk")
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG1_HOSTNAME/g" $ORDERER0_COMPOSE_PATH
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG1_HOSTNAME/g" $PEER_ORG1_COMPOSE_PATH
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG1_HOSTNAME/g" $SERVICE_ORG1_COMPOSE_PATH
sed "$FLAG" "s#- FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*#- FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$ORG1_CA_PATH#g" $SERVICE_ORG1_COMPOSE_PATH


# ORG2
ORG2_CA_PATH=$(ls /var/mynetwork/certs/crypto-config/peerOrganizations/org2.example.com/ca/ | grep "_sk")
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG2_HOSTNAME/g" $ORDERER1_COMPOSE_PATH
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG2_HOSTNAME/g" $PEER_ORG2_COMPOSE_PATH
sed "$FLAG" "s/- node.hostname == .*/- node.hostname == $ORG2_HOSTNAME/g" $SERVICE_ORG2_COMPOSE_PATH
sed "$FLAG" "s#- FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*#- FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$ORG2_CA_PATH#g" $SERVICE_ORG2_COMPOSE_PATH


if [ "$ARCH" == "Darwin" ]; then
  rm */**.ymlt
fi
