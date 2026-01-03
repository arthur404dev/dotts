package components

import (
	"strings"

	"github.com/arthur404dev/dotts/internal/tui/theme"
)

// SpacerDirection defines the direction of spacing.
type SpacerDirection int

const (
	// SpacerVertical creates vertical space (empty lines).
	SpacerVertical SpacerDirection = iota
	// SpacerHorizontal creates horizontal space (spaces).
	SpacerHorizontal
)

// Spacer is a component that renders empty space, useful for layout
// and visual separation between other components.
type Spacer struct {
	theme     *theme.Theme
	direction SpacerDirection
	size      int
}

// NewSpacer creates a vertical spacer with the given size in lines.
// A size of 1 represents a single empty line worth of space.
func NewSpacer(t *theme.Theme, size int) *Spacer {
	return &Spacer{
		theme:     t,
		direction: SpacerVertical,
		size:      size,
	}
}

// NewHSpacer creates a horizontal spacer with the given size in spaces.
// This is useful for creating inline spacing between components.
func NewHSpacer(t *theme.Theme, size int) *Spacer {
	return &Spacer{
		theme:     t,
		direction: SpacerHorizontal,
		size:      size,
	}
}

// SetSize sets the spacer size.
// For vertical spacers, size is in lines.
// For horizontal spacers, size is in space characters.
// Negative values are ignored. Returns the Spacer for method chaining.
func (s *Spacer) SetSize(size int) *Spacer {
	if size >= 0 {
		s.size = size
	}
	return s
}

// View renders the spacer as a string.
// For horizontal spacers, returns the specified number of space characters.
// For vertical spacers, returns newlines to create the specified number of
// empty lines. A size of 0 or less returns an empty string.
func (s *Spacer) View() string {
	if s.size <= 0 {
		return ""
	}

	if s.direction == SpacerHorizontal {
		return strings.Repeat(" ", s.size)
	}

	if s.size == 1 {
		return ""
	}
	return strings.Repeat("\n", s.size-1)
}
