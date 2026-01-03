package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// Modal is a generic modal overlay component.
type Modal struct {
	theme     *theme.Theme
	title     string
	content   string
	footer    string
	width     int
	closeable bool
}

// NewModal creates a new Modal component with the given theme.
func NewModal(t *theme.Theme) *Modal {
	return &Modal{
		theme:     t,
		width:     50,
		closeable: true,
	}
}

// SetTitle sets the modal title.
func (m *Modal) SetTitle(title string) *Modal {
	m.title = title
	return m
}

// SetContent sets the modal content.
func (m *Modal) SetContent(content string) *Modal {
	m.content = content
	return m
}

// SetFooter sets the modal footer.
func (m *Modal) SetFooter(footer string) *Modal {
	m.footer = footer
	return m
}

// SetWidth sets the modal width.
func (m *Modal) SetWidth(w int) *Modal {
	if w > 0 {
		m.width = w
	}
	return m
}

// SetCloseable sets whether the modal shows a close hint.
func (m *Modal) SetCloseable(closeable bool) *Modal {
	m.closeable = closeable
	return m
}

// View renders the modal.
func (m *Modal) View() string {
	t := m.theme

	var parts []string

	if m.title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true)
		parts = append(parts, titleStyle.Render(m.title))
		parts = append(parts, "")
	}

	if m.content != "" {
		contentStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		parts = append(parts, contentStyle.Render(m.content))
	}

	if m.footer != "" {
		parts = append(parts, "")
		footerStyle := lipgloss.NewStyle().Foreground(t.FgMuted)
		parts = append(parts, footerStyle.Render(m.footer))
	} else if m.closeable {
		parts = append(parts, "")
		footerStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
		parts = append(parts, footerStyle.Render("Press Esc to close"))
	}

	inner := lipgloss.JoinVertical(lipgloss.Left, parts...)

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus).
		Padding(1, 2).
		Width(m.width)

	return containerStyle.Render(inner)
}
