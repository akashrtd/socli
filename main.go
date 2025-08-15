package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"socli/config"
	"socli/content"
	"socli/crypto"
	"socli/internal"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui"

	"github.com/libp2p/go-libp2p/core/peer"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(internal.ConfigFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Set up the libp2p host
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	netManager, err := p2p.NewNetworkManager(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating network manager: %v\n", err)
		os.Exit(1)
	}
	defer netManager.Close()

	// --- Key Management ---
	var keyPair *crypto.KeyPair
	if _, err := os.Stat(cfg.Privacy.KeyPath); os.IsNotExist(err) {
		// Key file doesn't exist, generate a new one
		keyPair, err = crypto.GenerateKeyPair()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating key pair: %v\n", err)
			os.Exit(1)
		}
		err = crypto.SaveKeyPair(keyPair, cfg.Privacy.KeyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving key pair: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated new key pair and saved to", cfg.Privacy.KeyPath)
	} else {
		// Load existing key pair
		keyPair, err = crypto.LoadKeyPair(cfg.Privacy.KeyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading key pair: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Loaded key pair from", cfg.Privacy.KeyPath)
	}
	// --- End Key Management ---

	if err := netManager.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting network manager: %v\n", err)
		os.Exit(1)
	}

	// Set up pubsub
	psManager, err := p2p.NewPubSubManager(ctx, netManager.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pubsub manager: %v\n", err)
		os.Exit(1)
	}

	// Create a new in-memory store
	store := storage.NewMemoryStore()

	// Create a new markdown renderer
	renderer, err := content.NewMarkdownRenderer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating markdown renderer: %v\n", err)
		os.Exit(1)
	}

	// Create a new broadcaster, passing the keyPair
	broadcaster := messaging.NewBroadcaster(psManager, cfg, keyPair)

	// Initialize the main application model from the tui package, passing the keyPair
	appModel, err := tui.NewApp(netManager, store, renderer, broadcaster, keyPair)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Create a new BubbleTea program
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	// Set the peer connected callback after components are initialized
	netManager.SetPeerConnectedCallback(func(id peer.ID) {
		// Send a message to the TUI about the new peer
		log.Printf("Main: Notifying TUI of peer connected: %s\n", id.String()) // Use log for background processes
		p.Send(tui.PeerConnectedMsg{PeerID: id.String()})
	})

	// Subscribe to a default topic and start listening for messages
	// In a more advanced version, this would be dynamic based on user subscriptions.
	defaultTopicName := messaging.GetTopicForHashtag("general")
	topic, err := psManager.JoinTopic(defaultTopicName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error joining topic: %v\n", err)
		os.Exit(1)
	}

	sub, err := psManager.SubscribeToTopic(topic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error subscribing to topic: %v\n", err)
		os.Exit(1)
	}

	// --- TEMPORARY TEST CODE ---
	// Send a test message directly to verify P2P flow
	// This bypasses the TUI and Broadcaster for a direct test.
	go func() {
		// Give the network a second to fully initialize and discover peers
		time.Sleep(2 * time.Second)
		testMsg := &messaging.Message{
			ID:        "test-12345",
			Author:    netManager.Host.ID().String() + "-TESTER",
			Content:   "**Test Message**\n\nThis is a test message sent directly via `main.go` to verify P2P messaging. #test",
			Hashtags:  []string{"test", "general"},
			Timestamp: time.Now(),
			Type:      messaging.PostMsg,
		}
		testData, _ := json.Marshal(testMsg)

		// Encrypt if needed (simplified, uses own keys like broadcaster)
		var finalData []byte = testData
		if cfg.Privacy.EncryptMessages {
			encryptedData, err := crypto.Encrypt(testData, keyPair.PublicKey, keyPair.PrivateKey)
			if err == nil {
				finalData = encryptedData
				log.Println("Debug (Main Test Send): Test message encrypted.")
			} else {
				log.Printf("Debug (Main Test Send): Failed to encrypt test message: %v\n", err)
			}
		}

		log.Println("Debug (Main Test Send): Publishing test message to topic 'socli/hashtag/general'...")
		psManager.PublishMessage(ctx, topic, finalData)
	}()
	// --- END TEMPORARY TEST CODE ---

	go func() {
		for {
			fmt.Println("Debug (Subscription): Waiting for next message...") // Debug print
			msg, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., context cancellation)
				fmt.Printf("Debug (Subscription): Error receiving message: %v\n", err) // Debug print
				return
			}
			fmt.Printf("Debug (Subscription): Received a message from %s\n", msg.ReceivedFrom.String()) // Debug print
			// Don't display our own messages
			if msg.ReceivedFrom == netManager.Host.ID() {
				fmt.Println("Debug (Subscription): Ignoring own message.") // Debug print
				continue
			}

			// --- Message Processing ---
			data := msg.Data

			// Check if decryption is needed
			if cfg.Privacy.EncryptMessages {
				// Attempt to decrypt the message data using our private key.
				// We assume the sender encrypted it with our public key.
				// This is a simplification for prototype purposes.
				// A real implementation would require proper key exchange.
				decryptedData, ok := crypto.Decrypt(data, keyPair.PublicKey, keyPair.PrivateKey)
				if ok {
					data = decryptedData
				} else {
					// If decryption fails, the message might not be encrypted for us,
					// or it's not encrypted at all. Log and try to parse as plaintext JSON.
					// This is a simplification for the prototype.
					fmt.Println("Warning: Failed to decrypt message, attempting to parse as plaintext.")
				}
			}

			var receivedMsg messaging.Message
			if err := json.Unmarshal(data, &receivedMsg); err != nil {
				// Handle unmarshalling error
				fmt.Printf("Error unmarshalling message: %v\n", err)
				continue
			}
			fmt.Printf("Debug (Subscription): Successfully unmarshalled message ID %s from %s\n", receivedMsg.ID, receivedMsg.Author) // Debug print

			// Apply filters (placeholder)
			if !internal.ApplyFilters(&receivedMsg) {
				continue // Message was filtered out
			}

			// Send the processed message to the TUI
			p.Send(tui.PostReceivedMsg{Post: &receivedMsg})
			// --- End Message Processing ---
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}