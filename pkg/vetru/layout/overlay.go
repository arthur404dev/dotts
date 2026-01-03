package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Overlay provides utilities for rendering modal overlays
type Overlay struct {
	width  int
	height int
}

// NewOverlay creates a new overlay helper
func NewOverlay(width, height int) *Overlay {
	return &Overlay{
		width:  width,
		height: height,
	}
}

// SetSize updates the overlay dimensions
func (o *Overlay) SetSize(width, height int) {
	o.width = width
	o.height = height
}

// Center renders content centered on the screen with a semi-transparent background effect
func (o *Overlay) Center(content string) string {
	contentHeight := lipgloss.Height(content)
	contentWidth := lipgloss.Width(content)

	// Calculate padding
	topPadding := (o.height - contentHeight) / 2
	leftPadding := (o.width - contentWidth) / 2

	if topPadding < 0 {
		topPadding = 0
	}
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Build the overlay
	var result strings.Builder

	// Top padding (empty lines)
	for range topPadding {
		result.WriteString(strings.Repeat(" ", o.width) + "\n")
	}

	// Content with left padding
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		padding := strings.Repeat(" ", leftPadding)
		result.WriteString(padding + line + "\n")
	}

	return result.String()
}

// PlaceOver places the overlay content on top of a base view
// This creates a modal effect where the overlay appears over existing content
func (o *Overlay) PlaceOver(base, overlay string) string {
	return lipgloss.Place(
		o.width,
		o.height,
		lipgloss.Center,
		lipgloss.Center,
		overlay,
		lipgloss.WithWhitespaceChars(" "),
	)
}
