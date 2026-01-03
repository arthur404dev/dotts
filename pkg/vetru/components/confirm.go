// Package components provides reusable TUI components for the dotts application.
package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// ConfirmFocus indicates which option is currently focused in the confirmation dialog.
type ConfirmFocus int

const (
	// ConfirmYes indicates the Yes option is focused.
	ConfirmYes ConfirmFocus = iota
	// ConfirmNo indicates the No option is focused.
	ConfirmNo
)

// Confirm is an inline yes/no confirmation component with keyboard and mouse support.
// It displays a question with two options that can be toggled with arrow keys or clicked.
type Confirm struct {
	theme    *theme.Theme
	question string
	yesLabel string
	noLabel  string
	focus    ConfirmFocus
	focused  bool
	zoneID   string
}

// NewConfirm creates a new Confirm component with the given theme and question.
// The default selection is No for safety.
func NewConfirm(t *theme.Theme, question string) *Confirm {
	return &Confirm{
		theme:    t,
		question: question,
		yesLabel: "Yes",
		noLabel:  "No",
		focus:    ConfirmNo,
		zoneID:   "confirm",
	}
}

// SetQuestion sets the question text to display.
func (c *Confirm) SetQuestion(q string) *Confirm {
	c.question = q
	return c
}

// SetLabels sets custom labels for the Yes and No options.
func (c *Confirm) SetLabels(yes, no string) *Confirm {
	c.yesLabel = yes
	c.noLabel = no
	return c
}

// SetDefaultYes sets the default selection to Yes.
func (c *Confirm) SetDefaultYes() *Confirm {
	c.focus = ConfirmYes
	return c
}

// SetDefaultNo sets the default selection to No.
func (c *Confirm) SetDefaultNo() *Confirm {
	c.focus = ConfirmNo
	return c
}

// Value returns the current selection: true if Yes is selected, false if No.
func (c *Confirm) Value() bool {
	return c.focus == ConfirmYes
}

// Focus sets the component to focused state.
func (c *Confirm) Focus() tea.Cmd {
	c.focused = true
	return nil
}

// Blur removes focus from the component.
func (c *Confirm) Blur() {
	c.focused = false
}

// Focused returns whether the component is currently focused.
func (c *Confirm) Focused() bool {
	return c.focused
}

// Init implements tea.Model. Returns nil as no initialization is needed.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update handles keyboard and mouse events.
// Keyboard: left/h/y (select Yes), right/l/n (select No), tab (toggle).
// Mouse: click on Yes or No to select it.
func (c *Confirm) Update(msg tea.Msg) (*Confirm, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "y", "Y":
			c.focus = ConfirmYes
		case "right", "l", "n", "N":
			c.focus = ConfirmNo
		case "tab", "shift+tab":
			if c.focus == ConfirmYes {
				c.focus = ConfirmNo
			} else {
				c.focus = ConfirmYes
			}
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease {
			yesZone := c.zoneID + "-yes"
			noZone := c.zoneID + "-no"
			if zone.Get(yesZone).InBounds(msg) {
				c.focus = ConfirmYes
			} else if zone.Get(noZone).InBounds(msg) {
				c.focus = ConfirmNo
			}
		}
	}
	return c, nil
}

// View renders the confirmation dialog with the question and Yes/No buttons.
// The selected option is highlighted with the appropriate status color.
func (c *Confirm) View() string {
	t := c.theme

	questionStyle := lipgloss.NewStyle().Foreground(t.FgBase)

	yesStyle := lipgloss.NewStyle().Padding(0, 1)
	if c.focus == ConfirmYes {
		yesStyle = yesStyle.
			Background(t.Success).
			Foreground(t.BgBase).
			Bold(true)
	} else {
		yesStyle = yesStyle.
			Foreground(t.FgMuted)
	}

	noStyle := lipgloss.NewStyle().Padding(0, 1)
	if c.focus == ConfirmNo {
		noStyle = noStyle.
			Background(t.Error).
			Foreground(t.BgBase).
			Bold(true)
	} else {
		noStyle = noStyle.
			Foreground(t.FgMuted)
	}

	yesView := zone.Mark(c.zoneID+"-yes", yesStyle.Render(c.yesLabel))
	noView := zone.Mark(c.zoneID+"-no", noStyle.Render(c.noLabel))

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesView, " ", noView)

	return questionStyle.Render(c.question) + "  " + buttons
}
