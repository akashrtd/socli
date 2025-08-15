package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Input is a wrapper around the bubbles textarea component.
type Input struct {
	textarea textarea.Model
}

// NewInput creates a new input component.
func NewInput() *Input {
	ti := textarea.New()
	ti.Placeholder = "What's on your mind?"
	ti.Focus()
	return &Input{textarea: ti}
}

// Init initializes the component.
func (i *Input) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the component.
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	var cmd tea.Cmd
	i.textarea, cmd = i.textarea.Update(msg)
	return i, cmd
}

// View renders the component.
func (i *Input) View() string {
	return i.textarea.View()
}

// Value returns the value of the input.
func (i *Input) Value() string {
	return i.textarea.Value()
}
