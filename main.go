package main

import (
	"fmt"
	"os"

	"context"
	"socli/config"
	"socli/content"
	"socli/internal"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui"

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

	// Initialize the main application model from the tui package
	appModel, err := tui.NewApp(netManager, store, renderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Create a new BubbleTea program
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	// Subscribe to a default topic and start listening for messages
	topic, err := psManager.JoinTopic("socli/hashtag/general")
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
				// Handle error
				return
			}
			// Don't display our own messages
			if msg.ReceivedFrom == netManager.Host.ID() {
				continue
			}
			p.Send(tui.PostReceivedMsg{Post: &messaging.Message{
				Author:  msg.ReceivedFrom.String(),
				Content: string(msg.Data),
			}})
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
