package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// Card is an elevated container component with optional title and footer.
// It provides a visually distinct surface for grouping related content.
type Card struct {
	theme   *theme.Theme
	title   string
	content string
	footer  string
	width   int
	focused bool
}

// NewCard creates a new Card with the given theme.
func NewCard(t *theme.Theme) *Card {
	return &Card{
		theme: t,
	}
}

// SetTitle sets the card title.
func (c *Card) SetTitle(title string) *Card {
	c.title = title
	return c
}

// SetContent sets the card content.
func (c *Card) SetContent(content string) *Card {
	c.content = content
	return c
}

// SetFooter sets the card footer.
func (c *Card) SetFooter(footer string) *Card {
	c.footer = footer
	return c
}

// SetWidth sets the card width (0 for auto).
func (c *Card) SetWidth(w int) *Card {
	c.width = w
	return c
}

// SetFocused sets whether the card appears focused (highlighted border).
func (c *Card) SetFocused(focused bool) *Card {
	c.focused = focused
	return c
}

// View renders the card.
func (c *Card) View() string {
	t := c.theme

	borderColor := t.Border
	if c.focused {
		borderColor = t.BorderFocus
	}

	var parts []string

	if c.title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1)
		parts = append(parts, titleStyle.Render(c.title))
	}

	if c.content != "" {
		contentStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		parts = append(parts, contentStyle.Render(c.content))
	}

	if c.footer != "" {
		footerStyle := lipgloss.NewStyle().
			Foreground(t.FgMuted).
			MarginTop(1)
		parts = append(parts, footerStyle.Render(c.footer))
	}

	inner := lipgloss.JoinVertical(lipgloss.Left, parts...)

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)

	if c.width > 0 {
		containerStyle = containerStyle.Width(c.width)
	}

	return containerStyle.Render(inner)
}
