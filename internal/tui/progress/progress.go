package progress

import (
	"fmt"
	"strings"
	"sync"

	"github.com/arthur404dev/dotts/internal/tui/styles"
)

type Step struct {
	Name    string
	Status  StepStatus
	Message string
	Current int
	Total   int
}

type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepSuccess
	StepWarning
	StepFailed
	StepSkipped
)

type MultiProgress struct {
	steps []Step
	mu    sync.RWMutex
}

func New() *MultiProgress {
	return &MultiProgress{
		steps: make([]Step, 0),
	}
}

func (m *MultiProgress) AddStep(name string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.steps = append(m.steps, Step{
		Name:   name,
		Status: StepPending,
	})
	return len(m.steps) - 1
}

func (m *MultiProgress) SetStatus(index int, status StepStatus, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if index >= 0 && index < len(m.steps) {
		m.steps[index].Status = status
		m.steps[index].Message = message
	}
}

func (m *MultiProgress) SetProgress(index, current, total int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if index >= 0 && index < len(m.steps) {
		m.steps[index].Current = current
		m.steps[index].Total = total
	}
}

func (m *MultiProgress) Render() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lines []string
	for i, step := range m.steps {
		lines = append(lines, m.renderStep(i, step))
	}
	return strings.Join(lines, "\n")
}

func (m *MultiProgress) renderStep(index int, step Step) string {
	prefix := "├─"
	if index == len(m.steps)-1 {
		prefix = "└─"
	}

	icon := m.statusIcon(step.Status)
	name := step.Name

	var detail string
	switch step.Status {
	case StepRunning:
		if step.Total > 0 {
			bar := m.progressBar(step.Current, step.Total, 20)
			detail = fmt.Sprintf("%s %d/%d", bar, step.Current, step.Total)
		}
		if step.Message != "" {
			detail += " " + styles.Mute(step.Message)
		}
	case StepSuccess:
		if step.Message != "" {
			detail = styles.Mute(step.Message)
		}
	case StepFailed:
		if step.Message != "" {
			detail = styles.Err(step.Message)
		}
	case StepWarning:
		if step.Message != "" {
			detail = styles.Warn(step.Message)
		}
	case StepSkipped:
		detail = styles.Mute("skipped")
	}

	if detail != "" {
		return fmt.Sprintf("  %s [%s] %s  %s", prefix, icon, name, detail)
	}
	return fmt.Sprintf("  %s [%s] %s", prefix, icon, name)
}

func (m *MultiProgress) statusIcon(status StepStatus) string {
	switch status {
	case StepPending:
		return styles.Mute("○")
	case StepRunning:
		return styles.Info("●")
	case StepSuccess:
		return styles.Success("✓")
	case StepWarning:
		return styles.Warn("!")
	case StepFailed:
		return styles.Err("✗")
	case StepSkipped:
		return styles.Mute("-")
	default:
		return " "
	}
}

func (m *MultiProgress) progressBar(current, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}

	filled := (current * width) / total
	if filled > width {
		filled = width
	}

	return styles.Info(strings.Repeat("█", filled)) +
		styles.Mute(strings.Repeat("░", width-filled))
}

func PrintHeader(title string) {
	fmt.Println()
	fmt.Println(styles.Title(title))
	fmt.Println(strings.Repeat("─", 40))
}

func PrintSuccess(message string) {
	fmt.Println(styles.Success("✓ " + message))
}

func PrintError(message string) {
	fmt.Println(styles.Err("✗ " + message))
}

func PrintWarning(message string) {
	fmt.Println(styles.Warn("! " + message))
}

func PrintInfo(message string) {
	fmt.Println(styles.Info("→ " + message))
}

func PrintMuted(message string) {
	fmt.Println(styles.Mute("  " + message))
}
