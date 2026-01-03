package components

import (
	"fmt"
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// ProgressBar is a determinate progress indicator with gradient and percentage support.
type ProgressBar struct {
	theme       *theme.Theme
	progress    float64
	width       int
	showPercent bool
	showLabel   bool
	label       string
	gradient    bool
	filledChar  string
	emptyChar   string
}

// NewProgressBar creates a new ProgressBar with default settings.
func NewProgressBar(t *theme.Theme) *ProgressBar {
	return &ProgressBar{
		theme:       t,
		progress:    0,
		width:       40,
		showPercent: true,
		gradient:    true,
		filledChar:  "█",
		emptyChar:   "░",
	}
}

// SetProgress sets the progress value (0.0 to 1.0, clamped).
func (p *ProgressBar) SetProgress(progress float64) *ProgressBar {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.progress = progress
	return p
}

// SetProgressInt sets progress as an integer percentage (0 to 100).
func (p *ProgressBar) SetProgressInt(percent int) *ProgressBar {
	return p.SetProgress(float64(percent) / 100.0)
}

// SetWidth sets the progress bar width in characters.
func (p *ProgressBar) SetWidth(w int) *ProgressBar {
	if w > 0 {
		p.width = w
	}
	return p
}

// SetShowPercent enables or disables the percentage display.
func (p *ProgressBar) SetShowPercent(show bool) *ProgressBar {
	p.showPercent = show
	return p
}

// SetLabel sets a label to display before the progress bar.
func (p *ProgressBar) SetLabel(label string) *ProgressBar {
	p.label = label
	p.showLabel = label != ""
	return p
}

// SetGradient enables or disables gradient coloring on the filled portion.
func (p *ProgressBar) SetGradient(gradient bool) *ProgressBar {
	p.gradient = gradient
	return p
}

// SetChars sets the filled and empty characters for the bar.
func (p *ProgressBar) SetChars(filled, empty string) *ProgressBar {
	if filled != "" {
		p.filledChar = filled
	}
	if empty != "" {
		p.emptyChar = empty
	}
	return p
}

// Progress returns the current progress value (0.0 to 1.0).
func (p *ProgressBar) Progress() float64 {
	return p.progress
}

// ProgressInt returns the current progress as an integer percentage (0 to 100).
func (p *ProgressBar) ProgressInt() int {
	return int(p.progress * 100)
}

// IsComplete returns true if progress has reached 100%.
func (p *ProgressBar) IsComplete() bool {
	return p.progress >= 1.0
}

// View renders the progress bar as a styled string.
func (p *ProgressBar) View() string {
	t := p.theme

	filledWidth := int(p.progress * float64(p.width))
	emptyWidth := p.width - filledWidth

	var bar string
	if p.gradient && filledWidth > 0 {
		filledStr := strings.Repeat(p.filledChar, filledWidth)
		bar = theme.ApplyGradient(filledStr, t.Primary, t.Secondary, t.Tertiary)
	} else if filledWidth > 0 {
		filledStyle := lipgloss.NewStyle().Foreground(t.Primary)
		bar = filledStyle.Render(strings.Repeat(p.filledChar, filledWidth))
	}

	if emptyWidth > 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
		bar += emptyStyle.Render(strings.Repeat(p.emptyChar, emptyWidth))
	}

	var parts []string

	if p.showLabel && p.label != "" {
		labelStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		parts = append(parts, labelStyle.Render(p.label))
	}

	parts = append(parts, bar)

	if p.showPercent {
		percentStyle := lipgloss.NewStyle().Foreground(t.FgMuted)
		percent := fmt.Sprintf("%3.0f%%", p.progress*100)
		parts = append(parts, percentStyle.Render(percent))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, strings.Join(parts, " "))
}
