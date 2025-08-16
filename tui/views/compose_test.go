package views

import (
	"socli/config"
	"testing"
)

// TestComposeView tests the ComposeView component.
func TestComposeView(t *testing.T) {
	// Create a dummy config with a small max post length for testing
	cfg := &config.Config{
		UI: struct {
			Theme         string `yaml:"theme"`
			RefreshRate   int    `yaml:"refresh_rate_ms"`
			MaxPostLength int    `yaml:"max_post_length"`
		}{
			MaxPostLength: 10, // Set a small limit for easy testing
		},
	}

	// Create the ComposeView
	composeView := NewComposeView(cfg)

	// The Value() of a new ComposeView should be an empty string
	if composeView.Value() != "" {
		t.Errorf("New ComposeView Value() = %q, want %q", composeView.Value(), "")
	}

	// The underlying Input component should enforce the max length.
	// However, testing the actual input behavior (typing characters) is complex
	// with the bubbletea/bubbles components, as they are designed to work
	// within the bubbletea runtime.
	// We can at least verify that the ComposeView was created with the correct
	// max length by inspecting its behavior indirectly or by testing the
	// components.Input directly.
	
	// For now, let's just test that the ComposeView can be created and that
	// its Value() method works. The character counting logic is tested
	// in the components package.
}