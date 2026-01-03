// Package components provides reusable TUI building blocks.
package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// TabItem represents a single tab in a tab bar.
type TabItem struct {
	ID    string // Unique identifier for the tab
	Label string // Display label for the tab
}

// Tabs is a horizontal tab bar component with keyboard and mouse support.
// It provides navigation between different views or sections.
type Tabs struct {
	theme   *theme.Theme
	tabs    []TabItem
	active  int
	focused bool
	zoneID  string
}

// NewTabs creates a new Tabs component with the given theme and tab items.
func NewTabs(t *theme.Theme, tabs []TabItem) *Tabs {
	return &Tabs{
		theme:  t,
		tabs:   tabs,
		active: 0,
		zoneID: "tabs",
	}
}

// SetActive sets the active tab by index.
// If the index is out of bounds, it is ignored.
func (t *Tabs) SetActive(index int) *Tabs {
	if index >= 0 && index < len(t.tabs) {
		t.active = index
	}
	return t
}

// SetActiveByID sets the active tab by its ID.
// If no tab with the given ID exists, the active tab remains unchanged.
func (t *Tabs) SetActiveByID(id string) *Tabs {
	for i, tab := range t.tabs {
		if tab.ID == id {
			t.active = i
			return t
		}
	}
	return t
}

// Active returns the index of the currently active tab.
func (t *Tabs) Active() int {
	return t.active
}

// ActiveID returns the ID of the currently active tab.
// Returns an empty string if there are no tabs.
func (t *Tabs) ActiveID() string {
	if t.active >= 0 && t.active < len(t.tabs) {
		return t.tabs[t.active].ID
	}
	return ""
}

// ActiveTab returns a pointer to the currently active TabItem.
// Returns nil if there are no tabs.
func (t *Tabs) ActiveTab() *TabItem {
	if t.active >= 0 && t.active < len(t.tabs) {
		return &t.tabs[t.active]
	}
	return nil
}

// Next moves to the next tab, wrapping around to the first tab if at the end.
func (t *Tabs) Next() {
	if len(t.tabs) > 0 {
		t.active = (t.active + 1) % len(t.tabs)
	}
}

// Prev moves to the previous tab, wrapping around to the last tab if at the beginning.
func (t *Tabs) Prev() {
	if len(t.tabs) == 0 {
		return
	}
	t.active--
	if t.active < 0 {
		t.active = len(t.tabs) - 1
	}
}

// Focus sets the tabs as focused, enabling keyboard navigation.
func (t *Tabs) Focus() tea.Cmd {
	t.focused = true
	return nil
}

// Blur removes focus from the tabs.
func (t *Tabs) Blur() {
	t.focused = false
}

// Focused returns whether the tabs are currently focused.
func (t *Tabs) Focused() bool {
	return t.focused
}

// Init implements tea.Model and returns nil (no initialization command needed).
func (t *Tabs) Init() tea.Cmd {
	return nil
}

// Update handles keyboard and mouse events for tab navigation.
// Keyboard: left/h and right/l to navigate, tab/shift+tab to cycle.
// Mouse: click on a tab to select it.
func (t *Tabs) Update(msg tea.Msg) (*Tabs, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if t.focused {
			switch msg.String() {
			case "left", "h":
				t.Prev()
			case "right", "l":
				t.Next()
			}
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease {
			for i := range t.tabs {
				tabZone := t.zoneID + "-" + t.tabs[i].ID
				if zone.Get(tabZone).InBounds(msg) {
					t.active = i
					t.focused = true
					break
				}
			}
		}
	}
	return t, nil
}

// View renders the tabs as a horizontal tab bar.
// The active tab is highlighted with the primary color and a bottom border.
func (tb *Tabs) View() string {
	t := tb.theme

	var tabs []string

	for i, tab := range tb.tabs {
		tabZone := tb.zoneID + "-" + tab.ID

		var style lipgloss.Style
		if i == tb.active {
			// Active tab: primary color with bottom border indicator
			style = lipgloss.NewStyle().
				Foreground(t.Primary).
				Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(t.Primary).
				Padding(0, 2)
		} else {
			// Inactive tab: muted color without border
			style = lipgloss.NewStyle().
				Foreground(t.FgMuted).
				Padding(0, 2)
		}

		tabView := style.Render(tab.Label)
		tabs = append(tabs, zone.Mark(tabZone, tabView))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
}
