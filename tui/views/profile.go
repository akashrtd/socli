package views

import (
	"fmt"
	"socli/config"
	"socli/p2p"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProfileView displays user profile information and application settings.
type ProfileView struct {
	netManager *p2p.NetworkManager
	cfg        *config.Config
	// For editable fields, we might store temporary values or a copy of the config
	// that gets applied on save/exit.
}

// NewProfileView creates a new profile view.
func NewProfileView(netManager *p2p.NetworkManager, cfg *config.Config) *ProfileView {
	return &ProfileView{
		netManager: netManager,
		cfg:        cfg,
	}
}

// View renders the profile information.
func (v *ProfileView) View(width, height int) string {
	var b strings.Builder

	// --- Header ---
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")). // Purple
		MarginBottom(1)
	b.WriteString(headerStyle.Render("User Profile"))
	b.WriteString("\n")

	// --- Identity Section ---
	identityTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("33")). // Blue
		MarginTop(1).
		MarginBottom(1)
	b.WriteString(identityTitleStyle.Render("Identity"))
	b.WriteString("\n")

	if v.netManager != nil && v.netManager.Host != nil {
		peerID := v.netManager.Host.ID().String()
		b.WriteString(fmt.Sprintf("Peer ID: %s\n", peerID))
	} else {
		b.WriteString("Peer ID: Not available\n")
	}

	b.WriteString(fmt.Sprintf("Key File Path: %s\n", v.cfg.Privacy.KeyPath))
	b.WriteString("\n")

	// --- Network Section ---
	b.WriteString(identityTitleStyle.Render("Network"))
	b.WriteString("\n")

	if v.netManager != nil && v.netManager.Host != nil {
		b.WriteString("Listening Addresses:\n")
		for _, addr := range v.netManager.Host.Addrs() {
			b.WriteString(fmt.Sprintf("  - %s\n", addr.String()))
		}
	} else {
		b.WriteString("Listening Addresses: Not available\n")
	}

	// Connected Peers
	// This overlaps with the sidebar, but could show more detail here.
	// For now, we'll just note it's available in the main feed sidebar.
	b.WriteString("\nConnected peers are listed in the main feed sidebar.\n")
	b.WriteString("\n")

	// --- Configuration Section ---
	b.WriteString(identityTitleStyle.Render("Configuration"))
	b.WriteString("\n")

	// UI Settings
	b.WriteString("UI Settings:\n")
	b.WriteString(fmt.Sprintf("  Theme: %s\n", v.cfg.UI.Theme))
	b.WriteString(fmt.Sprintf("  Refresh Rate (ms): %d\n", v.cfg.UI.RefreshRate))
	b.WriteString(fmt.Sprintf("  Max Post Length: %d\n", v.cfg.UI.MaxPostLength))
	b.WriteString("\n")

	// Privacy Settings
	b.WriteString("Privacy Settings:\n")
	b.WriteString(fmt.Sprintf("  Encrypt Messages: %t\n", v.cfg.Privacy.EncryptMessages))
	b.WriteString(fmt.Sprintf("  Auto Clear on Exit: %t\n", v.cfg.Privacy.AutoClear))
	b.WriteString("\n")

	// Network Settings
	b.WriteString("Network Settings:\n")
	b.WriteString(fmt.Sprintf("  Listen Port: %d\n", v.cfg.Network.ListenPort))
	b.WriteString(fmt.Sprintf("  Enable mDNS: %t\n", v.cfg.Network.EnableMDNS))
	b.WriteString(fmt.Sprintf("  Enable DHT: %t\n", v.cfg.Network.EnableDHT))
	b.WriteString("\n")

	// --- Footer ---
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // Grey
		Align(lipgloss.Center).
		MarginTop(1)
	b.WriteString(footerStyle.Render("Press 'q' or 'Esc' to return to the main feed."))

	return b.String()
}

// TODO: Add methods to handle user input for editing settings if needed.
// This might involve adding fields to ProfileView to track temporary changes
// and then applying them on a 'Save' action.