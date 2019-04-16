docker kill $(docker ps -aq)
docker rm $(docker ps -aq)
rm -r crypto-config
rm -r channel-artifacts
