package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type StepStatus int

const (
	StepPending StepStatus = iota
	StepCurrent
	StepComplete
)

type Step struct {
	ID     string
	Label  string
	Status StepStatus
}

type Progress struct {
	steps         []Step
	currentIndex  int
	completeStyle lipgloss.Style
	currentStyle  lipgloss.Style
	pendingStyle  lipgloss.Style
	titleStyle    lipgloss.Style
	width         int
}

func NewProgress(steps []Step, completeStyle, currentStyle, pendingStyle, titleStyle lipgloss.Style) *Progress {
	return &Progress{
		steps:         steps,
		currentIndex:  0,
		completeStyle: completeStyle,
		currentStyle:  currentStyle,
		pendingStyle:  pendingStyle,
		titleStyle:    titleStyle,
		width:         20,
	}
}

func (p *Progress) SetWidth(w int) {
	p.width = w
}

func (p *Progress) SetCurrent(index int) {
	if index >= 0 && index < len(p.steps) {
		for i := range p.steps {
			if i < index {
				p.steps[i].Status = StepComplete
			} else if i == index {
				p.steps[i].Status = StepCurrent
			} else {
				p.steps[i].Status = StepPending
			}
		}
		p.currentIndex = index
	}
}

func (p *Progress) SetCurrentByID(id string) {
	for i, step := range p.steps {
		if step.ID == id {
			p.SetCurrent(i)
			return
		}
	}
}

func (p *Progress) Complete(id string) {
	for i, step := range p.steps {
		if step.ID == id {
			p.steps[i].Status = StepComplete
			return
		}
	}
}

func (p *Progress) Render() string {
	var items []string

	for i, step := range p.steps {
		var icon string
		var style lipgloss.Style

		switch step.Status {
		case StepComplete:
			icon = "✓"
			style = p.completeStyle
		case StepCurrent:
			icon = "●"
			style = p.currentStyle
		case StepPending:
			icon = "○"
			style = p.pendingStyle
		}

		item := style.Render(fmt.Sprintf("%s %s", icon, step.Label))
		items = append(items, item)

		if i < len(p.steps)-1 {
			connector := p.pendingStyle.Render(" → ")
			if step.Status == StepComplete {
				connector = p.completeStyle.Render(" → ")
			}
			items = append(items, connector)
		}
	}

	timeline := lipgloss.JoinHorizontal(lipgloss.Center, items...)

	container := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return container.Render(timeline)
}
