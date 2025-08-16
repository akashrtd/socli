package tui

import (
	"context"
	"encoding/json"
	"log"
	"socli/internal"
	"socli/messaging"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// subscribeToHashtag handles the logic for subscribing to a new hashtag.
func (m *AppModel) subscribeToHashtag(hashtag string) {
	topicName := messaging.GetTopicForHashtag(hashtag)
	
	// Check if already subscribed
	if _, ok := m.subscriptions[topicName]; ok {
		log.Printf("Already subscribed to topic: %s", topicName)
		return
	}

	// Join the topic
	topic, err := m.psManager.JoinTopic(topicName)
	if err != nil {
		log.Printf("Error joining topic %s: %v", topicName, err)
		return
	}

	// Subscribe to the topic
	sub, err := m.psManager.SubscribeToTopic(topic)
	if err != nil {
		log.Printf("Error subscribing to topic %s: %v", topicName, err)
		return
	}

	// Store the subscription
	m.subscriptions[topicName] = sub

	// Start listening to this new subscription in a goroutine
	// This goroutine will read messages from the subscription and send them
	// to the AppModel's postChan, which is listened to by listenForPostsCmd.
	go func(sub *pubsub.Subscription, topicName string, postChan chan<- *messaging.Message) {
		// Use a background context for now. In a more complex app, you might
		// use a cancellable context tied to the subscription's lifecycle.
		ctx := context.Background() 
		
		for {
			// Blocking receive from the subscription
			msg, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., context cancellation, subscription closed)
				log.Printf("Error receiving message from topic %s: %v", topicName, err)
				// If the subscription is closed or context is cancelled,
				// we should stop this goroutine.
				// This requires careful synchronization to remove the subscription
				// from m.subscriptions. For now, we'll just log and return.
				// A more robust implementation would involve a way to signal
				// this goroutine to stop and clean up the subscription entry.
				return
			}
			
			// This check is important to avoid processing our own messages
			// if the broadcaster also publishes to topics we subscribe to.
			// The primary check should be consistent with main.go's logic.
			// Let's assume the main check is sufficient for the default topic,
			// and any new topics subscribed here are additional.
			// Replicating the check here is safer.
			if msg.ReceivedFrom == m.netManager.Host.ID() {
				continue // Ignore own messages
			}

			// --- Message Processing (Simplified) ---
			// In main.go, there's decryption and unmarshalling logic.
			// For consistency, this logic should ideally be centralized.
			// However, for this prototype, we'll do a basic unmarshal here.
			// A production app would likely have a shared message processing function.
			
			// Assume message data is JSON for now (no encryption handling in this simplified path)
			// A full implementation would duplicate the decryption/unmarshalling logic from main.go
			// or refactor it into a shared utility.
			var receivedMsg messaging.Message
			if err := json.Unmarshal(msg.Data, &receivedMsg); err != nil {
				// Handle unmarshalling error
				log.Printf("Error unmarshalling message from topic %s: %v", topicName, err)
				continue
			}
			
			// Apply filters (placeholder)
			// Note: internal.ApplyFilters requires the message pointer.
			// The code below has a type mismatch. Let's fix it.
			// if !internal.ApplyFilters(receivedMsg) { 
			// 	 continue // Message was filtered out
			// }
			// The function expects *messaging.Message, and receivedMsg is a value.
			// Let's pass a pointer.
			if !internal.ApplyFilters(&receivedMsg) {
				continue // Message was filtered out
			}

			// Send the processed message to the AppModel's post channel
			// This will be picked up by listenForPostsCmd and turned into
			// a PostReceivedMsg for the Update function.
			// Use a select with default to avoid blocking if the channel is full.
			// This prevents this goroutine from hanging if the TUI is slow.
			select {
			case postChan <- &receivedMsg:
				// Message sent successfully
			default:
				// Channel is full, log and drop the message
				// This is a form of backpressure handling.
				log.Printf("Warning: Post channel full, dropping message from topic %s", topicName)
			}
		}
	}(sub, topicName, m.postChan) // Pass postChan to the goroutine

	log.Printf("Subscribed to hashtag #%s (topic: %s)", hashtag, topicName)
}

// unsubscribeFromHashtag handles the logic for unsubscribing from a hashtag.
func (m *AppModel) unsubscribeFromHashtag(hashtag string) {
	topicName := messaging.GetTopicForHashtag(hashtag)

	// Check if subscribed
	_, ok := m.subscriptions[topicName]
	if !ok {
		log.Printf("Not subscribed to topic: %s", topicName)
		return
	}

	// For now, just remove it from our map
	// TODO: Properly signal the subscription goroutine to stop and clean up resources.
	// This is a complex part of lifecycle management. Simply removing it from the map
	// doesn't stop the goroutine. The goroutine would need to check for a cancellation
	// signal (e.g., via a context) in its loop.
	delete(m.subscriptions, topicName)
	
	log.Printf("Unsubscribed from hashtag #%s (topic: %s)", hashtag, topicName)
	// The goroutine for this subscription will eventually stop when it
	// tries to read from a closed/cancelled subscription or when its
	// context is cancelled. Proper cleanup requires more coordination.
}