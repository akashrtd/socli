package components

import (
	"fmt"
	"strings" // Add strings import

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Input is a wrapper around the bubbles textarea component.
type Input struct {
	textarea   textarea.Model
	maxLength  int
	charCountStyle lipgloss.Style
}

// NewInput creates a new input component.
func NewInput(maxLength int) *Input {
	ti := textarea.New()
	ti.Placeholder = "What's on your mind?"
	ti.Focus()
	
	// Set the maximum length if provided
	if maxLength > 0 {
		// Note: bubbles/textarea doesn't enforce max length by default.
		// We will handle this in the Update method.
	}
	
	charCountStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // Grey color
		Align(lipgloss.Right)

	return &Input{
		textarea:       ti,
		maxLength:      maxLength,
		charCountStyle: charCountStyle,
	}
}

// Init initializes the component.
func (i *Input) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the component.
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	var cmd tea.Cmd
	
	// Handle key messages to enforce max length
	if keyMsg, ok := msg.(tea.KeyMsg); ok && i.maxLength > 0 {
		currentLen := len([]rune(i.textarea.Value()))
		// Allow backspace, delete, arrow keys etc. even at max length
		// Only restrict character input.
		if currentLen >= i.maxLength {
			switch keyMsg.Type {
			case tea.KeyRunes:
				// Ignore keypress if at max length
				return i, cmd
			}
		}
	}

	i.textarea, cmd = i.textarea.Update(msg)
	return i, cmd
}

// View renders the component, including the character count.
func (i *Input) View() string {
	var b strings.Builder
	b.WriteString(i.textarea.View())
	
	if i.maxLength > 0 {
		currentLen := len([]rune(i.textarea.Value()))
		countStr := fmt.Sprintf("%d/%d", currentLen, i.maxLength)
		b.WriteString("\n")
		b.WriteString(i.charCountStyle.Render(countStr))
	}
	
	return b.String()
}

// Value returns the value of the input.
func (i *Input) Value() string {
	return i.textarea.Value()
}
