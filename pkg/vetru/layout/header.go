package layout

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

const diag = "╱"

type Header struct {
	theme   *theme.Theme
	brand   string
	version string
	sysInfo string
	width   int
}

func NewHeader(t *theme.Theme, brand, version string) *Header {
	return &Header{
		theme:   t,
		brand:   brand,
		version: version,
	}
}

func (h *Header) SetWidth(w int) {
	h.width = w
}

func (h *Header) SetSysInfo(info string) {
	h.sysInfo = info
}

func (h *Header) SetBrand(brand string) {
	h.brand = brand
}

func (h *Header) View() string {
	t := h.theme
	if h.width < 40 {
		return h.viewCompact()
	}

	brand := theme.ApplyBoldGradient(h.brand, t.Primary, t.Secondary)

	version := lipgloss.NewStyle().
		Foreground(t.Tertiary).
		Render(h.version)

	logo := brand + " " + version
	logoWidth := lipgloss.Width(logo)

	sysInfo := ""
	sysInfoWidth := 0
	if h.sysInfo != "" {
		sysInfo = lipgloss.NewStyle().
			Foreground(t.FgMuted).
			Render(h.sysInfo)
		sysInfoWidth = lipgloss.Width(sysInfo)
	}

	leftPadding := 4
	leftField := theme.ApplyGradient(
		strings.Repeat(diag, leftPadding),
		t.Primary, t.Secondary,
	)

	rightWidth := h.width - leftPadding - logoWidth - sysInfoWidth - 4
	if rightWidth < 10 {
		rightWidth = 10
	}
	rightField := theme.ApplyGradient(
		strings.Repeat(diag, rightWidth),
		t.Secondary, t.BgSubtle,
	)

	gap := " "
	header := leftField + gap + logo + gap + rightField
	if sysInfo != "" {
		header += gap + sysInfo
	}

	return header
}

func (h *Header) viewCompact() string {
	t := h.theme
	brand := theme.ApplyBoldGradient(h.brand, t.Primary, t.Secondary)
	version := lipgloss.NewStyle().Foreground(t.FgMuted).Render(h.version)
	return brand + " " + version
}
