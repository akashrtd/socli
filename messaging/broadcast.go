package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"socli/config"
	"socli/crypto"
	"socli/p2p"
)

// Broadcaster handles the broadcasting of messages to the network.
type Broadcaster struct {
	psm     *p2p.PubSubManager
	cfg     *config.Config
	keyPair *crypto.KeyPair // Dedicated key pair for application-level encryption/signing
}

// NewBroadcaster creates a new message broadcaster.
// It requires a PubSubManager, Config, and the application's KeyPair.
func NewBroadcaster(psm *p2p.PubSubManager, cfg *config.Config, keyPair *crypto.KeyPair) *Broadcaster {
	return &Broadcaster{psm: psm, cfg: cfg, keyPair: keyPair}
}

// Broadcast sends a message to all relevant topics.
func (b *Broadcaster) Broadcast(ctx context.Context, msg *Message) error {
	fmt.Printf("Debug (Broadcaster): Preparing to broadcast message ID %s\n", msg.ID) // Debug print
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Check if encryption is enabled in the configuration
	if b.cfg.Privacy.EncryptMessages {
		fmt.Printf("Debug (Broadcaster): Encrypting message ID %s\n", msg.ID) // Debug print
		// Encrypt the JSON message data using our own keys as an example.
		// In a real-world P2P broadcast scenario, true E2E encryption for specific recipients
		// is complex. This encrypts the payload before it's put on the wire,
		// adding a layer even on top of libp2p's transport security (Noise/TLS).
		encryptedData, err := crypto.Encrypt(data, b.keyPair.PublicKey, b.keyPair.PrivateKey)
		if err != nil {
			return err // Handle encryption error
		}
		data = encryptedData
	}

	for _, hashtag := range msg.Hashtags {
		topicName := GetTopicForHashtag(hashtag)
		fmt.Printf("Debug (Broadcaster): Joining topic %s for message ID %s\n", topicName, msg.ID) // Debug print
		topic, err := b.psm.JoinTopic(topicName)
		if err != nil {
			// Log or handle the error, but continue trying other hashtags
			fmt.Printf("Debug (Broadcaster): Error joining topic %s: %v\n", topicName, err) // Debug print
			continue
		}
		fmt.Printf("Debug (Broadcaster): Publishing to topic %s for message ID %s\n", topicName, msg.ID) // Debug print
		if err := b.psm.PublishMessage(ctx, topic, data); err != nil {
			// Log or handle the error
			fmt.Printf("Debug (Broadcaster): Error publishing to topic %s: %v\n", topicName, err) // Debug print
		}
	}

	return nil
}