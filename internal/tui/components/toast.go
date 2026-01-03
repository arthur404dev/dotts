package components

import (
	"time"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToastVariant defines the type of toast.
type ToastVariant int

const (
	ToastInfo ToastVariant = iota
	ToastSuccess
	ToastWarning
	ToastError
)

// ToastDismissMsg is sent when a toast should be dismissed.
type ToastDismissMsg struct{}

// Toast is a temporary notification component.
type Toast struct {
	theme    *theme.Theme
	variant  ToastVariant
	message  string
	duration time.Duration
	visible  bool
}

// NewToast creates a new Toast component.
func NewToast(t *theme.Theme, variant ToastVariant, message string) *Toast {
	return &Toast{
		theme:    t,
		variant:  variant,
		message:  message,
		duration: 3 * time.Second,
		visible:  true,
	}
}

// InfoToast creates an info toast.
func InfoToast(t *theme.Theme, message string) *Toast {
	return NewToast(t, ToastInfo, message)
}

// SuccessToast creates a success toast.
func SuccessToast(t *theme.Theme, message string) *Toast {
	return NewToast(t, ToastSuccess, message)
}

// WarningToast creates a warning toast.
func WarningToast(t *theme.Theme, message string) *Toast {
	return NewToast(t, ToastWarning, message)
}

// ErrorToast creates an error toast.
func ErrorToast(t *theme.Theme, message string) *Toast {
	return NewToast(t, ToastError, message)
}

// SetDuration sets the auto-dismiss duration.
func (t *Toast) SetDuration(d time.Duration) *Toast {
	t.duration = d
	return t
}

// SetMessage sets the toast message.
func (t *Toast) SetMessage(msg string) *Toast {
	t.message = msg
	return t
}

// Show makes the toast visible and returns a dismiss command.
func (t *Toast) Show() tea.Cmd {
	t.visible = true
	return tea.Tick(t.duration, func(time.Time) tea.Msg {
		return ToastDismissMsg{}
	})
}

// Dismiss hides the toast.
func (t *Toast) Dismiss() {
	t.visible = false
}

// Visible returns whether the toast is visible.
func (t *Toast) Visible() bool {
	return t.visible
}

// Init implements tea.Model.
func (t *Toast) Init() tea.Cmd {
	return t.Show()
}

// Update handles events.
func (t *Toast) Update(msg tea.Msg) (*Toast, tea.Cmd) {
	switch msg.(type) {
	case ToastDismissMsg:
		t.Dismiss()
	}
	return t, nil
}

// View renders the toast.
func (ts *Toast) View() string {
	if !ts.visible {
		return ""
	}

	t := ts.theme

	var icon string
	var color lipgloss.TerminalColor

	switch ts.variant {
	case ToastSuccess:
		icon = theme.Icons.Success
		color = t.Success
	case ToastWarning:
		icon = theme.Icons.Warning
		color = t.Warning
	case ToastError:
		icon = theme.Icons.Error
		color = t.Error
	default: // ToastInfo
		icon = theme.Icons.Info
		color = t.Info
	}

	iconStyle := lipgloss.NewStyle().Foreground(color)
	msgStyle := lipgloss.NewStyle().Foreground(t.FgBase)

	content := iconStyle.Render(icon) + " " + msgStyle.Render(ts.message)

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1)

	return containerStyle.Render(content)
}
