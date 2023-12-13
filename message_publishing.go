package main

import (
	"context"
	"fmt"
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

	topic, err := pubsubService.Join("my-topic")
	if err != nil {
		panic(err)
	}

	// Publish a message to the topic
	message := "Hello, world!"
	err = topic.Publish(ctx, []byte(message))
	if err != nil {
		fmt.Println("Error publishing message:", err)
	}

	// Use pubsubService for subscribing to topics and handling messages
}