package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// BadgeVariant defines the visual style of a badge.
type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSuccess
	BadgeWarning
	BadgeError
	BadgeInfo
	BadgePrimary
	BadgeSecondary
)

// Badge is a small status indicator pill component.
type Badge struct {
	theme   *theme.Theme
	label   string
	variant BadgeVariant
	icon    string
}

// NewBadge creates a new Badge component.
func NewBadge(t *theme.Theme, label string) *Badge {
	return &Badge{
		theme:   t,
		label:   label,
		variant: BadgeDefault,
	}
}

// SetVariant sets the badge variant.
func (b *Badge) SetVariant(v BadgeVariant) *Badge {
	b.variant = v
	return b
}

// SetIcon sets an optional icon prefix.
func (b *Badge) SetIcon(icon string) *Badge {
	b.icon = icon
	return b
}

// SuccessBadge is a convenience constructor for success badges.
func SuccessBadge(t *theme.Theme, label string) *Badge {
	return NewBadge(t, label).SetVariant(BadgeSuccess).SetIcon(theme.Icons.Success)
}

// ErrorBadge is a convenience constructor for error badges.
func ErrorBadge(t *theme.Theme, label string) *Badge {
	return NewBadge(t, label).SetVariant(BadgeError).SetIcon(theme.Icons.Error)
}

// WarningBadge is a convenience constructor for warning badges.
func WarningBadge(t *theme.Theme, label string) *Badge {
	return NewBadge(t, label).SetVariant(BadgeWarning).SetIcon(theme.Icons.Warning)
}

// InfoBadge is a convenience constructor for info badges.
func InfoBadge(t *theme.Theme, label string) *Badge {
	return NewBadge(t, label).SetVariant(BadgeInfo).SetIcon(theme.Icons.Info)
}

// View renders the badge.
func (b *Badge) View() string {
	th := b.theme

	var fg, bg lipgloss.TerminalColor

	switch b.variant {
	case BadgeSuccess:
		fg = th.BgBase
		bg = th.Success
	case BadgeWarning:
		fg = th.BgBase
		bg = th.Warning
	case BadgeError:
		fg = th.BgBase
		bg = th.Error
	case BadgeInfo:
		fg = th.BgBase
		bg = th.Info
	case BadgePrimary:
		fg = th.BgBase
		bg = th.Primary
	case BadgeSecondary:
		fg = th.BgBase
		bg = th.Secondary
	default: // BadgeDefault
		fg = th.FgBase
		bg = th.BgSubtle
	}

	style := lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 1).
		Bold(true)

	content := b.label
	if b.icon != "" {
		content = b.icon + " " + b.label
	}

	return style.Render(content)
}
