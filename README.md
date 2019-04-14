# cloudledger
Blockchain network built on Hyperledger Fabric for a secure cloud platform

### What is it?
In the present day it is very hard to find a reliable cloud platform to backup your important documents. This project aims to build a platform based on blockchain technology, where you can store encrypted data which can be accessed by only you.
This platform exploits the security of blockchain to provide an interface between the user and the distributed data store to ensure that only the user can have access to their files.

### How does it work?
- There are two types of users: 1) Provides cloud storage 2) uses the cloud to store their data.
- The first type of users are the service providers and the list of such users with their credentials(IP address) is stored in a list.
- The second type of users  use the cloud to store data.
- We maintain two ledgers, one that can be queried by everyone and the other that can be queried only by the chaincode that maintains the first ledger.
- Each entry in the first ledger contains composite key(tuple of username and private key), the user's public key and the data of the files that can be accessed by that person. The file data consists of the secret key and the list of the service providers.
- Each entry in the second ledger has a mapping from the public key of a user to their composite key(the one mentioned above).
- The user submits the file which has to uploaded to the cloud.
- Each file is identified by a secret key(eg. concatenate the username and a random number) which can be accessed only by the owner of the file.
- The file is broken down into many parts(say 35) and is encrypted using this secret key.
- A large number of service providers(say 100) are chosen at random and their credentials are committed onto the ledger. We make sure that we upload each fragment of the file to at least 2 service providers so that we can retrieve all fragments even if some of them is not currently on the network.
- The chaincode is queried using the composite key and it returns the list of files that can be accessed by them. The user can go on and download each file using the information about the service providers that have the fragments.
- If A wants to share a file with B, then A invokes the chaincode with the public key of B. This chaincode queries the second ledger(mentioned above) and gets the composite key of B and updates the ledger with relevant information of the file so that B gets access rights to the file.

### Useful resources to understand hyperledger fabric
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/whatis.html
- https://medium.freecodecamp.org/how-to-build-a-blockchain-network-using-hyperledger-fabric-and-composer-e06644ff801d
- https://www.youtube.com/watch?v=MPNkUqOKhVE&list=PLjsqymUqgpSTGC4L6ULHCB_Mqmy43OcIh
