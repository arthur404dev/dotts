// Package components provides reusable UI components for the TUI.
package components

import (
	"fmt"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// SegmentItem represents an option in a segmented control
type SegmentItem struct {
	ID    string
	Label string
}

// SegmentedControl is a pill-style selector for choosing between options.
// It displays all options horizontally in a bordered container.
type SegmentedControl struct {
	theme   *theme.Theme
	zoneID  string
	items   []SegmentItem
	cursor  int
	focused bool
	hovered int
	width   int
}

// NewSegmentedControl creates a new segmented control with the given items.
func NewSegmentedControl(t *theme.Theme, items []SegmentItem) *SegmentedControl {
	id := "segmented"
	if len(items) > 0 {
		id = fmt.Sprintf("segmented-%s", items[0].ID)
	}

	return &SegmentedControl{
		theme:   t,
		zoneID:  id,
		items:   items,
		cursor:  0,
		hovered: -1,
		width:   0, // auto
	}
}

// SetWidth sets a fixed width for the control (0 = auto)
func (s *SegmentedControl) SetWidth(w int) *SegmentedControl {
	s.width = w
	return s
}

// Selected returns the currently selected item
func (s *SegmentedControl) Selected() *SegmentItem {
	if s.cursor >= 0 && s.cursor < len(s.items) {
		return &s.items[s.cursor]
	}
	return nil
}

// SelectedID returns the ID of the selected item
func (s *SegmentedControl) SelectedID() string {
	if item := s.Selected(); item != nil {
		return item.ID
	}
	return ""
}

// SelectedIndex returns the index of the selected item
func (s *SegmentedControl) SelectedIndex() int {
	return s.cursor
}

// SetCursor sets the cursor position
func (s *SegmentedControl) SetCursor(index int) *SegmentedControl {
	if index >= 0 && index < len(s.items) {
		s.cursor = index
	}
	return s
}

// SetCursorByID sets the cursor by item ID
func (s *SegmentedControl) SetCursorByID(id string) *SegmentedControl {
	for i, item := range s.items {
		if item.ID == id {
			s.cursor = i
			return s
		}
	}
	return s
}

// Focus focuses the control
func (s *SegmentedControl) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur removes focus
func (s *SegmentedControl) Blur() {
	s.focused = false
}

// Focused returns whether the control is focused
func (s *SegmentedControl) Focused() bool {
	return s.focused
}

// ZoneID returns the zone ID for mouse tracking
func (s *SegmentedControl) ZoneID() string {
	return s.zoneID
}

// Update handles keyboard and mouse events
func (s *SegmentedControl) Update(msg tea.Msg) (*SegmentedControl, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !s.focused {
			return s, nil
		}

		switch msg.String() {
		case "left", "h":
			if s.cursor > 0 {
				s.cursor--
			}
		case "right", "l", " ":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			} else if msg.String() == " " {
				s.cursor = 0
			}
		case "home":
			s.cursor = 0
		case "end":
			s.cursor = len(s.items) - 1
		}

	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionMotion:
			s.hovered = -1
			for i := range s.items {
				itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)
				if zone.Get(itemZone).InBounds(msg) {
					s.hovered = i
					break
				}
			}

		case tea.MouseActionRelease:
			for i := range s.items {
				itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)
				if zone.Get(itemZone).InBounds(msg) {
					s.cursor = i
					s.focused = true
					break
				}
			}
		}
	}

	return s, nil
}

// View renders the segmented control
func (s *SegmentedControl) View() string {
	t := s.theme

	if len(s.items) == 0 {
		return ""
	}

	// Calculate segment width
	maxLabelLen := 0
	for _, item := range s.items {
		if len(item.Label) > maxLabelLen {
			maxLabelLen = len(item.Label)
		}
	}
	segmentWidth := maxLabelLen + 4 // padding

	// Build segments
	var segments []string
	for i, item := range s.items {
		itemZone := fmt.Sprintf("%s-item-%d", s.zoneID, i)

		var style lipgloss.Style
		label := item.Label

		if i == s.cursor {
			// Selected segment - filled background
			style = lipgloss.NewStyle().
				Width(segmentWidth).
				Align(lipgloss.Center).
				Foreground(t.BgBase).
				Background(t.Primary).
				Bold(true)
			label = "● " + label
		} else if i == s.hovered {
			// Hovered segment
			style = lipgloss.NewStyle().
				Width(segmentWidth).
				Align(lipgloss.Center).
				Foreground(t.Primary)
		} else {
			// Inactive segment
			style = lipgloss.NewStyle().
				Width(segmentWidth).
				Align(lipgloss.Center).
				Foreground(t.FgMuted)
		}

		segment := style.Render(label)
		segments = append(segments, zone.Mark(itemZone, segment))
	}

	// Join segments horizontally with separator
	separator := lipgloss.NewStyle().
		Foreground(t.Border).
		Render("│")

	inner := ""
	for i, seg := range segments {
		inner += seg
		if i < len(segments)-1 {
			inner += separator
		}
	}

	// Wrap in border
	borderColor := t.Border
	if s.focused {
		borderColor = t.Primary
	}

	wrapper := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	if s.width > 0 {
		wrapper = wrapper.Width(s.width)
	}

	content := wrapper.Render(inner)
	return zone.Mark(s.zoneID, content)
}
