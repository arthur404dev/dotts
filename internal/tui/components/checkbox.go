// Package components provides reusable UI components for the TUI.
package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// Checkbox is an interactive checkbox component with label, hover, and focus states.
// It supports mouse and keyboard interaction for toggling the checked state.
type Checkbox struct {
	theme   *theme.Theme
	label   string
	checked bool
	focused bool
	hovered bool
	zoneID  string
}

// NewCheckbox creates a new Checkbox component with the given theme and label.
func NewCheckbox(t *theme.Theme, label string) *Checkbox {
	return &Checkbox{
		theme:  t,
		label:  label,
		zoneID: "checkbox-" + label,
	}
}

// SetLabel sets the checkbox label.
func (c *Checkbox) SetLabel(label string) *Checkbox {
	c.label = label
	return c
}

// SetChecked sets the checked state.
func (c *Checkbox) SetChecked(checked bool) *Checkbox {
	c.checked = checked
	return c
}

// Checked returns the current checked state.
func (c *Checkbox) Checked() bool {
	return c.checked
}

// Toggle toggles the checked state.
func (c *Checkbox) Toggle() {
	c.checked = !c.checked
}

// Focus focuses the checkbox.
func (c *Checkbox) Focus() tea.Cmd {
	c.focused = true
	return nil
}

// Blur removes focus from the checkbox.
func (c *Checkbox) Blur() {
	c.focused = false
}

// Focused returns whether the checkbox is focused.
func (c *Checkbox) Focused() bool {
	return c.focused
}

// Hovered returns whether the checkbox is hovered.
func (c *Checkbox) Hovered() bool {
	return c.hovered
}

// ZoneID returns the zone ID for mouse tracking.
func (c *Checkbox) ZoneID() string {
	return c.zoneID
}

// Init implements tea.Model.
func (c *Checkbox) Init() tea.Cmd {
	return nil
}

// Update handles keyboard and mouse events for the checkbox.
func (c *Checkbox) Update(msg tea.Msg) (*Checkbox, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.focused {
			switch msg.String() {
			case " ", "enter":
				c.Toggle()
			}
		}
	case tea.MouseMsg:
		inBounds := zone.Get(c.zoneID).InBounds(msg)

		switch msg.Action {
		case tea.MouseActionMotion:
			c.hovered = inBounds
		case tea.MouseActionRelease:
			if inBounds {
				c.Toggle()
				c.focused = true
			}
		}
	}
	return c, nil
}

func (c *Checkbox) View() string {
	t := c.theme

	var box string
	var boxStyle lipgloss.Style

	if c.checked {
		box = "[" + theme.Icons.Success + "]"
		boxStyle = lipgloss.NewStyle().Foreground(t.Primary)
	} else {
		box = "[ ]"
		switch {
		case c.focused:
			boxStyle = lipgloss.NewStyle().Foreground(t.Primary)
		case c.hovered:
			boxStyle = lipgloss.NewStyle().Foreground(t.FgBase)
		default:
			boxStyle = lipgloss.NewStyle().Foreground(t.FgMuted)
		}
	}

	labelStyle := lipgloss.NewStyle().Foreground(t.FgBase)
	switch {
	case c.focused:
		labelStyle = labelStyle.Bold(true)
	case c.hovered:
		labelStyle = labelStyle.Foreground(t.FgBase)
	default:
		labelStyle = labelStyle.Foreground(t.FgMuted)
	}

	content := boxStyle.Render(box) + " " + labelStyle.Render(c.label)

	return zone.Mark(c.zoneID, content)
}
