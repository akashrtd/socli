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
// The width and height parameters are provided by the main AppModel.View
// for potential future use (e.g., with a viewport).
func (v *FeedView) View(width, height int) string {
	var b strings.Builder

	posts := v.store.GetAllPosts()
	if len(posts) == 0 {
		b.WriteString("No posts yet. Be the first to post!\n")
		b.WriteString("Press 'c' to compose a new message.\n")
		return b.String()
	}

	// Display posts in reverse chronological order (newest first)
	// Simplified scrolling logic for now.
	// A more robust implementation would use a viewport or calculate
	// rendered content height.
	
	// For this prototype, let's display a fixed number of recent posts
	// or all posts if fewer than the limit.
	const maxPostsToDisplay = 50
	displayedCount := 0
	startIndex := len(posts) - 1 - v.offset
	if startIndex < 0 {
		startIndex = 0
	}

	// Adjust startIndex to ensure we don't exceed maxPostsToDisplay
	// unless offset pushes us beyond that.
	endIndex := startIndex - maxPostsToDisplay + 1
	if endIndex < 0 {
		endIndex = 0
	}

	for i := startIndex; i >= endIndex && i >= 0 && displayedCount < maxPostsToDisplay; i-- {
		post := posts[i]
		postLines := v.renderPost(post)
		b.WriteString(postLines)
		b.WriteString("\n---\n") // Separator between posts
		displayedCount++
	}

	// Indicate if scrolled
	if v.offset > 0 {
		b.WriteString(fmt.Sprintf("\n(Scrolled up by %d posts. Press 'j' to scroll down, 'k' to scroll up)\n", v.offset))
	}

	return b.String()
}

// ScrollUp moves the view up by one post (if possible).
func (v *FeedView) ScrollUp() {
	posts := v.store.GetAllPosts()
	// Allow scrolling up until we're showing the newest post at the bottom
	// of the display limit.
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