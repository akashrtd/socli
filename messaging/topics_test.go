package messaging

import (
	"testing"
)

// TestGetTopicForHashtag tests the GetTopicForHashtag function.
func TestGetTopicForHashtag(t *testing.T) {
	tests := []struct {
		name     string
		hashtag  string
		expected string
	}{
		{
			name:     "Simple hashtag",
			hashtag:  "golang",
			expected: "socli/hashtag/golang",
		},
		{
			name:     "Hashtag with numbers",
			hashtag:  "go123",
			expected: "socli/hashtag/go123",
		},
		{
			name:     "Hashtag with underscore",
			hashtag:  "go_lang",
			expected: "socli/hashtag/go_lang",
		},
		{
			name:     "Hashtag with hyphen",
			hashtag:  "go-lang",
			expected: "socli/hashtag/go-lang",
		},
		{
			name:     "Empty hashtag",
			hashtag:  "",
			expected: "socli/hashtag/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTopicForHashtag(tt.hashtag)
			if got != tt.expected {
				t.Errorf("GetTopicForHashtag(%q) = %q, want %q", tt.hashtag, got, tt.expected)
			}
		})
	}
}