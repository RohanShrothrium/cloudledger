docker kill $(docker ps -aq)
docker-compose -f docker-compose-cli.yaml down --volumes
docker rm $(docker ps -aq)
rm -r crypto-config
rm -r channel-artifacts
docker rmi $(docker image ls | awk '/mycc/{print $3}')
