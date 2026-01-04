package input

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type PathSuggestion struct {
	Path  string
	Name  string
	IsDir bool
}

type PathInput struct {
	model   textinput.Model
	theme   *theme.Theme
	zoneID  string
	label   string
	help    string
	focused bool
	hovered bool
	width   int

	suggestions     []PathSuggestion
	suggestionIdx   int
	showSuggestions bool
	columns         int
	rows            int
	dirsOnly        bool
}

func NewPathInput(t *theme.Theme, label, placeholder string) *PathInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 512
	ti.Width = 40

	zoneID := "path-input-" + strings.ReplaceAll(label, " ", "-")

	return &PathInput{
		model:         ti,
		theme:         t,
		zoneID:        zoneID,
		label:         label,
		width:         50,
		columns:       3,
		rows:          5,
		dirsOnly:      true,
		suggestionIdx: 0,
	}
}

func (p *PathInput) itemsPerPage() int {
	return p.columns * p.rows
}

func (p *PathInput) currentPage() int {
	if len(p.suggestions) == 0 {
		return 0
	}
	return p.suggestionIdx / p.itemsPerPage()
}

func (p *PathInput) totalPages() int {
	if len(p.suggestions) == 0 {
		return 0
	}
	return (len(p.suggestions) + p.itemsPerPage() - 1) / p.itemsPerPage()
}

func (p *PathInput) SetWidth(w int) *PathInput {
	p.width = w
	p.model.Width = w - 4
	return p
}

func (p *PathInput) SetLabel(label string) *PathInput {
	p.label = label
	return p
}

func (p *PathInput) SetHelp(help string) *PathInput {
	p.help = help
	return p
}

func (p *PathInput) SetPlaceholder(placeholder string) *PathInput {
	p.model.Placeholder = placeholder
	return p
}

func (p *PathInput) SetValue(value string) *PathInput {
	p.model.SetValue(value)
	p.updateSuggestions()
	return p
}

func (p *PathInput) SetDirsOnly(dirsOnly bool) *PathInput {
	p.dirsOnly = dirsOnly
	return p
}

func (p *PathInput) SetGridSize(cols, rows int) *PathInput {
	p.columns = cols
	p.rows = rows
	return p
}

func (p *PathInput) Value() string {
	return p.model.Value()
}

func (p *PathInput) ExpandedValue() string {
	return expandHome(p.model.Value())
}

func (p *PathInput) ZoneID() string {
	return p.zoneID
}

func (p *PathInput) Focus() tea.Cmd {
	p.focused = true
	p.updateSuggestions()
	return p.model.Focus()
}

func (p *PathInput) Blur() {
	p.focused = false
	p.showSuggestions = false
	p.model.Blur()
}

func (p *PathInput) Focused() bool {
	return p.focused
}

func (p *PathInput) Hovered() bool {
	return p.hovered
}

func (p *PathInput) PathExists() bool {
	expanded := p.ExpandedValue()
	_, err := os.Stat(expanded)
	return err == nil
}

func (p *PathInput) IsDirectory() bool {
	expanded := p.ExpandedValue()
	info, err := os.Stat(expanded)
	return err == nil && info.IsDir()
}

func (p *PathInput) HasSuggestions() bool {
	return p.showSuggestions && len(p.suggestions) > 0
}

func (p *PathInput) acceptSuggestion() {
	if len(p.suggestions) > 0 && p.suggestionIdx < len(p.suggestions) {
		p.model.SetValue(p.suggestions[p.suggestionIdx].Path)
		p.model.CursorEnd()
		p.updateSuggestions()
	}
}

func (p *PathInput) navNext() {
	if len(p.suggestions) == 0 {
		return
	}
	p.suggestionIdx++
	if p.suggestionIdx >= len(p.suggestions) {
		p.suggestionIdx = 0
	}
}

func (p *PathInput) navPrev() {
	if len(p.suggestions) == 0 {
		return
	}
	p.suggestionIdx--
	if p.suggestionIdx < 0 {
		p.suggestionIdx = len(p.suggestions) - 1
	}
}

func (p *PathInput) navDown() {
	p.navNext()
}

func (p *PathInput) navUp() {
	p.navPrev()
}

func (p *PathInput) navRight() {
	if len(p.suggestions) == 0 {
		return
	}

	page := p.currentPage()
	pageStart := page * p.itemsPerPage()
	pageEnd := pageStart + p.itemsPerPage()
	if pageEnd > len(p.suggestions) {
		pageEnd = len(p.suggestions)
	}
	itemsOnPage := pageEnd - pageStart

	localIdx := p.suggestionIdx - pageStart
	localIdx += p.rows

	if localIdx >= itemsOnPage {
		currentRow := localIdx % p.rows
		localIdx = currentRow
	}

	p.suggestionIdx = pageStart + localIdx
}

