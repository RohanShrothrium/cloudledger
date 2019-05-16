for i in {1..5}; do
     export BYFN_CA"${i}"_PRIVATE_KEY=$(ls /var/mynetwork/certs/crypto-config/peerOrganizations/org${i}.example.com/ca/ | grep _sk)
     echo $BYFN_CA"${i}"_PRIVATE_KEY
done
./scripts/network/deploy_services_kafka.sh
./scripts/network/deploy_services_org1.sh
./scripts/network/deploy_services_org2.sh
./scripts/network/deploy_services_org3.sh
./scripts/network/deploy_services_org4.sh
./scripts/network/deploy_services_org5.sh



