package tui

import (
	"context"
	"fmt"
	"socli/content"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"socli/tui/views"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// PostReceivedMsg is a message that is sent when a new post is received.
type PostReceivedMsg struct{ Post *messaging.Message }

// AppModel represents the main application model.
type AppModel struct {
	netManager  *p2p.NetworkManager
	store       *storage.MemoryStore
	renderer    *content.MarkdownRenderer
	composeView *views.ComposeView
	broadcaster *messaging.Broadcaster
	currentView string
}

// NewApp creates and returns a new application model.
func NewApp(netManager *p2p.NetworkManager, store *storage.MemoryStore, renderer *content.MarkdownRenderer, broadcaster *messaging.Broadcaster) (*AppModel, error) {
	return &AppModel{
		netManager:  netManager,
		store:       store,
		renderer:    renderer,
		composeView: views.NewComposeView(),
		broadcaster: broadcaster,
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
				msg := &messaging.Message{
					Author:    m.netManager.Host.ID().String(),
					Content:   m.composeView.Value(),
					Hashtags:  []string{"general"},
					Timestamp: time.Now(),
				}
				m.broadcaster.Broadcast(context.Background(), msg)
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

	var s strings.Builder
	s.WriteString("Welcome to socli!\n\n")

	posts := m.store.GetAllPosts()
	for _, post := range posts {
		rendered, err := m.renderer.Render(post.Content)
		if err != nil {
			rendered = post.Content
		}
		s.WriteString(fmt.Sprintf("From: %s\n%s\n\n", post.Author, rendered))
	}

	s.WriteString("---\n")
	s.WriteString(fmt.Sprintf("My Peer ID is: %s\n", m.netManager.Host.ID()))
	s.WriteString("Listening on:\n")
	for _, addr := range m.netManager.Host.Addrs() {
		s.WriteString(fmt.Sprintf("- %s\n", addr))
	}
	s.WriteString("\nPress 'c' to compose, 'q' to quit.")
	return s.String()
}
