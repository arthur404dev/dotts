package components

import (
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

type Center struct {
	theme      *theme.Theme
	content    string
	width      int
	height     int
	horizontal bool
	vertical   bool
}

func NewCenter(t *theme.Theme) *Center {
	return &Center{
		theme:      t,
		horizontal: true,
		vertical:   true,
	}
}

func (c *Center) SetContent(content string) *Center {
	c.content = content
	return c
}

func (c *Center) SetWidth(w int) *Center {
	c.width = w
	return c
}

func (c *Center) SetHeight(h int) *Center {
	c.height = h
	return c
}

func (c *Center) SetSize(w, h int) *Center {
	c.width = w
	c.height = h
	return c
}

func (c *Center) SetHorizontal(enabled bool) *Center {
	c.horizontal = enabled
	return c
}

func (c *Center) SetVertical(enabled bool) *Center {
	c.vertical = enabled
	return c
}

func (c *Center) View() string {
	if c.content == "" {
		return ""
	}

	hPos := lipgloss.Left
	vPos := lipgloss.Top

	if c.horizontal {
		hPos = lipgloss.Center
	}
	if c.vertical {
		vPos = lipgloss.Center
	}

	return lipgloss.Place(
		c.width,
		c.height,
		hPos,
		vPos,
		c.content,
	)
}
