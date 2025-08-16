// main.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"socli/config"
	"socli/content"
	"socli/crypto"
	"socli/internal"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui"
	"syscall"

	"github.com/libp2p/go-libp2p/core/peer"
	tea "github.com/charmbracelet/bubbletea"
)

// Version is the application version, set at build time.
var Version = "dev"

func main() {
	// Define command-line flags
	versionFlag := flag.Bool("version", false, "Print the version number and exit")
	flag.Parse()

	// If the version flag is set, print the version and exit
	if *versionFlag {
		fmt.Printf("SOCLI version %s\n", Version)
		os.Exit(0)
	}

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

	// Initialize the main application model from the tui package, passing the keyPair and config
	appModel, err := tui.NewApp(netManager, store, renderer, psManager, broadcaster, keyPair, cfg)
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

	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., context cancellation)
				return
			}
			// Don't display our own messages
			if msg.ReceivedFrom == netManager.Host.ID() {
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
					log.Println("Warning: Failed to decrypt message, attempting to parse as plaintext.")
				}
			}

			var receivedMsg messaging.Message
			if err := json.Unmarshal(data, &receivedMsg); err != nil {
				// Handle unmarshalling error
				log.Printf("Error unmarshalling message: %v\n", err)
				continue
			}

			// Apply filters (placeholder)
			if !internal.ApplyFilters(&receivedMsg) {
				continue // Message was filtered out
			}

			// Send the processed message to the TUI
			p.Send(tui.PostReceivedMsg{Post: &receivedMsg})
			// --- End Message Processing ---
		}
	}()

	// Set up a channel to listen for OS interrupt signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run the program in a goroutine so we can also listen for shutdown signals
	go func() {
		// Run the program
		// _, err := p.Run() // Capture the final model and error
		_, err := p.Run() // Just capture the error, ignore the final model
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}

		// --- Graceful Shutdown ---
		// The program has exited. Perform cleanup based on the final state.
		// We can check the final model's state if needed, though for a simple quit,
		// the context cancellation is usually sufficient.

		// Cancel the main context to signal shutdown to any background goroutines
		// (like the primary subscription listener)
		cancel()

		// Close the libp2p network manager
		fmt.Println("Shutting down libp2p host...")
		if err := netManager.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing network manager: %v\n", err)
			// Don't exit here, just log, as we're already shutting down
		}

		// Check configuration for auto-clear
		if cfg.Privacy.AutoClear {
			fmt.Println("Clearing in-memory data...")
			// Explicitly clear the in-memory store on shutdown if configured.
			store.Clear()
		}

		fmt.Println("Shutdown complete.")
		// --- End Graceful Shutdown ---

		// Exit the application
		os.Exit(0)
	}()

	// Wait for an interrupt signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal, shutting down...")

	// Cancel the context to signal shutdown to background goroutines
	cancel()

	// The program will exit when the BubbleTea program finishes
	// (which should happen quickly after context cancellation)
}