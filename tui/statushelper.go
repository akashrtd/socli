package tui

import (
	"socli/tui/types"
)

// setStatus updates the model's status message.
func (m *AppModel) setStatus(msg types.StatusMsg) {
	m.statusMsg = &msg
	// TODO: Add logic to display status messages in the UI
	// For now, we just store the message in the model.
	// A full implementation would involve displaying the message in a dedicated area
	// of the TUI, such as a status bar or a notifications panel.
	// This might involve updating a field in the model that the View() method
	// uses to render the status message.
}