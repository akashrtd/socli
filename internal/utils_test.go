package internal

import (
	"reflect"
	"socli/messaging"
	"testing"
	"time"
)

// TestExtractHashtags tests the ExtractHashtags function.
func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "Single hashtag",
			content:  "Hello #world",
			expected: []string{"world"},
		},
		{
			name:     "Multiple hashtags",
			content:  "Learning #Go and #BubbleTea",
			expected: []string{"Go", "BubbleTea"},
		},
		{
			name:     "Hashtag with punctuation",
			content:  "Check this out #test-case!",
			expected: []string{"test"},
		},
		{
			name:     "Hashtag at start",
			content:  "#start of the message",
			expected: []string{"start"},
		},
		{
			name:     "Hashtag at end",
			content:  "The end #finish",
			expected: []string{"finish"},
		},
		{
			name:     "No hashtags",
			content:  "Just a plain message",
			expected: []string{"general"}, // Default hashtag
		},
		{
			name:     "Empty string",
			content:  "",
			expected: []string{"general"}, // Default hashtag
		},
		{
			name:     "Hashtag with numbers",
			content:  "Version #v2.0 is out",
			expected: []string{"v2"},
		},
		{
			name:     "Consecutive hashtags",
			content:  "#first#second word",
			expected: []string{"first", "second"},
		},
		{
			name:     "Hashtag with underscore",
			content:  "Check #my_test case",
			expected: []string{"my_test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractHashtags(tt.content)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractHashtags(%q) = %v, want %v", tt.content, got, tt.expected)
			}
		})
	}
}

// TestApplyFilters tests the ApplyFilters function.
func TestApplyFilters(t *testing.T) {
	// Create a dummy message for testing
	msg := &messaging.Message{
		ID:        "test-id",
		Author:    "test-author",
		Content:   "This is a test message",
		Hashtags:  []string{"test"},
		Timestamp: time.Now(),
		Type:      messaging.PostMsg,
	}

	// Since ApplyFilters is a placeholder that always returns true,
	// we just test that it doesn't panic and returns true.
	got := ApplyFilters(msg)
	want := true

	if got != want {
		t.Errorf("ApplyFilters() = %v, want %v", got, want)
	}
}