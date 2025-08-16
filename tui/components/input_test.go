package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestInput tests the Input component.
func TestInput(t *testing.T) {
	// Create an Input with a small max length for testing
	maxLength := 5
	input := NewInput(maxLength)

	// The Value() of a new Input should be an empty string
	if input.Value() != "" {
		t.Errorf("New Input Value() = %q, want %q", input.Value(), "")
	}

	// Test that the Input enforces the max length.
	// To do this, we need to simulate key presses.
	// The bubbles/textarea component handles key presses through tea.Msg.
	// We can send tea.KeyMsg to the Input's Update method.

	// Simulate typing "hello" (5 characters, which is the limit)
	for _, char := range "hello" {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		var cmd tea.Cmd
		input, cmd = input.Update(keyMsg)
		// We don't need to check the cmd for now
		_ = cmd
	}

	// The Value() should now be "hello"
	if input.Value() != "hello" {
		t.Errorf("After typing 'hello', Input Value() = %q, want %q", input.Value(), "hello")
	}

	// Simulate typing one more character, which should be ignored
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}}
	var cmd tea.Cmd
	input, cmd = input.Update(keyMsg)
	_ = cmd

	// The Value() should still be "hello"
	if input.Value() != "hello" {
		t.Errorf("After typing 'hello' and then '!', Input Value() = %q, want %q", input.Value(), "hello")
	}

	// Test backspace
	// Simulate pressing backspace
	backspaceMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	input, cmd = input.Update(backspaceMsg)
	_ = cmd

	// The Value() should now be "hell"
	if input.Value() != "hell" {
		t.Errorf("After backspace, Input Value() = %q, want %q", input.Value(), "hell")
	}

	// Test that we can type again after backspace
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}
	input, cmd = input.Update(keyMsg)
	_ = cmd

	// The Value() should now be "hello" again
	if input.Value() != "hello" {
		t.Errorf("After backspace and typing 'o', Input Value() = %q, want %q", input.Value(), "hello")
	}
	
	// Test with maxLength = 0 (no limit)
	inputNoLimit := NewInput(0)
	keyMsgNoLimit := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	var cmdNoLimit tea.Cmd
	inputNoLimit, cmdNoLimit = inputNoLimit.Update(keyMsgNoLimit)
	_ = cmdNoLimit
	if inputNoLimit.Value() != "a" {
		t.Errorf("Input with maxLength=0 Value() = %q, want %q", inputNoLimit.Value(), "a")
	}
}