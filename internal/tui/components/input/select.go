package input

import (
	"fmt"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// SelectItem represents an item in a select list
type SelectItem struct {
	ID          string
	Label       string
	Description string
}

// Select is a select/dropdown component with hover states
type Select struct {
	theme      *theme.Theme
	zoneID     string
	label      string
	items      []SelectItem
	cursor     int
	hoverIndex int
	focused    bool
	width      int
	maxVisible int
	offset     int
}

// NewSelect creates a new select component
func NewSelect(t *theme.Theme, label string, items []SelectItem) *Select {
	return &Select{
		theme:      t,
		zoneID:     "select-" + label,
		label:      label,
		items:      items,
		cursor:     0,
		hoverIndex: -1,
		width:      50,
		maxVisible: 8,
	}
}

// SetWidth sets the select width
func (s *Select) SetWidth(w int) {
	s.width = w
}

// SetMaxVisible sets the maximum visible items
func (s *Select) SetMaxVisible(n int) {
	s.maxVisible = n
}

// Selected returns the currently selected item
func (s *Select) Selected() *SelectItem {
	if s.cursor >= 0 && s.cursor < len(s.items) {
		return &s.items[s.cursor]
	}
	return nil
}

// SelectedID returns the ID of the selected item
func (s *Select) SelectedID() string {
	if item := s.Selected(); item != nil {
		return item.ID
	}
	return ""
}

// SetCursor sets the cursor position
func (s *Select) SetCursor(index int) {
	if index >= 0 && index < len(s.items) {
		s.cursor = index
		s.updateOffset()
	}
}

// SetCursorByID sets the cursor by item ID
func (s *Select) SetCursorByID(id string) {
	for i, item := range s.items {
		if item.ID == id {
			s.SetCursor(i)
			return
		}
	}
}

// Focus focuses the select
func (s *Select) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur removes focus
func (s *Select) Blur() {
	s.focused = false
}

// Focused returns whether the select is focused
func (s *Select) Focused() bool {
	return s.focused
}

func (s *Select) updateOffset() {
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if s.cursor >= s.offset+s.maxVisible {
		s.offset = s.cursor - s.maxVisible + 1
	}
}

// Update handles events
func (s *Select) Update(msg tea.Msg) (*Select, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !s.focused {
			return s, nil
		}

		switch msg.String() {
		case "up", "k", "ctrl+p":
			if len(s.items) > 0 {
				if s.cursor > 0 {
					s.cursor--
				} else {
					s.cursor = len(s.items) - 1
				}
				s.updateOffset()
			}
		case "down", "j", "ctrl+n":
			if len(s.items) > 0 {
				if s.cursor < len(s.items)-1 {
					s.cursor++
				} else {
					s.cursor = 0
				}
				s.updateOffset()
			}
		case "home":
			s.cursor = 0
			s.offset = 0
		case "end":
			s.cursor = len(s.items) - 1
			s.updateOffset()
		}

	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionMotion:
			// Check hover on each item
			s.hoverIndex = -1
			for i := range s.items {
				itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)
				if zone.Get(itemZone).InBounds(msg) {
					s.hoverIndex = i
					break
				}
			}

		case tea.MouseActionRelease:
			// Click to select
			for i := range s.items {
				itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)
				if zone.Get(itemZone).InBounds(msg) {
					s.cursor = i
					s.updateOffset()
					break
				}
			}
		}
	}

	return s, nil
}

// View renders the select
func (s *Select) View() string {
	t := s.theme
	var rows []string

	// Label
	if s.label != "" {
		rows = append(rows, t.S().TextInput.Label.Render(s.label))
	}

	// Calculate visible range
	end := s.offset + s.maxVisible
	if end > len(s.items) {
		end = len(s.items)
	}

	// Render items
	for i := s.offset; i < end; i++ {
		item := s.items[i]
		itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)

		var row string
		var style lipgloss.Style
		prefix := "  "

		if i == s.cursor && s.focused {
			// Selected and focused
			style = lipgloss.NewStyle().
				Foreground(t.Primary).
				Bold(true)
			prefix = theme.Icons.ArrowRight + " "
		} else if i == s.hoverIndex {
			// Hovered
			style = lipgloss.NewStyle().
				Foreground(t.FgBase)
			prefix = theme.Icons.Chevron + " "
		} else {
			// Normal
			style = lipgloss.NewStyle().
				Foreground(t.FgMuted)
		}

		row = style.Render(fmt.Sprintf("%s%s", prefix, item.Label))
		row = zone.Mark(itemZone, row)
		rows = append(rows, row)

		// Show description for selected item
		if item.Description != "" && i == s.cursor {
			desc := t.S().Subtle.PaddingLeft(4).Render(item.Description)
			rows = append(rows, desc)
		}
	}

	// Scroll indicator
	if len(s.items) > s.maxVisible {
		scrollInfo := t.S().Subtle.
			PaddingLeft(2).
			Render(fmt.Sprintf("(%d/%d)", s.cursor+1, len(s.items)))
		rows = append(rows, scrollInfo)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return zone.Mark(s.zoneID, content)
}
