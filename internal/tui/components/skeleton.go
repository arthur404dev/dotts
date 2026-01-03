package components

import (
	"strings"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// Skeleton is a loading placeholder component.
type Skeleton struct {
	theme  *theme.Theme
	width  int
	height int
	char   string
}

// NewSkeleton creates a new Skeleton component.
func NewSkeleton(t *theme.Theme) *Skeleton {
	return &Skeleton{
		theme:  t,
		width:  20,
		height: 1,
		char:   "░",
	}
}

// SetWidth sets the skeleton width.
func (s *Skeleton) SetWidth(w int) *Skeleton {
	if w > 0 {
		s.width = w
	}
	return s
}

// SetHeight sets the skeleton height (in lines).
func (s *Skeleton) SetHeight(h int) *Skeleton {
	if h > 0 {
		s.height = h
	}
	return s
}

// SetSize sets both width and height.
func (s *Skeleton) SetSize(w, h int) *Skeleton {
	return s.SetWidth(w).SetHeight(h)
}

// SetChar sets the skeleton character.
func (s *Skeleton) SetChar(char string) *Skeleton {
	if char != "" {
		s.char = char
	}
	return s
}

// View renders the skeleton.
func (s *Skeleton) View() string {
	t := s.theme

	style := lipgloss.NewStyle().
		Foreground(t.FgSubtle).
		Background(t.BgSubtle)

	line := strings.Repeat(s.char, s.width)

	if s.height == 1 {
		return style.Render(line)
	}

	var lines []string
	for i := 0; i < s.height; i++ {
		lines = append(lines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// TextSkeleton creates a text-sized skeleton.
func TextSkeleton(t *theme.Theme, width int) *Skeleton {
	return NewSkeleton(t).SetWidth(width).SetHeight(1)
}

// BlockSkeleton creates a block skeleton.
func BlockSkeleton(t *theme.Theme, width, height int) *Skeleton {
	return NewSkeleton(t).SetSize(width, height)
}
