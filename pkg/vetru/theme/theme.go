// Package theme provides the visual theming system for the TUI.
// It supports gradient colors, adaptive terminal colors, and cached styles.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the application color scheme
type Theme struct {
	Name   string
	IsDark bool

	// Primary colors (for gradients)
	Primary   lipgloss.TerminalColor // Main accent (pink/magenta)
	Secondary lipgloss.TerminalColor // Secondary accent (purple/blue)
	Tertiary  lipgloss.TerminalColor // Tertiary accent (cyan)

	// Background colors
	BgBase    lipgloss.TerminalColor // Main background
	BgSubtle  lipgloss.TerminalColor // Elevated surfaces
	BgOverlay lipgloss.TerminalColor // Dialogs/overlays

	// Foreground colors
	FgBase   lipgloss.TerminalColor // Primary text
	FgMuted  lipgloss.TerminalColor // Secondary text
	FgSubtle lipgloss.TerminalColor // Hints, placeholders

	// Border colors
	Border      lipgloss.TerminalColor
	BorderFocus lipgloss.TerminalColor

	// Status colors
	Success lipgloss.TerminalColor
	Warning lipgloss.TerminalColor
	Error   lipgloss.TerminalColor
	Info    lipgloss.TerminalColor

	// Cached styles
	styles *Styles
}

// S returns the cached styles for this theme, building them if needed
func (t *Theme) S() *Styles {
	if t.styles == nil {
		t.styles = t.buildStyles()
	}
	return t.styles
}

// DefaultTheme returns the dotts brand theme (Crush-inspired)
func DefaultTheme() *Theme {
	return &Theme{
		Name:   "dotts",
		IsDark: true,

		// Gradient colors (pink -> purple -> blue)
		Primary:   lipgloss.Color("#f5a9b8"), // Soft pink
		Secondary: lipgloss.Color("#b4a7d6"), // Lavender
		Tertiary:  lipgloss.Color("#89b4fa"), // Sky blue

		// Backgrounds
		BgBase:    lipgloss.Color("#1a1b26"), // Deep navy
		BgSubtle:  lipgloss.Color("#24283b"), // Elevated
		BgOverlay: lipgloss.Color("#1f2335"), // Dialogs

		// Foregrounds
		FgBase:   lipgloss.Color("#c0caf5"), // Light text
		FgMuted:  lipgloss.Color("#9aa5ce"), // Muted text
		FgSubtle: lipgloss.Color("#565f89"), // Very muted

		// Borders
		Border:      lipgloss.Color("#3b4261"),
		BorderFocus: lipgloss.Color("#f5a9b8"), // Primary

		// Status
		Success: lipgloss.Color("#9ece6a"),
		Warning: lipgloss.Color("#e0af68"),
		Error:   lipgloss.Color("#f7768e"),
		Info:    lipgloss.Color("#7aa2f7"),
	}
}

// TerminalAdaptive returns a theme using terminal's native ANSI colors
func TerminalAdaptive() *Theme {
	return &Theme{
		Name:   "terminal",
		IsDark: true,

		// Use ANSI color numbers - these adapt to terminal theme
		Primary:   lipgloss.ANSIColor(13), // Bright magenta
		Secondary: lipgloss.ANSIColor(12), // Bright blue
		Tertiary:  lipgloss.ANSIColor(14), // Bright cyan

		BgBase:    lipgloss.ANSIColor(0),
		BgSubtle:  lipgloss.ANSIColor(8),
		BgOverlay: lipgloss.ANSIColor(0),

		FgBase:   lipgloss.ANSIColor(15),
		FgMuted:  lipgloss.ANSIColor(7),
		FgSubtle: lipgloss.ANSIColor(8),

		Border:      lipgloss.ANSIColor(8),
		BorderFocus: lipgloss.ANSIColor(13),

		Success: lipgloss.ANSIColor(10),
		Warning: lipgloss.ANSIColor(11),
		Error:   lipgloss.ANSIColor(9),
		Info:    lipgloss.ANSIColor(12),
	}
}

// current holds the active theme
var current *Theme

// Current returns the current active theme
func Current() *Theme {
	if current == nil {
		current = DefaultTheme()
	}
	return current
}

// SetCurrent sets the active theme
func SetCurrent(t *Theme) {
	current = t
}
