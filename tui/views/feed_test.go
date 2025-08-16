package views

import (
	"socli/content"
	"socli/messaging"
	"socli/storage"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

// TestFeedViewScrolling tests the scrolling functionality of FeedView.
func TestFeedViewScrolling(t *testing.T) {
	// Create a mock store and renderer for the FeedView
	store := storage.NewMemoryStore()
	renderer, err := content.NewMarkdownRenderer()
	if err != nil {
		t.Fatalf("Failed to create MarkdownRenderer: %v", err)
	}

	// Create the FeedView
	feedView := NewFeedView(store, renderer)

	// Initially, there should be no posts, so scrolling should have no effect
	// and the offset should remain 0.
	initialOffset := feedView.offset
	feedView.ScrollDown()
	if feedView.offset != initialOffset {
		t.Errorf("ScrollDown() with no posts changed offset from %d to %d, want no change", initialOffset, feedView.offset)
	}
	feedView.ScrollUp()
	if feedView.offset != initialOffset {
		t.Errorf("ScrollUp() with no posts changed offset from %d to %d, want no change", initialOffset, feedView.offset)
	}

	// Add some posts to the store
	for i := 1; i <= 5; i++ {
		post := &messaging.Message{
			ID:        string(rune(i)),
			Author:    peer.ID("test-author").String(),
			Content:   "Test post content",
			Hashtags:  []string{"test"},
			Timestamp: time.Now(),
			Type:      messaging.PostMsg,
		}
		store.AddPost(post)
	}

	// Now test scrolling with posts
	// Initial offset should still be 0
	if feedView.offset != 0 {
		t.Errorf("Initial offset is %d, want 0", feedView.offset)
	}

	// Scrolling down should not change the offset from 0
	feedView.ScrollDown()
	if feedView.offset != 0 {
		t.Errorf("ScrollDown() from offset 0 changed offset to %d, want 0", feedView.offset)
	}

	// Scrolling up should increase the offset
	feedView.ScrollUp()
	if feedView.offset != 1 {
		t.Errorf("ScrollUp() from offset 0 changed offset to %d, want 1", feedView.offset)
	}

	// Scroll up again
	feedView.ScrollUp()
	if feedView.offset != 2 {
		t.Errorf("ScrollUp() from offset 1 changed offset to %d, want 2", feedView.offset)
	}

	// Scroll down
	feedView.ScrollDown()
	if feedView.offset != 1 {
		t.Errorf("ScrollDown() from offset 2 changed offset to %d, want 1", feedView.offset)
	}

	// Test scrolling beyond the number of posts
	// There are 5 posts, so the maximum offset should be 4 (to show the oldest post at the top)
	// Let's scroll up a few more times to reach the limit
	feedView.ScrollUp() // offset 2
	feedView.ScrollUp() // offset 3
	feedView.ScrollUp() // offset 4
	if feedView.offset != 4 {
		t.Errorf("After scrolling up 3 more times, offset is %d, want 4", feedView.offset)
	}

	// One more scroll up should not change the offset
	feedView.ScrollUp()
	if feedView.offset != 4 {
		t.Errorf("ScrollUp() from offset 4 changed offset to %d, want 4", feedView.offset)
	}

	// Now scroll all the way back down
	feedView.ScrollDown() // offset 3
	feedView.ScrollDown() // offset 2
	feedView.ScrollDown() // offset 1
	feedView.ScrollDown() // offset 0
	if feedView.offset != 0 {
		t.Errorf("After scrolling down 4 times, offset is %d, want 0", feedView.offset)
	}

	// One more scroll down should not change the offset
	feedView.ScrollDown()
	if feedView.offset != 0 {
		t.Errorf("ScrollDown() from offset 0 changed offset to %d, want 0", feedView.offset)
	}
}