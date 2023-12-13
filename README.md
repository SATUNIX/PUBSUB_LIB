# PUBSUB_LIB
This repository is created to facilitate the development life cycle of messaging in the Dangerous Net program. 

## Decentralised Data Transfer Service

1. Use the libp2p pubsub mechanisms to build a chat messaging system, functionality can be thought of like an IRC.
2. Completely decentralised using the libp2p framework
3. Library should integrate into the Dangerous Net repository smoothly

## Note
1. All of the above files are just iterations Ive gone through, not complete programs at all
2. Alot of work to be done still though 80% of the code we need is in here, just need to organise it and use these functions to build a basic working model.

## Final Product Criteria 
1. Complete Decentralisation
2. P2P Messaging based off of peer ID's (peer ro peer messaging)
3. Daemon backround service (/etc/systemd/system/*) creation to handle listening of messages
4. Network saving all chats, encrypted or not over a cluster network (Clustering to be implemented in Dangerous Net in parralell)
5. This save of chats allows for previous messages to be downloaded and read by a client even if their node was previously offline
6. Recommendations for encryption / Group chats secured using dual encryption / name of new group + machine phrase = input + phrase salt = AES cipher key for group chats
    - AES Key treated as invite code to a group
    - See the dangerous net main.go encryption functions for inspiration 
7. Intuitive TUI interface (like terminal forms, bubble tea etc) (already started work on this in oldmain.go
8. Able to integrate directly into Dangerous Net and future Dangerous Net clustering service. 




## Acknowledgements 
1. The libp2p framework and libraries:
![image](https://github.com/SATUNIX/PUBSUB_LIB/assets/111553838/5fd76b0d-6a1a-4472-ac2e-b06abe5457ef)

2. IPFS
![image](https://github.com/SATUNIX/PUBSUB_LIB/assets/111553838/04bc7e41-6923-4d4f-99f8-ca6b33bfc3e5)

3. RSA
![image](https://github.com/SATUNIX/PUBSUB_LIB/assets/111553838/ac1a6815-50bb-48ee-acaa-18291e6eb137)

