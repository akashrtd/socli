package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// PubSubManagerInterface defines the interface for PubSubManager.
// This is useful for mocking in tests.
type PubSubManagerInterface interface {
	JoinTopic(topicName string) (*pubsub.Topic, error)
	PublishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) error
	SubscribeToTopic(topic *pubsub.Topic) (*pubsub.Subscription, error)
}