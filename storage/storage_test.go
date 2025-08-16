package storage

import (
	"socli/messaging"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/test"
	"github.com/multiformats/go-multiaddr"
)

// TestMemoryStorePosts tests adding, getting, and getting all posts.
func TestMemoryStorePosts(t *testing.T) {
	store := NewMemoryStore()

	// Test adding a post
	post1 := &messaging.Message{
		ID:        "1",
		Author:    "author1",
		Content:   "This is post 1",
		Hashtags:  []string{"test"},
		Timestamp: time.Now(),
		Type:      messaging.PostMsg,
	}
	store.AddPost(post1)

	// Test getting the post
	gotPost, found := store.GetPost("1")
	if !found {
		t.Fatal("Post not found")
	}
	if gotPost != post1 {
		// Since we're storing pointers, this should be the same object
		t.Errorf("GetPost() = %v, want %v", gotPost, post1)
	}

	// Test getting a non-existent post
	_, found = store.GetPost("2")
	if found {
		t.Error("GetPost() should return false for non-existent post")
	}

	// Test adding another post
	post2 := &messaging.Message{
		ID:        "2",
		Author:    "author2",
		Content:   "This is post 2",
		Hashtags:  []string{"test", "golang"},
		Timestamp: time.Now(),
		Type:      messaging.PostMsg,
	}
	store.AddPost(post2)

	// Test getting all posts
	allPosts := store.GetAllPosts()
	if len(allPosts) != 2 {
		t.Errorf("GetAllPosts() returned %d posts, want 2", len(allPosts))
	}

	// The order of posts in GetAllPosts is not guaranteed, so we just check if both are present
	foundPost1 := false
	foundPost2 := false
	for _, p := range allPosts {
		if p.ID == "1" {
			foundPost1 = true
		}
		if p.ID == "2" {
			foundPost2 = true
		}
	}
	if !foundPost1 {
		t.Error("Post 1 not found in GetAllPosts()")
	}
	if !foundPost2 {
		t.Error("Post 2 not found in GetAllPosts()")
	}
}

// TestMemoryStorePeers tests adding, getting, and getting all peers.
func TestMemoryStorePeers(t *testing.T) {
	store := NewMemoryStore()

	// Create a test peer ID and address
	peerID, err := test.RandPeerID()
	if err != nil {
		t.Fatalf("Failed to generate test peer ID: %v", err)
	}
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	if err != nil {
		t.Fatalf("Failed to create multiaddr: %v", err)
	}

	peerInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{addr},
	}

	// Test adding a peer
	store.AddPeer(peerInfo)

	// Test getting the peer
	gotPeer, found := store.GetPeer(peerID)
	if !found {
		t.Fatal("Peer not found")
	}
	if gotPeer.ID != peerInfo.ID {
		t.Errorf("GetPeer() ID = %v, want %v", gotPeer.ID, peerInfo.ID)
	}
	// For simplicity, we'll just check the ID. A more thorough test would check addresses too.

	// Test getting a non-existent peer
	peerID2, err := test.RandPeerID()
	if err != nil {
		t.Fatalf("Failed to generate test peer ID: %v", err)
	}
	_, found = store.GetPeer(peerID2)
	if found {
		t.Error("GetPeer() should return false for non-existent peer")
	}

	// Test adding another peer
	peerID3, err := test.RandPeerID()
	if err != nil {
		t.Fatalf("Failed to generate test peer ID: %v", err)
	}
	addr3, err := multiaddr.NewMultiaddr("/ip4/192.168.1.1/tcp/4002")
	if err != nil {
		t.Fatalf("Failed to create multiaddr: %v", err)
	}
	peerInfo3 := peer.AddrInfo{
		ID:    peerID3,
		Addrs: []multiaddr.Multiaddr{addr3},
	}
	store.AddPeer(peerInfo3)

	// Test getting all peers
	allPeers := store.GetAllPeers()
	if len(allPeers) != 2 {
		t.Errorf("GetAllPeers() returned %d peers, want 2", len(allPeers))
	}

	// Check if both peers are present
	foundPeer1 := false
	foundPeer3 := false
	for _, p := range allPeers {
		if p.ID == peerID {
			foundPeer1 = true
		}
		if p.ID == peerID3 {
			foundPeer3 = true
		}
	}
	if !foundPeer1 {
		t.Error("Peer 1 not found in GetAllPeers()")
	}
	if !foundPeer3 {
		t.Error("Peer 3 not found in GetAllPeers()")
	}
}

// TestMemoryStoreClear tests the Clear method.
func TestMemoryStoreClear(t *testing.T) {
	store := NewMemoryStore()

	// Add a post and a peer
	post := &messaging.Message{
		ID:        "1",
		Author:    "author1",
		Content:   "This is a test post",
		Hashtags:  []string{"test"},
		Timestamp: time.Now(),
		Type:      messaging.PostMsg,
	}
	store.AddPost(post)

	peerID, err := test.RandPeerID()
	if err != nil {
		t.Fatalf("Failed to generate test peer ID: %v", err)
	}
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	if err != nil {
		t.Fatalf("Failed to create multiaddr: %v", err)
	}
	peerInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{addr},
	}
	store.AddPeer(peerInfo)

	// Verify they are there
	if _, found := store.GetPost("1"); !found {
		t.Error("Post should be found before clear")
	}
	if _, found := store.GetPeer(peerID); !found {
		t.Error("Peer should be found before clear")
	}

	// Clear the store
	store.Clear()

	// Verify they are gone
	if _, found := store.GetPost("1"); found {
		t.Error("Post should not be found after clear")
	}
	if _, found := store.GetPeer(peerID); found {
		t.Error("Peer should not be found after clear")
	}

	allPosts := store.GetAllPosts()
	if len(allPosts) != 0 {
		t.Errorf("GetAllPosts() returned %d posts after clear, want 0", len(allPosts))
	}

	allPeers := store.GetAllPeers()
	if len(allPeers) != 0 {
		t.Errorf("GetAllPeers() returned %d peers after clear, want 0", len(allPeers))
	}
}