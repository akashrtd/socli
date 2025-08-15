package messaging

import (
	"context"
	"encoding/json"
	"socli/p2p"
)

// Broadcaster handles the broadcasting of messages to the network.
type Broadcaster struct {
	psm *p2p.PubSubManager
}

// NewBroadcaster creates a new message broadcaster.
func NewBroadcaster(psm *p2p.PubSubManager) *Broadcaster {
	return &Broadcaster{psm: psm}
}

// Broadcast sends a message to all relevant topics.
func (b *Broadcaster) Broadcast(ctx context.Context, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, hashtag := range msg.Hashtags {
		topicName := GetTopicForHashtag(hashtag)
		topic, err := b.psm.JoinTopic(topicName)
		if err != nil {
			// Log or handle the error, but continue trying other hashtags
			continue
		}
		if err := b.psm.PublishMessage(ctx, topic, data); err != nil {
			// Log or handle the error
		}
	}

	return nil
}
