package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = protocol.ID("/my-protocol/1.0.0")

func main() {
	ctx := context.Background()
	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Host created. We are: %s\n", host.ID())

	// Set a stream handler on host
	host.SetStreamHandler(protocolID, handleStream)

	// Connect to a peer (replace "<peer_multiaddr>" with the multiaddress of a peer)
	peerAddr, err := multiaddr.NewMultiaddr("<peer_multiaddr>")
	if err != nil {
		panic(err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
	if err != nil {
		panic(err)
	}

	err = host.Connect(ctx, *peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to: %s\n", peerInfo.ID)
}

func handleStream(stream network.Stream) {
	// Handle the stream
}