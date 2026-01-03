package wizard

import (
	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/components/input"
	"github.com/arthur404dev/dotts/pkg/vetru/keys"
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Step int

const (
	StepSource Step = iota
	StepMachine
	StepPersonal
	StepSettings
	StepFeatures
	StepSummary
)

var stepLabels = []string{
	"Source",
	"Machine",
	"Personal",
	"Settings",
	"Features",
	"Summary",
}

type Wizard struct {
	theme  *theme.Theme
	width  int
	height int

	step       Step
	focusIndex int

	stepper      *components.Stepper
	sourceSelect *input.Select

	nameInput   *input.TextInput
	emailInput  *input.TextInput
	githubInput *input.TextInput
}

func New(t *theme.Theme) *Wizard {
	stepItems := make([]components.StepItem, len(stepLabels))
	for i, label := range stepLabels {
		stepItems[i] = components.StepItem{
			ID:    label,
			Label: label,
		}
	}
	stepper := components.NewStepper(t, stepItems)

	sourceItems := []input.SelectItem{
		{ID: "default", Label: "Use default configs (arthur404dev/dotfiles)", Description: "Great for getting started quickly"},
		{ID: "fork", Label: "Start fresh with a template", Description: "Create your own config repository"},
		{ID: "custom", Label: "Use my existing config repo", Description: "Connect to your own dotfiles"},
	}

	sourceSelect := input.NewSelect(t, "Where should dotts get your configurations from?", sourceItems)

	nameInput := input.New(t, "Your full name", "John Doe")
	nameInput.SetHelp("Used for git commits and config files")

	emailInput := input.New(t, "Your email", "you@example.com")
	emailInput.SetHelp("Used for git commits")

	githubInput := input.New(t, "GitHub username (optional)", "username")
	githubInput.SetHelp("Used for GitHub-related configs")

	return &Wizard{
		theme:        t,
		step:         StepSource,
		stepper:      stepper,
		sourceSelect: sourceSelect,
		nameInput:    nameInput,
		emailInput:   emailInput,
		githubInput:  githubInput,
	}
}

func (w *Wizard) ID() messages.PageID {
	return app.PageWizard
}

func (w *Wizard) Title() string {
	return "Init Wizard"
}

func (w *Wizard) SetSize(width, height int) {
	w.width = width
	w.height = height

	inputWidth := width - 10
	if inputWidth > 60 {
		inputWidth = 60
	}
	w.sourceSelect.SetWidth(inputWidth)
	w.nameInput.SetWidth(inputWidth)
	w.emailInput.SetWidth(inputWidth)
	w.githubInput.SetWidth(inputWidth)
}

func (w *Wizard) Focus() tea.Cmd {
	return w.focusCurrent()
}

func (w *Wizard) Blur() {
	w.blurAll()
}

func (w *Wizard) Init() tea.Cmd {
	return w.focusCurrent()
}

func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keys.MatchesNavLeft(msg):
			return w.prevStep()

		case keys.MatchesNavRight(msg):
			return w.nextStep()

		case keys.MatchesNextField(msg):
			return w.nextField()

		case keys.MatchesPrevField(msg):
			return w.prevField()

		case keys.MatchesConfirm(msg):
			if w.step == StepSummary {
				return w, messages.WizardComplete(true, nil)
			}
			if w.isLastField() {
				return w.nextStep()
			}
			return w.nextField()

		case keys.MatchesEscape(msg):
			if w.step > StepSource {
				return w.prevStep()
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease || msg.Action == tea.MouseActionMotion {
			return w.handleMouse(msg)
		}
	}

	return w.updateCurrentInputs(msg)
}

func (w *Wizard) handleMouse(msg tea.MouseMsg) (*Wizard, tea.Cmd) {
	switch w.step {
	case StepSource:
		var cmd tea.Cmd
		w.sourceSelect, cmd = w.sourceSelect.Update(msg)
		return w, cmd

	case StepPersonal:
		inputs := []*input.TextInput{w.nameInput, w.emailInput, w.githubInput}
		for i, inp := range inputs {
			if zone.Get(inp.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
				w.blurAll()
				w.focusIndex = i
				return w, w.focusCurrent()
			}
		}
		return w.updateCurrentInputs(msg)
	}

	return w, nil
}

func (w *Wizard) nextField() (*Wizard, tea.Cmd) {
	maxFields := w.maxFieldsForStep()
	if maxFields <= 1 {
		return w, nil
	}

	w.blurAll()
	w.focusIndex = (w.focusIndex + 1) % maxFields
	return w, w.focusCurrent()
}

func (w *Wizard) prevField() (*Wizard, tea.Cmd) {
	maxFields := w.maxFieldsForStep()
	if maxFields <= 1 {
		return w, nil
	}

	w.blurAll()
	w.focusIndex--
	if w.focusIndex < 0 {
		w.focusIndex = maxFields - 1
	}
	return w, w.focusCurrent()
}

func (w *Wizard) maxFieldsForStep() int {
	switch w.step {
	case StepSource:
		return 1
	case StepPersonal:
		return 3
	default:
		return 1
	}
}

