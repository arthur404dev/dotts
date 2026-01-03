// Package components provides reusable TUI components for the dotts application.
package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerStyle defines the animation style for a spinner.
type SpinnerStyle int

const (
	// SpinnerDots displays a series of dots animation.
	SpinnerDots SpinnerStyle = iota
	// SpinnerLine displays a rotating line animation.
	SpinnerLine
	// SpinnerPulse displays a pulsing animation.
	SpinnerPulse
	// SpinnerPoints displays rotating points animation.
	SpinnerPoints
	// SpinnerGlobe displays a rotating globe animation.
	SpinnerGlobe
	// SpinnerMoon displays moon phase animation.
	SpinnerMoon
	// SpinnerMonkey displays a monkey animation.
	SpinnerMonkey
	// SpinnerMeter displays a meter animation.
	SpinnerMeter
	// SpinnerHamburger displays a hamburger menu animation.
	SpinnerHamburger
	// SpinnerEllipsis displays an ellipsis animation.
	SpinnerEllipsis
)

// Spinner is a loading indicator component that wraps bubbles/spinner
// with theme support and additional configuration options.
type Spinner struct {
	theme   *theme.Theme
	spinner spinner.Model
	label   string
}

// NewSpinner creates a new Spinner component with the default dots style.
func NewSpinner(t *theme.Theme) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(t.Primary)

	return &Spinner{
		theme:   t,
		spinner: s,
	}
}

// SetStyle sets the spinner animation style and returns the Spinner
// for method chaining.
func (s *Spinner) SetStyle(style SpinnerStyle) *Spinner {
	var sp spinner.Spinner
	switch style {
	case SpinnerDots:
		sp = spinner.Dot
	case SpinnerLine:
		sp = spinner.Line
	case SpinnerPulse:
		sp = spinner.Pulse
	case SpinnerPoints:
		sp = spinner.Points
	case SpinnerGlobe:
		sp = spinner.Globe
	case SpinnerMoon:
		sp = spinner.Moon
	case SpinnerMonkey:
		sp = spinner.Monkey
	case SpinnerMeter:
		sp = spinner.Meter
	case SpinnerHamburger:
		sp = spinner.Hamburger
	case SpinnerEllipsis:
		sp = spinner.Ellipsis
	default:
		sp = spinner.Dot
	}
	s.spinner.Spinner = sp
	return s
}

// SetLabel sets the text displayed next to the spinner and returns
// the Spinner for method chaining.
func (s *Spinner) SetLabel(label string) *Spinner {
	s.label = label
	return s
}

// SetColor overrides the spinner color and returns the Spinner
// for method chaining.
func (s *Spinner) SetColor(color lipgloss.TerminalColor) *Spinner {
	s.spinner.Style = lipgloss.NewStyle().Foreground(color)
	return s
}

// Init implements tea.Model and returns the initial tick command
// to start the spinner animation.
func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update handles spinner animation updates and returns the updated
// Spinner along with any commands to execute.
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View renders the spinner with its optional label.
func (s *Spinner) View() string {
	if s.label == "" {
		return s.spinner.View()
	}

	labelStyle := lipgloss.NewStyle().Foreground(s.theme.FgBase)
	return s.spinner.View() + " " + labelStyle.Render(s.label)
}

// Tick returns the tick command for animation. This can be used
// to manually trigger spinner updates.
func (s *Spinner) Tick() tea.Msg {
	return s.spinner.Tick()
}
