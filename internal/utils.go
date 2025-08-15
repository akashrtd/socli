package internal

import (
	"regexp"
	"socli/messaging"
)

// ExtractHashtags finds all hashtags in a message content.
// Hashtags are defined as # followed by alphanumeric characters or underscores.
func ExtractHashtags(content string) []string {
	re := regexp.MustCompile(`#(\w+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	hashtags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			hashtags = append(hashtags, match[1])
		}
	}

	// If no hashtags found, default to 'general'
	if len(hashtags) == 0 {
		hashtags = append(hashtags, "general")
	}

	return hashtags
}

// ApplyFilters applies any defined message filters.
// This is a placeholder for future filtering logic.
func ApplyFilters(msg *messaging.Message) bool {
	// Placeholder: Always allow messages for now.
	// Future implementation could filter based on content, author, etc.
	return true
}