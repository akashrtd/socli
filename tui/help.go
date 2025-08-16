package tui

// renderHelpView creates and renders the help screen.
// This function is implemented in tui/help_content.go
func (m *AppModel) renderHelpView() string {
	return m.renderHelpViewInternal()
}