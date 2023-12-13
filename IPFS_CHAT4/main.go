package main 

import(
	"bufio"
	"context"

	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"encoding/json"
	"github.com/multiformats/go-multiaddr"
	"path/filepath"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
/*
 * FOR BROKEN ENCRYPTION / USER ACCOUNT AUTH LOGIC
	"golang.org/x/crypto/ssh/terminal"
        "golang.org/x/crypto/scrypt"
        "crypto/aes"
        "crypto/cipher"
	"errors"
	"io"
*/

)


func main() {
    sourcePort := flag.Int("sp", 0, "Source port number")
    flag.Parse()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    h, err := MakeHost(*sourcePort)
    if err != nil {
        log.Fatal(err)
    }

    ps, err := pubsub.NewGossipSub(ctx, h)
    if err != nil {
        log.Fatal(err)
    }

    var chatRoom *ChatRoom

    // Display the main menu
    for {
	ChatMenuDisplay()

        var choice int
        fmt.Print("Enter your choice: ")
        fmt.Scanln(&choice)

        switch choice {
        case 1:
            // Join Chat room / subscribe to topic
            fmt.Print("Enter chat room name: ")
            var roomName string
            fmt.Scanln(&roomName)

            fmt.Print("Enter your nickname: ")
            var nickname string
            fmt.Scanln(&nickname)

            chatRoom, err = JoinChatRoom(ctx, ps, h.ID(), nickname, roomName)
            if err != nil {
                log.Println("Error joining chat room:", err)
                continue
            }
            fmt.Println("Joined chat room:", roomName)

        case 2:
            // Publish message
            if chatRoom == nil {
                fmt.Println("Please join a chat room first.")
                continue
            }

            fmt.Print("Enter message: ")
            var message string
            fmt.Scanln(&message)

            err := chatRoom.Publish(message)
            if err != nil {
                log.Println("Error publishing message:", err)
            }

        case 3:
            // Start Interactive Chat
            if chatRoom == nil {
                fmt.Println("Please join a chat room first.")
                continue
            }
            startChatInterface(ctx, chatRoom)

        case 0:
            // Exit
            fmt.Println("Exiting application.")
            return

        default:
            fmt.Println("Invalid choice. Please enter a valid option.")
        }
    }
}

func ChatMenuDisplay() {

    // Yellow lines
    fmt.Println("\033[1;33m---------------------------------------------\033[0m")
    // Bold white title
    fmt.Println("\033[1;37mDangerous Net | LIBP2P Chat Application\033[0m")
    // Yellow lines
    fmt.Println("\033[1;33m=============================================\033[0m")

    // Normal text for the options
    fmt.Println("\033[1;32m>\033[0;32m 1.\033[0m Join Chat Room")
    fmt.Println("\033[1;32m>\033[0;32m 2.\033[0m Publish Message to Chat Room")
    fmt.Println("\033[1;32m>\033[0;32m 3.\033[0m Interactive Chat \033[1;32m(EXPERIMENTAL)\033[0m")

    fmt.Println("\033[1;32m>\033[0;32m 0.\033[0m Exit")

    // Yellow lines
    fmt.Println("\033[1;33m=============================================\033[0m")
}

func GetConfigDir() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    return filepath.Join(homeDir, ".config", "DangerousNet", "Chat", "Keys")
}

/*
BROKEN ENCRYPTION LOGIC FOR USER ACCOUNT (EXPERIMENTAL)
func EncryptKey(key []byte, password string, saltPath string) ([]byte, error) {
    salt := make([]byte, 8)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, err
    }

    // Save salt to a file
    err := os.WriteFile(saltPath, salt, 0600)
    if err != nil {
        return nil, err
    }

    derivedKey, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(derivedKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    encryptedData := gcm.Seal(nil, nonce, key, nil)
    return encryptedData, nil // Return only the encrypted data
}
func DecryptKey(encryptedKey []byte, password string, saltPath string) ([]byte, error) {

    // Read salt from file
    salt, err := os.ReadFile(saltPath)
    if err != nil {
        return nil, err
    }
    if len(encryptedKey) < 8 {
        return nil, errors.New("encrypted key data is too short")
    }
    encryptedData := encryptedKey[8:]

    derivedKey, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(derivedKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedData) < nonceSize {
        return nil, errors.New("encrypted data is too short")
    }
    nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// Broken - Fix then merge with Dangerous Net then add keycard verification

BROKEN KEY ENCRYPTION LOGIC
func LoadOrCreateKey(configDir string) (crypto.PrivKey, error) {
    keyFilePath := filepath.Join(configDir, "private.key")
    saltFilePath := filepath.Join(configDir, "salt") // Path for the salt file

    // Check if the key file exists
    _, err := os.Stat(keyFilePath)
    if os.IsNotExist(err) {
        // Key file does not exist, generate a new key
        privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 4096, rand.Reader)
        if err != nil {
            return nil, err
        }

        // Create config directory if it does not exist
        err = os.MkdirAll(configDir, os.ModePerm)
        if err != nil {
            return nil, err
        }

        // Save the key to the file
        keyBytes, err := crypto.MarshalPrivateKey(privateKey)
        if err != nil {
            return nil, err
        }

        fmt.Print("Create a password: ")
        passwordBytes, err := terminal.ReadPassword(0)
        if err != nil {
            return nil, err
        }
        fmt.Println() // Print a newline after the password input

        encryptedKey, err := EncryptKey(keyBytes, string(passwordBytes), saltFilePath)
        if err != nil {
            return nil, err
        }

        err = os.WriteFile(keyFilePath, encryptedKey, 0640)
        if err != nil {
            return nil, err
        }

        return privateKey, nil
    } else if err != nil {
        return nil, err
    }

    // Key file exists, load the key
    fmt.Print("Enter your password: ")
    passwordBytes, err := terminal.ReadPassword(0)
    if err != nil {
        return nil, err
    }
    fmt.Println() // Print a newline after the password input

    encryptedKey, err := os.ReadFile(keyFilePath)
    if err != nil {
        return nil, err
    }

    decryptedKey, err := DecryptKey(encryptedKey, string(passwordBytes), saltFilePath)
    if err != nil {
        return nil, err
    }

    return crypto.UnmarshalPrivateKey(decryptedKey)
}
*/

