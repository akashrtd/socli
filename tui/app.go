package tui

import (
	"context"
	"fmt"
	"log"
	"socli/config"
	"socli/content"
	"socli/crypto"
	"socli/internal"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui/types"
	"socli/tui/views"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PostReceivedMsg is a message that is sent when a new post is received.
type PostReceivedMsg struct{ Post *messaging.Message }

// PeerConnectedMsg is a message sent when a new peer is discovered and connected.
type PeerConnectedMsg struct {
	PeerID string
}

// AppModel represents the main application model.
type AppModel struct {
	netManager      *p2p.NetworkManager
	store           *storage.MemoryStore
	renderer        *content.MarkdownRenderer
	composeView     *views.ComposeView
	feedView        *views.FeedView
	broadcaster     *messaging.Broadcaster
	psManager       *p2p.PubSubManager // Store PubSubManager for dynamic subscriptions
	keyPair         *crypto.KeyPair
	cfg             *config.Config
	currentView     string // "feed", "compose", or "help"
	subscriptions   map[string]*pubsub.Subscription // Map of topic names to subscriptions
	postChan        chan *messaging.Message         // Channel for receiving posts from dynamic subscriptions
	broadcastResultChan chan types.StatusMsg                  // Channel for receiving broadcast results
	terminalWidth   int                             // Store terminal width
	terminalHeight  int                             // Store terminal height
	statusMsg       *types.StatusMsg                      // For displaying status messages to the user
	helpScrollOffset int                           // Scroll offset for the help view
}

// NewApp creates and returns a new application model.
func NewApp(netManager *p2p.NetworkManager, store *storage.MemoryStore, renderer *content.MarkdownRenderer, psManager *p2p.PubSubManager, broadcaster *messaging.Broadcaster, keyPair *crypto.KeyPair, cfg *config.Config) (*AppModel, error) {
	// Initialize with default "general" subscription
	// Note: The primary subscription logic for the default topic will remain in main.go for now
	// to avoid duplication. AppModel will handle dynamic subscriptions.

	// Initialize the channel for receiving posts from dynamic subscriptions
	postChan := make(chan *messaging.Message, 10) // Buffered channel
	
	// Initialize the broadcast result channel
	broadcastResultChan := make(chan types.StatusMsg, 1) // Buffered to prevent blocking

	// Initialize with default terminal size, will be updated by tea.WindowSizeMsg
	width, height := 80, 24

	// Get our own Peer ID as a string for the UI
	// ownPeerID := netManager.Host.ID().String()

	return &AppModel{
		netManager:         netManager,
		store:              store,
		renderer:           renderer,
		composeView:        views.NewComposeView(cfg), // Pass config for max length
		feedView:        views.NewFeedView(store, renderer), // Pass store and renderer
		broadcaster:        broadcaster,
		psManager:          psManager, // Store psManager
		keyPair:            keyPair,
		cfg:                cfg,
		currentView:        "feed",
		subscriptions:      make(map[string]*pubsub.Subscription), // Initialize empty subscriptions map
		postChan:           postChan,                           // Initialize post channel
		broadcastResultChan: broadcastResultChan,              // Initialize broadcast result channel
		terminalWidth:      width,
		terminalHeight:     height,
		statusMsg:          nil, // No initial status message
		helpScrollOffset:   0,   // Initialize help scroll offset
	}, nil
}

// Init is the first function that will be called. It returns a command.
func (m *AppModel) Init() tea.Cmd {
	// Start the commands to listen for posts and broadcast results
	return tea.Batch(m.listenForPostsCmd(), m.listenForBroadcastResults())
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
				// Clear any status message when switching views
				m.statusMsg = nil
				return m, nil
			case "j":
				// Scroll down in feed (older posts)
				m.feedView.ScrollDown()
				return m, nil
			case "k":
				// Scroll up in feed (newer posts)
				m.feedView.ScrollUp()
				return m, nil
			case "?":
				// Toggle help view
				m.currentView = "help"
				return m, nil
			}
		case "compose":
			switch msg.String() {
			case "enter":
				content := m.composeView.Value()
				// Check if it's a command
				if strings.HasPrefix(strings.TrimSpace(content), "/") {
					// Parse and handle command
					command, args := ParseCommand(content)
					switch command {
					case "subscribe":
						if len(args) > 0 {
							hashtag := args[0]
							// Set status message
							m.statusMsg = &types.SubscribingMsg
							// Handle subscription logic
							m.subscribeToHashtag(hashtag)
							// Clear the input after command
							m.composeView = views.NewComposeView(m.cfg)
						}
						// Fall through to switch back to feed
					case "unsubscribe":
						if len(args) > 0 {
							hashtag := args[0]
							// Set status message
							m.statusMsg = &types.UnsubscribingMsg
							// Handle unsubscription logic
							m.unsubscribeFromHashtag(hashtag)
							// Clear the input after command
							m.composeView = views.NewComposeView(m.cfg)
						}
						// Fall through to switch back to feed
					default:
						// Show unknown command message
						m.statusMsg = &types.UnknownCmdMsg
						// Clear the input
						m.composeView = views.NewComposeView(m.cfg)
					}
					m.currentView = "feed"
					return m, nil
				}

				// Regular post
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
						log.Printf("Error signing message: %v\n", err)
					} else {
						msg.Signature = signature
					}

					// 3. Show "Publishing..." status
					m.statusMsg = &types.PostingMsg
					// Trigger a refresh to show the status immediately
					// We can return a simple 'nil' command or a custom one if needed
					// For now, the view will update on the next refresh cycle.
					
					// 4. Broadcast the message asynchronously
					// To provide feedback after broadcasting, we perform the broadcast
					// in a goroutine and then send a message back to the Update loop.
					// This requires defining new message types for success/failure.
					// Let's define them.
					
					// Perform broadcast in a goroutine to avoid blocking the UI
					// and send a result message back to the Update function.
					// We capture 'm' (the model) to access its fields, but we must be
					// careful not to mutate it directly in the goroutine.
					// Sending a message back is the safe way.
					go func() {
						err := m.broadcaster.Broadcast(context.Background(), msg)
						// Send a message back to the Update loop with the result
						// We need to get the program instance to send the message.
						// One way is to pass it to the model or have a callback.
						// BubbleTea programs can send messages using p.Send().
						// However, we don't have direct access to 'p' here.
						// A common pattern is to pass the program to the model
						// or use a channel. Let's re-evaluate this.
						
						// Placeholder for now, just log the result.
						// A full implementation would send a PostPublishedMsg or PostPublishFailedMsg
						// back to the Update loop.
						if err != nil {
							log.Printf("Error broadcasting message: %v", err)
							// Ideally, send PostPublishFailedMsg
						} else {
							log.Println("Message broadcast successfully")
							// Ideally, send PostSentMsg
						}
					}()

					// 5. Clear compose view and switch back to feed
					m.currentView = "feed"
					m.composeView = views.NewComposeView(m.cfg)
					return m, nil
				}
				// If content was empty, just go back to feed
				m.currentView = "feed"
				m.composeView = views.NewComposeView(m.cfg)
				return m, nil
			case "esc":
				m.currentView = "feed"
				// Reset compose view to clear content
				m.composeView = views.NewComposeView(m.cfg)
				// Clear status message
				m.statusMsg = nil
				return m, nil
			case "?":
				// Toggle help view from compose as well
				m.currentView = "help"
				return m, nil
			}
		case "help":
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "?", "esc":
				// Toggle help view off, go back to feed
				m.currentView = "feed"
				return m, nil
			case "j":
				// Scroll down in help view (older content)
				m.helpScrollOffset++
				return m, nil
			case "k":
				// Scroll up in help view (newer content)
				if m.helpScrollOffset > 0 {
					m.helpScrollOffset--
				}
				return m, nil
			}
			// For other keys in help view, do nothing
			return m, nil
		}
	case PostReceivedMsg:
		m.store.AddPost(msg.Post)
		return m, nil
	case PeerConnectedMsg:
		// A new peer has connected. Add it to our local store.
		// The peer ID is in msg.PeerID. We need to get AddrInfo.
		// We can try to get it from the host's peerstore.
		// This requires converting the string PeerID back to peer.ID.
		peerID, err := peer.Decode(msg.PeerID)
		if err != nil {
			log.Printf("Error decoding peer ID %s: %v", msg.PeerID, err)
			return m, nil
		}
		
		// Get AddrInfo from the host's peerstore
		// The peerstore should have been updated by the discovery process.
		pi := m.netManager.Host.Peerstore().PeerInfo(peerID)
		
		// Add the peer to our in-memory store
		m.store.AddPeer(pi)
		
		// We don't need to return a specific command here for the peer list update
		// as the View() method will read from the store directly.
		// If we had a more complex UI update mechanism, we might return a cmd.
		return m, nil
	case tea.WindowSizeMsg:
		// Handle terminal resize events
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		return m, nil
	// --- New Status Message Handlers ---
	case types.StatusMsg:
		// A status message has been sent to the model from an external source
		// (e.g., via p.Send() or another goroutine with direct access).
		// Update the model's statusMsg field.
		m.statusMsg = &msg
		return m, nil
	case types.BroadcastResultMsg:
		// A broadcast result has been received from the listenForBroadcastResults command.
		// Update the model's statusMsg field with the result.
		m.statusMsg = &msg.Status
		// Restart the listener command to continue receiving results.
		return m, m.listenForBroadcastResults()
	// --- End Status Message Handlers ---
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
	switch m.currentView {
	case "compose":
		return appStyle.Render(m.composeView.View())
	case "help":
		return m.renderHelpView()
	default: // "feed" view
		// --- Main Layout Construction ---
		// For simplicity, let's create a basic layout with a header, main content (feed),
		// a sidebar (peers/topics), and a status bar.
		// A more complex layout could use lipgloss.Layout or similar.

		// 1. Header
		header := headerStyle.Render("SOCLI - Decentralized P2P Social Platform")

		// 2. Main Content Area (Feed)
		// We need to calculate the space available for the main content
		// after accounting for header and status bar.
		// Let's assume header and status are 1 line each.
		// Sidebar will be a fixed width column on the right.
		// This is a simplified calculation.
		availableHeight := m.terminalHeight - 2 // Subtract header and status bar height
		if availableHeight < 0 {
			availableHeight = 0
		}
		availableWidth := m.terminalWidth
		if availableWidth < 0 {
			availableWidth = 0
		}

		// 3. Sidebar (Peers and Topics)
		sidebarWidth := 20
		mainContentWidth := availableWidth - sidebarWidth - 1 // -1 for potential spacing
		if mainContentWidth < 0 {
			mainContentWidth = availableWidth
			sidebarWidth = 0 // Hide sidebar if not enough space
		}

		// Get data for sidebar
		peers := m.store.GetAllPeers()
		topics := make([]string, 0, len(m.subscriptions))
		for topicName := range m.subscriptions {
			// Extract hashtag name from topic name, e.g., "socli/hashtag/general" -> "general"
			parts := strings.Split(topicName, "/")
			if len(parts) > 0 {
				topics = append(topics, "#"+parts[len(parts)-1])
			} else {
				topics = append(topics, topicName) // Fallback
			}
		}

		// Render Sidebar components
		peerList := m.renderPeerList(peers)
		topicList := m.renderTopicList(topics)

		// Combine peer and topic lists for the sidebar
		// Using lipgloss to join them vertically
		sidebarContent := lipgloss.JoinVertical(lipgloss.Left, peerList, topicList)
		sidebar := sidebarStyle.Width(sidebarWidth).Render(sidebarContent)

		// 4. Main Feed Content
		// Pass calculated dimensions to the feed view
		// Note: The feed view's scrolling logic might need adjustment based on actual rendered height.
		// For now, we pass the available width and a guess for height.
		// A full implementation would involve more complex height calculation or viewport management.
		feedContent := m.feedView.View(mainContentWidth, availableHeight-2) // Subtract some for internal feed padding/margins

		// 5. Status Bar
		// Determine the status message to display
		statusText := fmt.Sprintf("My Peer ID: %s | Press 'c' to compose, '?' for help, 'q' to quit", m.netManager.Host.ID().String())
		if m.statusMsg != nil {
			statusText = m.statusMsg.Message
		}
		statusBar := statusStyle.Width(m.terminalWidth).Render(statusText)

		// 6. Assemble the main view
		// Use lipgloss to lay out the main content and sidebar horizontally
		mainAreaView := lipgloss.JoinHorizontal(lipgloss.Top, feedContent, sidebar)

		// Join header, main area, and status bar vertically
		// appStyle adds overall padding
		view := lipgloss.JoinVertical(lipgloss.Left, header, mainAreaView, statusBar)

		return appStyle.Render(view)
	}
}

