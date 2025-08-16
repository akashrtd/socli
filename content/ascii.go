package content

import (
	"regexp"
	"strings"
)

// ASCII represents a detected ASCII art block within a post.
type ASCII struct {
	Content   string  // The raw text content of the ASCII block
	StartLine int     // The starting line number (0-indexed) of the block in the original content
	EndLine   int     // The ending line number (0-indexed, inclusive) of the block in the original content
	MinWidth  int     // Minimum width of lines in the block
	MaxWidth  int     // Maximum width of lines in the block
	AvgWidth  float64 // Average width of lines in the block
}

// ExtractASCII scans the provided Markdown content and finds code blocks
// with the language identifier "ascii".
//
// It returns a slice of pointers to ASCII structs, each representing a detected block.
// Blocks are returned in the order they appear in the content.
func ExtractASCII(markdown string) []*ASCII {
	lines := strings.Split(markdown, "\n")
	var asciiBlocks []*ASCII

	// Regular expression to match opening code fence with optional language identifier
	// Matches ``` followed by optional whitespace and then the word "ascii" (case-insensitive)
	// followed by optional whitespace and end of line.
	// We construct the pattern without using raw string literals to avoid escaping issues.
	// Pattern breakdown:
	// ^ - Start of line
	// [\s]* - Zero or more whitespace characters
	// ``` - Literal backticks
	// [\s]* - Zero or more whitespace characters
	// (?i:ascii) - Case-insensitive group matching "ascii"
	// [\s]* - Zero or more whitespace characters
	// $ - End of line
	pattern := "^" + `[` + `\s` + `]*` + regexp.QuoteMeta("```") + `[` + `\s` + `]*(?i:ascii)[` + `\s` + `]*` + "$"
	codeFenceRegex := regexp.MustCompile(pattern)

	inAsciiBlock := false
	blockStartLine := -1
	var currentBlockLines []string

	for i, line := range lines {
		isCodeFenceStart := codeFenceRegex.MatchString(line)
		// Check for closing code fence (any line starting with ```)
		// but not the opening ```ascii fence itself if it's on the same line (unlikely but handled)
		isCodeFenceEnd := strings.HasPrefix(strings.TrimSpace(line), "```") && !codeFenceRegex.MatchString(line)

		if isCodeFenceStart && !inAsciiBlock {
			// Start of an ASCII code block
			inAsciiBlock = true
			blockStartLine = i
			currentBlockLines = []string{}
		} else if isCodeFenceEnd && inAsciiBlock {
			// End of the current ASCII code block
			inAsciiBlock = false
			
			// Process the collected lines
			if len(currentBlockLines) > 0 {
				// Join the lines to get the content
				content := strings.Join(currentBlockLines, "\n")
				
				// Calculate min/max/avg width
				minWidth, maxWidth, avgWidth := calculateDimensions(currentBlockLines)
				
				// Create the ASCII struct
				asciiBlock := &ASCII{
					Content:   content,
					StartLine: blockStartLine,
					EndLine:   i,
					MinWidth:  minWidth,
					MaxWidth:  maxWidth,
					AvgWidth:  avgWidth,
				}
				asciiBlocks = append(asciiBlocks, asciiBlock)
			}
			
			// Reset for the next potential block
			blockStartLine = -1
			currentBlockLines = []string{}
		} else if inAsciiBlock {
			// Inside an ASCII code block, collect the line
			// Do not include the closing ```
			currentBlockLines = append(currentBlockLines, line)
		}
	}

	// Handle the case where the content ends with an unclosed ASCII block
	// This is generally an error in Markdown, but we can handle it gracefully
	if inAsciiBlock && len(currentBlockLines) > 0 {
		content := strings.Join(currentBlockLines, "\n")
		minWidth, maxWidth, avgWidth := calculateDimensions(currentBlockLines)
		asciiBlock := &ASCII{
			Content:   content,
			StartLine: blockStartLine,
			EndLine:   len(lines) - 1, // Last line index
			MinWidth:  minWidth,
			MaxWidth:  maxWidth,
			AvgWidth:  avgWidth,
		}
		asciiBlocks = append(asciiBlocks, asciiBlock)
	}

	return asciiBlocks
}

// calculateDimensions calculates min, max, and average width of lines.
func calculateDimensions(lines []string) (min, max int, avg float64) {
	if len(lines) == 0 {
		return 0, 0, 0.0
	}
	
	min = len(lines[0])
	max = len(lines[0])
	total := 0
	
	for _, line := range lines {
		l := len(line)
		total += l
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
	}
	
	avg = float64(total) / float64(len(lines))
	return min, max, avg
}