package status

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status struct {
	theme  *theme.Theme
	width  int
	height int
}

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
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		t.S().Muted.Render("Coming soon..."),
	)

	return components.NewPageContainer(t).
		SetSize(s.width, s.height).
		SetContent(content).
		View()
}
