package messaging

import (
	"testing"
	"time"
)

// TestMessageCreation tests the creation of a Message struct.
func TestMessageCreation(t *testing.T) {
	// Create a new message
	msg := &Message{
		ID:        "test-id-123",
		Author:    "test-author-peer-id",
		Content:   "**Hello, World!** This is a *test* message.",
		Hashtags:  []string{"test", "golang"},
		Timestamp: time.Now(),
		Signature: []byte("test-signature"),
		Type:      PostMsg,
		ReplyTo:   "",
	}

	// Verify the fields
	if msg.ID != "test-id-123" {
		t.Errorf("Message.ID = %s, want test-id-123", msg.ID)
	}

	if msg.Author != "test-author-peer-id" {
		t.Errorf("Message.Author = %s, want test-author-peer-id", msg.Author)
	}

	if msg.Content != "**Hello, World!** This is a *test* message." {
		t.Errorf("Message.Content = %s, want **Hello, World!** This is a *test* message.", msg.Content)
	}

	if len(msg.Hashtags) != 2 {
		t.Errorf("len(Message.Hashtags) = %d, want 2", len(msg.Hashtags))
	}

	if msg.Hashtags[0] != "test" {
		t.Errorf("Message.Hashtags[0] = %s, want test", msg.Hashtags[0])
	}

	if msg.Hashtags[1] != "golang" {
		t.Errorf("Message.Hashtags[1] = %s, want golang", msg.Hashtags[1])
	}

	if msg.Type != PostMsg {
		t.Errorf("Message.Type = %s, want PostMsg", msg.Type)
	}

	if msg.ReplyTo != "" {
		t.Errorf("Message.ReplyTo = %s, want empty string", msg.ReplyTo)
	}
}