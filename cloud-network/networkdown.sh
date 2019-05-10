docker kill $(docker ps -aq)
docker-compose -f docker-compose-cli.yaml down --volumes
docker rm $(docker ps -aq)
rm -r crypto-config
rm -r channel-artifacts
docker rmi $(docker image ls | awk '/mycc/{print $3}')
cd /home/rohan/Documents/cloudledger/SDK/wallet
rm -rf *
cd ~
rm -rf .hfc-key-store
