// Package components provides reusable TUI building blocks for dotts.
package components

import (
	"strings"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// SeparatorOrientation defines the direction of the separator.
type SeparatorOrientation int

const (
	// SeparatorHorizontal renders a horizontal line (─────).
	SeparatorHorizontal SeparatorOrientation = iota
	// SeparatorVertical renders a vertical line (│).
	SeparatorVertical
)

// Separator is a visual divider component that can be rendered horizontally
// or vertically. It supports optional centered labels for horizontal separators.
type Separator struct {
	theme       *theme.Theme
	orientation SeparatorOrientation
	label       string
	length      int
	char        string
}

// NewSeparator creates a new horizontal Separator with default settings.
// The default character is "─" and the default length is 40.
func NewSeparator(t *theme.Theme) *Separator {
	return &Separator{
		theme:       t,
		orientation: SeparatorHorizontal,
		char:        "─",
		length:      40,
	}
}

// NewVerticalSeparator creates a vertical Separator with default settings.
// The default character is "│" and the default length is 1.
func NewVerticalSeparator(t *theme.Theme) *Separator {
	return &Separator{
		theme:       t,
		orientation: SeparatorVertical,
		char:        "│",
		length:      1,
	}
}

// SetLabel sets an optional centered label for horizontal separators.
// The label will be displayed in the middle of the separator line.
// Returns the Separator for method chaining.
func (s *Separator) SetLabel(label string) *Separator {
	s.label = label
	return s
}

// SetLength sets the length of the separator.
// For horizontal separators, this is the width in characters.
// For vertical separators, this is the height in lines.
// Values <= 0 are ignored. Returns the Separator for method chaining.
func (s *Separator) SetLength(length int) *Separator {
	if length > 0 {
		s.length = length
	}
	return s
}

// SetChar sets the character used for the separator line.
// Empty strings are ignored. Returns the Separator for method chaining.
func (s *Separator) SetChar(char string) *Separator {
	if char != "" {
		s.char = char
	}
	return s
}

// View renders the separator as a string.
// For vertical separators, it returns multiple lines joined by newlines.
// For horizontal separators with labels, the label is centered within the line.
func (s *Separator) View() string {
	t := s.theme
	style := lipgloss.NewStyle().Foreground(t.FgSubtle)

	if s.orientation == SeparatorVertical {
		var lines []string
		for i := 0; i < s.length; i++ {
			lines = append(lines, s.char)
		}
		return style.Render(strings.Join(lines, "\n"))
	}

	if s.label == "" {
		return style.Render(strings.Repeat(s.char, s.length))
	}

	labelStyle := lipgloss.NewStyle().Foreground(t.FgMuted).Padding(0, 1)
	labelView := labelStyle.Render(s.label)
	labelWidth := lipgloss.Width(labelView)

	sideLength := (s.length - labelWidth) / 2
	if sideLength < 0 {
		sideLength = 0
	}

	left := style.Render(strings.Repeat(s.char, sideLength))
	right := style.Render(strings.Repeat(s.char, sideLength))

	return left + labelView + right
}
