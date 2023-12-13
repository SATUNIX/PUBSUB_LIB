## Implementing Pubsub for Real-time Messaging

### Overview

To enable real-time messaging in the chat Dapp, we need to implement IPFS pubsub. Pubsub allows nodes to subscribe to topics and receive messages published to those topics.

### Steps

1. Enable pubsub in the IPFS configuration
2. Subscribe to a chat topic
3. Publish messages to the chat topic
4. Receive messages from the chat topic
5. Update the view model to handle sending and receiving messages through IPFS pubsub