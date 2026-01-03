package vetru

import (
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/palette"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	PageID   = messages.PageID
	ActionID = messages.ActionID
	Theme    = theme.Theme
	Command  = palette.Command
)

func DefaultTheme() *Theme {
	return theme.DefaultTheme()
}

func Run(cfg Config) error {
	model := NewModel(cfg)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)
	_, err := p.Run()
	return err
}
