package layout

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// Footer renders the help bar at the bottom of the TUI
type Footer struct {
	theme  *theme.Theme
	help   help.Model
	keyMap help.KeyMap
	width  int
}

// NewFooter creates a new footer component
func NewFooter(t *theme.Theme) *Footer {
	h := help.New()
	h.Styles = t.S().Help
	h.ShowAll = false

	return &Footer{
		theme: t,
		help:  h,
	}
}

// SetWidth sets the footer width
func (f *Footer) SetWidth(w int) {
	f.width = w
	f.help.Width = w - 4
}

// SetKeyMap sets the key bindings to display
func (f *Footer) SetKeyMap(km help.KeyMap) {
	f.keyMap = km
}

// ToggleFullHelp toggles between short and full help
func (f *Footer) ToggleFullHelp() {
	f.help.ShowAll = !f.help.ShowAll
}

// View renders the footer
func (f *Footer) View() string {
	if f.keyMap == nil {
		return ""
	}

	container := lipgloss.NewStyle().
		Padding(0, 1).
		Width(f.width)

	return container.Render(f.help.View(f.keyMap))
}
