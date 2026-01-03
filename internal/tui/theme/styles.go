package theme

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the pre-built lipgloss styles for the theme
type Styles struct {
	// Base styles
	Base lipgloss.Style

	// Text styles
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Text     lipgloss.Style
	Muted    lipgloss.Style
	Subtle   lipgloss.Style

	// Status styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Interactive states
	Focused  lipgloss.Style
	Hovered  lipgloss.Style
	Disabled lipgloss.Style

	// Border styles
	BorderNormal lipgloss.Style
	BorderFocus  lipgloss.Style

	// Input styles
	TextInput TextInputStyles

	// Help styles
	Help help.Styles
}

// TextInputStyles contains styles for text input components
type TextInputStyles struct {
	Label    lipgloss.Style
	Help     lipgloss.Style
	Normal   lipgloss.Style
	Hovered  lipgloss.Style
	FocusedB lipgloss.Style // Border style when focused
}

func (t *Theme) buildStyles() *Styles {
	base := lipgloss.NewStyle().Foreground(t.FgBase)

	return &Styles{
		Base: base,

		Title: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Bold(true),

		Text: base,

		Muted: lipgloss.NewStyle().Foreground(t.FgMuted),

		Subtle: lipgloss.NewStyle().Foreground(t.FgSubtle),

		// Status
		Success: lipgloss.NewStyle().Foreground(t.Success),
		Warning: lipgloss.NewStyle().Foreground(t.Warning),
		Error:   lipgloss.NewStyle().Foreground(t.Error),
		Info:    lipgloss.NewStyle().Foreground(t.Info),

		// Interactive
		Focused:  lipgloss.NewStyle().Foreground(t.Primary).Bold(true),
		Hovered:  lipgloss.NewStyle().Foreground(t.FgBase),
		Disabled: lipgloss.NewStyle().Foreground(t.FgSubtle),

		// Borders
		BorderNormal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),

		BorderFocus: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.BorderFocus),

		// Text input
		TextInput: TextInputStyles{
			Label: lipgloss.NewStyle().Foreground(t.FgMuted),
			Help:  lipgloss.NewStyle().Foreground(t.FgSubtle).Italic(true),
			Normal: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(t.Border).
				Padding(0, 1),
			Hovered: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(t.FgMuted).
				Padding(0, 1),
			FocusedB: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(t.BorderFocus).
				Padding(0, 1),
		},

		// Help
		Help: help.Styles{
			ShortKey:       lipgloss.NewStyle().Foreground(t.Primary),
			ShortDesc:      lipgloss.NewStyle().Foreground(t.FgMuted),
			ShortSeparator: lipgloss.NewStyle().Foreground(t.FgSubtle),
			FullKey:        lipgloss.NewStyle().Foreground(t.Primary),
			FullDesc:       lipgloss.NewStyle().Foreground(t.FgMuted),
			FullSeparator:  lipgloss.NewStyle().Foreground(t.FgSubtle),
		},
	}
}
