// Package components provides reusable TUI components for the dotts application.
package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// StepperStatus represents the state of a step in the stepper.
type StepperStatus int

const (
	// StepperPending indicates the step has not been reached yet.
	StepperPending StepperStatus = iota
	// StepperCurrent indicates the step is currently active.
	StepperCurrent
	// StepperComplete indicates the step has been completed.
	StepperComplete
)

// StepItem represents a single step in the stepper.
type StepItem struct {
	// ID is the unique identifier for the step.
	ID string
	// Label is the display text for the step.
	Label string
}

// Stepper is a horizontal step indicator component that displays progress
// through a sequence of steps. It is theme-aware and provides visual
// differentiation between completed, current, and pending steps.
//
// Visual design:
//
//	Current at step 0:  → Source  ──  ○ Machine  ──  ○ Personal
//	Current at step 2:  ✓ Source  ──  ✓ Machine  ──  → Personal
type Stepper struct {
	theme      *theme.Theme
	steps      []StepItem
	current    int
	connector  string
	gap        int
	showLabels bool
}

// NewStepper creates a new Stepper component with the given theme and steps.
// By default, the stepper shows labels, uses "──" as the connector, and has
// a gap of 2 spaces around connectors.
func NewStepper(t *theme.Theme, steps []StepItem) *Stepper {
	return &Stepper{
		theme:      t,
		steps:      steps,
		current:    0,
		connector:  "──",
		gap:        2,
		showLabels: true,
	}
}

// SetConnector sets the connector string displayed between steps.
// Returns the Stepper for method chaining.
func (s *Stepper) SetConnector(c string) *Stepper {
	s.connector = c
	return s
}

// SetGap sets the number of spaces on each side of the connector.
// Returns the Stepper for method chaining.
func (s *Stepper) SetGap(g int) *Stepper {
	if g >= 0 {
		s.gap = g
	}
	return s
}

// SetShowLabels toggles whether step labels are displayed.
// When false, only the status icons are shown.
// Returns the Stepper for method chaining.
func (s *Stepper) SetShowLabels(show bool) *Stepper {
	s.showLabels = show
	return s
}

// SetCurrent sets the current step by index.
// If the index is out of bounds, no change is made.
func (s *Stepper) SetCurrent(index int) {
	if index >= 0 && index < len(s.steps) {
		s.current = index
	}
}

// SetCurrentByID sets the current step by its ID.
// If no step with the given ID exists, no change is made.
func (s *Stepper) SetCurrentByID(id string) {
	for i, step := range s.steps {
		if step.ID == id {
			s.current = i
			return
		}
	}
}

// Next advances to the next step.
// Returns true if successful, false if already at the last step.
func (s *Stepper) Next() bool {
	if s.current < len(s.steps)-1 {
		s.current++
		return true
	}
	return false
}

// Prev moves to the previous step.
// Returns true if successful, false if already at the first step.
func (s *Stepper) Prev() bool {
	if s.current > 0 {
		s.current--
		return true
	}
	return false
}

// Current returns the index of the current step.
func (s *Stepper) Current() int {
	return s.current
}

// CurrentID returns the ID of the current step.
// Returns an empty string if there are no steps.
func (s *Stepper) CurrentID() string {
	if s.current >= 0 && s.current < len(s.steps) {
		return s.steps[s.current].ID
	}
	return ""
}

// IsFirst returns true if the current step is the first step.
func (s *Stepper) IsFirst() bool {
	return s.current == 0
}

// IsLast returns true if the current step is the last step.
func (s *Stepper) IsLast() bool {
	return s.current == len(s.steps)-1
}

// Len returns the total number of steps.
func (s *Stepper) Len() int {
	return len(s.steps)
}

// Steps returns a copy of the step items.
func (s *Stepper) Steps() []StepItem {
	result := make([]StepItem, len(s.steps))
	copy(result, s.steps)
	return result
}

// getStatus returns the status of the step at the given index.
func (s *Stepper) getStatus(index int) StepperStatus {
	if index < s.current {
		return StepperComplete
	} else if index == s.current {
		return StepperCurrent
	}
	return StepperPending
}

// View renders the stepper as a styled string.
// The output displays steps horizontally with connectors between them:
//
//	Completed: ✓ Label (green)
//	Current:   → Label (primary color, bold)
//	Pending:   ○ Label (muted)
func (s *Stepper) View() string {
	if len(s.steps) == 0 {
		return ""
	}

	t := s.theme
	gap := strings.Repeat(" ", s.gap)

	var parts []string

	for i, step := range s.steps {
		status := s.getStatus(i)

		var icon string
		var style lipgloss.Style

		switch status {
		case StepperComplete:
			icon = theme.Icons.Success
			style = lipgloss.NewStyle().Foreground(t.Success)
		case StepperCurrent:
			icon = theme.Icons.ArrowRight
			style = lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
		case StepperPending:
			icon = theme.Icons.Pending
			style = lipgloss.NewStyle().Foreground(t.FgMuted)
		}

		var stepView string
		if s.showLabels {
			stepView = style.Render(icon + " " + step.Label)
		} else {
			stepView = style.Render(icon)
		}

		parts = append(parts, stepView)

		if i < len(s.steps)-1 {
			connectorStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
			parts = append(parts, gap+connectorStyle.Render(s.connector)+gap)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
