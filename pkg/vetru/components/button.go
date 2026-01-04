package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type ButtonVariant string

const (
	ButtonPrimary   ButtonVariant = "primary"
	ButtonSecondary ButtonVariant = "secondary"
	ButtonGhost     ButtonVariant = "ghost"
)

type Button struct {
	theme   *theme.Theme
	zoneID  string
	label   string
	icon    string
	iconPos string // "left" or "right"
	variant ButtonVariant
	focused bool
	hovered bool
	width   int
}

func NewButton(t *theme.Theme, label string) *Button {
	zoneID := "btn-" + strings.ReplaceAll(strings.ToLower(label), " ", "-")
	return &Button{
		theme:   t,
		zoneID:  zoneID,
		label:   label,
		variant: ButtonPrimary,
		iconPos: "left",
	}
}

func (b *Button) SetLabel(label string) *Button {
	b.label = label
	return b
}

func (b *Button) SetIcon(icon string) *Button {
	b.icon = icon
	return b
}

func (b *Button) SetIconRight() *Button {
	b.iconPos = "right"
	return b
}

func (b *Button) SetVariant(v ButtonVariant) *Button {
	b.variant = v
	return b
}

func (b *Button) SetWidth(w int) *Button {
	b.width = w
	return b
}

func (b *Button) ZoneID() string {
	return b.zoneID
}

func (b *Button) Focus() tea.Cmd {
	b.focused = true
	return nil
}

func (b *Button) Blur() {
	b.focused = false
}

func (b *Button) Focused() bool {
	return b.focused
}

func (b *Button) SetHovered(h bool) {
	b.hovered = h
}

func (b *Button) Update(msg tea.Msg) (*Button, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if zone.Get(b.zoneID).InBounds(msg) {
			b.hovered = true
		} else {
			b.hovered = false
		}
	}
	return b, nil
}

func (b *Button) View() string {
	t := b.theme

	var fg, bg lipgloss.TerminalColor
	var borderColor lipgloss.TerminalColor
	var bold bool

	switch b.variant {
	case ButtonPrimary:
		if b.focused {
			fg = t.BgBase
			bg = t.Primary
			borderColor = t.Primary
			bold = true
		} else if b.hovered {
			fg = t.BgBase
			bg = t.Secondary
			borderColor = t.Secondary
			bold = true
		} else {
			fg = t.Primary
			bg = t.BgSubtle
			borderColor = t.Border
			bold = false
		}
	case ButtonSecondary:
		if b.focused {
			fg = t.BgBase
			bg = t.Secondary
			borderColor = t.Secondary
			bold = true
		} else if b.hovered {
			fg = t.FgBase
			bg = t.BgOverlay
			borderColor = t.FgMuted
			bold = false
		} else {
			fg = t.FgMuted
			bg = t.BgSubtle
			borderColor = t.Border
			bold = false
		}
	case ButtonGhost:
		bg = nil
		borderColor = nil
		if b.focused {
			fg = t.Primary
			bold = true
		} else if b.hovered {
			fg = t.FgBase
			bold = false
		} else {
			fg = t.FgMuted
			bold = false
		}
	}

	style := lipgloss.NewStyle().
		Foreground(fg).
		Bold(bold).
		Padding(0, 2)

	if bg != nil {
		style = style.Background(bg)
	}

	if b.variant != ButtonGhost && borderColor != nil {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	if b.width > 0 {
		style = style.Width(b.width).Align(lipgloss.Center)
	}

	content := b.label
	if b.icon != "" {
		if b.iconPos == "right" {
			content = b.label + " " + b.icon
		} else {
			content = b.icon + " " + b.label
		}
	}

	return zone.Mark(b.zoneID, style.Render(content))
}
