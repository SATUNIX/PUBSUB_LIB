package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"
    "log"
    "encoding/json"
    tea "github.com/charmbracelet/bubbletea"
    libp2p "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-core/host"
    "github.com/libp2p/go-libp2p-core/peer"
    pubsub "github.com/libp2p/go-libp2p-pubsub"
    "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
    "time"
)

const DiscoveryServiceTag = "librum-pubsub"
const DiscoveryInterval = time.Hour

type model struct {
    host        host.Host
    gossipSub   *pubsub.PubSub
    topic       *pubsub.Topic
    subscriber  *pubsub.Subscription
    chatClient  *ipfschat.IPFSChat // Replace with actual implementation
    currentView string
    input       string
    messages    []string
    errorMessage string
    messageChan  chan string // Channel for incoming messages
}

// ... (include other types and functions from your original code)

func main() {
    ctx := context.Background()

    // Create a new libp2p Host
    h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
    if err != nil {
        panic(err)
    }

    // Set up mDNS discovery
    setupDiscovery(h)

    // Create a new PubSub service
    gossipSub, err := pubsub.NewGossipSub(ctx, h)
    if err != nil {
        panic(err)
    }

    // Join the pubsub topic
    topic, err := gossipSub.Join("librum")
    if err != nil {
        panic(err)
    }

    // Subscribe to the topic
    subscriber, err := topic.Subscribe()
    if err != nil {
        panic(err)
    }

    // Initialize the TUI model
    m := initialModel()
    m.host = h
    m.gossipSub = gossipSub
    m.topic = topic
    m.subscriber = subscriber

    // Start the TUI program
    p := tea.NewProgram(m)
    if err := p.Start(); err != nil {
        fmt.Printf("Error running program: %v", err)
        os.Exit(1)
    }
}

// ... (include other methods and functions from your original code)

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
func setupDiscovery(h host.Host) error {
    // setup mDNS discovery to find local peers
    s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
    return s.Start()
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
    h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
    fmt.Printf("Discovered new peer %s\n", pi.ID.Pretty())
    err := n.h.Connect(context.Background(), pi)
    if err != nil {
        fmt.Printf("Error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
    }
}