// renderPeerList creates a styled list of connected peers.
func (m *AppModel) renderPeerList(peers []peer.AddrInfo) string {
	title := headerStyle.Render("Peers")
	items := make([]string, 0, len(peers)+1)
	items = append(items, title)

	if len(peers) == 0 {
		items = append(items, listItemStyle.Render("No peers connected"))
	} else {
		for _, p := range peers {
			// Truncate or format peer ID for display
			peerIDStr := p.ID.String()
			if len(peerIDStr) > 15 {
				peerIDStr = peerIDStr[:8] + "..." + peerIDStr[len(peerIDStr)-6:]
			}
			items = append(items, listItemStyle.Render(peerIDStr))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderTopicList creates a styled list of subscribed topics.
func (m *AppModel) renderTopicList(topics []string) string {
	title := headerStyle.Render("Topics")
	items := make([]string, 0, len(topics)+1)
	items = append(items, title)

	if len(topics) == 0 {
		// This should ideally show the default "general" topic if not dynamically subscribed
		// For now, we just show an empty list if map is empty
		items = append(items, listItemStyle.Render("No topics subscribed"))
	} else {
		for _, t := range topics {
			items = append(items, listItemStyle.Render(t))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderHelpView creates and renders the help screen.
// (This function is defined in tui/help.go)