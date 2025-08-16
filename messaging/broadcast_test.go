package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"socli/config"
	"socli/crypto"
	"testing"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// mockPubSubManager is a mock implementation of p2p.PubSubManagerInterface for testing.
type mockPubSubManager struct {
	joinTopicFunc      func(topicName string) (*pubsub.Topic, error)
	publishMessageFunc func(ctx context.Context, topic *pubsub.Topic, data []byte) error
	subscribeToTopicFunc func(topic *pubsub.Topic) (*pubsub.Subscription, error)
}

// JoinTopic mocks the JoinTopic method.
func (m *mockPubSubManager) JoinTopic(topicName string) (*pubsub.Topic, error) {
	if m.joinTopicFunc != nil {
		return m.joinTopicFunc(topicName)
	}
	// Default mock behavior: return a nil topic and no error
	return nil, nil
}

// PublishMessage mocks the PublishMessage method.
func (m *mockPubSubManager) PublishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) error {
	if m.publishMessageFunc != nil {
		return m.publishMessageFunc(ctx, topic, data)
	}
	// Default mock behavior: return no error
	return nil
}

// SubscribeToTopic mocks the SubscribeToTopic method.
func (m *mockPubSubManager) SubscribeToTopic(topic *pubsub.Topic) (*pubsub.Subscription, error) {
	if m.subscribeToTopicFunc != nil {
		return m.subscribeToTopicFunc(topic)
	}
	// Default mock behavior: return a nil subscription and no error
	return nil, nil
}

