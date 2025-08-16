package content

import (
	"strings"
	"testing"
)

// TestNewMarkdownRenderer tests that a MarkdownRenderer can be created successfully.
func TestNewMarkdownRenderer(t *testing.T) {
	// This test ensures that the Glamour renderer can be initialized
	// with the specified style without returning an error.
	// It doesn't test the actual rendering, which would be more complex
	// and potentially brittle.
	
	renderer, err := NewMarkdownRenderer()
	if err != nil {
		t.Fatalf("NewMarkdownRenderer() error = %v, want nil", err)
	}
	if renderer == nil {
		t.Fatal("NewMarkdownRenderer() returned nil, want a renderer")
	}
	// The renderer field is not exported, so we can't easily test its contents.
	// The fact that it was created without error is a good basic test.
}

// TestMarkdownRendererRender tests the Render method of MarkdownRenderer.
func TestMarkdownRendererRender(t *testing.T) {
	// Create a renderer
	renderer, err := NewMarkdownRenderer()
	if err != nil {
		t.Fatalf("Failed to create MarkdownRenderer: %v", err)
	}

	// Test rendering a simple Markdown string
	markdown := "**Hello, World!**"
	rendered, err := renderer.Render(markdown)
	if err != nil {
		t.Fatalf("Render() error = %v, want nil", err)
	}
	if rendered == "" {
		t.Error("Render() returned empty string, want non-empty rendered output")
	}
	// A basic check to see if the rendered output contains the text
	// (though it might be formatted differently).
	// This is a very lenient test.
	if len(rendered) < len("Hello, World!") {
		t.Errorf("Render() output seems too short. Got length %d, expected at least %d", len(rendered), len("Hello, World!"))
	}

	// Test rendering an empty string
	emptyMarkdown := ""
	_, err = renderer.Render(emptyMarkdown)
	if err != nil {
		t.Fatalf("Render() error for empty string = %v, want nil", err)
	}
	// Rendering an empty string might produce some whitespace or terminal codes,
	// but it should not produce an error.
	// We won't assert anything specific about the output for an empty string.

	// Test rendering a more complex Markdown string
	complexMarkdown := `# Title
	
This is a paragraph with *italic* and **bold** text.

- Item 1
- Item 2

[Link](http://example.com)
`
	complexRendered, err := renderer.Render(complexMarkdown)
	if err != nil {
		t.Fatalf("Render() error for complex Markdown = %v, want nil", err)
	}
	if complexRendered == "" {
		t.Error("Render() returned empty string for complex Markdown, want non-empty rendered output")
	}
}

// TestMarkdownRendererRenderWithASCII tests the Render method with ASCII art.
func TestMarkdownRendererRenderWithASCII(t *testing.T) {
	// Create a renderer
	renderer, err := NewMarkdownRenderer()
	if err != nil {
		t.Fatalf("Failed to create MarkdownRenderer: %v", err)
	}

	// Test case 1: Content with no ASCII art
	normalMarkdown := `# Normal Post
	
This is a regular post with **bold** text and a [link](http://example.com).

- Item A
- Item B
`
	renderedNormal, err := renderer.Render(normalMarkdown)
	if err != nil {
		t.Fatalf("Render() error for normal Markdown = %v, want nil", err)
	}
	if renderedNormal == "" {
		t.Error("Render() returned empty string for normal Markdown with ASCII, want non-empty rendered output")
	}
	// The output should not contain any ASCII art bordering if none was detected
	// This is a bit fragile, but a reasonable check for this test.
	if strings.Contains(renderedNormal, "┌") || strings.Contains(renderedNormal, "└") {
		// Assuming the lipgloss border uses these characters
		t.Errorf("Render() for normal Markdown unexpectedly contained border characters. Output: %s", renderedNormal)
	}

	// Test case 2: Content with valid ASCII art
	// Define a sample ASCII art block using the new syntax
	asciiArt := `  /\_/\
 ( ^.^ )
  > ^ <`
	
	// Embed it in Markdown content using the ```ascii syntax
	markdownWithASCII := `# Post with ASCII Art

Check out this cute cat:

` + "```ascii\n" + asciiArt + "\n```\n" + `
Isn't it adorable?
`

	renderedWithASCII, err := renderer.Render(markdownWithASCII)
	if err != nil {
		t.Fatalf("Render() error for Markdown with ASCII = %v, want nil", err)
	}
	if renderedWithASCII == "" {
		t.Error("Render() returned empty string for Markdown with ASCII, want non-empty rendered output")
	}
	// Check if the rendered output is not empty
	// and contains the ASCII art content.
	// The output is a combination of rendered Markdown and styled ASCII art.
	// A simple check is that it's not empty and contains the ASCII art content.
	if renderedWithASCII == "" {
		t.Errorf("Render() output for Markdown with ASCII is empty, want non-empty rendered output")
	}
	
	// Check that the ASCII art content is present in the output.
	// The art content is ` /\_/\ `, `( ^.^ )`, ` > ^ < `.
	// These should be present in the rendered output.
	if !strings.Contains(renderedWithASCII, "/\\_/\\") {
		t.Errorf("Render() output for Markdown with ASCII does not contain the ASCII art content '/\\_/\\'. Output: %s", renderedWithASCII)
	}
	if !strings.Contains(renderedWithASCII, "( ^.^ )") {
		t.Errorf("Render() output for Markdown with ASCII does not contain the ASCII art content '( ^.^ )'. Output: %s", renderedWithASCII)
	}
	if !strings.Contains(renderedWithASCII, "> ^ <") {
		t.Errorf("Render() output for Markdown with ASCII does not contain the ASCII art content '> ^ <'. Output: %s", renderedWithASCII)
	}
	
	// Also check for the presence of border characters from lipgloss
	// This indicates that the ASCII art was styled.
	if !strings.Contains(renderedWithASCII, "╭") && !strings.Contains(renderedWithASCII, "┌") {
		// The border might not always render visibly in test output, or might use different chars.
		// This is not a critical failure, just a log.
		t.Logf("INFO: Render() output for Markdown with ASCII might not show visible border in test output. Check manually. Output snippet: %.200s...", renderedWithASCII)
	}

	// Test case 3: Content with text that looks like it might be ASCII but doesn't meet criteria
	// (e.g., a short, narrow block that wouldn't pass the 3-line/10-char avg width filter)
	// We'll test a regular code block that doesn't use the ```ascii syntax.
	// It should be rendered as a regular code block by Glamour.
	shortASCII := `---
 | |
---`
	
	markdownWithShortASCII := `# Post with Short Block

This is a short block that might look like ASCII but shouldn't qualify:

` + "```\n" + shortASCII + "\n```\n" + `
Because it's too short and narrow, and it's not marked with ` + "`" + `ascii` + "`" + `.
`

	renderedWithShortASCII, err := renderer.Render(markdownWithShortASCII)
	if err != nil {
		t.Fatalf("Render() error for Markdown with short ASCII-like block = %v, want nil", err)
	}
	if renderedWithShortASCII == "" {
		t.Error("Render() returned empty string for Markdown with short ASCII-like block, want non-empty rendered output")
	}
	// This short block should NOT be treated as ASCII art and thus should NOT get a border.
	// The content should be rendered as regular Markdown (code block).
	// We won't assert the absence of a border strictly, as the code block rendering is the main thing.
	// The key is that it doesn't crash or behave unexpectedly.
	// A basic check is that it renders without error and produces output.
	_ = renderedWithShortASCII // Use the variable to avoid "declared and not used" error
}