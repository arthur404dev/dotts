package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type KeyBinding struct {
	Key  string
	Desc string
}

type Help struct {
	bindings  []KeyBinding
	keyStyle  lipgloss.Style
	descStyle lipgloss.Style
	sepStyle  lipgloss.Style
	width     int
}

func NewHelp(keyStyle, descStyle lipgloss.Style) *Help {
	return &Help{
		bindings:  []KeyBinding{},
		keyStyle:  keyStyle,
		descStyle: descStyle,
		sepStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#585b70")),
		width:     80,
	}
}

func (h *Help) SetWidth(w int) {
	h.width = w
}

func (h *Help) SetBindings(bindings []KeyBinding) {
	h.bindings = bindings
}

func (h *Help) Render() string {
	if len(h.bindings) == 0 {
		return ""
	}

	var parts []string
	for _, b := range h.bindings {
		key := h.keyStyle.Render("[" + b.Key + "]")
		desc := h.descStyle.Render(b.Desc)
		parts = append(parts, key+" "+desc)
	}

	sep := h.sepStyle.Render("  ")
	content := strings.Join(parts, sep)

	container := lipgloss.NewStyle().
		Width(h.width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return container.Render(content)
}
