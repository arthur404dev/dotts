package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// AlertVariant defines the type and visual style of an alert.
type AlertVariant int

const (
	// AlertInfo displays an informational message with blue styling.
	AlertInfo AlertVariant = iota
	// AlertSuccess displays a success message with green styling.
	AlertSuccess
	// AlertWarning displays a warning message with yellow/orange styling.
	AlertWarning
	// AlertError displays an error message with red styling.
	AlertError
)

// Alert is an inline alert message component that displays contextual
// feedback to users with appropriate visual styling based on the message type.
type Alert struct {
	theme   *theme.Theme
	variant AlertVariant
	title   string
	message string
	icon    bool
	width   int
}

// NewAlert creates a new Alert component with the specified variant and message.
func NewAlert(t *theme.Theme, variant AlertVariant, message string) *Alert {
	return &Alert{
		theme:   t,
		variant: variant,
		message: message,
		icon:    true,
	}
}

// InfoAlert is a convenience constructor for info alerts.
func InfoAlert(t *theme.Theme, message string) *Alert {
	return NewAlert(t, AlertInfo, message)
}

// SuccessAlert is a convenience constructor for success alerts.
func SuccessAlert(t *theme.Theme, message string) *Alert {
	return NewAlert(t, AlertSuccess, message)
}

// WarningAlert is a convenience constructor for warning alerts.
func WarningAlert(t *theme.Theme, message string) *Alert {
	return NewAlert(t, AlertWarning, message)
}

// ErrorAlert is a convenience constructor for error alerts.
func ErrorAlert(t *theme.Theme, message string) *Alert {
	return NewAlert(t, AlertError, message)
}

// SetTitle sets an optional title for the alert and returns the Alert
// for method chaining. The title is displayed above the message in bold.
func (a *Alert) SetTitle(title string) *Alert {
	a.title = title
	return a
}

// SetMessage updates the alert message and returns the Alert
// for method chaining.
func (a *Alert) SetMessage(message string) *Alert {
	a.message = message
	return a
}

// SetIcon enables or disables the icon display and returns the Alert
// for method chaining. Icons are enabled by default.
func (a *Alert) SetIcon(show bool) *Alert {
	a.icon = show
	return a
}

// SetWidth sets the alert width and returns the Alert for method chaining.
// A width of 0 enables auto-sizing based on content.
func (a *Alert) SetWidth(w int) *Alert {
	a.width = w
	return a
}

// View renders the alert with appropriate styling based on its variant.
func (a *Alert) View() string {
	t := a.theme

	var iconStr string
	var color lipgloss.TerminalColor
	var borderColor lipgloss.TerminalColor

	switch a.variant {
	case AlertSuccess:
		iconStr = theme.Icons.Success
		color = t.Success
		borderColor = t.Success
	case AlertWarning:
		iconStr = theme.Icons.Warning
		color = t.Warning
		borderColor = t.Warning
	case AlertError:
		iconStr = theme.Icons.Error
		color = t.Error
		borderColor = t.Error
	default: // AlertInfo
		iconStr = theme.Icons.Info
		color = t.Info
		borderColor = t.Info
	}

	// Build content parts
	var parts []string

	// Title line (if present)
	if a.title != "" {
		titleStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
		titleLine := titleStyle.Render(a.title)
		if a.icon {
			titleLine = lipgloss.NewStyle().Foreground(color).Render(iconStr) + " " + titleLine
		}
		parts = append(parts, titleLine)
	}

	// Message line
	if a.message != "" {
		msgStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		msgLine := a.message
		// If no title but icon enabled, prepend icon to message
		if a.title == "" && a.icon {
			msgLine = lipgloss.NewStyle().Foreground(color).Render(iconStr) + " " + msgLine
		}
		parts = append(parts, msgStyle.Render(msgLine))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Container with left border accent
	containerStyle := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(borderColor).
		PaddingLeft(1)

	if a.width > 0 {
		containerStyle = containerStyle.Width(a.width)
	}

	return containerStyle.Render(content)
}
