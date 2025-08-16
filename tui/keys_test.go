package tui

import (
	"reflect"
	"testing"
)

// TestParseCommand tests the ParseCommand function.
func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "Simple command",
			input:    "/subscribe tech",
			wantCmd:  "subscribe",
			wantArgs: []string{"tech"},
		},
		{
			name:     "Command with multiple arguments",
			input:    "/subscribe tech news",
			wantCmd:  "subscribe",
			wantArgs: []string{"tech", "news"},
		},
		{
			name:     "Command with no arguments",
			input:    "/help",
			wantCmd:  "help",
			wantArgs: []string{},
		},
		{
			name:     "Not a command",
			input:    "This is a normal message",
			wantCmd:  "",
			wantArgs: []string(nil),
		},
		{
			name:     "Empty string",
			input:    "",
			wantCmd:  "",
			wantArgs: []string(nil),
		},
		{
			name:     "Command with leading spaces",
			input:    "  /subscribe tech",
			wantCmd:  "subscribe",
			wantArgs: []string{"tech"},
		},
		{
			name:     "Command with trailing spaces",
			input:    "/subscribe tech  ",
			wantCmd:  "subscribe",
			wantArgs: []string{"tech"},
		},
		{
			name:     "Command with leading and trailing spaces",
			input:    "  /subscribe tech  ",
			wantCmd:  "subscribe",
			wantArgs: []string{"tech"},
		},
		{
			name:     "Command with only spaces",
			input:    "   ",
			wantCmd:  "",
			wantArgs: []string(nil),
		},
		{
			name:     "Just a slash",
			input:    "/",
			wantCmd:  "",
			wantArgs: []string(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs := ParseCommand(tt.input)
			if gotCmd != tt.wantCmd {
				t.Errorf("ParseCommand(%q) got command %q, want %q", tt.input, gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("ParseCommand(%q) got args %v, want %v", tt.input, gotArgs, tt.wantArgs)
			}
		})
	}
}