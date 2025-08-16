package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// listenForPostsCmd returns a tea.Cmd that listens for posts on the postChan.
// When a post is received, it sends a PostReceivedMsg to the model's Update function.
func (m *AppModel) listenForPostsCmd() tea.Cmd {
	return func() tea.Msg {
		// This function will be called by BubbleTea in a goroutine
		// It blocks until a message is received on the channel
		post := <-m.postChan
		// Wrap the post in a PostReceivedMsg and return it
		// This msg will be passed to the Update function
		return PostReceivedMsg{Post: post}
	}
}