// Package components provides reusable TUI building blocks for the dotts application.
package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// Empty is an empty state placeholder component.
// It displays a centered message with optional icon, title, description, and action hint.
// Use this component when there is no content to display, such as empty lists or initial states.
type Empty struct {
	theme   *theme.Theme
	icon    string
	title   string
	message string
	action  string
}

// NewEmpty creates a new Empty component with the given theme.
// Default settings: mailbox icon (📭), no title, message, or action.
func NewEmpty(t *theme.Theme) *Empty {
	return &Empty{
		theme: t,
		icon:  "📭",
	}
}

// SetIcon sets the icon or emoji displayed at the top of the empty state.
func (e *Empty) SetIcon(icon string) *Empty {
	e.icon = icon
	return e
}

// SetTitle sets the main title text displayed below the icon.
func (e *Empty) SetTitle(title string) *Empty {
	e.title = title
	return e
}

// SetMessage sets the description message displayed below the title.
func (e *Empty) SetMessage(message string) *Empty {
	e.message = message
	return e
}

// SetAction sets the action hint displayed at the bottom (e.g., "Press Enter to create").
func (e *Empty) SetAction(action string) *Empty {
	e.action = action
	return e
}

// View renders the empty state as a vertically centered stack of elements.
func (e *Empty) View() string {
	t := e.theme

	var parts []string

	if e.icon != "" {
		iconStyle := lipgloss.NewStyle().
			Foreground(t.FgSubtle).
			MarginBottom(1)
		parts = append(parts, iconStyle.Render(e.icon))
	}

	if e.title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(t.FgMuted).
			Bold(true)
		parts = append(parts, titleStyle.Render(e.title))
	}

	if e.message != "" {
		msgStyle := lipgloss.NewStyle().
			Foreground(t.FgSubtle)
		parts = append(parts, msgStyle.Render(e.message))
	}

	if e.action != "" {
		actionStyle := lipgloss.NewStyle().
			Foreground(t.Primary).
			MarginTop(1)
		parts = append(parts, actionStyle.Render(e.action))
	}

	content := lipgloss.JoinVertical(lipgloss.Center, parts...)

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(content)
}
