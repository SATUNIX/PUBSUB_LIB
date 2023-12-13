package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-pubsub"
)

func main() {
	ctx := context.Background()
	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

	pubsubService, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// Use pubsubService for subscribing to topics and publishing messages
}