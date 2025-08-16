package messaging

import (
	"context"
	"encoding/json"
	"log"
	"socli/config"
	"socli/crypto"
	"socli/p2p" // This import is needed for the interface
)

// Broadcaster handles the broadcasting of messages to the network.
type Broadcaster struct {
	psm     p2p.PubSubManagerInterface // Use the interface
	cfg     *config.Config
	keyPair *crypto.KeyPair // Dedicated key pair for application-level encryption/signing
}

// NewBroadcaster creates a new message broadcaster.
// It requires a PubSubManager, Config, and the application's KeyPair.
func NewBroadcaster(psm p2p.PubSubManagerInterface, cfg *config.Config, keyPair *crypto.KeyPair) *Broadcaster {
	return &Broadcaster{psm: psm, cfg: cfg, keyPair: keyPair}
}

// Broadcast sends a message to all relevant topics.
func (b *Broadcaster) Broadcast(ctx context.Context, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Check if encryption is enabled in the configuration
	if b.cfg.Privacy.EncryptMessages {
		// Encrypt the JSON message data using our own keys as an example.
		// In a real-world P2P broadcast scenario, true E2E encryption for specific recipients
		// is complex. This encrypts the payload before it's put on the wire,
		// adding a layer even on top of libp2p's transport security (Noise/TLS).
		encryptedData, err := crypto.Encrypt(data, b.keyPair.PublicKey, b.keyPair.PrivateKey)
		if err != nil {
			log.Printf("Error encrypting message ID %s: %v", msg.ID, err)
			return err // Handle encryption error
		}
		data = encryptedData
	}

	for _, hashtag := range msg.Hashtags {
		topicName := GetTopicForHashtag(hashtag)
		// Attempt to join the topic. If it's already joined, pubsub implementations
		// usually handle this gracefully or return a specific error.
		// We will log all errors for now, but not stop broadcasting.
		// A more advanced implementation might cache joined topics.
		topic, err := b.psm.JoinTopic(topicName)
		if err != nil {
			// Log the error but continue with other hashtags.
			// The "topic already exists" error is likely benign and can be ignored
			// if the goal is simply to ensure the topic is joined before publishing.
			// However, without knowing the specific error type, we log it for visibility.
			log.Printf("Warning: Could not join topic '%s' for message ID %s: %v", topicName, msg.ID, err)
			// Depending on the pubsub implementation, we might still be able to publish
			// to a topic we think we failed to join. This is a potential area for refinement.
			// For now, we skip publishing to this specific topic if join failed.
			continue
		}
		
		if err := b.psm.PublishMessage(ctx, topic, data); err != nil {
			// Log publish errors
			log.Printf("Error publishing to topic '%s' for message ID %s: %v", topicName, msg.ID, err)
			// Continue with other hashtags even if one fails
		}
	}

	return nil
}