package views

import (
	"socli/tui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// ComposeView is the view for composing new messages.
type ComposeView struct {
	input *components.Input
}

// NewComposeView creates a new compose view.
func NewComposeView() *ComposeView {
	return &ComposeView{
		input: components.NewInput(),
	}
}

// Init initializes the component.
func (v *ComposeView) Init() tea.Cmd {
	return v.input.Init()
}

// Update handles messages for the component.
func (v *ComposeView) Update(msg tea.Msg) (*ComposeView, tea.Cmd) {
	var cmd tea.Cmd
	v.input, cmd = v.input.Update(msg)
	return v, cmd
}

// View renders the component.
func (v *ComposeView) View() string {
	return v.input.View()
}

// Value returns the value of the input.
func (v *ComposeView) Value() string {
	return v.input.Value()
}
