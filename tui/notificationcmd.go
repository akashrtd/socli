package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// listenForNotifications returns a tea.Cmd that listens for status messages
// on the AppModel's broadcastResultChan.
// When a status is received, it sends a StatusMsg to the model's Update function.
func (m *AppModel) listenForNotifications() tea.Cmd {
	return func() tea.Msg {
		// This function will be called by BubbleTea in a goroutine
		// It blocks until a message is received on the channel
		status := <-m.broadcastResultChan
		// Wrap the status in a StatusMsg and return it
		// This msg will be passed to the Update function
		return status
	}
}