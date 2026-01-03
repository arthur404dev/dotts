// Package palette provides the command palette component
package palette

import tea "github.com/charmbracelet/bubbletea"

// Command represents a command in the palette
type Command struct {
	ID          string
	Label       string
	Description string
	Category    string
	Keywords    []string // For fuzzy search
	Shortcut    string   // Display hint (e.g., "ctrl+s")
	Handler     func() tea.Cmd
}

// Category constants
const (
	CategoryNavigation = "Navigation"
	CategoryActions    = "Actions"
	CategorySystem     = "System"
)
