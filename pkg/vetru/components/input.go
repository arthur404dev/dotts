package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Input struct {
	textinput    textinput.Model
	label        string
	help         string
	zoneID       string
	labelStyle   lipgloss.Style
	helpStyle    lipgloss.Style
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
	width        int
	focused      bool
}

func NewInput(label, placeholder string) *Input {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 40

	return &Input{
		textinput:    ti,
		label:        label,
		labelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")),
		helpStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Italic(true),
		focusedStyle: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#cba6f7")).Padding(0, 1),
		blurredStyle: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#45475a")).Padding(0, 1),
		width:        50,
	}
}

func (i *Input) SetWidth(w int) {
	i.width = w
	i.textinput.Width = w - 4
}

func (i *Input) SetLabel(label string) {
	i.label = label
}

func (i *Input) SetHelp(help string) {
	i.help = help
}

func (i *Input) SetValue(value string) {
	i.textinput.SetValue(value)
}

func (i *Input) SetZoneID(id string) {
	i.zoneID = id
}

func (i *Input) ZoneID() string {
	return i.zoneID
}

func (i *Input) Value() string {
	return i.textinput.Value()
}

func (i *Input) Focus() tea.Cmd {
	i.focused = true
	return i.textinput.Focus()
}

func (i *Input) Blur() {
	i.focused = false
	i.textinput.Blur()
}

func (i *Input) Focused() bool {
	return i.focused
}

func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	var cmd tea.Cmd
	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

func (i *Input) View() string {
	label := i.labelStyle.Render(i.label)

	var fieldStyle lipgloss.Style
	if i.focused {
		fieldStyle = i.focusedStyle.Width(i.width)
	} else {
		fieldStyle = i.blurredStyle.Width(i.width)
	}

	field := fieldStyle.Render(i.textinput.View())

	content := lipgloss.JoinVertical(lipgloss.Left, label, field)

	if i.help != "" {
		help := i.helpStyle.Render(i.help)
		content = lipgloss.JoinVertical(lipgloss.Left, content, help)
	}

	if i.zoneID != "" {
		return zone.Mark(i.zoneID, content)
	}
	return content
}

func (i *Input) Blink() tea.Msg {
	return textinput.Blink()
}
