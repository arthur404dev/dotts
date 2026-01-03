// Package layout provides layout primitives for the TUI
package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// PlaceOverlay renders fg on top of bg at position (x, y) with ANSI-aware compositing.
func PlaceOverlay(x, y int, fg, bg string) string {
	if fg == "" {
		return bg
	}
	if bg == "" {
		return fg
	}

	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	x = clamp(x, 0, max(0, bgWidth-fgWidth))
	y = clamp(y, 0, max(0, bgHeight-fgHeight))

	var b strings.Builder
	b.Grow(len(bg) + len(fg))

	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}

		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		bgLineWidth := ansi.StringWidth(bgLine)

		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			leftWidth := ansi.StringWidth(left)
			b.WriteString(left)
			if leftWidth < x {
				b.WriteString(strings.Repeat(" ", x-leftWidth))
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)

		fgLineWidth := ansi.StringWidth(fgLine)
		rightStart := x + fgLineWidth

		if rightStart < bgLineWidth {
			right := ansi.TruncateLeft(bgLine, rightStart, "")
			rightWidth := ansi.StringWidth(right)
			expectedRightWidth := bgLineWidth - rightStart
			if rightWidth < expectedRightWidth {
				b.WriteString(strings.Repeat(" ", expectedRightWidth-rightWidth))
			}
			b.WriteString(right)
		}
	}

	return b.String()
}

// CenterOverlay centers fg on bg within the given dimensions.
func CenterOverlay(fg, bg string, width, height int) string {
	if fg == "" {
		return bg
	}

	fgLines, fgWidth := getLines(fg)
	fgHeight := len(fgLines)

	x := (width - fgWidth) / 2
	y := (height - fgHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return PlaceOverlay(x, y, fg, bg)
}

// DimBackground applies faint styling to dim the content for overlay backgrounds.
func DimBackground(content string) string {
	if content == "" {
		return content
	}
	return lipgloss.NewStyle().Faint(true).Render(content)
}

func getLines(s string) ([]string, int) {
	if s == "" {
		return []string{}, 0
	}

	lines := strings.Split(s, "\n")
	widest := 0

	for _, line := range lines {
		w := ansi.StringWidth(line)
		if w > widest {
			widest = w
		}
	}

	return lines, widest
}

func clamp(v, lower, upper int) int {
	if upper < lower {
		upper = lower
	}
	return min(max(v, lower), upper)
}
