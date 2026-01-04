package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

type PageContainer struct {
	theme  *theme.Theme
	width  int
	height int

	header  string
	content string
	footer  string

	widthRatio  float64
	heightRatio float64

	hAlign lipgloss.Position
	vAlign lipgloss.Position

	contentAlign lipgloss.Position

	gap int
}

func NewPageContainer(t *theme.Theme) *PageContainer {
	return &PageContainer{
		theme:        t,
		widthRatio:   0.8,
		heightRatio:  0.8,
		hAlign:       lipgloss.Center,
		vAlign:       lipgloss.Center,
		contentAlign: lipgloss.Center,
		gap:          2,
	}
}

func (p *PageContainer) SetGap(gap int) *PageContainer {
	p.gap = gap
	return p
}

func (p *PageContainer) SetSize(width, height int) *PageContainer {
	p.width = width
	p.height = height
	return p
}

func (p *PageContainer) SetRatio(widthRatio, heightRatio float64) *PageContainer {
	p.widthRatio = widthRatio
	p.heightRatio = heightRatio
	return p
}

func (p *PageContainer) SetAlign(h, v lipgloss.Position) *PageContainer {
	p.hAlign = h
	p.vAlign = v
	return p
}

func (p *PageContainer) SetContentAlign(align lipgloss.Position) *PageContainer {
	p.contentAlign = align
	return p
}

func (p *PageContainer) SetHeader(header string) *PageContainer {
	p.header = header
	return p
}

func (p *PageContainer) SetContent(content string) *PageContainer {
	p.content = content
	return p
}

func (p *PageContainer) SetFooter(footer string) *PageContainer {
	p.footer = footer
	return p
}

func (p *PageContainer) InnerWidth() int {
	return int(float64(p.width) * p.widthRatio)
}

func (p *PageContainer) InnerHeight() int {
	return int(float64(p.height) * p.heightRatio)
}

func (p *PageContainer) ContentHeight() int {
	innerHeight := p.InnerHeight()

	headerHeight := 0
	footerHeight := 0
	gapCount := 0

	if p.header != "" {
		headerHeight = lipgloss.Height(p.header)
		gapCount++
	}
	if p.footer != "" {
		footerHeight = lipgloss.Height(p.footer)
		gapCount++
	}

	totalGapHeight := gapCount * p.gap
	contentHeight := innerHeight - headerHeight - footerHeight - totalGapHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	return contentHeight
}

func (p *PageContainer) View() string {
	innerWidth := p.InnerWidth()
	innerHeight := p.InnerHeight()

	headerHeight := 0
	footerHeight := 0
	gapCount := 0

	if p.header != "" {
		headerHeight = lipgloss.Height(p.header)
		gapCount++
	}
	if p.footer != "" {
		footerHeight = lipgloss.Height(p.footer)
		gapCount++
	}

	totalGapHeight := gapCount * p.gap
	contentHeight := innerHeight - headerHeight - footerHeight - totalGapHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	var parts []string

	if p.header != "" {
		headerView := lipgloss.NewStyle().
			Width(innerWidth).
			Align(lipgloss.Center).
			Render(p.header)
		parts = append(parts, headerView)
		parts = append(parts, strings.Repeat("\n", p.gap-1))
	}

	if p.content != "" {
		contentView := lipgloss.Place(
			innerWidth,
			contentHeight,
			p.contentAlign,
			lipgloss.Top,
			p.content,
		)
		parts = append(parts, contentView)
	} else {
		spacer := lipgloss.NewStyle().
			Width(innerWidth).
			Height(contentHeight).
			Render("")
		parts = append(parts, spacer)
	}

	if p.footer != "" {
		parts = append(parts, strings.Repeat("\n", p.gap-1))
		footerView := lipgloss.NewStyle().
			Width(innerWidth).
			Align(lipgloss.Center).
			Render(p.footer)
		parts = append(parts, footerView)
	}

	inner := lipgloss.JoinVertical(lipgloss.Center, parts...)

	return lipgloss.Place(
		p.width,
		p.height,
		p.hAlign,
		p.vAlign,
		inner,
	)
}
