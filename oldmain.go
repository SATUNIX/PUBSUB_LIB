package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "log"
    "json"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    chatClient   *ipfschat.IPFSChat
    currentView  string
    input        string
    messages     []string
    errorMessage string
    messageChan   chan string // Channel for incoming messages
}


type errMsg struct{ err error }

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Error running program: %v", err)
        os.Exit(1)
    }
}

type Config struct {
    DefaultSubscription string `json:"defaultSubscription"`
}

func readConfig() Config {
    configFile, err := os.Open("config.json")
    if err != nil {
        log.Println("Error opening config file:", err)
        return Config{}
    }
    defer configFile.Close()

    var config Config
    json.NewDecoder(configFile).Decode(&config)
    return config
}

func setupLogger() {
    // Create log file
    file, err := os.OpenFile("ipfschat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }

    // Set log output
    log.SetOutput(file)
}
func initialModel() model {
    return model{
        chatClient:  ipfschat.NewIPFSChat("localhost:5001"),
        currentView: "menu",
        messageChan: make(chan string),
    }
}


func (m model) Init() tea.Cmd {
    return nil
}

func (m *model) Subscribe(topic string) {
    // Assuming you have a method to subscribe to a topic and receive messages
    go func() {
        for {
            message, err := m.chatClient.ReceiveMessage(topic) // Replace with actual method to receive messages
            if err != nil {
                // Handle error, perhaps by sending an error message to the channel
                m.messageChan <- fmt.Sprintf("Error receiving message: %v", err)
                continue
            }
            m.messageChan <- message
        }
    }()
}

func (m *model) SubscribeToDangerousNet() {
    specialTopic := "1337DangerousNet1337" // Define a special topic name for Dangerous Net
    log.Println("Subscribing to:", specialTopic)

    go func() {
        // Implement the subscription logic. This is a placeholder.
        // Replace it with actual IPFS subscription logic.
        for {
            message, err := m.chatClient.ReceiveMessage(specialTopic)
            if err != nil {
                m.messageChan <- fmt.Sprintf("Error receiving message: %v", err)
                continue
            }
            m.messageChan <- message
        }
    }()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "1", "2", "3", "4":
            m.currentView = msg.String()

	case "enter":
	    if m.currentView == "subscribe" {
                if m.input == "1337" {
            // Special handling for Dangerous Net
                log.Println("Accessing Dangerous Net")
                m.SubscribeToDangerousNet()
                m.currentView = "viewMessages"
            } else {
                m.Subscribe(m.input)
                m.currentView = "viewMessages"
            }
            m.input = ""
        } else if m.currentView == "publish" {
                err := m.Publish("your_topic", m.input)  // Use the Publish method
                if err != nil {
                    return m, func() tea.Msg {
                        return errMsg{err}
                    }
                }
                m.input = ""
            }

        case "esc":
            // Return to the main menu
            m.currentView = "menu"

        case "backspace":
            // Handle backspace for input
            if len(m.input) > 0 {
                m.input = m.input[:len(m.input)-1]
            }

        default:
            if m.currentView == "subscribe" || m.currentView == "publish" {
                // Append typed characters to input
                m.input += msg.String()
            }
        }

    case errMsg:
        // Handle error messages
        m.errorMessage = msg.err.Error()

    case string:
        // Handle incoming messages
        m.messages = append(m.messages, msg)

    default:
        // Handle other cases or ignore
    }
    return m, nil
}


func (m model) View() string {
    var s strings.Builder
    switch m.currentView {
    case "menu":
        s.WriteString("Dangerous Net | IPFS Chat Menu \n")
        s.WriteString("1. Subscribe to a topic\n")
        s.WriteString("2. Publish a message\n")
        s.WriteString("3. List topics\n")
        s.WriteString("4. List peers\n")
        s.WriteString("\nPress number to select, Esc to return to menu")

    case "subscribe":
        // Input for subscribing to a topic
	s.WriteString("Enter topic to subscribe or type '1337' for the Dangerous Net: " + m.input + "\n")

        // Display chat messages if subscribed
        s.WriteString("\n--- Chat Messages ---\n")
        for _, message := range m.messages {
            s.WriteString(message + "\n")
        }

    case "publish":
        // Input for publishing a message
        s.WriteString("Enter message to publish: " + m.input + "\n")
    }

    // Display error messages if any
    if m.errorMessage != "" {
        s.WriteString("\nError: " + m.errorMessage)
    }

    return s.String()
}

