#!/bin/bash

# Setting up important environment variables
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
PEER0_ORG3_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt

# Helper function to set up the environment variables
setGlobals() {
  PEER=$1
  ORG=$2
  echo "Setting global variables for Org{$ORG} Peer{$PEER}"
  if [ $ORG -eq 1 ]; then
    CORE_PEER_LOCALMSPID="Org1MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org1.example.com:7051
    else
      CORE_PEER_ADDRESS=peer1.org1.example.com:8051
    fi
  elif [ $ORG -eq 2 ]; then
    CORE_PEER_LOCALMSPID="Org2MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org2.example.com:9051
    else
      CORE_PEER_ADDRESS=peer1.org2.example.com:10051
    fi

  elif [ $ORG -eq 3 ]; then
    CORE_PEER_LOCALMSPID="Org3MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org3.example.com:11051
    else
      CORE_PEER_ADDRESS=peer1.org3.example.com:12051
    fi
  else
    echo "================== ERROR !!! ORG Unknown =================="
  fi

}

# Create the channel
createChannel(){
    export CHANNEL_NAME=mychannel
    setGlobals 0 1
    echo "================================================="
    echo "Creating channel with channel name" $CHANNEL_NAME
    echo "================================================="
    echo "peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile $ORDERER_CA >&log.txt"
    peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile $ORDERER_CA >&log.txt
	cat log.txt  
}

# Joing peers to the channel
joinChannel(){
  export CHANNEL_NAME=mychannel
  for org in 1 2 3; do
    for peer in 0 1; do
      echo "================================================="
      echo "PEER" $peer "ORG" $org "joining the channel"
      echo "================================================="
      setGlobals $peer $org
      peer channel join -b $CHANNEL_NAME.block
      res=$?
      if [ $res -ne 0 ]; then
      sleep 3
        peer channel join -b $CHANNEL_NAME.block
        res=$?
        if [ $res -ne 0 ]; then
          sleep 3
          peer channel join -b $CHANNEL_NAME.block
          res=$?
          if [ $res -ne 0 ]; then
            sleep 3
            peer channel join -b $CHANNEL_NAME.block
          fi
        fi
      fi
      sleep 5
      echo $org$peer "joined the channel."
    done
	done
}

# Updating Anchor Peers
updateAnchorPeers(){
  export CHANNEL_NAME=mychannel
  echo "UPDATING ANCHOR PEERS"
  echo ""
  for org in 1 2 3; do
    for peer in 0; do
      echo "================================================="
      echo "PEER" $peer "ORG" $org "set as anchor peer."
      echo "================================================="
      setGlobals $peer $org
      peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA
      echo $org$peer "set as Anchor Peer."
    done                        
  done
}

# Installing Chaincodes on each peer 
installChaincode(){
    export CHANNEL_NAME=mychannel
    echo "INSTALLING CHAINCODE"
    echo ""
    for org in 1 2 3; do
        for peer in 0 1; do
            echo "================================================="
            echo "Installing chaincode on PEER" $peer "ORG" $org "."
            echo "================================================="
            setGlobals $peer $org
            peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/go
        done
    done
}

# Instantiating installed chaincode
instantiateChaincode(){
  export CHANNEL_NAME=mychannel
	setGlobals 0 1
    echo "================================================="
    echo "            INSTANTIATING CHAINCODE              "
    echo "================================================="
    peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -l node -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')"            
}

# setting globals
setGlobals 0 1
createChannel
joinChannel
updateAnchorPeers
installChaincode
instantiateChaincode