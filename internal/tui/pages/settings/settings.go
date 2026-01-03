// Package settings provides the settings page
package settings

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/internal/tui/messages"
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Settings is the settings page
type Settings struct {
	theme  *theme.Theme
	width  int
	height int
}

// New creates a new settings page
func New(t *theme.Theme) *Settings {
	return &Settings{theme: t}
}

func (s *Settings) ID() messages.PageID                     { return app.PageSettings }
func (s *Settings) Title() string                           { return "Settings" }
func (s *Settings) SetSize(w, h int)                        { s.width, s.height = w, h }
func (s *Settings) Focus() tea.Cmd                          { return nil }
func (s *Settings) Blur()                                   {}
func (s *Settings) Init() tea.Cmd                           { return nil }
func (s *Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }

func (s *Settings) View() string {
	t := s.theme
	title := t.S().Title.Render("Settings")
	content := t.S().Muted.Render("Coming soon...")

	container := lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Padding(1, 2)

	return container.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", content))
}
