package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectItem struct {
	ID          string
	Label       string
	Description string
}

type SelectList struct {
	items         []SelectItem
	cursor        int
	label         string
	labelStyle    lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
	descStyle     lipgloss.Style
	width         int
	maxVisible    int
	offset        int
}

func NewSelectList(label string, items []SelectItem) *SelectList {
	return &SelectList{
		items:         items,
		cursor:        0,
		label:         label,
		labelStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).MarginBottom(1),
		itemStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Padding(0, 2),
		selectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")).Bold(true).Padding(0, 2),
		descStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Padding(0, 4),
		width:         50,
		maxVisible:    8,
		offset:        0,
	}
}

func (s *SelectList) SetWidth(w int) {
	s.width = w
}

func (s *SelectList) SetMaxVisible(n int) {
	s.maxVisible = n
}

func (s *SelectList) Selected() *SelectItem {
	if s.cursor < len(s.items) {
		return &s.items[s.cursor]
	}
	return nil
}

func (s *SelectList) SelectedID() string {
	if item := s.Selected(); item != nil {
		return item.ID
	}
	return ""
}

func (s *SelectList) SetCursor(index int) {
	if index >= 0 && index < len(s.items) {
		s.cursor = index
		s.updateOffset()
	}
}

func (s *SelectList) SetCursorByID(id string) {
	for i, item := range s.items {
		if item.ID == id {
			s.SetCursor(i)
			return
		}
	}
}

func (s *SelectList) updateOffset() {
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if s.cursor >= s.offset+s.maxVisible {
		s.offset = s.cursor - s.maxVisible + 1
	}
}

func (s *SelectList) Update(msg tea.Msg) (*SelectList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
				s.updateOffset()
			}
		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
				s.updateOffset()
			}
		case "home":
			s.cursor = 0
			s.offset = 0
		case "end":
			s.cursor = len(s.items) - 1
			s.updateOffset()
		}
	}
	return s, nil
}

func (s *SelectList) View() string {
	var rows []string

	if s.label != "" {
		rows = append(rows, s.labelStyle.Render(s.label))
	}

	end := s.offset + s.maxVisible
	if end > len(s.items) {
		end = len(s.items)
	}

	for i := s.offset; i < end; i++ {
		item := s.items[i]

		var row string
		if i == s.cursor {
			row = s.selectedStyle.Render(fmt.Sprintf("→ %s", item.Label))
		} else {
			row = s.itemStyle.Render(fmt.Sprintf("  %s", item.Label))
		}
		rows = append(rows, row)

		if item.Description != "" && i == s.cursor {
			desc := s.descStyle.Render(item.Description)
			rows = append(rows, desc)
		}
	}

	if len(s.items) > s.maxVisible {
		scrollInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#585b70")).
			Padding(0, 2).
			Render(fmt.Sprintf("(%d/%d)", s.cursor+1, len(s.items)))
		rows = append(rows, scrollInfo)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
