# cloudledger
Blockchain network built on Hyperledger Fabric for a secure cloud platform

### What is it?
In the present day it is very hard to find a reliable cloud platform to backup your important documents. This project aims to build a platform based on blockchain technology, where you can store encrypted data which can be accessed by only you.

### How it works
- There are two types of users: 1) Provides cloud storage 2) uses the cloud to store their data.
- The first type of users are the service providers and the list of such users with their credentials(IP address) is stored in a list.
- The second type of user submits the file which has to be stored on the cloud.
- Each file has a secret key(eg. concatenate the username and a random number) which can be accessed only by the owner of the file.
- File is encoded into a string using base64 and the string is returned.
- A large number of service providers(say 50) are chosen at random and their credentials are committed onto the user's ledger.
- The returned string is then broken down into large number of components(50 in this case) and each sub-string is then uploaded to the peer server of the service providers who were chosen in the previous step. The file is uploaded using a P2P network.
- When the peer wants to download the file he queries the ledger by the secret key to get the credentials of the service providers and downloads the respective substrings from each of them.
- The substrings are then concatenated and decrypted at the client's PC.
- When the client wants to share the file to other peers, the secret key and the information of the service providers are committed onto that peers ledger.

### Useful resources to understand hyperledger fabric
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/whatis.html
- https://medium.freecodecamp.org/how-to-build-a-blockchain-network-using-hyperledger-fabric-and-composer-e06644ff801d
- https://www.youtube.com/watch?v=MPNkUqOKhVE&list=PLjsqymUqgpSTGC4L6ULHCB_Mqmy43OcIh
