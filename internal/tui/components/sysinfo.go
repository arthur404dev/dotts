package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type SystemInfo struct {
	OS       string
	Hostname string
	Arch     string
	Distro   string
	Manager  string

	titleStyle lipgloss.Style
	labelStyle lipgloss.Style
	valueStyle lipgloss.Style
	width      int
}

func NewSystemInfo(titleStyle, labelStyle, valueStyle lipgloss.Style) *SystemInfo {
	return &SystemInfo{
		titleStyle: titleStyle,
		labelStyle: labelStyle,
		valueStyle: valueStyle,
		width:      20,
	}
}

func (s *SystemInfo) SetWidth(w int) {
	s.width = w
}

func (s *SystemInfo) SetInfo(os, hostname, arch, distro, manager string) {
	s.OS = os
	s.Hostname = hostname
	s.Arch = arch
	s.Distro = distro
	s.Manager = manager
}

func (s *SystemInfo) Render() string {
	title := s.titleStyle.Render("System")

	rows := []string{title, ""}

	if s.Distro != "" {
		rows = append(rows, s.valueStyle.Render(s.Distro))
	} else if s.OS != "" {
		rows = append(rows, s.valueStyle.Render(s.OS))
	}

	if s.Hostname != "" {
		rows = append(rows, s.labelStyle.Render(s.Hostname))
	}

	if s.Arch != "" {
		rows = append(rows, s.labelStyle.Render(s.Arch))
	}

	if s.Manager != "" {
		rows = append(rows, s.labelStyle.Render(fmt.Sprintf("pkg: %s", s.Manager)))
	}

	container := lipgloss.NewStyle().
		Width(s.width).
		Padding(1, 1)

	return container.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}
