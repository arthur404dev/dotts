package messages

import tea "github.com/charmbracelet/bubbletea"

type PageID string

type NavigateMsg struct {
	Page PageID
}

func Navigate(page PageID) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{Page: page}
	}
}

type ActionID string

type ActionMsg struct {
	Action ActionID
}

func Action(action ActionID) tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{Action: action}
	}
}

type OpenPaletteMsg struct{}

func OpenPalette() tea.Cmd {
	return func() tea.Msg {
		return OpenPaletteMsg{}
	}
}

type ClosePaletteMsg struct{}

func ClosePalette() tea.Cmd {
	return func() tea.Msg {
		return ClosePaletteMsg{}
	}
}

type StartWizardMsg struct{}

func StartWizard() tea.Cmd {
	return func() tea.Msg {
		return StartWizardMsg{}
	}
}

type WizardCompleteMsg struct {
	Success bool
	Error   error
}

func WizardComplete(success bool, err error) tea.Cmd {
	return func() tea.Msg {
		return WizardCompleteMsg{Success: success, Error: err}
	}
}

type ExitWizardMsg struct{}

func ExitWizard() tea.Cmd {
	return func() tea.Msg {
		return ExitWizardMsg{}
	}
}

type DialogType string

const (
	DialogConfirm DialogType = "confirm"
	DialogAlert   DialogType = "alert"
	DialogError   DialogType = "error"
)

type DialogMsg struct {
	Type     DialogType
	Title    string
	Message  string
	OnOK     tea.Cmd
	OnCancel tea.Cmd
}

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

type CloseDialogMsg struct{}

func CloseDialog() tea.Cmd {
	return func() tea.Msg {
		return CloseDialogMsg{}
	}
}

type StatusType string

const (
	StatusSuccess StatusType = "success"
	StatusWarning StatusType = "warning"
	StatusError   StatusType = "error"
	StatusInfo    StatusType = "info"
)

type StatusInfoMsg struct {
	Message string
	Type    StatusType
}

func ShowStatus(message string, stype StatusType) tea.Cmd {
	return func() tea.Msg {
		return StatusInfoMsg{Message: message, Type: stype}
	}
}

type ClearStatusMsg struct{}

func ClearStatus() tea.Cmd {
	return func() tea.Msg {
		return ClearStatusMsg{}
	}
}
