// Package components provides reusable TUI building blocks for the dotts application.
package components

import (
	"fmt"
	"strings"

	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// Code is a code block display component.
// It renders source code with optional line numbers, a subtle background,
// and configurable width. Useful for displaying configuration files,
// shell commands, or any monospace content.
type Code struct {
	theme       *theme.Theme
	content     string
	language    string
	lineNumbers bool
	wrap        bool
	width       int
}

// NewCode creates a new Code component with the given theme and content.
// Default settings: no line numbers, no wrapping, auto width.
func NewCode(t *theme.Theme, content string) *Code {
	return &Code{
		theme:       t,
		content:     content,
		lineNumbers: false,
		wrap:        false,
	}
}

// SetContent sets the code content to display.
func (c *Code) SetContent(content string) *Code {
	c.content = content
	return c
}

// SetLanguage sets the language hint (for future syntax highlighting support).
func (c *Code) SetLanguage(lang string) *Code {
	c.language = lang
	return c
}

// SetLineNumbers enables or disables line number display.
func (c *Code) SetLineNumbers(show bool) *Code {
	c.lineNumbers = show
	return c
}

// SetWrap enables or disables line wrapping (reserved for future use).
func (c *Code) SetWrap(wrap bool) *Code {
	c.wrap = wrap
	return c
}

// SetWidth sets the code block width (0 for auto).
func (c *Code) SetWidth(w int) *Code {
	c.width = w
	return c
}

// View renders the code block with optional line numbers and styled background.
func (c *Code) View() string {
	if c.content == "" {
		return ""
	}

	t := c.theme
	lines := strings.Split(c.content, "\n")
	lineNumWidth := len(fmt.Sprintf("%d", len(lines)))

	var outputLines []string

	for i, line := range lines {
		var lineView string

		if c.lineNumbers {
			numStyle := lipgloss.NewStyle().
				Foreground(t.FgSubtle).
				Width(lineNumWidth).
				Align(lipgloss.Right)
			lineNum := numStyle.Render(fmt.Sprintf("%d", i+1))

			sepStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
			sep := sepStyle.Render(" │ ")

			codeStyle := lipgloss.NewStyle().Foreground(t.Tertiary)
			code := codeStyle.Render(line)

			lineView = lineNum + sep + code
		} else {
			codeStyle := lipgloss.NewStyle().Foreground(t.Tertiary)
			lineView = codeStyle.Render(line)
		}

		outputLines = append(outputLines, lineView)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, outputLines...)

	containerStyle := lipgloss.NewStyle().
		Background(t.BgSubtle).
		Padding(1, 2)

	if c.width > 0 {
		containerStyle = containerStyle.Width(c.width)
	}

	return containerStyle.Render(content)
}
