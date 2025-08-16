package types

// StatusMsg represents a status message to be displayed to the user.
// It can be an informational message, a warning, or an error.
type StatusMsg struct {
	Type    StatusType
	Message string
}

// StatusType defines the type of status message.
type StatusType int

const (
	Info StatusType = iota
	Warning
	Error
	Success
)

// StatusMsg implementations for common actions
var (
	// Posting
	PostingMsg     = StatusMsg{Info, "Publishing your post..."}
	PostSentMsg    = StatusMsg{Success, "Post sent successfully!"}
	PostFailedMsg  = StatusMsg{Error, "Failed to send post. Please try again."}

	// Subscription
	SubscribingMsg    = StatusMsg{Info, "Subscribing to topic..."}
	SubscribedMsg     = StatusMsg{Success, "Successfully subscribed!"}
	SubscribeFailedMsg = StatusMsg{Error, "Failed to subscribe. Please try again."}
	
	UnsubscribingMsg    = StatusMsg{Info, "Unsubscribing from topic..."}
	UnsubscribedMsg     = StatusMsg{Success, "Successfully unsubscribed!"}
	UnsubscribeFailedMsg = StatusMsg{Error, "Failed to unsubscribe. Please try again."}

	// Peer Connection
	PeerConnectedMsg = StatusMsg{Info, "New peer connected!"}
	PeerDisconnectedMsg = StatusMsg{Info, "Peer disconnected!"}

	// General
	UnknownCmdMsg = StatusMsg{Warning, "Unknown command. Type /help for a list of commands."}
)

// BroadcastResultMsg is a message sent when a post broadcast operation completes.
// It carries the result status of the broadcast.
type BroadcastResultMsg struct {
	Status StatusMsg
}