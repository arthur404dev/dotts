// Package components provides reusable TUI building blocks for the dotts application.
package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// BorderStyle defines the type of border for a Box.
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderNormal
	BorderRounded
	BorderDouble
	BorderHidden
	BorderThick
)

// Box is a container component with configurable padding and borders.
// It provides a flexible way to wrap content with consistent styling.
type Box struct {
	theme       *theme.Theme
	content     string
	border      BorderStyle
	borderColor lipgloss.TerminalColor
	padding     theme.Padding
	width       int
	height      int
	title       string
}

// NewBox creates a new Box with the given theme.
// Default settings: rounded border, theme border color, horizontal padding of 1.
func NewBox(t *theme.Theme) *Box {
	return &Box{
		theme:       t,
		border:      BorderRounded,
		borderColor: t.Border,
		padding:     theme.PaddingXY(1, 0),
	}
}

// SetContent sets the box content.
func (b *Box) SetContent(content string) *Box {
	b.content = content
	return b
}

// SetBorder sets the border style.
func (b *Box) SetBorder(style BorderStyle) *Box {
	b.border = style
	return b
}

// SetBorderColor sets the border color.
func (b *Box) SetBorderColor(color lipgloss.TerminalColor) *Box {
	b.borderColor = color
	return b
}

// SetPadding sets the padding.
func (b *Box) SetPadding(p theme.Padding) *Box {
	b.padding = p
	return b
}

// SetWidth sets a fixed width (0 for auto).
func (b *Box) SetWidth(w int) *Box {
	b.width = w
	return b
}

// SetHeight sets a fixed height (0 for auto).
func (b *Box) SetHeight(h int) *Box {
	b.height = h
	return b
}

// SetTitle sets an optional title for the box.
func (b *Box) SetTitle(title string) *Box {
	b.title = title
	return b
}

// View renders the box.
func (b *Box) View() string {
	style := lipgloss.NewStyle().
		Padding(b.padding.Top, b.padding.Right, b.padding.Bottom, b.padding.Left)

	switch b.border {
	case BorderNone:
	case BorderNormal:
		style = style.Border(lipgloss.NormalBorder()).BorderForeground(b.borderColor)
	case BorderRounded:
		style = style.Border(lipgloss.RoundedBorder()).BorderForeground(b.borderColor)
	case BorderDouble:
		style = style.Border(lipgloss.DoubleBorder()).BorderForeground(b.borderColor)
	case BorderHidden:
		style = style.Border(lipgloss.HiddenBorder())
	case BorderThick:
		style = style.Border(lipgloss.ThickBorder()).BorderForeground(b.borderColor)
	}

	if b.width > 0 {
		style = style.Width(b.width)
	}
	if b.height > 0 {
		style = style.Height(b.height)
	}

	return style.Render(b.content)
}
