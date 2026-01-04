// Package keys provides centralized key binding definitions for the TUI.
package keys

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	OpenPalette = key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "commands"),
	)

	ForceQuit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "force quit"),
	)

	Escape = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back/close"),
	)

	Help = key.NewBinding(
		key.WithKeys("ctrl+g", "?"),
		key.WithHelp("ctrl+g", "help"),
	)

	Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)
)

var (
	ListUp = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	)

	ListDown = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	)

	ListUpAlt = key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑", "up"),
	)

	ListDownAlt = key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓", "down"),
	)

	ListHome = key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "first"),
	)

	ListEnd = key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "last"),
	)

	Confirm = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	)
)

var (
	NavLeft = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev"),
	)

	NavRight = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next"),
	)

	StepPrev = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "prev step"),
	)

	StepNext = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next step"),
	)
)

var (
	NextField = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	)

	PrevField = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	)
)

var (
	ScrollUp = key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup", "scroll up"),
	)

	ScrollDown = key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn", "scroll down"),
	)

	ScrollTop = key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("home", "scroll top"),
	)

	ScrollBottom = key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("end", "scroll bottom"),
	)
)

func MatchesListUp(msg tea.KeyMsg) bool {
	return key.Matches(msg, ListUp)
}

func MatchesListDown(msg tea.KeyMsg) bool {
	return key.Matches(msg, ListDown)
}

func MatchesListUpAlt(msg tea.KeyMsg) bool {
	return key.Matches(msg, ListUpAlt)
}

func MatchesListDownAlt(msg tea.KeyMsg) bool {
	return key.Matches(msg, ListDownAlt)
}

func MatchesNavLeft(msg tea.KeyMsg) bool {
	return key.Matches(msg, NavLeft)
}

func MatchesNavRight(msg tea.KeyMsg) bool {
	return key.Matches(msg, NavRight)
}

func MatchesNextField(msg tea.KeyMsg) bool {
	return key.Matches(msg, NextField)
}

func MatchesPrevField(msg tea.KeyMsg) bool {
	return key.Matches(msg, PrevField)
}

func MatchesConfirm(msg tea.KeyMsg) bool {
	return key.Matches(msg, Confirm)
}

func MatchesEscape(msg tea.KeyMsg) bool {
	return key.Matches(msg, Escape)
}

func MatchesStepPrev(msg tea.KeyMsg) bool {
	return key.Matches(msg, StepPrev)
}

func MatchesStepNext(msg tea.KeyMsg) bool {
	return key.Matches(msg, StepNext)
}

func MatchesScrollUp(msg tea.KeyMsg) bool {
	return key.Matches(msg, ScrollUp)
}

func MatchesScrollDown(msg tea.KeyMsg) bool {
	return key.Matches(msg, ScrollDown)
}

func MatchesScrollTop(msg tea.KeyMsg) bool {
	return key.Matches(msg, ScrollTop)
}

func MatchesScrollBottom(msg tea.KeyMsg) bool {
	return key.Matches(msg, ScrollBottom)
}