// TestBroadcasterBroadcast tests the Broadcast method of Broadcaster.
func TestBroadcasterBroadcast(t *testing.T) {
	// Create a test key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create a test config
	cfg := &config.Config{
		Privacy: struct {
			EncryptMessages bool   `yaml:"encrypt_messages"`
			KeyPath         string `yaml:"key_path"`
			AutoClear       bool   `yaml:"auto_clear_on_exit"`
		}{
			EncryptMessages: false, // Start with encryption off for simplicity
		},
	}

	// Create a test message
	msg := &Message{
		ID:        "test-id",
		Author:    "test-author",
		Content:   "This is a test message",
		Hashtags:  []string{"test"},
		Timestamp: time.Now(),
		Type:      PostMsg,
	}

	// Test successful broadcast without encryption
	t.Run("SuccessfulBroadcastWithoutEncryption", func(t *testing.T) {
		// Create a mock PubSubManager
		mockPSM := &mockPubSubManager{
			joinTopicFunc: func(topicName string) (*pubsub.Topic, error) {
				// Verify the topic name is correct
				expectedTopicName := GetTopicForHashtag("test")
				if topicName != expectedTopicName {
					t.Errorf("JoinTopic called with topic %s, want %s", topicName, expectedTopicName)
				}
				// Return a nil topic and no error for simplicity
				return nil, nil
			},
			publishMessageFunc: func(ctx context.Context, topic *pubsub.Topic, data []byte) error {
				// Verify the data is the marshaled message
				var publishedMsg Message
				if err := json.Unmarshal(data, &publishedMsg); err != nil {
					t.Errorf("Failed to unmarshal published data: %v", err)
				}
				if publishedMsg.ID != msg.ID {
					t.Errorf("Published message ID = %s, want %s", publishedMsg.ID, msg.ID)
				}
				// Return no error for success
				return nil
			},
		}

		// Create the broadcaster
		broadcaster := NewBroadcaster(mockPSM, cfg, keyPair)

		// Perform the broadcast
		err := broadcaster.Broadcast(context.Background(), msg)
		if err != nil {
			t.Errorf("Broadcast() error = %v, want nil", err)
		}
	})

	// Test successful broadcast with encryption
	t.Run("SuccessfulBroadcastWithEncryption", func(t *testing.T) {
		// Modify config to enable encryption
		cfg.Privacy.EncryptMessages = true

		// Create a mock PubSubManager
		mockPSM := &mockPubSubManager{
			joinTopicFunc: func(topicName string) (*pubsub.Topic, error) {
				// Verify the topic name is correct
				expectedTopicName := GetTopicForHashtag("test")
				if topicName != expectedTopicName {
					t.Errorf("JoinTopic called with topic %s, want %s", topicName, expectedTopicName)
				}
				// Return a nil topic and no error for simplicity
				return nil, nil
			},
			publishMessageFunc: func(ctx context.Context, topic *pubsub.Topic, data []byte) error {
				// With encryption enabled, the data should be encrypted
				// We can't easily verify the content, but we can check that it's not the original JSON
				originalData, _ := json.Marshal(msg)
				if string(data) == string(originalData) {
					t.Error("Published data is not encrypted, want encrypted data")
				}
				// Return no error for success
				return nil
			},
		}

		// Create the broadcaster
		broadcaster := NewBroadcaster(mockPSM, cfg, keyPair)

		// Perform the broadcast
		err := broadcaster.Broadcast(context.Background(), msg)
		if err != nil {
			t.Errorf("Broadcast() error = %v, want nil", err)
		}

		// Reset config for other tests
		cfg.Privacy.EncryptMessages = false
	})

	// Test broadcast with multiple hashtags
	t.Run("BroadcastWithMultipleHashtags", func(t *testing.T) {
		// Create a message with multiple hashtags
		multiHashtagMsg := &Message{
			ID:        "multi-hashtag-id",
			Author:    "test-author",
			Content:   "This is a test message with multiple hashtags",
			Hashtags:  []string{"test", "golang", "decentralized"},
			Timestamp: time.Now(),
			Type:      PostMsg,
		}

		// Counter for joinTopic calls
		joinTopicCalls := 0
		expectedTopics := []string{
			GetTopicForHashtag("test"),
			GetTopicForHashtag("golang"),
			GetTopicForHashtag("decentralized"),
		}

		// Create a mock PubSubManager
		mockPSM := &mockPubSubManager{
			joinTopicFunc: func(topicName string) (*pubsub.Topic, error) {
				// Verify the topic name is one of the expected ones
				found := false
				for _, expected := range expectedTopics {
					if topicName == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("JoinTopic called with unexpected topic %s", topicName)
				}
				joinTopicCalls++
				// Return a nil topic and no error for simplicity
				return nil, nil
			},
			publishMessageFunc: func(ctx context.Context, topic *pubsub.Topic, data []byte) error {
				// Verify the data is the marshaled message
				var publishedMsg Message
				if err := json.Unmarshal(data, &publishedMsg); err != nil {
					t.Errorf("Failed to unmarshal published data: %v", err)
				}
				if publishedMsg.ID != multiHashtagMsg.ID {
					t.Errorf("Published message ID = %s, want %s", publishedMsg.ID, multiHashtagMsg.ID)
				}
				// Return no error for success
				return nil
			},
		}

		// Create the broadcaster
		broadcaster := NewBroadcaster(mockPSM, cfg, keyPair)

		// Perform the broadcast
		err := broadcaster.Broadcast(context.Background(), multiHashtagMsg)
		if err != nil {
			t.Errorf("Broadcast() error = %v, want nil", err)
		}

		// Verify that JoinTopic was called for each hashtag
		if joinTopicCalls != len(expectedTopics) {
			t.Errorf("JoinTopic was called %d times, want %d", joinTopicCalls, len(expectedTopics))
		}
	})

	// Test broadcast with join topic error
	t.Run("BroadcastWithJoinTopicError", func(t *testing.T) {
		// Create a mock PubSubManager that returns an error on JoinTopic
		mockPSM := &mockPubSubManager{
			joinTopicFunc: func(topicName string) (*pubsub.Topic, error) {
				// Return an error for any topic
				return nil, errors.New("join topic error")
			},
			publishMessageFunc: func(ctx context.Context, topic *pubsub.Topic, data []byte) error {
				// This should not be called if JoinTopic fails
				t.Error("PublishMessage was called unexpectedly")
				return nil
			},
		}

		// Create the broadcaster
		broadcaster := NewBroadcaster(mockPSM, cfg, keyPair)

		// Perform the broadcast
		// We expect the broadcast to succeed even if joining a topic fails,
		// as it should continue with other hashtags.
		err := broadcaster.Broadcast(context.Background(), msg)
		if err != nil {
			t.Errorf("Broadcast() error = %v, want nil", err)
		}
	})

	// Test broadcast with publish message error
	t.Run("BroadcastWithPublishMessageError", func(t *testing.T) {
		// Create a mock PubSubManager that returns an error on PublishMessage
		mockPSM := &mockPubSubManager{
			joinTopicFunc: func(topicName string) (*pubsub.Topic, error) {
				// Return a nil topic and no error for simplicity
				return nil, nil
			},
			publishMessageFunc: func(ctx context.Context, topic *pubsub.Topic, data []byte) error {
				// Return an error for any publish
				return errors.New("publish message error")
			},
		}

		// Create the broadcaster
		broadcaster := NewBroadcaster(mockPSM, cfg, keyPair)

		// Perform the broadcast
		// We expect the broadcast to succeed even if publishing fails,
		// as it should continue with other hashtags and log the error.
		err := broadcaster.Broadcast(context.Background(), msg)
		if err != nil {
			t.Errorf("Broadcast() error = %v, want nil", err)
		}
	})
}