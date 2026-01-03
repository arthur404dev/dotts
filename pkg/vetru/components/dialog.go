package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// DialogType defines the type of dialog.
type DialogType int

const (
	DialogConfirm DialogType = iota
	DialogAlert
	DialogError
)

// DialogFocus indicates which button is focused.
type DialogFocus int

const (
	DialogFocusConfirm DialogFocus = iota
	DialogFocusCancel
)

// Dialog is a pre-built dialog component.
type Dialog struct {
	theme       *theme.Theme
	dialogType  DialogType
	title       string
	message     string
	confirmText string
	cancelText  string
	focus       DialogFocus
	width       int
	zoneID      string
}

// NewDialog creates a new Dialog component.
func NewDialog(t *theme.Theme, dtype DialogType) *Dialog {
	return &Dialog{
		theme:       t,
		dialogType:  dtype,
		confirmText: "OK",
		cancelText:  "Cancel",
		focus:       DialogFocusConfirm,
		width:       50,
		zoneID:      "dialog",
	}
}

// ConfirmDialog creates a confirmation dialog.
func ConfirmDialog(t *theme.Theme, title, message string) *Dialog {
	return NewDialog(t, DialogConfirm).
		SetTitle(title).
		SetMessage(message).
		SetButtons("Yes", "No")
}

// AlertDialog creates an alert dialog.
func AlertDialog(t *theme.Theme, title, message string) *Dialog {
	return NewDialog(t, DialogAlert).
		SetTitle(title).
		SetMessage(message).
		SetButtons("OK", "")
}

// ErrorDialog creates an error dialog.
func ErrorDialog(t *theme.Theme, title, message string) *Dialog {
	return NewDialog(t, DialogError).
		SetTitle(title).
		SetMessage(message).
		SetButtons("OK", "")
}

// SetTitle sets the dialog title.
func (d *Dialog) SetTitle(title string) *Dialog {
	d.title = title
	return d
}

// SetMessage sets the dialog message.
func (d *Dialog) SetMessage(message string) *Dialog {
	d.message = message
	return d
}

// SetButtons sets the confirm and cancel button labels.
func (d *Dialog) SetButtons(confirm, cancel string) *Dialog {
	d.confirmText = confirm
	d.cancelText = cancel
	return d
}

// SetWidth sets the dialog width.
func (d *Dialog) SetWidth(w int) *Dialog {
	if w > 0 {
		d.width = w
	}
	return d
}

// Focus returns the current focus.
func (d *Dialog) Focus() DialogFocus {
	return d.focus
}

// IsConfirmed returns true if confirm button is focused.
func (d *Dialog) IsConfirmed() bool {
	return d.focus == DialogFocusConfirm
}

// Init implements tea.Model.
func (d *Dialog) Init() tea.Cmd {
	return nil
}

// Update handles events.
func (d *Dialog) Update(msg tea.Msg) (*Dialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "tab", "shift+tab":
			if d.cancelText != "" {
				if d.focus == DialogFocusConfirm {
					d.focus = DialogFocusCancel
				} else {
					d.focus = DialogFocusConfirm
				}
			}
		case "right", "l":
			if d.cancelText != "" {
				if d.focus == DialogFocusCancel {
					d.focus = DialogFocusConfirm
				} else {
					d.focus = DialogFocusCancel
				}
			}
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease {
			confirmZone := d.zoneID + "-confirm"
			cancelZone := d.zoneID + "-cancel"
			if zone.Get(confirmZone).InBounds(msg) {
				d.focus = DialogFocusConfirm
			} else if zone.Get(cancelZone).InBounds(msg) {
				d.focus = DialogFocusCancel
			}
		}
	}
	return d, nil
}

// View renders the dialog.
func (d *Dialog) View() string {
	t := d.theme

	var parts []string

	var icon string
	var titleColor lipgloss.TerminalColor
	switch d.dialogType {
	case DialogError:
		icon = theme.Icons.Error + " "
		titleColor = t.Error
	case DialogAlert:
		icon = theme.Icons.Warning + " "
		titleColor = t.Warning
	default:
		icon = theme.Icons.Info + " "
		titleColor = t.Info
	}

	if d.title != "" {
		titleStyle := lipgloss.NewStyle().Foreground(titleColor).Bold(true)
		parts = append(parts, titleStyle.Render(icon+d.title))
		parts = append(parts, "")
	}

	if d.message != "" {
		msgStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		parts = append(parts, msgStyle.Render(d.message))
		parts = append(parts, "")
	}

	var buttons []string

	confirmStyle := lipgloss.NewStyle().Padding(0, 2)
	if d.focus == DialogFocusConfirm {
		confirmStyle = confirmStyle.
			Background(t.Primary).
			Foreground(t.BgBase).
			Bold(true)
	} else {
		confirmStyle = confirmStyle.
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border).
			Foreground(t.FgBase)
	}
	confirmView := zone.Mark(d.zoneID+"-confirm", confirmStyle.Render(d.confirmText))
	buttons = append(buttons, confirmView)

	if d.cancelText != "" {
		cancelStyle := lipgloss.NewStyle().Padding(0, 2)
		if d.focus == DialogFocusCancel {
			cancelStyle = cancelStyle.
				Background(t.FgMuted).
				Foreground(t.BgBase).
				Bold(true)
		} else {
			cancelStyle = cancelStyle.
				Border(lipgloss.NormalBorder()).
				BorderForeground(t.Border).
				Foreground(t.FgMuted)
		}
		cancelView := zone.Mark(d.zoneID+"-cancel", cancelStyle.Render(d.cancelText))
		buttons = append(buttons, "  ")
		buttons = append(buttons, cancelView)
	}

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)
	parts = append(parts, buttonRow)

	inner := lipgloss.JoinVertical(lipgloss.Left, parts...)

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus).
		Padding(1, 2).
		Width(d.width)

	return containerStyle.Render(inner)
}
