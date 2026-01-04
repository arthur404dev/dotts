package settings

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	theme  *theme.Theme
	width  int
	height int
}

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
