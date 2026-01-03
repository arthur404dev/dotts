// Package components provides reusable TUI components for the dotts application.
package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// MenuItem represents a single item in the menu.
type MenuItem struct {
	// ID is the unique identifier for this menu item.
	ID string
	// Label is the display text for this menu item.
	Label string
	// Icon is an optional icon to display before the label.
	Icon string
	// Enabled determines if the item can be selected.
	Enabled bool
}

// Menu is a vertical navigation menu component with keyboard and mouse support.
// It supports item icons, disabled items, and focus states.
type Menu struct {
	theme   *theme.Theme
	items   []MenuItem
	cursor  int
	focused bool
	zoneID  string
	width   int
}

// NewMenu creates a new Menu component with the given theme and items.
// All items are enabled by default if not explicitly set.
func NewMenu(t *theme.Theme, items []MenuItem) *Menu {
	for i := range items {
		if !items[i].Enabled {
			items[i].Enabled = true
		}
	}
	return &Menu{
		theme:  t,
		items:  items,
		cursor: 0,
		zoneID: "menu",
		width:  0,
	}
}

// SetWidth sets the menu width. Use 0 for auto-width.
func (m *Menu) SetWidth(w int) *Menu {
	m.width = w
	return m
}

// SetCursor sets the cursor position by index.
// If the index is out of bounds, the cursor position is unchanged.
func (m *Menu) SetCursor(index int) *Menu {
	if index >= 0 && index < len(m.items) {
		m.cursor = index
	}
	return m
}

// SetCursorByID sets the cursor position by item ID.
// If no item with the given ID exists, the cursor position is unchanged.
func (m *Menu) SetCursorByID(id string) *Menu {
	for i, item := range m.items {
		if item.ID == id {
			m.cursor = i
			return m
		}
	}
	return m
}

// Cursor returns the current cursor position.
func (m *Menu) Cursor() int {
	return m.cursor
}

// Selected returns the currently selected menu item.
// Returns nil if the cursor is out of bounds.
func (m *Menu) Selected() *MenuItem {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return &m.items[m.cursor]
	}
	return nil
}

// SelectedID returns the ID of the currently selected item.
// Returns an empty string if no item is selected.
func (m *Menu) SelectedID() string {
	if item := m.Selected(); item != nil {
		return item.ID
	}
	return ""
}

// Focus sets the menu to focused state.
func (m *Menu) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// Blur removes focus from the menu.
func (m *Menu) Blur() {
	m.focused = false
}

// Focused returns whether the menu is currently focused.
func (m *Menu) Focused() bool {
	return m.focused
}

// Init implements tea.Model. Returns nil as no initialization is needed.
func (m *Menu) Init() tea.Cmd {
	return nil
}

// Update handles keyboard and mouse events.
// Keyboard: up/k (move up), down/j (move down), home (first item), end (last item).
// Mouse: click on an item to select it.
func (m *Menu) Update(msg tea.Msg) (*Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focused {
			switch msg.String() {
			case "up", "k", "ctrl+p":
				m.moveCursorUp()
			case "down", "j", "ctrl+n":
				m.moveCursorDown()
			case "home":
				m.cursor = 0
			case "end":
				m.cursor = len(m.items) - 1
			}
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease {
			for i := range m.items {
				itemZone := m.zoneID + "-" + m.items[i].ID
				if zone.Get(itemZone).InBounds(msg) {
					if m.items[i].Enabled {
						m.cursor = i
						m.focused = true
					}
					break
				}
			}
		}
	}
	return m, nil
}

func (m *Menu) moveCursorUp() {
	if len(m.items) == 0 {
		return
	}
	start := m.cursor
	for {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
		}
		if m.items[m.cursor].Enabled || m.cursor == start {
			return
		}
	}
}

func (m *Menu) moveCursorDown() {
	if len(m.items) == 0 {
		return
	}
	start := m.cursor
	for {
		m.cursor++
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		if m.items[m.cursor].Enabled || m.cursor == start {
			return
		}
	}
}

// View renders the menu as a vertical list of items.
// The selected item is highlighted with an arrow prefix.
// Disabled items are rendered in a muted style.
func (mn *Menu) View() string {
	t := mn.theme

	var rows []string

	for i, item := range mn.items {
		itemZone := mn.zoneID + "-" + item.ID

		var style lipgloss.Style
		prefix := "  "

		if !item.Enabled {
			style = lipgloss.NewStyle().Foreground(t.FgSubtle)
		} else if i == mn.cursor && mn.focused {
			style = lipgloss.NewStyle().
				Foreground(t.Primary).
				Bold(true)
			prefix = theme.Icons.ArrowRight + " "
		} else if i == mn.cursor {
			style = lipgloss.NewStyle().Foreground(t.FgBase)
			prefix = theme.Icons.Chevron + " "
		} else {
			style = lipgloss.NewStyle().Foreground(t.FgMuted)
		}

		label := item.Label
		if item.Icon != "" {
			label = item.Icon + " " + label
		}

		row := style.Render(prefix + label)
		if mn.width > 0 {
			row = lipgloss.NewStyle().Width(mn.width).Render(row)
		}

		rows = append(rows, zone.Mark(itemZone, row))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
