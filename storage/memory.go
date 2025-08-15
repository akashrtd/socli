package storage

import (
	"socli/messaging"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
)

// MemoryStore provides in-memory storage for posts and peers.
type MemoryStore struct {
	posts map[string]*messaging.Message
	peers map[peer.ID]peer.AddrInfo
	mu    sync.RWMutex
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		posts: make(map[string]*messaging.Message),
		peers: make(map[peer.ID]peer.AddrInfo),
	}
}

// AddPost stores a new post in memory.
func (s *MemoryStore) AddPost(post *messaging.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[post.ID] = post
}

// GetPost retrieves a post by its ID.
func (s *MemoryStore) GetPost(id string) (*messaging.Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, found := s.posts[id]
	return post, found
}

// GetAllPosts returns all stored posts.
func (s *MemoryStore) GetAllPosts() []*messaging.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	posts := make([]*messaging.Message, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts
}

// AddPeer stores a new peer in memory.
func (s *MemoryStore) AddPeer(pi peer.AddrInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers[pi.ID] = pi
}

// GetPeer retrieves a peer by its ID.
func (s *MemoryStore) GetPeer(id peer.ID) (peer.AddrInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peer, found := s.peers[id]
	return peer, found
}

// GetAllPeers returns all stored peers.
func (s *MemoryStore) GetAllPeers() []peer.AddrInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers := make([]peer.AddrInfo, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	return peers
}
