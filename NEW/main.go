package main

import (
    "bufio"
    "context"
    "crypto/rand"
    "fmt"
    "os"
    "io"
    "strings"
    "log"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p/core/crypto"
    "github.com/libp2p/go-libp2p/core/host"
    "github.com/libp2p/go-libp2p/core/network"
    "github.com/libp2p/go-libp2p/core/peer"
    "github.com/libp2p/go-libp2p/core/peerstore"
    "github.com/multiformats/go-multiaddr"
    tea "github.com/charmbracelet/bubbletea"
    pubsub "github.com/libp2p/go-libp2p-pubsub"

)

type model struct {
    host         host.Host
    currentView  string
    input        string
    messages     []string
    errorMessage string
    messageChan  chan string
    ps *pubsub.PubSub
    selectedMenuItem int
    subscribedTopics []string
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m *model) updateMessages() {
    m.host.SetStreamHandler("/chat/1.0.0", func(s network.Stream) {
        rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
        go func() {
            for {
                str, err := rw.ReadString('\n')
                if err != nil {
                    log.Println("Error reading from buffer")
                    break
                }
                m.messageChan <- str
            }
        }()
    })
}


func (m *model) publishMessage(topic string, msg string) error {
    t, err := m.ps.Join(topic)
    if err != nil {
        return err
    }

    sub, err := t.Subscribe()
    if err != nil {
        return err
    }
    defer sub.Cancel()

    return t.Publish(context.Background(), []byte(msg))
}


func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp, tea.KeyDown:
            if msg.Type == tea.KeyUp {
                m.selectedMenuItem--
            } else {
                m.selectedMenuItem++
            }

            menuItemsCount := 4
            if m.selectedMenuItem > menuItemsCount {
                m.selectedMenuItem = 1
            } else if m.selectedMenuItem < 1 {
                m.selectedMenuItem = menuItemsCount
            }

        case tea.KeyEnter:
            switch m.currentView {
            case "publish":
                if m.input != "" {
                    err := m.publishMessage("your_topic", m.input)
                    if err != nil {
                        m.errorMessage = err.Error()
                    }
                    m.input = ""
                }
            case "menu":
                switch m.selectedMenuItem {
                case 1:
                    m.currentView = "subscribe"
                case 2:
                    m.currentView = "publish"
                case 3:
                    m.currentView = "listTopics"
                case 4:
                    m.currentView = "listPeers"
                }
            }

        case tea.KeyCtrlC, tea.KeyEsc:
            return m, tea.Quit
        }

    case string:
        m.messages = append(m.messages, msg)
    }

    return m, nil
}



func (m model) View() string {
    var s strings.Builder

    switch m.currentView {
    case "menu":
        // Dynamically display menu items with selection
        menuItems := []string{"Subscribe to a topic", "Publish a message", "List topics", "List peers"}

        s.WriteString("Dangerous Net | IPFS Chat Menu\n")
        for i, item := range menuItems {
            if m.selectedMenuItem == i+1 {
                s.WriteString(fmt.Sprintf("-> %d. %s\n", i+1, item)) // Highlight selected item
            } else {
                s.WriteString(fmt.Sprintf("   %d. %s\n", i+1, item))
            }
        }
        s.WriteString("\nPress number to select, Esc to return to menu\n")

    case "subscribe":
        // Handle the subscribe view here
        // Example:
        s.WriteString("Enter topic to subscribe: " + m.input + "\n")

    case "publish":
        // Handle the publish view here
        // Example:
        s.WriteString("Enter message to publish: " + m.input + "\n")

    case "listTopics":
        s.WriteString("List of Topics:\n")
        for _, topic := range listTopics() { // Assuming listTopics is a function that returns []string
            s.WriteString(fmt.Sprintf("- %s\n", topic))
        }

    case "listPeers":
        s.WriteString("List of Peers:\n")
        for _, peerID := range listPeers(m.host) { // Assuming listPeers is a function that returns []string
            s.WriteString(fmt.Sprintf("- %s\n", peerID))
        }

    // Add other cases if necessary
    }

    // Error message display
    if m.errorMessage != "" {
        s.WriteString("\nError: " + m.errorMessage + "\n")
    }

    return s.String()
}



func (m *model) subscribeToTopic(topicName string) error {
    // Join the topic
    t, err := m.ps.Join(topicName)
    if err != nil {
        return err
    }

    // Subscribe to the topic
    _, err = t.Subscribe()
    if err != nil {
        return err
    }

    // Optionally, keep track of subscribed topics
    m.subscribedTopics = append(m.subscribedTopics, topicName)
    return nil
}


func makeHost(port int, randomness io.Reader) (host.Host, error) {
    var prvKey crypto.PrivKey
    var err error

    // Check if a specific randomness source is provided
    if randomness != nil {
        // Use the provided randomness source
        prvKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 4096, randomness)
    } else {
        // Use default randomness source
        prvKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 4096, rand.Reader)
    }

    if err != nil {
        log.Println(err)
        return nil, err
    }

    // 0.0.0.0 will listen on any interface device.
    sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // Constructs a new libp2p Host with the given options.
    return libp2p.New(
        libp2p.ListenAddrs(sourceMultiAddr),
        libp2p.Identity(prvKey),
    )
}




func main() {
    // Initialize libp2p host and other necessary components
    h, err := makeHost(0, nil) // nil indicates default randomness
    if err != nil {
        log.Println(err)
        return
    }
    defer h.Close()

    // Initialize the PubSub service
    ps, err := pubsub.NewGossipSub(context.Background(), h)
    if err != nil {
        log.Fatal(err)
    }

    // Initialize the model with the host, PubSub service, and set initial view
    m := model{
        host:        h,
        messageChan: make(chan string),
        ps:          ps,
        currentView: "menu", // Set initial view to "menu"
	selectedMenuItem: 1,
    }

    m.updateMessages()

    p := tea.NewProgram(&m)
    if err := p.Start(); err != nil {
        fmt.Printf("Error running program: %v", err)
        os.Exit(1)
    }
}

func startPeer(ctx context.Context, h host.Host, streamHandler network.StreamHandler) {
	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	h.SetStreamHandler("/chat/1.0.0", streamHandler)

	// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		log.Println("was not able to find actual local port")
		return
	}

	log.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, h.ID())
	log.Println("You can replace 127.0.0.1 with public IP as well.")
	log.Println("Waiting for incoming connection")
	log.Println()
}

func listPeers(h host.Host) []string {
    peers := h.Network().Peers()
    peerIDs := make([]string, len(peers))
    for i, peer := range peers {
        peerIDs[i] = peer.String()
    }
    return peerIDs
}


var topics map[string]*pubsub.Topic

func addTopic(topicName string, topic *pubsub.Topic) {
    if topics == nil {
        topics = make(map[string]*pubsub.Topic)
    }
    topics[topicName] = topic
}

func listTopics() []string {
    var topicNames []string
    for name := range topics {
        topicNames = append(topicNames, name)
    }
    return topicNames
}


func startPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
	log.Println("This node's multiaddresses:")
	for _, la := range h.Addrs() {
		log.Printf(" - %v\n", la)
	}
	log.Println()

	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := h.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Established connection to destination")

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return rw, nil
}
