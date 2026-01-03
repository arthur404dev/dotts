package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

type StackDirection int

const (
	Vertical StackDirection = iota
	Horizontal
)

type StackAlign int

const (
	AlignStart StackAlign = iota
	AlignCenter
	AlignEnd
)

type Stack struct {
	theme     *theme.Theme
	direction StackDirection
	gap       int
	align     StackAlign
	children  []string
}

func NewVStack(t *theme.Theme, gap int) *Stack {
	return &Stack{
		theme:     t,
		direction: Vertical,
		gap:       gap,
		align:     AlignStart,
		children:  []string{},
	}
}

func NewHStack(t *theme.Theme, gap int) *Stack {
	return &Stack{
		theme:     t,
		direction: Horizontal,
		gap:       gap,
		align:     AlignStart,
		children:  []string{},
	}
}

func (s *Stack) SetGap(gap int) *Stack {
	if gap >= 0 {
		s.gap = gap
	}
	return s
}

func (s *Stack) SetAlign(align StackAlign) *Stack {
	s.align = align
	return s
}

func (s *Stack) Add(children ...string) *Stack {
	s.children = append(s.children, children...)
	return s
}

func (s *Stack) Clear() *Stack {
	s.children = []string{}
	return s
}

func (s *Stack) View() string {
	if len(s.children) == 0 {
		return ""
	}

	filtered := make([]string, 0, len(s.children))
	for _, child := range s.children {
		if child != "" {
			filtered = append(filtered, child)
		}
	}

	if len(filtered) == 0 {
		return ""
	}

	pos := s.alignToPosition()

	if s.direction == Vertical {
		return s.renderVertical(filtered, pos)
	}
	return s.renderHorizontal(filtered, pos)
}

func (s *Stack) alignToPosition() lipgloss.Position {
	switch s.align {
	case AlignStart:
		if s.direction == Horizontal {
			return lipgloss.Top
		}
		return lipgloss.Left
	case AlignCenter:
		return lipgloss.Center
	case AlignEnd:
		if s.direction == Horizontal {
			return lipgloss.Bottom
		}
		return lipgloss.Right
	default:
		return lipgloss.Left
	}
}

func (s *Stack) renderVertical(children []string, pos lipgloss.Position) string {
	if s.gap <= 0 {
		return lipgloss.JoinVertical(pos, children...)
	}

	gapStr := strings.Repeat("\n", s.gap)
	withGaps := make([]string, 0, len(children)*2-1)

	for i, child := range children {
		withGaps = append(withGaps, child)
		if i < len(children)-1 {
			withGaps = append(withGaps, gapStr)
		}
	}

	return lipgloss.JoinVertical(pos, withGaps...)
}

func (s *Stack) renderHorizontal(children []string, pos lipgloss.Position) string {
	if s.gap <= 0 {
		return lipgloss.JoinHorizontal(pos, children...)
	}

	gapStr := strings.Repeat(" ", s.gap)
	withGaps := make([]string, 0, len(children)*2-1)

	for i, child := range children {
		withGaps = append(withGaps, child)
		if i < len(children)-1 {
			withGaps = append(withGaps, gapStr)
		}
	}

	return lipgloss.JoinHorizontal(pos, withGaps...)
}
