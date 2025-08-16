package tui

import (
	"socli/config"
	"socli/content"
	"socli/crypto"
	"socli/messaging"
	"socli/p2p"
	"socli/storage"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

// TestAppModelPostIntegration tests the integration of posting a message
// through the AppModel and checking if it appears in the store.
// This test simulates the receipt of a PostReceivedMsg by the AppModel.
func TestAppModelPostIntegration(t *testing.T) {
	// Create dependencies for AppModel
	cfg := config.DefaultConfig()
	store := storage.NewMemoryStore()
	renderer, err := content.NewMarkdownRenderer()
	if err != nil {
		t.Fatalf("Failed to create MarkdownRenderer: %v", err)
	}
	
	// For the test, we'll pass nil for the P2P dependencies since we're not testing P2P functionality here.
	// In a real application, these would be initialized.
	var netManager *p2p.NetworkManager = nil
	var psManager *p2p.PubSubManager = nil
	var broadcaster *messaging.Broadcaster = nil
	
	// Create a dummy key pair for signing
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create the AppModel
	appModel, err := NewApp(netManager, store, renderer, psManager, broadcaster, keyPair, cfg)
	if err != nil {
		t.Fatalf("Failed to create AppModel: %v", err)
	}

	// Create a test message
	testPost := &messaging.Message{
		ID:        "integration-test-id",
		Author:    peer.ID("integration-test-author").String(),
		Content:   "This is an integration test post",
		Hashtags:  []string{"integration", "test"},
		Timestamp: time.Now(),
		Type:      messaging.PostMsg,
	}

	// Create a PostReceivedMsg
	postReceivedMsg := PostReceivedMsg{Post: testPost}

	// Send the PostReceivedMsg to the AppModel's Update method
	// We don't need to check the returned model or command for this test
	_, _ = appModel.Update(postReceivedMsg)

	// Check if the post was added to the store
	allPosts := store.GetAllPosts()
	found := false
	for _, post := range allPosts {
		if post.ID == testPost.ID {
			found = true
			// Check if the post content matches
			if post.Content != testPost.Content {
				t.Errorf("Post content mismatch. Got %q, want %q", post.Content, testPost.Content)
			}
			break
		}
	}
	if !found {
		t.Errorf("Post with ID %q was not found in the store after PostReceivedMsg", testPost.ID)
	}
}