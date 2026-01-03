package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type NavItem struct {
	ID    string
	Label string
	Icon  string
}

type Nav struct {
	items       []NavItem
	activeIndex int
	hoverIndex  int
	activeStyle lipgloss.Style
	normalStyle lipgloss.Style
	hoverStyle  lipgloss.Style
	titleStyle  lipgloss.Style
	width       int
}

func NewNav(items []NavItem, activeStyle, normalStyle, hoverStyle, titleStyle lipgloss.Style) *Nav {
	return &Nav{
		items:       items,
		activeIndex: 0,
		hoverIndex:  -1,
		activeStyle: activeStyle,
		normalStyle: normalStyle,
		hoverStyle:  hoverStyle,
		titleStyle:  titleStyle,
		width:       20,
	}
}

func (n *Nav) SetWidth(w int) {
	n.width = w
}

func (n *Nav) SetActive(index int) {
	if index >= 0 && index < len(n.items) {
		n.activeIndex = index
	}
}

func (n *Nav) SetActiveByID(id string) {
	for i, item := range n.items {
		if item.ID == id {
			n.activeIndex = i
			return
		}
	}
}

func (n *Nav) ActiveID() string {
	if n.activeIndex < len(n.items) {
		return n.items[n.activeIndex].ID
	}
	return ""
}

func (n *Nav) Items() []NavItem {
	return n.items
}

func (n *Nav) Len() int {
	return len(n.items)
}

func (n *Nav) SetHover(index int) {
	n.hoverIndex = index
}

func (n *Nav) ClearHover() {
	n.hoverIndex = -1
}

func (n *Nav) Next() {
	n.activeIndex = (n.activeIndex + 1) % len(n.items)
}

func (n *Nav) Prev() {
	n.activeIndex--
	if n.activeIndex < 0 {
		n.activeIndex = len(n.items) - 1
	}
}

func (n *Nav) Render() string {
	var rows []string

	title := n.titleStyle.Render("Navigation")
	rows = append(rows, title)
	rows = append(rows, "")

	for i, item := range n.items {
		var row string
		var style lipgloss.Style
		prefix := "   "

		if i == n.activeIndex {
			style = n.activeStyle
			prefix = " → "
		} else if i == n.hoverIndex {
			style = n.hoverStyle
			prefix = " ▸ "
		} else {
			style = n.normalStyle
		}

		row = style.Render(fmt.Sprintf("%s%s", prefix, item.Label))
		rows = append(rows, zone.Mark(item.ID, row))
	}

	container := lipgloss.NewStyle().
		Width(n.width).
		Padding(1, 1)

	return container.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}