func (p *PathInput) navLeft() {
	if len(p.suggestions) == 0 {
		return
	}

	page := p.currentPage()
	pageStart := page * p.itemsPerPage()
	pageEnd := pageStart + p.itemsPerPage()
	if pageEnd > len(p.suggestions) {
		pageEnd = len(p.suggestions)
	}
	itemsOnPage := pageEnd - pageStart

	localIdx := p.suggestionIdx - pageStart
	localIdx -= p.rows

	if localIdx < 0 {
		currentRow := (p.suggestionIdx - pageStart) % p.rows
		numCols := (itemsOnPage + p.rows - 1) / p.rows
		lastColStart := (numCols - 1) * p.rows
		localIdx = lastColStart + currentRow
		if localIdx >= itemsOnPage {
			localIdx -= p.rows
		}
		if localIdx < 0 {
			localIdx = 0
		}
	}

	p.suggestionIdx = pageStart + localIdx
}

func (p *PathInput) Update(msg tea.Msg) (*PathInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !p.focused {
			return p, nil
		}

		switch msg.String() {
		case "tab", " ":
			if p.showSuggestions && len(p.suggestions) > 0 {
				p.acceptSuggestion()
				return p, nil
			}
			if msg.String() == " " {
				break
			}
			return p, nil

		case "ctrl+n", "down":
			if p.showSuggestions && len(p.suggestions) > 0 {
				p.navNext()
				return p, nil
			}

		case "ctrl+p", "up":
			if p.showSuggestions && len(p.suggestions) > 0 {
				p.navPrev()
				return p, nil
			}

		case "right":
			if p.showSuggestions && len(p.suggestions) > 0 {
				p.navRight()
				return p, nil
			}

		case "left":
			if p.showSuggestions && len(p.suggestions) > 0 {
				p.navLeft()
				return p, nil
			}

		case "esc":
			if p.showSuggestions {
				p.showSuggestions = false
				return p, nil
			}
		}

	case tea.MouseMsg:
		inBounds := zone.Get(p.zoneID).InBounds(msg)

		switch msg.Action {
		case tea.MouseActionMotion:
			p.hovered = inBounds

			if p.showSuggestions {
				page := p.currentPage()
				pageStart := page * p.itemsPerPage()
				pageEnd := pageStart + p.itemsPerPage()
				if pageEnd > len(p.suggestions) {
					pageEnd = len(p.suggestions)
				}

				for i := pageStart; i < pageEnd; i++ {
					sugZone := fmt.Sprintf("%s-sug-%d", p.zoneID, i)
					if zone.Get(sugZone).InBounds(msg) {
						p.suggestionIdx = i
						break
					}
				}
			}

		case tea.MouseActionRelease:
			if inBounds && !p.focused {
				p.focused = true
				return p, p.model.Focus()
			}

			if p.showSuggestions {
				page := p.currentPage()
				pageStart := page * p.itemsPerPage()
				pageEnd := pageStart + p.itemsPerPage()
				if pageEnd > len(p.suggestions) {
					pageEnd = len(p.suggestions)
				}

				for i := pageStart; i < pageEnd; i++ {
					sugZone := fmt.Sprintf("%s-sug-%d", p.zoneID, i)
					if zone.Get(sugZone).InBounds(msg) {
						p.acceptSuggestion()
						return p, nil
					}
				}
			}
		}
	}

	if p.focused {
		prevValue := p.model.Value()
		var cmd tea.Cmd
		p.model, cmd = p.model.Update(msg)

		if p.model.Value() != prevValue {
			p.updateSuggestions()
		}

		return p, cmd
	}

	return p, nil
}

