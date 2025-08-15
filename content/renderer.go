package content

import (
	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer wraps the Glamour renderer.
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
}

// NewMarkdownRenderer creates a new Markdown renderer with a default style.
func NewMarkdownRenderer() (*MarkdownRenderer, error) {
	// Create a new Glamour renderer with a dark style
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
	)
	if err != nil {
		return nil, err
	}
	return &MarkdownRenderer{renderer: r}, nil
}

// Render renders a Markdown string to a terminal-friendly format.
func (r *MarkdownRenderer) Render(in string) (string, error) {
	out, err := r.renderer.Render(in)
	if err != nil {
		return "", err
	}
	return out, nil
}
