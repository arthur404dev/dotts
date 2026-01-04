package components

import (
	"fmt"
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Scrollable struct {
	viewport viewport.Model
	theme    *theme.Theme
	width    int
	height   int
	content  string
	focused  bool

	showIndicator bool
}

func NewScrollable(t *theme.Theme) *Scrollable {
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle()

	return &Scrollable{
		viewport:      vp,
		theme:         t,
		showIndicator: true,
	}
}

func (s *Scrollable) SetSize(width, height int) *Scrollable {
	s.width = width
	s.height = height
	s.viewport.Width = width
	s.viewport.Height = height
	return s
}

func (s *Scrollable) SetContent(content string) *Scrollable {
	s.content = content
	s.viewport.SetContent(content)
	return s
}

func (s *Scrollable) SetShowIndicator(show bool) *Scrollable {
	s.showIndicator = show
	return s
}

func (s *Scrollable) Focus() tea.Cmd {
	s.focused = true
	return nil
}

func (s *Scrollable) Blur() {
	s.focused = false
}

func (s *Scrollable) Focused() bool {
	return s.focused
}

func (s *Scrollable) AtTop() bool {
	return s.viewport.AtTop()
}

func (s *Scrollable) AtBottom() bool {
	return s.viewport.AtBottom()
}

func (s *Scrollable) ScrollPercent() float64 {
	return s.viewport.ScrollPercent()
}

func (s *Scrollable) YOffset() int {
	return s.viewport.YOffset
}

func (s *Scrollable) TotalLines() int {
	return s.viewport.TotalLineCount()
}

func (s *Scrollable) VisibleLines() int {
	return s.viewport.VisibleLineCount()
}

func (s *Scrollable) NeedsScrolling() bool {
	return s.TotalLines() > s.VisibleLines()
}

func (s *Scrollable) ScrollUp(lines int) {
	s.viewport.LineUp(lines)
}

func (s *Scrollable) ScrollDown(lines int) {
	s.viewport.LineDown(lines)
}

func (s *Scrollable) GotoTop() {
	s.viewport.GotoTop()
}

func (s *Scrollable) GotoBottom() {
	s.viewport.GotoBottom()
}

func (s *Scrollable) Update(msg tea.Msg) (*Scrollable, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *Scrollable) View() string {
	if !s.NeedsScrolling() {
		return s.content
	}

	content := s.viewport.View()

	if !s.showIndicator {
		return content
	}

	return s.viewWithIndicator(content)
}

func (s *Scrollable) viewWithIndicator(content string) string {
	t := s.theme

	lines := strings.Split(content, "\n")
	indicatorWidth := 1

	contentWidth := s.width - indicatorWidth - 1

	scrollbarHeight := s.height
	if scrollbarHeight < 1 {
		scrollbarHeight = 1
	}

	thumbSize := max(1, int(float64(scrollbarHeight)*float64(s.VisibleLines())/float64(s.TotalLines())))
	thumbPos := int(float64(scrollbarHeight-thumbSize) * s.ScrollPercent())

	var scrollbar []string
	for i := 0; i < scrollbarHeight; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			scrollbar = append(scrollbar, lipgloss.NewStyle().Foreground(t.Primary).Render("┃"))
		} else {
			scrollbar = append(scrollbar, lipgloss.NewStyle().Foreground(t.Border).Render("│"))
		}
	}

	var result []string
	for i := 0; i < len(lines) && i < s.height; i++ {
		line := lines[i]
		lineWidth := lipgloss.Width(line)

		padding := contentWidth - lineWidth
		if padding > 0 {
			line += strings.Repeat(" ", padding)
		}

		scrollChar := " "
		if i < len(scrollbar) {
			scrollChar = scrollbar[i]
		}

		result = append(result, line+" "+scrollChar)
	}

	return strings.Join(result, "\n")
}

func (s *Scrollable) ViewWithStatus() string {
	content := s.View()

	if !s.NeedsScrolling() {
		return content
	}

	t := s.theme
	status := fmt.Sprintf(" %d/%d ", s.YOffset()+s.VisibleLines(), s.TotalLines())
	statusView := lipgloss.NewStyle().
		Foreground(t.FgMuted).
		Background(t.BgSubtle).
		Render(status)

	return lipgloss.JoinVertical(lipgloss.Right, content, statusView)
}
