package styles

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	Primary   = lipgloss.Color("#89b4fa")
	Secondary = lipgloss.Color("#a6e3a1")
	Accent    = lipgloss.Color("#f9e2af")
	Error     = lipgloss.Color("#f38ba8")
	Warning   = lipgloss.Color("#fab387")
	Muted     = lipgloss.Color("#6c7086")
	Surface   = lipgloss.Color("#313244")
	Base      = lipgloss.Color("#1e1e2e")
	Text      = lipgloss.Color("#cdd6f4")
	Subtext   = lipgloss.Color("#a6adc8")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Subtext).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Accent)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Surface).
			Padding(0, 2).
			MarginBottom(1)
)

var (
	SuccessIcon = SuccessStyle.Render("✓")
	ErrorIcon   = ErrorStyle.Render("✗")
	WarningIcon = WarningStyle.Render("!")
	InfoIcon    = AccentStyle.Render("ℹ")
	PendingIcon = MutedStyle.Render("○")
	ActiveIcon  = AccentStyle.Render("●")
)

func Banner() string {
	banner := `
     _       _   _       
  __| | ___ | |_| |_ ___ 
 / _' |/ _ \| __| __/ __|
| (_| | (_) | |_| |_\__ \
 \__,_|\___/ \__|\__|___/
`
	return lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		Render(banner)
}

func Title(text string) string {
	return TitleStyle.Render(text)
}

func Subtitle(text string) string {
	return SubtitleStyle.Render(text)
}

func Success(text string) string {
	return SuccessStyle.Render(SuccessIcon + " " + text)
}

func Err(text string) string {
	return ErrorStyle.Render(ErrorIcon + " " + text)
}

func Warn(text string) string {
	return WarningStyle.Render(WarningIcon + " " + text)
}

func Info(text string) string {
	return AccentStyle.Render(InfoIcon + " " + text)
}

func Mute(text string) string {
	return MutedStyle.Render(text)
}

func Box(content string) string {
	return BoxStyle.Render(content)
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func GetHuhTheme() *huh.Theme {
	t := huh.ThemeCatppuccin()
	return t
}

func StatusLine(icon, label, value string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		icon+" ",
		lipgloss.NewStyle().Width(20).Render(label+":"),
		lipgloss.NewStyle().Foreground(Text).Render(value),
	)
}

func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	filled := int(float64(current) / float64(total) * float64(width))
	empty := width - filled

	bar := lipgloss.NewStyle().Foreground(Primary).Render(repeat("━", filled))
	bar += lipgloss.NewStyle().Foreground(Muted).Render(repeat("━", empty))

	return bar
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
