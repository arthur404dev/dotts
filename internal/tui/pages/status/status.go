// Package status provides the status page
package status

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Status is the status page
type Status struct {
	theme  *theme.Theme
	width  int
	height int
}

// New creates a new status page
func New(t *theme.Theme) *Status {
	return &Status{theme: t}
}

func (s *Status) ID() messages.PageID                     { return app.PageStatus }
func (s *Status) Title() string                           { return "Status" }
func (s *Status) SetSize(w, h int)                        { s.width, s.height = w, h }
func (s *Status) Focus() tea.Cmd                          { return nil }
func (s *Status) Blur()                                   {}
func (s *Status) Init() tea.Cmd                           { return nil }
func (s *Status) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }

func (s *Status) View() string {
	t := s.theme
	title := t.S().Title.Render("Status")
	content := t.S().Muted.Render("Coming soon...")

	container := lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Padding(1, 2)

	return container.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", content))
}
