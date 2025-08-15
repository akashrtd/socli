package views

import (
	"fmt"
	"socli/content"
	"socli/messaging"
	"socli/storage"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// FeedView displays the main timeline of posts.
type FeedView struct {
	store    *storage.MemoryStore
	renderer *content.MarkdownRenderer
	offset   int // Vertical scroll offset
}

// NewFeedView creates a new feed view.
func NewFeedView(store *storage.MemoryStore, renderer *content.MarkdownRenderer) *FeedView {
	return &FeedView{
		store:    store,
		renderer: renderer,
		offset:   0,
	}
}

// View renders the feed, considering the scroll offset.
func (v *FeedView) View(width, height int) string {
	var b strings.Builder

	posts := v.store.GetAllPosts()
	if len(posts) == 0 {
		b.WriteString("No posts yet. Be the first to post!\n")
		b.WriteString("Press 'c' to compose a new message.\n")
		return b.String()
	}

	// Display posts in reverse chronological order (newest first)
	// We need to calculate how many posts can fit in the given height.
	// This is a simplified calculation. A full implementation might
	// involve measuring rendered content height.
	availableLines := height - 2 // Reserve lines for header/footer if needed
	if availableLines <= 0 {
		availableLines = 1
	}

	displayedPosts := 0
	startIndex := len(posts) - 1 - v.offset
	if startIndex < 0 {
		startIndex = 0
	}

	for i := startIndex; i >= 0 && displayedPosts < availableLines; i-- {
		post := posts[i]
		postLines := v.renderPost(post)
		// A very rough estimate: assume post takes a few lines.
		// A more accurate way would be to split `postLines` and count.
		// For now, we'll use a simple heuristic or just display a fixed number.
		// Let's simplify and show a fixed number of posts for now,
		// and refine scrolling logic later.
		b.WriteString(postLines)
		b.WriteString("\n---\n") // Separator between posts
		displayedPosts++
	}

	// If we are scrolled, indicate it
	if v.offset > 0 {
		b.WriteString(fmt.Sprintf("\n(Scrolled up by %d posts. Press 'j' to scroll down, 'k' to scroll up)\n", v.offset))
	}

	return b.String()
}

// ScrollUp moves the view up by one post (if possible).
func (v *FeedView) ScrollUp() {
	posts := v.store.GetAllPosts()
	if v.offset < len(posts)-1 {
		v.offset++
	}
}

// ScrollDown moves the view down by one post (if possible).
func (v *FeedView) ScrollDown() {
	if v.offset > 0 {
		v.offset--
	}
}

// renderPost formats a single post for display.
func (v *FeedView) renderPost(post *messaging.Message) string {
	var b strings.Builder

	// Author and Timestamp
	authorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")) // Purple
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))           // Grey
	b.WriteString(authorStyle.Render(fmt.Sprintf("From: %s", post.Author)))
	b.WriteString(" ")
	b.WriteString(timeStyle.Render(post.Timestamp.Format(time.Stamp)))
	b.WriteString("\n")

	// Content with Markdown rendering
	renderedContent, err := v.renderer.Render(post.Content)
	if err != nil {
		// Fallback to plain text if rendering fails
		renderedContent = post.Content
	}
	b.WriteString(renderedContent)
	b.WriteString("\n")

	// Hashtags
	if len(post.Hashtags) > 0 {
		hashtagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33")) // Blue
		b.WriteString(hashtagStyle.Render(fmt.Sprintf("Tags: #%s", strings.Join(post.Hashtags, " #"))))
		b.WriteString("\n")
	}

	return b.String()
}