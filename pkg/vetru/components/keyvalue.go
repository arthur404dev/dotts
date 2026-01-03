package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// KVPair represents a single key-value pair for display.
type KVPair struct {
	Key   string
	Value string
}

// KeyValue is a component for displaying aligned key-value pairs.
type KeyValue struct {
	theme      *theme.Theme
	pairs      []KVPair
	separator  string
	labelWidth int
}

// NewKeyValue creates a new KeyValue component with the given theme.
func NewKeyValue(t *theme.Theme) *KeyValue {
	return &KeyValue{
		theme:     t,
		pairs:     []KVPair{},
		separator: ":",
	}
}

// Add adds a single key-value pair and returns the component for chaining.
func (k *KeyValue) Add(key, value string) *KeyValue {
	k.pairs = append(k.pairs, KVPair{Key: key, Value: value})
	return k
}

// AddPairs adds multiple key-value pairs and returns the component for chaining.
func (k *KeyValue) AddPairs(pairs ...KVPair) *KeyValue {
	k.pairs = append(k.pairs, pairs...)
	return k
}

// SetPairs replaces all existing pairs with the provided ones.
func (k *KeyValue) SetPairs(pairs []KVPair) *KeyValue {
	k.pairs = pairs
	return k
}

// Clear removes all key-value pairs from the component.
func (k *KeyValue) Clear() *KeyValue {
	k.pairs = []KVPair{}
	return k
}

// SetSeparator sets the separator string between key and value (default ":").
func (k *KeyValue) SetSeparator(sep string) *KeyValue {
	k.separator = sep
	return k
}

// SetLabelWidth sets a fixed width for labels (0 for auto-calculation).
func (k *KeyValue) SetLabelWidth(w int) *KeyValue {
	k.labelWidth = w
	return k
}

// View renders the key-value pairs as a vertically stacked, aligned list.
func (k *KeyValue) View() string {
	if len(k.pairs) == 0 {
		return ""
	}

	t := k.theme

	labelWidth := k.labelWidth
	if labelWidth == 0 {
		for _, pair := range k.pairs {
			if len(pair.Key) > labelWidth {
				labelWidth = len(pair.Key)
			}
		}
	}

	var rows []string

	for _, pair := range k.pairs {
		key := pair.Key
		if len(key) < labelWidth {
			key = key + strings.Repeat(" ", labelWidth-len(key))
		}

		keyStyle := lipgloss.NewStyle().Foreground(t.FgMuted)
		valStyle := lipgloss.NewStyle().Foreground(t.FgBase)
		sepStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)

		row := keyStyle.Render(key) + sepStyle.Render(k.separator) + " " + valStyle.Render(pair.Value)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
