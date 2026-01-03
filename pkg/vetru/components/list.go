package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem implements the list.Item interface for use with List.
type ListItem struct {
	id          string
	title       string
	description string
}

// NewListItem creates a new ListItem.
func NewListItem(id, title, description string) ListItem {
	return ListItem{
		id:          id,
		title:       title,
		description: description,
	}
}

// ID returns the item's unique identifier.
func (i ListItem) ID() string { return i.id }

// Title implements list.Item.
func (i ListItem) Title() string { return i.title }

// Description implements list.Item.
func (i ListItem) Description() string { return i.description }

// FilterValue implements list.Item.
func (i ListItem) FilterValue() string { return i.title }

// List is a themed wrapper around bubbles/list.
type List struct {
	theme  *theme.Theme
	list   list.Model
	title  string
	width  int
	height int
}

// NewList creates a new List component.
func NewList(t *theme.Theme, items []ListItem, width, height int) *List {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(t.Primary).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(t.FgMuted).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(t.Primary).
		Padding(0, 0, 0, 1)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(t.FgBase).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(t.FgSubtle).
		Padding(0, 0, 0, 2)
	delegate.Styles.DimmedTitle = lipgloss.NewStyle().
		Foreground(t.FgSubtle).
		Padding(0, 0, 0, 2)
	delegate.Styles.DimmedDesc = lipgloss.NewStyle().
		Foreground(t.FgSubtle).
		Padding(0, 0, 0, 2)

	l := list.New(listItems, delegate, width, height)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Padding(0, 0, 1, 0)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(t.Primary)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(t.Primary)

	return &List{
		theme:  t,
		list:   l,
		width:  width,
		height: height,
	}
}

// SetTitle sets the list title.
func (l *List) SetTitle(title string) *List {
	l.title = title
	l.list.Title = title
	return l
}

// SetItems replaces all items in the list.
func (l *List) SetItems(items []ListItem) *List {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}
	l.list.SetItems(listItems)
	return l
}

// SetSize sets the list dimensions.
func (l *List) SetSize(width, height int) *List {
	l.width = width
	l.height = height
	l.list.SetSize(width, height)
	return l
}

// SetFilteringEnabled enables or disables filtering.
func (l *List) SetFilteringEnabled(enabled bool) *List {
	l.list.SetFilteringEnabled(enabled)
	return l
}

// SetShowStatusBar shows or hides the status bar.
func (l *List) SetShowStatusBar(show bool) *List {
	l.list.SetShowStatusBar(show)
	return l
}

// SetShowHelp shows or hides the help text.
func (l *List) SetShowHelp(show bool) *List {
	l.list.SetShowHelp(show)
	return l
}

// Selected returns the currently selected item, or nil if none.
func (l *List) Selected() *ListItem {
	item, ok := l.list.SelectedItem().(ListItem)
	if !ok {
		return nil
	}
	return &item
}

// SelectedID returns the ID of the selected item, or empty string if none.
func (l *List) SelectedID() string {
	if item := l.Selected(); item != nil {
		return item.ID()
	}
	return ""
}

// Index returns the current cursor index.
func (l *List) Index() int {
	return l.list.Index()
}

// Init implements tea.Model.
func (l *List) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (l *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

// View implements tea.Model.
func (l *List) View() string {
	return l.list.View()
}

// Model returns the underlying list.Model for advanced usage.
func (l *List) Model() *list.Model {
	return &l.list
}
