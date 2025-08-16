package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHelpViewInternal creates and renders the help screen.
func (m *AppModel) renderHelpViewInternal() string {
	// Define styles for different parts of the help screen
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")). // Purple
		MarginBottom(1)

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("33")). // Blue
		MarginTop(1).
		MarginBottom(1)

	itemStyle := lipgloss.NewStyle().
		PaddingLeft(2)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")) // Pink

	// valueStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("240")) // Grey

	exampleStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("220")) // Yellow

	// Build the help content
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("SOCLI - Help"))
	b.WriteString("\n")

	// Introduction
	b.WriteString("Welcome to SOCLI, a decentralized, ephemeral, terminal-based social platform.\n")
	b.WriteString("\n")

	// Global Keybindings
	b.WriteString(sectionTitleStyle.Render("Global Keybindings"))
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Quit the application", keyStyle.Render("q, Ctrl+C"))) + "\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Switch to compose view", keyStyle.Render("c"))) + "\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Toggle this help screen", keyStyle.Render("?"))) + "\n")
	b.WriteString("\n")

	// Feed View Keybindings
	b.WriteString(sectionTitleStyle.Render("Feed View Keybindings"))
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Scroll down to older posts", keyStyle.Render("j"))) + "\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Scroll up to newer posts", keyStyle.Render("k"))) + "\n")
	b.WriteString("\n")

	// Compose View Keybindings
	b.WriteString(sectionTitleStyle.Render("Compose View Keybindings"))
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Send the typed message or execute the command", keyStyle.Render("Enter"))) + "\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Discard the current message/command and return to the feed view", keyStyle.Render("Esc"))) + "\n")
	b.WriteString("\n")

	// Commands
	b.WriteString(sectionTitleStyle.Render("Commands"))
	b.WriteString("While in the compose view, you can enter special commands prefixed with '/'.\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Join a new topic to start receiving posts tagged with #hashtag.", keyStyle.Render("/subscribe <hashtag>"))) + " Example: " + exampleStyle.Render("/subscribe tech") + "\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("%s : Leave a topic to stop receiving posts for #hashtag.", keyStyle.Render("/unsubscribe <hashtag>"))) + " Example: " + exampleStyle.Render("/unsubscribe tech") + "\n")
	b.WriteString("\n")

	// Features
	b.WriteString(sectionTitleStyle.Render("Features"))
	b.WriteString(itemStyle.Render("Real-time Messaging: Instantly broadcast and receive messages via GossipSub.") + "\n")
	b.WriteString(itemStyle.Render("Markdown Support: Format your posts with Markdown, rendered beautifully in the terminal.") + "\n")
	b.WriteString(itemStyle.Render("Dynamic Hashtag Subscriptions: Join and leave topics (#hashtags) on the fly.") + "\n")
	b.WriteString(itemStyle.Render("Peer Discovery: Automatically discover other SOCLI users on your local network (mDNS) and globally (DHT).") + "\n")
	b.WriteString(itemStyle.Render("End-to-End Encryption: Message payloads are encrypted using NaCl Box before being sent over the network.") + "\n")
	b.WriteString(itemStyle.Render("Message Signing: All posts are cryptographically signed for authenticity.") + "\n")
	b.WriteString(itemStyle.Render("Privacy First: All content is ephemeral, stored only in memory and vanishes on exit.") + "\n")
	b.WriteString("\n")

	// Configuration
	b.WriteString(sectionTitleStyle.Render("Configuration"))
	b.WriteString("SOCLI can be configured via a 'config.yaml' file in the application directory.\n")
	b.WriteString("See 'config/defaults.go' for default values and descriptions.\n")
	b.WriteString("\n")

	// Troubleshooting
	b.WriteString(sectionTitleStyle.Render("Troubleshooting"))
	b.WriteString(itemStyle.Render("If you encounter issues, check the application logs for error messages.") + "\n")
	b.WriteString(itemStyle.Render("Ensure your firewall allows traffic on the configured port (default random).") + "\n")
	b.WriteString(itemStyle.Render("For network connectivity issues, verify that other SOCLI instances are on the same network (mDNS) or reachable via DHT.") + "\n")
	b.WriteString("\n")

	// Advanced Usage
	b.WriteString(sectionTitleStyle.Render("Advanced Usage"))
	b.WriteString(itemStyle.Render("SOCLI is designed to be extensible. See the project repository for developer documentation.") + "\n")
	b.WriteString("\n")

	// Footer
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center).Render(
		fmt.Sprintf("Press 'q' or 'Ctrl+C' to quit, '?' or 'Esc' to return to the feed. Scroll with 'j'/'k'. My Peer ID: %s", m.netManager.Host.ID().String()),
	))

	helpText := b.String()

	// Render the final help view with styling
	helpStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")) // Purple border

	return helpStyle.Render(helpText)
}