// Package components provides reusable UI components for the TUI.
package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// Toggle is an on/off switch component with customizable labels.
// It supports mouse and keyboard interaction for toggling the enabled state.
type Toggle struct {
	theme    *theme.Theme
	label    string
	enabled  bool
	focused  bool
	hovered  bool
	zoneID   string
	onLabel  string
	offLabel string
}

// NewToggle creates a new Toggle component with the given theme and label.
// The default on/off labels are "ON" and "OFF".
func NewToggle(t *theme.Theme, label string) *Toggle {
	return &Toggle{
		theme:    t,
		label:    label,
		zoneID:   "toggle-" + label,
		onLabel:  "ON",
		offLabel: "OFF",
	}
}

// SetLabel sets the toggle label.
func (tg *Toggle) SetLabel(label string) *Toggle {
	tg.label = label
	return tg
}

// SetEnabled sets the enabled state.
func (tg *Toggle) SetEnabled(enabled bool) *Toggle {
	tg.enabled = enabled
	return tg
}

// SetLabels sets custom on/off labels for the toggle.
func (tg *Toggle) SetLabels(on, off string) *Toggle {
	tg.onLabel = on
	tg.offLabel = off
	return tg
}

// Enabled returns the current enabled state.
func (tg *Toggle) Enabled() bool {
	return tg.enabled
}

// Toggle toggles the enabled state.
func (tg *Toggle) Toggle() {
	tg.enabled = !tg.enabled
}

// Focus focuses the toggle.
func (tg *Toggle) Focus() tea.Cmd {
	tg.focused = true
	return nil
}

// Blur removes focus from the toggle.
func (tg *Toggle) Blur() {
	tg.focused = false
}

// Focused returns whether the toggle is focused.
func (tg *Toggle) Focused() bool {
	return tg.focused
}

// Hovered returns whether the toggle is hovered.
func (tg *Toggle) Hovered() bool {
	return tg.hovered
}

// ZoneID returns the zone ID for mouse tracking.
func (tg *Toggle) ZoneID() string {
	return tg.zoneID
}

// Init implements tea.Model.
func (tg *Toggle) Init() tea.Cmd {
	return nil
}

// Update handles keyboard and mouse events for the toggle.
func (tg *Toggle) Update(msg tea.Msg) (*Toggle, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if tg.focused {
			switch msg.String() {
			case " ", "enter":
				tg.Toggle()
			}
		}
	case tea.MouseMsg:
		inBounds := zone.Get(tg.zoneID).InBounds(msg)

		switch msg.Action {
		case tea.MouseActionMotion:
			tg.hovered = inBounds
		case tea.MouseActionRelease:
			if inBounds {
				tg.Toggle()
				tg.focused = true
			}
		}
	}
	return tg, nil
}

func (tg *Toggle) View() string {
	t := tg.theme

	labelStyle := lipgloss.NewStyle().Foreground(t.FgBase)
	switch {
	case tg.focused:
		labelStyle = labelStyle.Bold(true)
	case tg.hovered:
		labelStyle = labelStyle.Foreground(t.FgBase)
	default:
		labelStyle = labelStyle.Foreground(t.FgMuted)
	}

	var switchView string
	if tg.enabled {
		onStyle := lipgloss.NewStyle().Foreground(t.Success).Bold(true)
		trackStyle := lipgloss.NewStyle().Foreground(t.Success)
		switchView = trackStyle.Render("[") +
			onStyle.Render(theme.Icons.Current+"━") +
			trackStyle.Render("]") +
			" " + onStyle.Render(tg.onLabel)
	} else {
		offStyle := lipgloss.NewStyle().Foreground(t.FgMuted)
		switchView = offStyle.Render("[━" + theme.Icons.Pending + "] " + tg.offLabel)
	}

	content := labelStyle.Render(tg.label) + "  " + switchView

	return zone.Mark(tg.zoneID, content)
}
