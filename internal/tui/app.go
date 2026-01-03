// Package tui provides the terminal user interface for dotts
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the TUI application
func Run() error {
	model := NewModel()

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	_, err := p.Run()
	return err
}
