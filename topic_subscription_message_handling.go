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

	subscription, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			fmt.Printf("Received message from %s: %s\n", msg.GetFrom(), string(msg.GetData()))
		}
	}()

	// Use pubsubService for publishing messages
}