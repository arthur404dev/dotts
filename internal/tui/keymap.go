package tui

import (
	"github.com/arthur404dev/dotts/internal/tui/keys"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	OpenPalette key.Binding
	Quit        key.Binding
	ForceQuit   key.Binding
	Escape      key.Binding
	Help        key.Binding

	ListUp   key.Binding
	ListDown key.Binding
	NavLeft  key.Binding
	NavRight key.Binding

	NextField key.Binding
	PrevField key.Binding
	Confirm   key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		OpenPalette: keys.OpenPalette,
		Quit:        keys.Quit,
		ForceQuit:   keys.ForceQuit,
		Escape:      keys.Escape,
		Help:        keys.Help,

		ListUp:   keys.ListUp,
		ListDown: keys.ListDown,
		NavLeft:  keys.NavLeft,
		NavRight: keys.NavRight,

		NextField: keys.NextField,
		PrevField: keys.PrevField,
		Confirm:   keys.Confirm,
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.OpenPalette, k.Help, k.ForceQuit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ListUp, k.ListDown, k.NavLeft, k.NavRight},
		{k.NextField, k.PrevField, k.Confirm},
		{k.OpenPalette, k.Help, k.Escape, k.Quit, k.ForceQuit},
	}
}

var _ help.KeyMap = KeyMap{}
