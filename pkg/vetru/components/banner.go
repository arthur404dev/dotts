package components

import "github.com/charmbracelet/lipgloss"

const bannerArt = `
     _       _   _       
  __| | ___ | |_| |_ ___ 
 / _' |/ _ \| __| __/ __|
| (_| | (_) | |_| |_\__ \
 \__,_|\___/ \__|\__|___/
`

type Banner struct {
	style   lipgloss.Style
	version string
}

func NewBanner(color lipgloss.Color, version string) *Banner {
	return &Banner{
		style: lipgloss.NewStyle().
			Foreground(color).
			Bold(true),
		version: version,
	}
}

func (b *Banner) Render(width int) string {
	banner := b.style.Render(bannerArt)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086"))

	version := versionStyle.Render("v" + b.version)

	return lipgloss.JoinVertical(lipgloss.Left, banner, version)
}

func (b *Banner) Height() int {
	return 7
}