func (w *Wizard) isLastField() bool {
	return w.focusIndex >= w.maxFieldsForStep()-1
}

func (w *Wizard) nextStep() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	switch w.step {
	case StepSource:
		w.step = StepPersonal
	case StepPersonal:
		w.step = StepSummary
	case StepSummary:
		return w, messages.WizardComplete(true, nil)
	default:
		w.step++
		if w.step > StepSummary {
			return w, messages.WizardComplete(true, nil)
		}
	}

	return w, w.focusCurrent()
}

func (w *Wizard) prevStep() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	switch w.step {
	case StepPersonal:
		w.step = StepSource
	case StepSummary:
		w.step = StepPersonal
	default:
		if w.step > StepSource {
			w.step--
		}
	}

	return w, w.focusCurrent()
}

func (w *Wizard) blurAll() {
	w.sourceSelect.Blur()
	w.nameInput.Blur()
	w.emailInput.Blur()
	w.githubInput.Blur()
}

func (w *Wizard) focusCurrent() tea.Cmd {
	switch w.step {
	case StepSource:
		return w.sourceSelect.Focus()
	case StepPersonal:
		switch w.focusIndex {
		case 0:
			return w.nameInput.Focus()
		case 1:
			return w.emailInput.Focus()
		case 2:
			return w.githubInput.Focus()
		}
	}
	return nil
}

func (w *Wizard) updateCurrentInputs(msg tea.Msg) (*Wizard, tea.Cmd) {
	var cmd tea.Cmd

	switch w.step {
	case StepSource:
		w.sourceSelect, cmd = w.sourceSelect.Update(msg)
	case StepPersonal:
		switch w.focusIndex {
		case 0:
			w.nameInput, cmd = w.nameInput.Update(msg)
		case 1:
			w.emailInput, cmd = w.emailInput.Update(msg)
		case 2:
			w.githubInput, cmd = w.githubInput.Update(msg)
		}
	}

	return w, cmd
}

func (w *Wizard) View() string {
	t := w.theme

	progress := w.renderProgress()

	var content string
	switch w.step {
	case StepSource:
		content = w.renderSourceStep()
	case StepPersonal:
		content = w.renderPersonalStep()
	case StepSummary:
		content = w.renderSummaryStep()
	default:
		content = t.S().Title.Render("Step: " + stepLabels[w.step])
	}

	hint := w.renderHint()

	container := lipgloss.NewStyle().
		Width(w.width).
		Height(w.height).
		Padding(1, 2)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		progress,
		"",
		content,
		"",
		hint,
	)

	return container.Render(body)
}

func (w *Wizard) renderProgress() string {
	w.stepper.SetCurrent(int(w.step))
	return w.stepper.View()
}

func (w *Wizard) renderSourceStep() string {
	t := w.theme

	title := t.S().Title.Render("Configuration Source")
	subtitle := t.S().Subtle.Render("Choose where to get your dotfiles configuration from.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		w.sourceSelect.View(),
	)
}

func (w *Wizard) renderPersonalStep() string {
	t := w.theme

	title := t.S().Title.Render("Personal Settings")
	subtitle := t.S().Subtle.Render("These settings personalize your dotfiles.")

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		w.nameInput.View(),
		"",
		w.emailInput.View(),
		"",
		w.githubInput.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		form,
	)
}

func (w *Wizard) renderSummaryStep() string {
	t := w.theme

	title := t.S().Title.Render("Summary")
	subtitle := t.S().Subtle.Render("Review your settings before completing setup.")

	source := w.sourceSelect.Selected()
	sourceLabel := "Unknown"
	if source != nil {
		sourceLabel = source.Label
	}

	summaryItems := []string{
		t.S().Muted.Render("  Source:  ") + t.S().Text.Render(sourceLabel),
		t.S().Muted.Render("  Name:    ") + t.S().Text.Render(w.nameInput.Value()),
		t.S().Muted.Render("  Email:   ") + t.S().Text.Render(w.emailInput.Value()),
		t.S().Muted.Render("  GitHub:  ") + t.S().Text.Render(w.githubInput.Value()),
	}

	summary := lipgloss.JoinVertical(lipgloss.Left, summaryItems...)

	box := t.S().BorderNormal.
		Width(w.width-8).
		Padding(1, 2).
		Render(summary)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		box,
		"",
		t.S().Success.Render(theme.Icons.ArrowRight+" Press Enter to complete setup"),
	)
}

func (w *Wizard) renderHint() string {
	t := w.theme

	hints := []string{}

	switch w.step {
	case StepSource:
		hints = []string{"[↑/↓] select", "[←/→] step", "[enter] continue"}
	case StepPersonal:
		hints = []string{"[tab] next field", "[←/→] step", "[enter] continue"}
	case StepSummary:
		hints = []string{"[enter] complete", "[←] back"}
	}

	if len(hints) == 0 {
		return ""
	}

	return t.S().Muted.Render(lipgloss.JoinHorizontal(lipgloss.Top, hints...))
}
