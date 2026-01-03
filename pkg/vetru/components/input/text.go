// Package input provides text input components with hover/focus states
package input

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// TextInput is a text input component with hover and focus states
type TextInput struct {
	model   textinput.Model
	theme   *theme.Theme
	zoneID  string
	label   string
	help    string
	focused bool
	hovered bool
	width   int
}

// New creates a new text input component
func New(t *theme.Theme, label, placeholder string) *TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 40

	// Generate a unique zone ID
	zoneID := "input-" + label

	return &TextInput{
		model:  ti,
		theme:  t,
		zoneID: zoneID,
		label:  label,
		width:  50,
	}
}

// SetWidth sets the input width
func (t *TextInput) SetWidth(w int) {
	t.width = w
	t.model.Width = w - 4 // Account for border padding
}

// SetLabel sets the input label
func (t *TextInput) SetLabel(label string) {
	t.label = label
}

// SetHelp sets the help text
func (t *TextInput) SetHelp(help string) {
	t.help = help
}

// SetPlaceholder sets the placeholder text
func (t *TextInput) SetPlaceholder(placeholder string) {
	t.model.Placeholder = placeholder
}

// SetValue sets the input value
func (t *TextInput) SetValue(value string) {
	t.model.SetValue(value)
}

// Value returns the current input value
func (t *TextInput) Value() string {
	return t.model.Value()
}

// ZoneID returns the zone ID for mouse tracking
func (t *TextInput) ZoneID() string {
	return t.zoneID
}

// Focus focuses the input
func (t *TextInput) Focus() tea.Cmd {
	t.focused = true
	return t.model.Focus()
}

// Blur removes focus from the input
func (t *TextInput) Blur() {
	t.focused = false
	t.model.Blur()
}

// Focused returns whether the input is focused
func (t *TextInput) Focused() bool {
	return t.focused
}

// Hovered returns whether the input is hovered
func (t *TextInput) Hovered() bool {
	return t.hovered
}

// Update handles input events
func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		inBounds := zone.Get(t.zoneID).InBounds(msg)

		switch msg.Action {
		case tea.MouseActionMotion:
			// Update hover state
			t.hovered = inBounds

		case tea.MouseActionRelease:
			// Click to focus
			if inBounds && !t.focused {
				t.focused = true
				return t, t.model.Focus()
			}
		}
	}

	// Pass events to underlying model if focused
	if t.focused {
		var cmd tea.Cmd
		t.model, cmd = t.model.Update(msg)
		return t, cmd
	}

	return t, nil
}

// View renders the input
func (t *TextInput) View() string {
	th := t.theme

	// Label
	labelView := th.S().TextInput.Label.Render(t.label)

	// Determine border style based on state
	var borderStyle lipgloss.Style
	switch {
	case t.focused:
		borderStyle = th.S().TextInput.FocusedB.Width(t.width)
	case t.hovered:
		borderStyle = th.S().TextInput.Hovered.Width(t.width)
	default:
		borderStyle = th.S().TextInput.Normal.Width(t.width)
	}

	// Input field
	fieldView := borderStyle.Render(t.model.View())

	// Help text
	var helpView string
	if t.help != "" {
		helpView = th.S().TextInput.Help.Render(t.help)
	}

	// Compose
	content := lipgloss.JoinVertical(lipgloss.Left, labelView, fieldView)
	if helpView != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, helpView)
	}

	// Wrap in zone for mouse tracking
	return zone.Mark(t.zoneID, content)
}

// Blink returns the blink command for cursor animation
func (t *TextInput) Blink() tea.Msg {
	return textinput.Blink()
}
