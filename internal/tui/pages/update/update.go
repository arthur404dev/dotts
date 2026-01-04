package update

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Update struct {
	theme  *theme.Theme
	width  int
	height int
}

func New(t *theme.Theme) *Update {
	return &Update{theme: t}
}

func (u *Update) ID() messages.PageID                     { return app.PageUpdate }
func (u *Update) Title() string                           { return "Update" }
func (u *Update) SetSize(w, h int)                        { u.width, u.height = w, h }
func (u *Update) Focus() tea.Cmd                          { return nil }
func (u *Update) Blur()                                   {}
func (u *Update) Init() tea.Cmd                           { return nil }
func (u *Update) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return u, nil }

func (u *Update) View() string {
	t := u.theme
	title := t.S().Title.Render("Update")
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		t.S().Muted.Render("Coming soon..."),
	)

	return components.NewPageContainer(t).
		SetSize(u.width, u.height).
		SetContent(content).
		View()
}
