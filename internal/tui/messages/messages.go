// Package messages defines shared message types for TUI communication
package messages

import tea "github.com/charmbracelet/bubbletea"

// PageID identifies different pages in the TUI
type PageID string

const (
	PageDashboard PageID = "dashboard"
	PageStatus    PageID = "status"
	PageUpdate    PageID = "update"
	PageDoctor    PageID = "doctor"
	PageSettings  PageID = "settings"
	PageWizard    PageID = "wizard"
)

// ActionID identifies different actions that can be triggered
type ActionID string

const (
	ActionSync   ActionID = "sync"
	ActionUpdate ActionID = "update"
	ActionDoctor ActionID = "doctor"
)

// NavigateMsg requests navigation to a different page
type NavigateMsg struct {
	Page PageID
}

// Navigate returns a command that navigates to the specified page
func Navigate(page PageID) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{Page: page}
	}
}

// ActionMsg requests an action to be performed
type ActionMsg struct {
	Action ActionID
}

// Action returns a command that triggers the specified action
func Action(action ActionID) tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{Action: action}
	}
}

// OpenPaletteMsg requests the command palette to open
type OpenPaletteMsg struct{}

// OpenPalette returns a command that opens the command palette
func OpenPalette() tea.Cmd {
	return func() tea.Msg {
		return OpenPaletteMsg{}
	}
}

// ClosePaletteMsg requests the command palette to close
type ClosePaletteMsg struct{}

// ClosePalette returns a command that closes the command palette
func ClosePalette() tea.Cmd {
	return func() tea.Msg {
		return ClosePaletteMsg{}
	}
}

// StartWizardMsg requests the init wizard to start
type StartWizardMsg struct{}

// StartWizard returns a command that starts the init wizard
func StartWizard() tea.Cmd {
	return func() tea.Msg {
		return StartWizardMsg{}
	}
}

// WizardCompleteMsg indicates the wizard has completed
type WizardCompleteMsg struct {
	Success bool
	Error   error
}

// WizardComplete returns a command that signals wizard completion
func WizardComplete(success bool, err error) tea.Cmd {
	return func() tea.Msg {
		return WizardCompleteMsg{Success: success, Error: err}
	}
}

// ExitWizardMsg requests to exit the wizard (only valid when wizard allows it)
type ExitWizardMsg struct{}

// ExitWizard returns a command that exits the wizard
func ExitWizard() tea.Cmd {
	return func() tea.Msg {
		return ExitWizardMsg{}
	}
}

// DialogMsg manages dialog state
type DialogMsg struct {
	Type     DialogType
	Title    string
	Message  string
	OnOK     tea.Cmd
	OnCancel tea.Cmd
}

// DialogType identifies different dialog types
type DialogType string

const (
	DialogConfirm DialogType = "confirm"
	DialogAlert   DialogType = "alert"
	DialogError   DialogType = "error"
)

// ShowDialog returns a command that shows a dialog
func ShowDialog(dtype DialogType, title, message string, onOK, onCancel tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return DialogMsg{
			Type:     dtype,
			Title:    title,
			Message:  message,
			OnOK:     onOK,
			OnCancel: onCancel,
		}
	}
}

// CloseDialogMsg requests the current dialog to close
type CloseDialogMsg struct{}

// CloseDialog returns a command that closes the current dialog
func CloseDialog() tea.Cmd {
	return func() tea.Msg {
		return CloseDialogMsg{}
	}
}

// StatusInfoMsg displays a temporary status message
type StatusInfoMsg struct {
	Message string
	Type    StatusType
}

// StatusType identifies the type of status message
type StatusType string

const (
	StatusSuccess StatusType = "success"
	StatusWarning StatusType = "warning"
	StatusError   StatusType = "error"
	StatusInfo    StatusType = "info"
)

// ShowStatus returns a command that shows a status message
func ShowStatus(message string, stype StatusType) tea.Cmd {
	return func() tea.Msg {
		return StatusInfoMsg{Message: message, Type: stype}
	}
}

// ClearStatusMsg clears the current status message
type ClearStatusMsg struct{}

// ClearStatus returns a command that clears the status message
func ClearStatus() tea.Cmd {
	return func() tea.Msg {
		return ClearStatusMsg{}
	}
}