func LoadOrCreateKey(configDir string) (crypto.PrivKey, error) {
    keyFilePath := filepath.Join(configDir, "private.key")

    // Check if the key file exists
    if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
        // Key file does not exist, generate a new key
        privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
        if err != nil {
            return nil, err
        }

        // Create config directory if it does not exist
        err = os.MkdirAll(configDir, os.ModePerm)
        if err != nil {
            return nil, err
        }

        // Save the key to the file
        keyBytes, err := crypto.MarshalPrivateKey(privateKey)
        if err != nil {
            return nil, err
        }

        err = os.WriteFile(keyFilePath, keyBytes, 0600)
        if err != nil {
            return nil, err
        }

        return privateKey, nil
    }

    // Key file exists, load the key
    keyBytes, err := os.ReadFile(keyFilePath)
    if err != nil {
        return nil, err
    }

    return crypto.UnmarshalPrivateKey(keyBytes)
}

func MakeHost(port int) (host.Host, error) {
    configDir := GetConfigDir()
    prvKey, err := LoadOrCreateKey(configDir)
    if err != nil {
        log.Fatal(err)
    }

    sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
    return libp2p.New(
        libp2p.ListenAddrs(sourceMultiAddr),
        libp2p.Identity(prvKey),
    )
}


const ChatRoomBufSize = 128  // Adjust as needed

type ChatRoom struct {
    ctx      context.Context
    ps       *pubsub.PubSub
    topic    *pubsub.Topic
    sub      *pubsub.Subscription
    self     peer.ID
    nick     string
    roomName string
    Messages chan *ChatMessage
}

type ChatMessage struct {
    Message    string
    SenderID   string
    SenderNick string
}
// topic handler
func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nickname, roomName string) (*ChatRoom, error) {
    topic, err := ps.Join(TopicName(roomName))
    if err != nil {
        return nil, err
    }

    sub, err := topic.Subscribe()
    if err != nil {
        return nil, err
    }

    chatRoom := &ChatRoom{
        ctx:      ctx,
        ps:       ps,
        topic:    topic,
        sub:      sub,
        self:     selfID,
        nick:     nickname,
        roomName: roomName,
        Messages: make(chan *ChatMessage, ChatRoomBufSize),
    }

    go chatRoom.readLoop()
    return chatRoom, nil
}


// message handler for chatrooms
func (cr *ChatRoom) Publish(message string) error {
    m := ChatMessage{
        Message:    message,
        SenderID:   cr.self.String(),
        SenderNick: cr.nick,
    }
    msgBytes, err := json.Marshal(m)
    if err != nil {
        return err
    }
    return cr.topic.Publish(cr.ctx, msgBytes)
}


func (cr *ChatRoom) ListPeers() []peer.ID {
    return cr.ps.ListPeers(TopicName(cr.roomName))
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}
		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		cr.Messages <- cm
	}
}

func TopicName(roomName string) string {
	return "chat-room:" + roomName
}

func HandleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}


func ReadData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func WriteData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
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




func startChatInterface(ctx context.Context, chatRoom *ChatRoom) {
    // Start a goroutine to handle incoming messages
    go func() {
        for msg := range chatRoom.Messages {
            // Check if the message is from the current user
            if msg.SenderID == chatRoom.self.String() {
                continue // Skip the user's own messages
            }
            fmt.Printf("\r\x1b[32m%s\x1b[0m: %s\n> ", msg.SenderNick, msg.Message)
        }
    }()

    // Main loop for sending messages
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("> ")
    for scanner.Scan() {
        text := scanner.Text()
        if text == "/exit" {
            fmt.Println("Exiting chat room...")
            return // Exit command to leave the chat
        }

        // Send message
        if err := chatRoom.Publish(text); err != nil {
            fmt.Println("Error sending message:", err)
        }
        fmt.Print("> ") // Prompt for next message
    }

    if err := scanner.Err(); err != nil {
        log.Println("Error reading from stdin:", err)
    }
}
