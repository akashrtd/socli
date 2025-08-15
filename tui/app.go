package tui

import (
	"context"
	"fmt"
	"socli/content"
	"socli/crypto"
	"socli/internal"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui/views"
	"time"

	"github.com/google/uuid"
	tea "github.com/charmbracelet/bubbletea"
)

// PostReceivedMsg is a message that is sent when a new post is received.
type PostReceivedMsg struct{ Post *messaging.Message }

// PeerConnectedMsg is a message sent when a new peer is discovered and connected.
type PeerConnectedMsg struct {
	PeerID string
}

// AppModel represents the main application model.
type AppModel struct {
	netManager  *p2p.NetworkManager
	store       *storage.MemoryStore
	renderer    *content.MarkdownRenderer
	composeView *views.ComposeView
	feedView    *views.FeedView // Add FeedView
	broadcaster *messaging.Broadcaster
	keyPair     *crypto.KeyPair // Local node's keypair for signing
	currentView string
}

// NewApp creates and returns a new application model.
func NewApp(netManager *p2p.NetworkManager, store *storage.MemoryStore, renderer *content.MarkdownRenderer, broadcaster *messaging.Broadcaster, keyPair *crypto.KeyPair) (*AppModel, error) {
	return &AppModel{
		netManager:  netManager,
		store:       store,
		renderer:    renderer,
		composeView: views.NewComposeView(),
		feedView:    views.NewFeedView(store, renderer), // Initialize FeedView
		broadcaster: broadcaster,
		keyPair:     keyPair, // Pass the keypair for signing
		currentView: "feed",
	}, nil
}

// Init is the first function that will be called. It returns a command.
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentView {
		case "feed":
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "c":
				m.currentView = "compose"
				return m, nil
			}
		case "compose":
			switch msg.String() {
			case "enter":
				content := m.composeView.Value()
				fmt.Printf("Debug: Compose view content is '%s'\n", content) // Debug print
				if content != "" { // Only send non-empty messages
					// 1. Create Message struct
					msg := &messaging.Message{
						ID:        uuid.New().String(), // Generate unique ID
						Author:    m.netManager.Host.ID().String(),
						Content:   content,
						Hashtags:  internal.ExtractHashtags(content), // Extract hashtags
						Timestamp: time.Now(),
						Type:      messaging.PostMsg, // Set message type
					}

					// 2. Sign the message
					// For simplicity, we'll sign the content. A more robust approach
					// might sign a hash of the entire message or specific fields.
					signature, err := crypto.SignMessage([]byte(msg.Content), m.keyPair.PrivateKey)
					if err != nil {
						// TODO: Handle signing error in UI
						fmt.Printf("Error signing message: %v\n", err)
					} else {
						msg.Signature = signature
					}

					// 3. Broadcast the message
					fmt.Printf("Debug: Sending message - ID: %s, Content: '%s', Hashtags: %v\n", msg.ID, msg.Content, msg.Hashtags) // Debug print
					m.broadcaster.Broadcast(context.Background(), msg)
				}
				m.currentView = "feed"
				return m, nil
			case "esc":
				m.currentView = "feed"
				return m, nil
			}
		}
	case PostReceivedMsg:
		m.store.AddPost(msg.Post)
		return m, nil
	case PeerConnectedMsg:
		// TODO: Update store or UI with new peer info
		// For now, just log it
		fmt.Printf("TUI notified of peer connected: %s\n", msg.PeerID)
		return m, nil
	}

	if m.currentView == "compose" {
		var cmd tea.Cmd
		m.composeView, cmd = m.composeView.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the application's UI.
func (m *AppModel) View() string {
	if m.currentView == "compose" {
		return m.composeView.View()
	}

	// Default to feed view
	return m.feedView.View()
}