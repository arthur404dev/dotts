// Package components provides reusable UI primitives for the dotts TUI.
package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// TextVariant defines the visual style of text.
type TextVariant int

const (
	TextBody TextVariant = iota
	TextTitle
	TextSubtitle
	TextCaption
	TextCode
	TextMuted
	TextSuccess
	TextWarning
	TextError
	TextInfo
)

// Text is a styled text component with various variants.
type Text struct {
	theme   *theme.Theme
	content string
	variant TextVariant
	bold    bool
	italic  bool
}

// NewText creates a new Text component.
func NewText(t *theme.Theme, content string) *Text {
	return &Text{
		theme:   t,
		content: content,
		variant: TextBody,
	}
}

// SetVariant sets the text variant.
func (t *Text) SetVariant(v TextVariant) *Text {
	t.variant = v
	return t
}

// SetBold sets bold styling.
func (t *Text) SetBold(bold bool) *Text {
	t.bold = bold
	return t
}

// SetItalic sets italic styling.
func (t *Text) SetItalic(italic bool) *Text {
	t.italic = italic
	return t
}

// Title is a convenience constructor for title text.
func Title(th *theme.Theme, content string) *Text {
	return NewText(th, content).SetVariant(TextTitle)
}

// Subtitle is a convenience constructor for subtitle text.
func Subtitle(th *theme.Theme, content string) *Text {
	return NewText(th, content).SetVariant(TextSubtitle)
}

// Muted is a convenience constructor for muted text.
func Muted(th *theme.Theme, content string) *Text {
	return NewText(th, content).SetVariant(TextMuted)
}

// View renders the text.
func (t *Text) View() string {
	th := t.theme

	var style lipgloss.Style

	switch t.variant {
	case TextTitle:
		style = lipgloss.NewStyle().Foreground(th.Primary).Bold(true)
	case TextSubtitle:
		style = lipgloss.NewStyle().Foreground(th.Secondary).Bold(true)
	case TextCaption:
		style = lipgloss.NewStyle().Foreground(th.FgSubtle).Italic(true)
	case TextCode:
		style = lipgloss.NewStyle().Foreground(th.Tertiary).Background(th.BgSubtle)
	case TextMuted:
		style = lipgloss.NewStyle().Foreground(th.FgMuted)
	case TextSuccess:
		style = lipgloss.NewStyle().Foreground(th.Success)
	case TextWarning:
		style = lipgloss.NewStyle().Foreground(th.Warning)
	case TextError:
		style = lipgloss.NewStyle().Foreground(th.Error)
	case TextInfo:
		style = lipgloss.NewStyle().Foreground(th.Info)
	default: // TextBody
		style = lipgloss.NewStyle().Foreground(th.FgBase)
	}

	if t.bold {
		style = style.Bold(true)
	}
	if t.italic {
		style = style.Italic(true)
	}

	return style.Render(t.content)
}

// String allows Text to be used directly as a string.
func (t *Text) String() string {
	return t.View()
}
