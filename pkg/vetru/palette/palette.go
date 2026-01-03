package palette

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const diag = "╱"

// Palette is the command palette component
type Palette struct {
	theme    *theme.Theme
	input    textinput.Model
	commands []Command
	filtered []Command
	cursor   int
	width    int
}

// New creates a new command palette
func New(t *theme.Theme, commands []Command) *Palette {
	ti := textinput.New()
	ti.Placeholder = "Type to filter"
	ti.Prompt = "> "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Primary)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.FgBase)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.FgSubtle)

	return &Palette{
		theme:    t,
		input:    ti,
		commands: commands,
		filtered: commands,
		width:    60,
	}
}

// SetWidth sets the palette width
func (p *Palette) SetWidth(w int) {
	// Constrain width
	if w > 70 {
		w = 70
	}
	if w < 40 {
		w = 40
	}
	p.width = w
	p.input.Width = w - 6
}

// Focus focuses the palette input
func (p *Palette) Focus() tea.Cmd {
	p.cursor = 0
	p.input.SetValue("")
	p.filtered = p.commands
	return p.input.Focus()
}

// Update handles events
func (p *Palette) Update(msg tea.Msg) (*Palette, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Execute selected command
			if len(p.filtered) > 0 && p.cursor < len(p.filtered) {
				cmd := p.filtered[p.cursor]
				return p, tea.Batch(
					messages.ClosePalette(),
					cmd.Handler(),
				)
			}
			return p, messages.ClosePalette()

		case "esc":
			return p, messages.ClosePalette()

		case "up", "ctrl+p":
			if len(p.filtered) > 0 {
				if p.cursor > 0 {
					p.cursor--
				} else {
					p.cursor = len(p.filtered) - 1
				}
			}
			return p, nil

		case "down", "ctrl+n":
			if len(p.filtered) > 0 {
				if p.cursor < len(p.filtered)-1 {
					p.cursor++
				} else {
					p.cursor = 0
				}
			}
			return p, nil

		case "tab":
			// Move down
			if p.cursor < len(p.filtered)-1 {
				p.cursor++
			} else {
				p.cursor = 0
			}
			return p, nil

		case "shift+tab":
			// Move up
			if p.cursor > 0 {
				p.cursor--
			} else {
				p.cursor = len(p.filtered) - 1
			}
			return p, nil

		default:
			// Update input and filter
			var cmd tea.Cmd
			p.input, cmd = p.input.Update(msg)
			p.filter()
			return p, cmd
		}
	}

	return p, nil
}

// filter updates the filtered commands based on input
func (p *Palette) filter() {
	query := strings.ToLower(strings.TrimSpace(p.input.Value()))
	if query == "" {
		p.filtered = p.commands
		p.cursor = 0
		return
	}

	p.filtered = []Command{}
	for _, cmd := range p.commands {
		// Match against label, category, and keywords
		label := strings.ToLower(cmd.Label)
		category := strings.ToLower(cmd.Category)

		if strings.Contains(label, query) || strings.Contains(category, query) {
			p.filtered = append(p.filtered, cmd)
			continue
		}

		// Check keywords
		for _, kw := range cmd.Keywords {
			if strings.Contains(strings.ToLower(kw), query) {
				p.filtered = append(p.filtered, cmd)
				break
			}
		}
	}

	// Reset cursor if out of bounds
	if p.cursor >= len(p.filtered) {
		p.cursor = 0
	}
}

// View renders the palette
func (p *Palette) View() string {
	t := p.theme
	contentWidth := p.width - 4 // Account for border padding

	// Title with gradient stripes (Crush-style)
	title := "Commands"
	titleStyle := lipgloss.NewStyle().Foreground(t.Primary)
	remainingWidth := contentWidth - lipgloss.Width(title) - 1
	stripes := ""
	if remainingWidth > 0 {
		stripes = theme.ApplyGradient(
			strings.Repeat(diag, remainingWidth),
			t.Primary, t.Secondary,
		)
	}
	header := titleStyle.Render(title) + " " + stripes

	// Input field
	inputView := p.input.View()

	// Command list
	var listItems []string
	currentCategory := ""

	for i, cmd := range p.filtered {
		// Category header
		if cmd.Category != currentCategory {
			currentCategory = cmd.Category
			if len(listItems) > 0 {
				listItems = append(listItems, "") // Spacing
			}
			catStyle := lipgloss.NewStyle().
				Foreground(t.FgSubtle).
				Bold(true)
			listItems = append(listItems, catStyle.Render(currentCategory))
		}

		// Command row
		var rowStyle lipgloss.Style
		if i == p.cursor {
			// Highlighted - full width background
			rowStyle = lipgloss.NewStyle().
				Background(t.Primary).
				Foreground(t.BgBase).
				Width(contentWidth)
		} else {
			rowStyle = lipgloss.NewStyle().
				Foreground(t.FgBase)
		}

		// Build row content
		label := cmd.Label
		shortcut := ""
		if cmd.Shortcut != "" {
			shortcut = cmd.Shortcut
		}

		if shortcut != "" {
			// Right-align shortcut
			shortcutStyle := lipgloss.NewStyle().Foreground(t.FgMuted)
			if i == p.cursor {
				shortcutStyle = shortcutStyle.Foreground(t.BgSubtle)
			}
			shortcutView := shortcutStyle.Render(shortcut)
			labelWidth := contentWidth - lipgloss.Width(shortcutView) - 1
			row := rowStyle.Render(
				lipgloss.NewStyle().Width(labelWidth).Render(label) + " " + shortcutView,
			)
			listItems = append(listItems, row)
		} else {
			listItems = append(listItems, rowStyle.Render(label))
		}
	}

	listView := lipgloss.JoinVertical(lipgloss.Left, listItems...)

	hints := lipgloss.NewStyle().Foreground(t.FgSubtle).Render(
		"[↑/↓] select " + theme.Icons.Dot + " [enter] confirm " + theme.Icons.Dot + " [esc] close",
	)

	// Compose dialog content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		inputView,
		"",
		listView,
		"",
		hints,
	)

	// Dialog container with border
	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus).
		Padding(1, 1).
		Width(p.width)

	return container.Render(content)
}
