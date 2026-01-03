// Package doctor provides the doctor page
package doctor

import (
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/internal/tui/messages"
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Doctor is the doctor page
type Doctor struct {
	theme  *theme.Theme
	width  int
	height int
}

// New creates a new doctor page
func New(t *theme.Theme) *Doctor {
	return &Doctor{theme: t}
}

func (d *Doctor) ID() messages.PageID                     { return app.PageDoctor }
func (d *Doctor) Title() string                           { return "Doctor" }
func (d *Doctor) SetSize(w, h int)                        { d.width, d.height = w, h }
func (d *Doctor) Focus() tea.Cmd                          { return nil }
func (d *Doctor) Blur()                                   {}
func (d *Doctor) Init() tea.Cmd                           { return nil }
func (d *Doctor) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return d, nil }

func (d *Doctor) View() string {
	t := d.theme
	title := t.S().Title.Render("Doctor")
	content := t.S().Muted.Render("Coming soon...")

	container := lipgloss.NewStyle().
		Width(d.width).
		Height(d.height).
		Padding(1, 2)

	return container.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", content))
}
