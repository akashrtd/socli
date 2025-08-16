package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

// Ensure PubSubManager implements PubSubManagerInterface.
var _ PubSubManagerInterface = (*PubSubManager)(nil)

// PubSubManager manages the GossipSub protocol for real-time messaging.
type PubSubManager struct {
	ps *pubsub.PubSub
}

// NewPubSubManager creates a new GossipSub router.
func NewPubSubManager(ctx context.Context, h host.Host) (*PubSubManager, error) {
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}
	return &PubSubManager{ps: ps}, nil
}

// JoinTopic subscribes to a given topic (hashtag).
func (psm *PubSubManager) JoinTopic(topicName string) (*pubsub.Topic, error) {
	return psm.ps.Join(topicName)
}

// PublishMessage broadcasts a message to a topic.
func (psm *PubSubManager) PublishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) error {
	return topic.Publish(ctx, data)
}

// SubscribeToTopic creates a subscription to a topic.
func (psm *PubSubManager) SubscribeToTopic(topic *pubsub.Topic) (*pubsub.Subscription, error) {
	return topic.Subscribe()
}
