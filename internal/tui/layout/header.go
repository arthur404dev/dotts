// Package layout provides layout primitives for the TUI
package layout

import (
	"strings"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

const diag = "╱"

// Header renders the application header with logo and system info
type Header struct {
	theme   *theme.Theme
	version string
	sysInfo string
	width   int
}

// NewHeader creates a new header component
func NewHeader(t *theme.Theme, version string) *Header {
	return &Header{
		theme:   t,
		version: version,
	}
}

// SetWidth sets the header width
func (h *Header) SetWidth(w int) {
	h.width = w
}

// SetSysInfo sets the system info text (e.g., "linux · arch · x86_64")
func (h *Header) SetSysInfo(info string) {
	h.sysInfo = info
}

// View renders the header
func (h *Header) View() string {
	t := h.theme
	if h.width < 40 {
		return h.viewCompact()
	}

	// Brand with gradient
	brand := theme.ApplyBoldGradient("dotts", t.Primary, t.Secondary)

	// Version
	version := lipgloss.NewStyle().
		Foreground(t.Tertiary).
		Render(h.version)

	logo := brand + " " + version
	logoWidth := lipgloss.Width(logo)

	// System info on right (discrete)
	sysInfo := ""
	sysInfoWidth := 0
	if h.sysInfo != "" {
		sysInfo = lipgloss.NewStyle().
			Foreground(t.FgMuted).
			Render(h.sysInfo)
		sysInfoWidth = lipgloss.Width(sysInfo)
	}

	// Left diagonal field
	leftPadding := 4
	leftField := theme.ApplyGradient(
		strings.Repeat(diag, leftPadding),
		t.Primary, t.Secondary,
	)

	// Right diagonal field (fills remaining space, fades out)
	rightWidth := h.width - leftPadding - logoWidth - sysInfoWidth - 4
	if rightWidth < 10 {
		rightWidth = 10
	}
	rightField := theme.ApplyGradient(
		strings.Repeat(diag, rightWidth),
		t.Secondary, t.BgSubtle,
	)

	// Compose header
	gap := " "
	header := leftField + gap + logo + gap + rightField
	if sysInfo != "" {
		header += gap + sysInfo
	}

	return header
}

// viewCompact renders a compact header for narrow terminals
func (h *Header) viewCompact() string {
	t := h.theme
	brand := theme.ApplyBoldGradient("dotts", t.Primary, t.Secondary)
	version := lipgloss.NewStyle().Foreground(t.FgMuted).Render(h.version)
	return brand + " " + version
}
