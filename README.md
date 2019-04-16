# cloudledger
Blockchain network built on Hyperledger Fabric for a secure cloud platform

### What is it?
In the present day it is very hard to find a reliable cloud platform to backup your important documents. This project aims to build a platform based on blockchain technology, where you can store encrypted data which can be accessed by only you.
This platform exploits the security of blockchain to provide an interface between the user and the distributed data store to ensure that only the user can have access to their files.

### How does it work?
#### Overview
- There are two types of users:
  - Service providers: ones that provides cloud storage
  - Service availers: ones that use the cloud to store their data
- The first type of users are the service providers and the list of such users with their credentials(IP address) is stored.
- The second type of users use the cloud to store data.
- We maintain two ledgers, one that can be queried by everyone and the other that can be queried only by the chaincode that maintains the first ledger.
- Each entry in the first ledger contains a composite key(tuple of username and private key), the user's public key and the data of the files that can be accessed by that person. The file data consists of the secret key and the list of the service providers who have the fragments of the file.
- Each entry in the second ledger has a mapping from the public key of a user to their composite key(the one mentioned above).
#### Uploading a file
- The user submits the file which has to be uploaded on the cloud.
- Each file is identified by a secret key(eg. concatenate the username and a random number) which can be accessed only by the owner of the file.
- The file is broken down into many parts(say 35) and is encrypted using this secret key.
- A large number of service providers(say 100) are chosen at random and their credentials are committed onto the ledger. We make sure that we upload each fragment of the file to at least 2 service providers so that we can retrieve all fragments even if some of them are not currently on the network.
#### Downloading a file
- When the user wants to download a file the chaincode is queried using the composite key and it returns the list of files and the file data that can be accessed by them. The user can go on and download each file using the information about the service providers that have the fragments.
#### Sharing a file
- If A wants to share a file with B, then A invokes the chaincode with the public key of B. This chaincode queries the second ledger(mentioned above) and gets the composite key of B and updates the ledger with relevant information of the file so that B gets access rights to the file.
#### Deleting a file
- When a user wants to delete a file from his cloud storage, he invokes the chaincode to remove the file data(data regarding the secret key and the service providers who have the fragments) from his state. This ensures that the particular file is effectively deleted only for this user but not for the users to whom the file was shared.

### Known Issues
- The list of service providers is accessible by everyone on the UI and is centralised. This data can be manipulated such that the future uploads by a user are funneled into an adversarial computer. We plan to solve this issue by introducing a third ledger which is accessible to the first chaincode. 

### Useful resources to understand hyperledger fabric
#### Official documentation
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/whatis.html
#### Setting up the network
- https://medium.freecodecamp.org/how-to-build-a-blockchain-network-using-hyperledger-fabric-and-composer-e06644ff801d
- https://www.youtube.com/watch?v=MPNkUqOKhVE&list=PLjsqymUqgpSTGC4L6ULHCB_Mqmy43OcIh
#### For writing chiancode in Golang
- https://www.ibm.com/developerworks/cloud/library/cl-ibm-blockchain-chaincode-development-using-golang/cl-ibm-blockchain-chaincode-development-using-golang-pdf.pdf
