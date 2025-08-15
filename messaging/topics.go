package messaging

import "fmt"

const (
	// HashtagTopicPrefix is the prefix for all hashtag-based topics.
	HashtagTopicPrefix = "socli/hashtag/"
)

// GetTopicForHashtag returns the full topic string for a given hashtag.
func GetTopicForHashtag(hashtag string) string {
	return fmt.Sprintf("%s%s", HashtagTopicPrefix, hashtag)
}
