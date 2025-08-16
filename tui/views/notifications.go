package views

import (
	"fmt"
	"socli/tui/types" // Import tui/types for StatusMsg and StatusType
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// NotificationsView displays a log of real-time application events and status messages.
type NotificationsView struct {
	// A slice to store the log of messages.
	// In a full implementation, this might be a ring buffer or a more sophisticated structure.
	// For now, a simple slice is sufficient.
	messages []types.StatusMsg
	// Offset for scrolling through the log if it gets long
	offset   int
	maxLines int // Maximum number of lines to display based on terminal height
}

// NewNotificationsView creates a new notifications view.
func NewNotificationsView() *NotificationsView {
	// Initialize with a pre-allocated slice capacity.
	// The actual maxLines will be set dynamically based on terminal size.
	return &NotificationsView{
		messages: make([]types.StatusMsg, 0, 100), // Initial capacity of 100
		offset:   0,
		maxLines: 20, // Default, will be updated
	}
}

// AddMessage adds a new status message to the log.
func (v *NotificationsView) AddMessage(msg types.StatusMsg) {
	v.messages = append(v.messages, msg)
	// TODO: Implement logic to limit the number of stored messages (e.g., keep last 100)
	// TODO: Implement auto-scrolling to the latest message when a new one is added
	//       (unless the user has manually scrolled up).
}

// SetMaxLines updates the maximum number of lines to display,
// typically based on the terminal height.
func (v *NotificationsView) SetMaxLines(lines int) {
	v.maxLines = lines
	// Adjust offset if necessary to stay within bounds
	if v.offset > len(v.messages) {
		v.offset = len(v.messages)
	}
	if v.offset < 0 {
		v.offset = 0
	}
}

// ScrollUp scrolls the view up by one line (shows older messages).
func (v *NotificationsView) ScrollUp() {
	if v.offset < len(v.messages)-1 {
		v.offset++
	}
}

// ScrollDown scrolls the view down by one line (shows newer messages).
func (v *NotificationsView) ScrollDown() {
	if v.offset > 0 {
		v.offset--
	}
}

// View renders the notifications log.
func (v *NotificationsView) View(width, height int) string {
	var b strings.Builder

	// --- Header ---
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")). // Purple
		MarginBottom(1)
	b.WriteString(headerStyle.Render("Notifications & Activity Log"))
	b.WriteString("\n")

	// --- Notifications Content ---
	contentStyle := lipgloss.NewStyle().
		MaxWidth(width).
		Height(height - 2) // Account for header and footer

	// Determine the range of messages to display based on offset and maxLines
	numMessages := len(v.messages)
	if numMessages == 0 {
		b.WriteString("(No notifications yet)")
	} else {
		// Calculate start and end indices for slicing the messages
		// We want to display the most recent messages at the bottom of the view.
		// So, if offset is 0, we show the latest messages.
		// If offset is N, we show messages N steps back in history.
		endIndex := numMessages - v.offset
		startIndex := endIndex - v.maxLines
		if startIndex < 0 {
			startIndex = 0
		}
		if endIndex > numMessages {
			endIndex = numMessages
		}

		// Slice the messages to display
		messagesToShow := v.messages[startIndex:endIndex]

		// Display messages from oldest to newest within the slice
		for _, msg := range messagesToShow {
			// Style the message based on its type
			var msgStyle lipgloss.Style
			switch msg.Type {
			case types.Info:
				msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33")) // Blue
			case types.Success:
				msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")) // Green
			case types.Warning:
				msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220")) // Yellow
			case types.Error:
				msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
			default:
				msgStyle = lipgloss.NewStyle() // Default (white)
			}

			// Add a timestamp prefix
			// Using a fixed timestamp format for consistency
			timestamp := time.Now().Format("15:04:05")
			formattedMsg := fmt.Sprintf("[%s] %s", timestamp, msg.Message)
			b.WriteString(msgStyle.Render(formattedMsg))
			b.WriteString("\n")
		}
	}

	content := b.String()
	styledContent := contentStyle.Render(content)

	// --- Footer ---
	scrollInfo := ""
	if len(v.messages) > v.maxLines {
		scrollInfo = fmt.Sprintf(" (Showing %d-%d of %d)", len(v.messages)-v.offset-v.maxLines+1, len(v.messages)-v.offset, len(v.messages))
	}
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // Grey
		Align(lipgloss.Center).
		MarginTop(1)
	b.Reset() // Clear builder for footer
	b.WriteString(footerStyle.Render(fmt.Sprintf("Press 'j'/'k' to scroll, 'q'/'Esc' to return%s", scrollInfo)))

	footer := b.String()

	// --- Assemble Final View ---
	// Join content and footer vertically
	view := lipgloss.JoinVertical(lipgloss.Left, styledContent, footer)

	return view
}