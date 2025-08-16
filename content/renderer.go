package content

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// MarkdownRenderer wraps the Glamour renderer and handles ASCII art.
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

// Render renders a Markdown string to a terminal-friendly format, preserving ASCII art.
func (r *MarkdownRenderer) Render(in string) (string, error) {
	// 1. Extract ASCII art blocks
	asciiBlocks := ExtractASCII(in)

	// If no ASCII art, render normally with Glamour
	if len(asciiBlocks) == 0 {
		out, err := r.renderer.Render(in)
		if err != nil {
			return "", err
		}
		return out, nil
	}

	// 2. Split content and render segments
	lines := strings.Split(in, "\n")
	var renderedSegments []string

	lastEndIndex := 0
	for _, block := range asciiBlocks {
		// Render text before the ASCII block
		if block.StartLine > lastEndIndex {
			preText := strings.Join(lines[lastEndIndex:block.StartLine], "\n")
			preRendered, err := r.renderer.Render(preText)
			if err != nil {
				// If Glamour fails, fall back to plain text
				renderedSegments = append(renderedSegments, preText)
			} else {
				renderedSegments = append(renderedSegments, preRendered)
			}
		}

		// Render the ASCII block itself with special styling
		// Use lipgloss to preserve whitespace and apply a monospace style
		asciiStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")). // Grey border
			Padding(0, 1).
			Align(lipgloss.Left)

		renderedASCII := asciiStyle.Render(block.Content)
		renderedSegments = append(renderedSegments, renderedASCII)

		lastEndIndex = block.EndLine + 1
	}

	// Render any remaining text after the last ASCII block
	if lastEndIndex < len(lines) {
		postText := strings.Join(lines[lastEndIndex:], "\n")
		postRendered, err := r.renderer.Render(postText)
		if err != nil {
			// If Glamour fails, fall back to plain text
			renderedSegments = append(renderedSegments, postText)
		} else {
			renderedSegments = append(renderedSegments, postRendered)
		}
	}

	// 3. Concatenate all rendered segments
	finalOutput := strings.Join(renderedSegments, "\n")
	return finalOutput, nil
}
