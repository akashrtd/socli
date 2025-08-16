package tui

import (
	"strings"
	"unicode"
)

// ParseCommand checks if the input string is a command and parses it.
// Commands start with '/'.
// Returns the command name and a slice of arguments.
// If the input is not a command, it returns an empty command name.
func ParseCommand(input string) (command string, args []string) {
	// Trim leading whitespace
	input = strings.TrimLeftFunc(input, unicode.IsSpace)
	
	// Check if it starts with '/'
	if !strings.HasPrefix(input, "/") {
		return "", nil
	}
	
	// Remove the '/'
	content := input[1:]
	
	// Split by whitespace to get command and args
	parts := strings.Fields(content)
	
	if len(parts) == 0 {
		return "", nil
	}
	
	command = parts[0]
	args = parts[1:]
	
	return command, args
}