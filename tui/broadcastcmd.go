package tui

import (
	"socli/tui/types"
	tea "github.com/charmbracelet/bubbletea"
)

// listenForBroadcastResults returns a tea.Cmd that listens for StatusMsg
// on the AppModel's broadcastResultChan.
// When a status is received, it sends a BroadcastResultMsg to the model's Update function.
func (m *AppModel) listenForBroadcastResults() tea.Cmd {
	return func() tea.Msg {
		// This function will be called by BubbleTea in a goroutine
		// It blocks until a message is received on the channel
		status := <-m.broadcastResultChan
		// Wrap the status in a BroadcastResultMsg and return it
		// This msg will be passed to the Update function
		return types.BroadcastResultMsg{Status: status}
	}
}