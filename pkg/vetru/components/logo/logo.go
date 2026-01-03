// Package logo provides the ASCII art logo component
package logo

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// ASCII art for "dotts"
const asciiLogo = `     _       _   _       
  __| | ___ | |_| |_ ___ 
 / _' |/ _ \| __| __/ __|
| (_| | (_) | |_| |_\__ \
 \__,_|\___/ \__|\__|___/`

// Logo renders the ASCII art logo with gradient colors
type Logo struct {
	theme   *theme.Theme
	version string
}

// New creates a new logo component
func New(t *theme.Theme, version string) *Logo {
	return &Logo{
		theme:   t,
		version: version,
	}
}

// View renders the full ASCII art logo with gradient
func (l *Logo) View() string {
	t := l.theme
	lines := strings.Split(asciiLogo, "\n")

	var result strings.Builder
	for _, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}
		// Apply gradient to each line
		gradientLine := theme.ApplyGradient(line, t.Primary, t.Secondary, t.Tertiary)
		result.WriteString(gradientLine)
		result.WriteString("\n")
	}

	// Add version below
	versionLine := lipgloss.NewStyle().
		Foreground(t.FgMuted).
		Render("v" + l.version)

	return result.String() + versionLine
}

// Compact renders a single-line version of the logo
func (l *Logo) Compact() string {
	t := l.theme
	brand := theme.ApplyBoldGradient("dotts", t.Primary, t.Secondary)
	ver := lipgloss.NewStyle().Foreground(t.FgMuted).Render("v" + l.version)
	return brand + " " + ver
}

// Width returns the width of the full logo
func (l *Logo) Width() int {
	return lipgloss.Width(asciiLogo)
}

// Height returns the height of the full logo
func (l *Logo) Height() int {
	return lipgloss.Height(asciiLogo) + 1 // +1 for version line
}
