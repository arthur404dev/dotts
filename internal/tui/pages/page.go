// Package pages defines the page interface and page implementations
package pages

import (
	"github.com/arthur404dev/dotts/internal/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

// Page defines the interface for all TUI pages
type Page interface {
	tea.Model

	// Identity
	ID() messages.PageID
	Title() string

	// Lifecycle
	SetSize(width, height int)

	// Focus management
	Focus() tea.Cmd
	Blur()
}
