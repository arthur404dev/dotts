// Package components provides reusable TUI building blocks.
package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// BreadcrumbItem represents a single item in a breadcrumb trail.
type BreadcrumbItem struct {
	ID    string // Unique identifier for the breadcrumb item
	Label string // Display label for the breadcrumb item
}

// Breadcrumb is a navigation trail component that shows the current
// location within a hierarchical structure.
type Breadcrumb struct {
	theme     *theme.Theme
	items     []BreadcrumbItem
	separator string
}

// NewBreadcrumb creates a new Breadcrumb component with the given theme and items.
// The default separator is " › ".
func NewBreadcrumb(t *theme.Theme, items []BreadcrumbItem) *Breadcrumb {
	return &Breadcrumb{
		theme:     t,
		items:     items,
		separator: " › ",
	}
}

// SetSeparator sets the separator string between breadcrumb items.
// Common alternatives include " / ", " > ", or " → ".
func (b *Breadcrumb) SetSeparator(sep string) *Breadcrumb {
	b.separator = sep
	return b
}

// SetItems replaces all breadcrumb items with the given items.
func (b *Breadcrumb) SetItems(items []BreadcrumbItem) *Breadcrumb {
	b.items = items
	return b
}

// Push adds a new item to the end of the breadcrumb trail.
func (b *Breadcrumb) Push(item BreadcrumbItem) *Breadcrumb {
	b.items = append(b.items, item)
	return b
}

// Pop removes and returns the last item from the breadcrumb trail.
// Returns nil if the breadcrumb is empty.
func (b *Breadcrumb) Pop() *BreadcrumbItem {
	if len(b.items) == 0 {
		return nil
	}
	item := b.items[len(b.items)-1]
	b.items = b.items[:len(b.items)-1]
	return &item
}

// Current returns the last (current) item in the breadcrumb trail.
// Returns nil if the breadcrumb is empty.
func (b *Breadcrumb) Current() *BreadcrumbItem {
	if len(b.items) == 0 {
		return nil
	}
	return &b.items[len(b.items)-1]
}

// Len returns the number of items in the breadcrumb trail.
func (b *Breadcrumb) Len() int {
	return len(b.items)
}

// View renders the breadcrumb as a horizontal trail.
// The last item (current location) is highlighted with the primary color.
// Previous items are rendered in muted color with separators between them.
func (b *Breadcrumb) View() string {
	if len(b.items) == 0 {
		return ""
	}

	t := b.theme
	sepStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)

	var parts []string

	for i, item := range b.items {
		var style lipgloss.Style

		// Last item is current (highlighted)
		if i == len(b.items)-1 {
			style = lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
		} else {
			style = lipgloss.NewStyle().Foreground(t.FgMuted)
		}

		parts = append(parts, style.Render(item.Label))

		// Add separator (except after last item)
		if i < len(b.items)-1 {
			parts = append(parts, sepStyle.Render(b.separator))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