func (p *PathInput) updateSuggestions() {
	value := p.model.Value()
	if value == "" {
		p.suggestions = nil
		p.showSuggestions = false
		p.suggestionIdx = 0
		return
	}

	expanded := expandHome(value)

	var dir, prefix string
	if value == "~" || strings.HasSuffix(value, "/") {
		dir = expanded
		if !strings.HasSuffix(dir, string(os.PathSeparator)) {
			dir += string(os.PathSeparator)
		}
		prefix = ""
	} else {
		dir = filepath.Dir(expanded)
		prefix = filepath.Base(expanded)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		p.suggestions = nil
		p.showSuggestions = false
		return
	}

	var suggestions []PathSuggestion
	for _, entry := range entries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") && !strings.HasPrefix(prefix, ".") {
			continue
		}

		if prefix != "" && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			continue
		}

		if p.dirsOnly && !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, name)
		displayPath := contractHome(fullPath)
		if entry.IsDir() {
			displayPath += string(os.PathSeparator)
		}

		suggestions = append(suggestions, PathSuggestion{
			Path:  displayPath,
			Name:  name,
			IsDir: entry.IsDir(),
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Name < suggestions[j].Name
	})

	p.suggestions = suggestions
	p.showSuggestions = len(suggestions) > 0
	if p.suggestionIdx >= len(suggestions) {
		p.suggestionIdx = 0
	}
}

func (p *PathInput) View() string {
	th := p.theme
	var rows []string

	if p.label != "" {
		labelView := th.S().TextInput.Label.Render(p.label)
		rows = append(rows, labelView)
	}

	var borderStyle lipgloss.Style
	switch {
	case p.focused:
		borderStyle = th.S().TextInput.FocusedB.Width(p.width)
	case p.hovered:
		borderStyle = th.S().TextInput.Hovered.Width(p.width)
	default:
		borderStyle = th.S().TextInput.Normal.Width(p.width)
	}

	fieldView := borderStyle.Render(p.model.View())
	rows = append(rows, fieldView)

	if p.showSuggestions && len(p.suggestions) > 0 && p.focused {
		dropdown := p.renderDropdown()
		rows = append(rows, dropdown)
	}

	if p.help != "" && !p.showSuggestions {
		helpView := th.S().TextInput.Help.Render(p.help)
		rows = append(rows, helpView)
	} else if p.showSuggestions {
		hint := p.renderHint()
		rows = append(rows, hint)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return zone.Mark(p.zoneID, content)
}

func (p *PathInput) renderHint() string {
	th := p.theme

	var pageInfo string
	if p.totalPages() > 1 {
		pageInfo = fmt.Sprintf(" • Page %d/%d", p.currentPage()+1, p.totalPages())
	}

	return th.S().Muted.Render(fmt.Sprintf("↑/↓ navigate • ←/→ columns • Tab accept%s", pageInfo))
}

func (p *PathInput) renderDropdown() string {
	th := p.theme

	colWidth := (p.width - 4) / p.columns
	if colWidth < 15 {
		colWidth = 15
	}

	page := p.currentPage()
	pageStart := page * p.itemsPerPage()
	pageEnd := pageStart + p.itemsPerPage()
	if pageEnd > len(p.suggestions) {
		pageEnd = len(p.suggestions)
	}

	itemsOnPage := pageEnd - pageStart
	numCols := (itemsOnPage + p.rows - 1) / p.rows
	if numCols > p.columns {
		numCols = p.columns
	}

	var gridRows []string
	for row := 0; row < p.rows; row++ {
		var cols []string
		for col := 0; col < numCols; col++ {
			localIdx := col*p.rows + row
			globalIdx := pageStart + localIdx

			if globalIdx >= pageEnd {
				cols = append(cols, lipgloss.NewStyle().Width(colWidth).Render(""))
				continue
			}

			sug := p.suggestions[globalIdx]
			sugZone := fmt.Sprintf("%s-sug-%d", p.zoneID, globalIdx)

			icon := theme.Icons.File
			if sug.IsDir {
				icon = theme.Icons.Folder
			}

			name := sug.Name
			maxNameLen := colWidth - 5
			if len(name) > maxNameLen && maxNameLen > 3 {
				name = name[:maxNameLen-1] + "…"
			}

			var style lipgloss.Style
			if globalIdx == p.suggestionIdx {
				style = lipgloss.NewStyle().
					Width(colWidth).
					Foreground(th.BgBase).
					Background(th.Primary).
					Bold(true)
			} else {
				style = lipgloss.NewStyle().
					Width(colWidth).
					Foreground(th.FgBase)
			}

			cell := style.Render(" " + icon + " " + name)
			cols = append(cols, zone.Mark(sugZone, cell))
		}
		gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, cols...))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, gridRows...)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Border).
		BorderTop(false).
		Padding(0, 1).
		Render(content)
}

func (p *PathInput) Blink() tea.Msg {
	return textinput.Blink()
}

func expandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		if path == "~" {
			return home
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func contractHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+string(os.PathSeparator)) {
		return "~" + path[len(home):]
	}
	return path
}
