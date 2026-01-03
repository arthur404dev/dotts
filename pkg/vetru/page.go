package vetru

import (
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	tea.Model

	ID() messages.PageID
	Title() string

	SetSize(width, height int)

	Focus() tea.Cmd
	Blur()
}
