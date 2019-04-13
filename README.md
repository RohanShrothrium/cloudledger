# cloudledger
Blockchain network built on Hyperledger Fabric for a secure cloud platform

### What is it?
In the present day it is very hard to find a reliable cloud platform to backup your important documents. This project aims to build a platform based on blockchain technology to where you can store encrypted data which can be accessed by only you.

### How it works
- User submits file.
- File is encoded into a string using base64 and the string is returned.
- The returned string is then broken down into three components and each sub-string is committed onto three ledgers which do not belong to the same channel with a suitable key(eg. concatenate the username and a random number).

### Useful resources to understand hyperledger fabric
- https://hyperledger-fabric.readthedocs.io/en/release-1.4/whatis.html
- https://medium.freecodecamp.org/how-to-build-a-blockchain-network-using-hyperledger-fabric-and-composer-e06644ff801d
- https://www.youtube.com/watch?v=MPNkUqOKhVE&list=PLjsqymUqgpSTGC4L6ULHCB_Mqmy43OcIh
