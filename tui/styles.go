package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Main layout styles
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")). // Purple
			Background(lipgloss.Color("235")). // Dark grey
			Bold(true).
			Padding(0, 1)

	// View styles (Feed, Compose, etc.)
	viewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")) // Grey

	// Sidebar styles (for peer/topic list)
	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")). // Grey
			Padding(0, 1)

	// Status bar styles
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Grey
			Background(lipgloss.Color("235")). // Dark grey
			Padding(0, 1)

	// List item styles
	listItemStyle = lipgloss.NewStyle().PaddingLeft(1)

	// Active list item styles
	activeListItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("205")) // Pink

	// Help text styles
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Grey
			Align(lipgloss.Center)
)