package messaging

import "time"

// MsgType defines the type of a message.
type MsgType string

const (
	// PostMsg is a standard post.
	PostMsg MsgType = "post"
	// ReplyMsg is a reply to another post.
	ReplyMsg MsgType = "reply"
	// ShareMsg is a share of another post.
	ShareMsg MsgType = "share"
)

// Message represents a message sent over the p2p network.
type Message struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`   // Peer ID
	Content   string    `json:"content"`  // Markdown content
	Hashtags  []string  `json:"hashtags"` // Extracted hashtags
	Timestamp time.Time `json:"timestamp"`
	Signature []byte    `json:"signature"` // Message signature
	Type      MsgType   `json:"type"`      // Post, Reply, Share
	ReplyTo   string    `json:"reply_to,omitempty"`
}
