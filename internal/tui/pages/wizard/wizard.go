package wizard

import (
	"os"

	"github.com/arthur404dev/dotts/internal/setup"
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/components/input"
	"github.com/arthur404dev/dotts/pkg/vetru/keys"
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
	StepSetup
	StepComplete
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

	stepper *components.Stepper

	sourceStep *SourceStep

	hostnameInput    *input.TextInput
	descriptionInput *input.TextInput
	profileSelect    *input.Select

	nameInput   *input.TextInput
	emailInput  *input.TextInput
	githubInput *input.TextInput

	shellSelect  *input.Select
	editorSelect *input.Select
	themeSelect  *input.Select

	featureSSH    *components.Checkbox
	featureGPG    *components.Checkbox
	featureGitHub *components.Checkbox
	featureDocker *components.Checkbox
	featureAsdf   *components.Checkbox
	featureGUI    *components.Checkbox

	prevButton  *components.Button
	nextButton  *components.Button
	prevHovered bool
	nextHovered bool
	prevFocused bool
	nextFocused bool

	scrollable *components.Scrollable

	setupInProgress bool
	setupResult     *messages.SetupCompleteMsg
	setupError      error
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

	sourceStep := NewSourceStep(t)

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "mymachine"
	}
	hostnameInput := input.New(t, "Hostname", hostname)
	hostnameInput.SetHelp("Identifies this machine in your configs")

	descriptionInput := input.New(t, "Description", "My workstation")
	descriptionInput.SetHelp("Brief description of this machine")

	profileItems := []input.SelectItem{
		{ID: "desktop", Label: "Desktop", Description: "Multi-monitor workstation"},
		{ID: "notebook", Label: "Notebook", Description: "Laptop with battery optimization"},
		{ID: "server", Label: "Server", Description: "Headless, minimal setup"},
	}
	profileSelect := input.NewSelect(t, "Base profile for this machine:", profileItems)

	nameInput := input.New(t, "Your full name", "John Doe")
	nameInput.SetHelp("Used for git commits and config files")

	emailInput := input.New(t, "Your email", "you@example.com")
	emailInput.SetHelp("Used for git commits")

	githubInput := input.New(t, "GitHub username (optional)", "username")
	githubInput.SetHelp("Used for GitHub-related configs")

	shellItems := []input.SelectItem{
		{ID: "fish", Label: "fish", Description: "Friendly interactive shell"},
		{ID: "zsh", Label: "zsh", Description: "Z shell with plugins"},
		{ID: "bash", Label: "bash", Description: "Classic Bourne-again shell"},
	}
	shellSelect := input.NewSelect(t, "Preferred shell:", shellItems)

	editorItems := []input.SelectItem{
		{ID: "nvim", Label: "Neovim", Description: "Modern Vim fork"},
		{ID: "vim", Label: "Vim", Description: "Classic modal editor"},
		{ID: "code", Label: "VS Code", Description: "Visual Studio Code"},
		{ID: "emacs", Label: "Emacs", Description: "Extensible text editor"},
	}
	editorSelect := input.NewSelect(t, "Preferred editor:", editorItems)

	themeItems := []input.SelectItem{
		{ID: "catppuccin", Label: "Catppuccin", Description: "Soothing pastel theme"},
		{ID: "nord", Label: "Nord", Description: "Arctic, north-bluish palette"},
		{ID: "dracula", Label: "Dracula", Description: "Dark theme for vampires"},
		{ID: "gruvbox", Label: "Gruvbox", Description: "Retro groove colors"},
	}
	themeSelect := input.NewSelect(t, "Color theme:", themeItems)

	prevButton := components.NewButton(t, "Previous").
		SetIcon(theme.Icons.ArrowRight).
		SetVariant(components.ButtonSecondary)
	nextButton := components.NewButton(t, "Next").
		SetIcon(theme.Icons.ArrowRight).
		SetIconRight().
		SetVariant(components.ButtonPrimary)

	scrollable := components.NewScrollable(t)

	return &Wizard{
		theme:            t,
		step:             StepSource,
		stepper:          stepper,
		sourceStep:       sourceStep,
		hostnameInput:    hostnameInput,
		descriptionInput: descriptionInput,
		profileSelect:    profileSelect,
		nameInput:        nameInput,
		emailInput:       emailInput,
		githubInput:      githubInput,
		shellSelect:      shellSelect,
		editorSelect:     editorSelect,
		themeSelect:      themeSelect,
		featureSSH:       components.NewCheckbox(t, "SSH - Generate SSH keys").SetChecked(true),
		featureGPG:       components.NewCheckbox(t, "GPG - Setup GPG signing"),
		featureGitHub:    components.NewCheckbox(t, "GitHub - Authenticate GitHub CLI").SetChecked(true),
		featureDocker:    components.NewCheckbox(t, "Docker - Container tooling"),
		featureAsdf:      components.NewCheckbox(t, "asdf - Version manager"),
		featureGUI:       components.NewCheckbox(t, "GUI - Desktop applications"),
		prevButton:       prevButton,
		nextButton:       nextButton,
		scrollable:       scrollable,
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
	w.sourceStep.SetSize(width, height)
	w.hostnameInput.SetWidth(inputWidth)
	w.descriptionInput.SetWidth(inputWidth)
	w.profileSelect.SetWidth(inputWidth)
	w.nameInput.SetWidth(inputWidth)
	w.emailInput.SetWidth(inputWidth)
	w.githubInput.SetWidth(inputWidth)
	w.shellSelect.SetWidth(inputWidth)
	w.editorSelect.SetWidth(inputWidth)
	w.themeSelect.SetWidth(inputWidth)
}

func (w *Wizard) Focus() tea.Cmd {
	return w.focusCurrent()
}

func (w *Wizard) Blur() {
	w.blurAll()
}

func (w *Wizard) CapturesKey(msg tea.KeyMsg) bool {
	if w.step == StepSource {
		return w.sourceStep.CapturesKey(msg)
	}
	return false
}

func (w *Wizard) Init() tea.Cmd {
	return w.focusCurrent()
}

func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SetupCompleteMsg:
		w.setupInProgress = false
		w.setupResult = &msg
		if msg.Success {
			w.step = StepComplete
			w.setupError = nil
		} else {
			w.setupError = msg.Error
			w.step = StepSummary
			w.focusIndex = 1
		}
		return w, w.focusCurrent()

	case tea.KeyMsg:
		if w.step == StepSetup {
			return w, nil
		}

		if w.step == StepComplete {
			if keys.MatchesConfirm(msg) {
				return w, messages.WizardComplete(true, nil)
			}
			return w, nil
		}

		if w.step == StepSource {
			return w.handleSourceStepKey(msg)
		}

		if w.stepUsesScrollable() {
			if handled, cmd := w.handleScrollKey(msg); handled {
				return w, cmd
			}
		}

		switch {
		case keys.MatchesNextField(msg):
			return w.nextField()

		case keys.MatchesPrevField(msg):
			return w.prevField()

		case keys.MatchesConfirm(msg):
			return w.handleEnter()

		case keys.MatchesEscape(msg):
			return w.handleEscape()
		}

	case tea.MouseMsg:
		if w.step == StepSetup || w.step == StepComplete {
			return w, nil
		}
		return w.handleMouse(msg)
	}

	return w.updateCurrentInputs(msg)
}

func (w *Wizard) stepUsesScrollable() bool {
	switch w.step {
	case StepMachine, StepPersonal, StepSettings, StepFeatures, StepSummary:
		return true
	}
	return false
}

func (w *Wizard) handleScrollKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch {
	case keys.MatchesScrollUp(msg):
		w.scrollable.ScrollUp(5)
		return true, nil
	case keys.MatchesScrollDown(msg):
		w.scrollable.ScrollDown(5)
		return true, nil
	case keys.MatchesScrollTop(msg):
		w.scrollable.GotoTop()
		return true, nil
	case keys.MatchesScrollBottom(msg):
		w.scrollable.GotoBottom()
		return true, nil
	}
	return false, nil
}

func (w *Wizard) handleEnter() (*Wizard, tea.Cmd) {
	if w.nextFocused {
		return w.navigateForward()
	}
	if w.prevFocused {
		return w.navigateBack()
	}

	if w.step == StepSummary {
		return w.startSetup()
	}

	if w.step == StepSource {
		var cmd tea.Cmd
		w.sourceStep, cmd = w.sourceStep.Update(tea.KeyMsg{Type: tea.KeyEnter})
		return w, cmd
	}

	w.updateCurrentInputs(tea.KeyMsg{Type: tea.KeySpace})

	if w.isLastField() {
		w.blurAll()
		w.nextFocused = true
		return w, nil
	}
	return w.nextField()
}

func (w *Wizard) handleEscape() (*Wizard, tea.Cmd) {
	if w.nextFocused {
		w.nextFocused = false
		w.focusIndex = w.maxFieldsForStep() - 1
		if w.focusIndex < 0 {
			w.focusIndex = 0
		}
		return w, w.focusCurrent()
	}
	if w.prevFocused {
		w.prevFocused = false
		w.focusIndex = 0
		return w, w.focusCurrent()
	}

	if w.step == StepSource {
		var cmd tea.Cmd
		w.sourceStep, cmd = w.sourceStep.Update(tea.KeyMsg{Type: tea.KeyEscape})
		return w, cmd
	}

	if w.focusIndex == 0 && w.canNavigateBack() {
		w.blurAll()
		w.prevFocused = true
		return w, nil
	}

	return w.prevField()
}

func (w *Wizard) handleSourceStepKey(msg tea.KeyMsg) (*Wizard, tea.Cmd) {
	prevComplete := w.sourceStep.IsComplete()

	var cmd tea.Cmd
	w.sourceStep, cmd = w.sourceStep.Update(msg)

	if !prevComplete && w.sourceStep.IsComplete() {
		w.blurAll()
		w.step = StepMachine
		w.focusIndex = 0
		return w, w.focusCurrent()
	}

	return w, cmd
}

func (w *Wizard) startSetup() (*Wizard, tea.Cmd) {
	w.step = StepSetup
	w.setupInProgress = true
	w.blurAll()

	sourceConfig := w.sourceStep.Config()

	return w, func() tea.Msg {
		s := setup.New()

		cfg := setup.SourceConfig{
			CreateSymlink: sourceConfig.CreateSymlink,
			SymlinkPath:   sourceConfig.SymlinkPath,
		}

		switch sourceConfig.Mode {
		case SourceModeStarter:
			cfg.Type = setup.SourceTypeStarter
		case SourceModeBYOP:
			cfg.Type = setup.SourceTypeBYOP
			if sourceConfig.BYOPType == BYOPTypeRemote {
				cfg.BYOPSubType = setup.BYOPSubTypeRemote
				cfg.RemoteURL = sourceConfig.RemoteURL
			} else {
				cfg.BYOPSubType = setup.BYOPSubTypeLocal
				cfg.LocalPath = sourceConfig.LocalPath
			}
		case SourceModeForkUser:
			cfg.Type = setup.SourceTypeForkUser
			cfg.ForkFrom = sourceConfig.ForkFrom
		}

		result, err := s.Execute(cfg)
		if err != nil {
			return messages.SetupCompleteMsg{
				Success: false,
				Error:   err,
			}
		}

		return messages.SetupCompleteMsg{
			Success:     true,
			ConfigPath:  result.ConfigPath,
			SymlinkPath: result.SymlinkPath,
			BackupPath:  result.BackupPath,
			ClonedFrom:  result.ClonedFrom,
			CopiedFrom:  result.CopiedFrom,
		}
	}
}

func (w *Wizard) handleMouse(msg tea.MouseMsg) (*Wizard, tea.Cmd) {
	if w.stepUsesScrollable() && w.scrollable.NeedsScrolling() {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			w.scrollable.ScrollUp(3)
			return w, nil
		case tea.MouseButtonWheelDown:
			w.scrollable.ScrollDown(3)
			return w, nil
		}
	}

	w.prevHovered = zone.Get(w.prevButton.ZoneID()).InBounds(msg) && w.canNavigateBack()
	w.nextHovered = zone.Get(w.nextButton.ZoneID()).InBounds(msg) && w.canNavigateForward()

	if clicked, stepIndex := w.stepper.HandleClick(msg); clicked {
		return w.goToStep(Step(stepIndex))
	}

	if msg.Action == tea.MouseActionRelease {
		if w.prevHovered {
			return w.navigateBack()
		}
		if w.nextHovered {
			return w.navigateForward()
		}
	}

	switch w.step {
	case StepSource:
		var cmd tea.Cmd
		w.sourceStep, cmd = w.sourceStep.Update(msg)
		return w, cmd

	case StepMachine:
		inputs := []*input.TextInput{w.hostnameInput, w.descriptionInput}
		for i, inp := range inputs {
			if zone.Get(inp.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
				w.blurAll()
				w.focusIndex = i
				return w, w.focusCurrent()
			}
		}

	case StepPersonal:
		inputs := []*input.TextInput{w.nameInput, w.emailInput, w.githubInput}
		for i, inp := range inputs {
			if zone.Get(inp.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
				w.blurAll()
				w.focusIndex = i
				return w, w.focusCurrent()
			}
		}
	}

	return w.updateCurrentInputs(msg)
}

func (w *Wizard) navigateBack() (*Wizard, tea.Cmd) {
	if w.step == StepSource {
		cmd := w.sourceStep.GoBack()
		return w, cmd
	}
	return w.prevStep()
}

func (w *Wizard) navigateForward() (*Wizard, tea.Cmd) {
	if w.step == StepSource {
		if w.sourceStep.IsOnLastSubStep() {
			w.blurAll()
			w.step = StepMachine
			w.focusIndex = 0
			return w, w.focusCurrent()
		}
		cmd := w.sourceStep.GoNext()
		return w, cmd
	}
	if w.step == StepSummary {
		return w.startSetup()
	}
	return w.nextStep()
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
	case StepMachine:
		return 3
	case StepPersonal:
		return 3
	case StepSettings:
		return 3
	case StepFeatures:
		return 6
	case StepSummary:
		return 0
	default:
		return 1
	}
}

func (w *Wizard) inputCountForStep() int {
	return w.maxFieldsForStep()
}

func (w *Wizard) isLastField() bool {
	return w.focusIndex >= w.maxFieldsForStep()-1
}

func (w *Wizard) nextStep() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	if w.step < StepSummary {
		w.step++
	} else {
		return w, messages.WizardComplete(true, nil)
	}

	return w, w.focusCurrent()
}

func (w *Wizard) nextStepNav() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	w.step++
	if w.step > StepSummary {
		w.step = StepSource
	}

	return w, w.focusCurrent()
}

func (w *Wizard) prevStep() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	if w.step > StepSource {
		w.step--
	}

	return w, w.focusCurrent()
}

func (w *Wizard) goToStep(step Step) (*Wizard, tea.Cmd) {
	if step < StepSource || step > StepSummary {
		return w, nil
	}
	if step >= w.step {
		return w, nil
	}

	w.blurAll()
	w.step = step
	w.focusIndex = 0

	if step == StepSource {
		w.sourceStep = NewSourceStep(w.theme)
		w.sourceStep.SetSize(w.width, w.height)
	}

	return w, w.focusCurrent()
}

func (w *Wizard) prevStepNav() (*Wizard, tea.Cmd) {
	w.blurAll()
	w.focusIndex = 0

	w.step--
	if w.step < StepSource {
		w.step = StepSummary
	}

	return w, w.focusCurrent()
}

func (w *Wizard) blurAll() {
	w.prevFocused = false
	w.nextFocused = false
	w.sourceStep.Blur()
	w.hostnameInput.Blur()
	w.descriptionInput.Blur()
	w.profileSelect.Blur()
	w.nameInput.Blur()
	w.emailInput.Blur()
	w.githubInput.Blur()
	w.shellSelect.Blur()
	w.editorSelect.Blur()
	w.themeSelect.Blur()
	w.featureSSH.Blur()
	w.featureGPG.Blur()
	w.featureGitHub.Blur()
	w.featureDocker.Blur()
	w.featureAsdf.Blur()
	w.featureGUI.Blur()
	w.prevButton.Blur()
	w.nextButton.Blur()
}

func (w *Wizard) focusCurrent() tea.Cmd {
	switch w.step {
	case StepSource:
		return w.sourceStep.Focus()
	case StepMachine:
		switch w.focusIndex {
		case 0:
			return w.hostnameInput.Focus()
		case 1:
			return w.descriptionInput.Focus()
		case 2:
			return w.profileSelect.Focus()
		}
	case StepPersonal:
		switch w.focusIndex {
		case 0:
			return w.nameInput.Focus()
		case 1:
			return w.emailInput.Focus()
		case 2:
			return w.githubInput.Focus()
		}
	case StepSettings:
		switch w.focusIndex {
		case 0:
			return w.shellSelect.Focus()
		case 1:
			return w.editorSelect.Focus()
		case 2:
			return w.themeSelect.Focus()
		}
	case StepFeatures:
		switch w.focusIndex {
		case 0:
			return w.featureSSH.Focus()
		case 1:
			return w.featureGPG.Focus()
		case 2:
			return w.featureGitHub.Focus()
		case 3:
			return w.featureDocker.Focus()
		case 4:
			return w.featureAsdf.Focus()
		case 5:
			return w.featureGUI.Focus()
		}
	}
	return nil
}

func (w *Wizard) updateCurrentInputs(msg tea.Msg) (*Wizard, tea.Cmd) {
	var cmd tea.Cmd

	switch w.step {
	case StepSource:
		w.sourceStep, cmd = w.sourceStep.Update(msg)
	case StepMachine:
		switch w.focusIndex {
		case 0:
			w.hostnameInput, cmd = w.hostnameInput.Update(msg)
		case 1:
			w.descriptionInput, cmd = w.descriptionInput.Update(msg)
		case 2:
			w.profileSelect, cmd = w.profileSelect.Update(msg)
		}
	case StepPersonal:
		switch w.focusIndex {
		case 0:
			w.nameInput, cmd = w.nameInput.Update(msg)
		case 1:
			w.emailInput, cmd = w.emailInput.Update(msg)
		case 2:
			w.githubInput, cmd = w.githubInput.Update(msg)
		}
	case StepSettings:
		switch w.focusIndex {
		case 0:
			w.shellSelect, cmd = w.shellSelect.Update(msg)
		case 1:
			w.editorSelect, cmd = w.editorSelect.Update(msg)
		case 2:
			w.themeSelect, cmd = w.themeSelect.Update(msg)
		}
	case StepFeatures:
		switch w.focusIndex {
		case 0:
			w.featureSSH, cmd = w.featureSSH.Update(msg)
		case 1:
			w.featureGPG, cmd = w.featureGPG.Update(msg)
		case 2:
			w.featureGitHub, cmd = w.featureGitHub.Update(msg)
		case 3:
			w.featureDocker, cmd = w.featureDocker.Update(msg)
		case 4:
			w.featureAsdf, cmd = w.featureAsdf.Update(msg)
		case 5:
			w.featureGUI, cmd = w.featureGUI.Update(msg)
		}
	}

	return w, cmd
}

func (w *Wizard) View() string {
	t := w.theme

	var progress string
	if w.step <= StepSummary {
		progress = w.renderProgress()
	}

	hint := w.renderHint()

	arrowWidth := 5
	centerWidth := w.width - (arrowWidth * 2)

	container := components.NewPageContainer(t).
		SetSize(centerWidth, w.height).
		SetHeader(progress).
		SetFooter(hint)

	contentWidth := container.InnerWidth()
	contentHeight := container.ContentHeight()

	var content string
	var useScrollable bool

	switch w.step {
	case StepSource:
		content = w.renderSourceStep()
		useScrollable = false
	case StepMachine:
		content = w.renderMachineStep()
		useScrollable = true
	case StepPersonal:
		content = w.renderPersonalStep()
		useScrollable = true
	case StepSettings:
		content = w.renderSettingsStep()
		useScrollable = true
	case StepFeatures:
		content = w.renderFeaturesStep()
		useScrollable = true
	case StepSummary:
		content = w.renderSummaryStep()
		useScrollable = true
	case StepSetup:
		content = w.renderSetupStep()
		useScrollable = false
	case StepComplete:
		content = w.renderCompleteStep()
		useScrollable = false
	default:
		content = t.S().Title.Render("Unknown Step")
		useScrollable = false
	}

	if useScrollable {
		w.scrollable.SetSize(contentWidth, contentHeight)
		w.scrollable.SetContent(content)
		content = w.scrollable.View()
	}

	container.SetContent(content)

	centerContent := container.View()

	return w.renderWithSideNav(centerContent)
}

func (w *Wizard) renderWithSideNav(centerContent string) string {
	t := w.theme

	arrowWidth := 5
	showPrev := w.canNavigateBack()
	showNext := w.canNavigateForward()

	var leftArrow, rightArrow string

	prevHighlight := w.prevHovered || w.prevFocused || (w.step == StepSource && w.sourceStep.IsPrevFocused())
	nextHighlight := w.nextHovered || w.nextFocused || (w.step == StepSource && w.sourceStep.IsNextFocused())

	if showPrev {
		arrowStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
		if prevHighlight {
			arrowStyle = arrowStyle.
				Foreground(t.BgBase).
				Background(t.Primary)
		} else {
			arrowStyle = arrowStyle.Foreground(t.Primary)
		}
		leftArrow = zone.Mark(w.prevButton.ZoneID(), arrowStyle.Render("◀"))
	}

	if showNext {
		arrowStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
		if nextHighlight {
			arrowStyle = arrowStyle.
				Foreground(t.BgBase).
				Background(t.Primary)
		} else {
			arrowStyle = arrowStyle.Foreground(t.Primary)
		}
		rightArrow = zone.Mark(w.nextButton.ZoneID(), arrowStyle.Render("▶"))
	}

	leftCol := lipgloss.Place(
		arrowWidth,
		w.height,
		lipgloss.Center,
		lipgloss.Center,
		leftArrow,
	)

	rightCol := lipgloss.Place(
		arrowWidth,
		w.height,
		lipgloss.Center,
		lipgloss.Center,
		rightArrow,
	)

	centerWidth := w.width - (arrowWidth * 2)
	centerCol := lipgloss.Place(
		centerWidth,
		w.height,
		lipgloss.Center,
		lipgloss.Center,
		centerContent,
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, centerCol, rightCol)
}

func (w *Wizard) canNavigateBack() bool {
	if w.step == StepSource {
		return w.sourceStep.CanGoBack()
	}
	return w.step > StepSource && w.step <= StepSummary
}

func (w *Wizard) canNavigateForward() bool {
	return w.step <= StepSummary
}

func (w *Wizard) renderProgress() string {
	w.stepper.SetCurrent(int(w.step))
	return w.stepper.View()
}

func (w *Wizard) renderSourceStep() string {
	return w.sourceStep.View()
}

func (w *Wizard) renderMachineStep() string {
	t := w.theme

	title := t.S().Title.Render("Machine Configuration")
	subtitle := t.S().Subtle.Render("Configure your machine identity.")

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		w.hostnameInput.View(),
		"",
		w.descriptionInput.View(),
		"",
		w.profileSelect.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		form,
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

func (w *Wizard) renderSettingsStep() string {
	t := w.theme

	title := t.S().Title.Render("Preferences")
	subtitle := t.S().Subtle.Render("Customize your environment.")

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		w.shellSelect.View(),
		"",
		w.editorSelect.View(),
		"",
		w.themeSelect.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		form,
	)
}

func (w *Wizard) renderFeaturesStep() string {
	t := w.theme

	title := t.S().Title.Render("Features")
	subtitle := t.S().Subtle.Render("Select features to enable.")

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		w.featureSSH.View(),
		w.featureGPG.View(),
		w.featureGitHub.View(),
		w.featureDocker.View(),
		w.featureAsdf.View(),
		w.featureGUI.View(),
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

	sourceConfig := w.sourceStep.Config()
	sourceLabel := "Not configured"
	switch sourceConfig.Mode {
	case SourceModeStarter:
		sourceLabel = "dotts-starter"
	case SourceModeBYOP:
		if sourceConfig.BYOPType == BYOPTypeRemote {
			sourceLabel = "BYOP: " + sourceConfig.RemoteURL
		} else {
			sourceLabel = "BYOP: " + sourceConfig.LocalPath
		}
	case SourceModeForkUser:
		sourceLabel = "Fork: " + sourceConfig.ForkFrom
	}

	profile := w.profileSelect.Selected()
	profileLabel := "Not selected"
	if profile != nil {
		profileLabel = profile.Label
	}

	shell := w.shellSelect.Selected()
	shellLabel := "Not selected"
	if shell != nil {
		shellLabel = shell.Label
	}

	editor := w.editorSelect.Selected()
	editorLabel := "Not selected"
	if editor != nil {
		editorLabel = editor.Label
	}

	colorTheme := w.themeSelect.Selected()
	themeLabel := "Not selected"
	if colorTheme != nil {
		themeLabel = colorTheme.Label
	}

	var features []string
	if w.featureSSH.Checked() {
		features = append(features, "SSH")
	}
	if w.featureGPG.Checked() {
		features = append(features, "GPG")
	}
	if w.featureGitHub.Checked() {
		features = append(features, "GitHub")
	}
	if w.featureDocker.Checked() {
		features = append(features, "Docker")
	}
	if w.featureAsdf.Checked() {
		features = append(features, "asdf")
	}
	if w.featureGUI.Checked() {
		features = append(features, "GUI")
	}
	featuresStr := "None"
	if len(features) > 0 {
		featuresStr = lipgloss.JoinHorizontal(lipgloss.Top, features[0])
		for i := 1; i < len(features); i++ {
			featuresStr += ", " + features[i]
		}
	}

	labelStyle := t.S().Muted.Copy().Width(12)
	valueStyle := t.S().Text

	summaryItems := []string{
		labelStyle.Render("Source:") + " " + valueStyle.Render(sourceLabel),
		labelStyle.Render("Hostname:") + " " + valueStyle.Render(w.hostnameInput.Value()),
		labelStyle.Render("Profile:") + " " + valueStyle.Render(profileLabel),
		labelStyle.Render("Name:") + " " + valueStyle.Render(w.nameInput.Value()),
		labelStyle.Render("Email:") + " " + valueStyle.Render(w.emailInput.Value()),
		labelStyle.Render("GitHub:") + " " + valueStyle.Render(w.githubInput.Value()),
		labelStyle.Render("Shell:") + " " + valueStyle.Render(shellLabel),
		labelStyle.Render("Editor:") + " " + valueStyle.Render(editorLabel),
		labelStyle.Render("Theme:") + " " + valueStyle.Render(themeLabel),
		labelStyle.Render("Features:") + " " + valueStyle.Render(featuresStr),
	}

	summary := lipgloss.JoinVertical(lipgloss.Left, summaryItems...)

	boxWidth := 50
	if w.width > 0 && w.width < 70 {
		boxWidth = w.width - 20
	}

	box := t.S().BorderNormal.
		Width(boxWidth).
		Padding(1, 2).
		Render(summary)

	var errorView string
	if w.setupError != nil {
		errorView = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			t.S().Error.Render(theme.Icons.Error+" Previous setup attempt failed:"),
			t.S().Muted.PaddingLeft(2).Render(w.setupError.Error()),
			"",
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		box,
		errorView,
	)
}

func (w *Wizard) renderSetupStep() string {
	t := w.theme

	title := t.S().Title.Render("Setting Up...")
	subtitle := t.S().Subtle.Render("Please wait while dotts configures your system.")

	spinnerStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Padding(2, 0)

	spinner := spinnerStyle.Render(theme.Icons.Current + " Setting up configuration source...")

	var errorView string
	if w.setupError != nil {
		errorView = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			t.S().Error.Render(theme.Icons.Error+" Setup failed:"),
			t.S().Muted.PaddingLeft(2).Render(w.setupError.Error()),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		spinner,
		errorView,
	)
}

func (w *Wizard) renderCompleteStep() string {
	t := w.theme

	title := t.S().Title.Render(theme.Icons.Success + " Setup Complete!")
	subtitle := t.S().Subtle.Render("Your dotts configuration is ready.")

	var details []string

	if w.setupResult != nil {
		details = append(details, t.S().Success.Render(theme.Icons.Success+" Config stored at:"))
		details = append(details, t.S().Text.PaddingLeft(2).Render(w.setupResult.ConfigPath))

		if w.setupResult.SymlinkPath != "" {
			details = append(details, "")
			details = append(details, t.S().Success.Render(theme.Icons.Success+" Symlink created at:"))
			details = append(details, t.S().Text.PaddingLeft(2).Render(w.setupResult.SymlinkPath))
		}

		if w.setupResult.BackupPath != "" {
			details = append(details, "")
			details = append(details, t.S().Warning.Render(theme.Icons.Info+" Previous config backed up to:"))
			details = append(details, t.S().Muted.PaddingLeft(2).Render(w.setupResult.BackupPath))
		}

		if w.setupResult.ClonedFrom != "" {
			details = append(details, "")
			details = append(details, t.S().Info.Render(theme.Icons.Info+" Cloned from:"))
			details = append(details, t.S().Muted.PaddingLeft(2).Render(w.setupResult.ClonedFrom))
		}

		if w.setupResult.CopiedFrom != "" {
			details = append(details, "")
			details = append(details, t.S().Info.Render(theme.Icons.Info+" Copied from:"))
			details = append(details, t.S().Muted.PaddingLeft(2).Render(w.setupResult.CopiedFrom))
		}
	}

	boxWidth := 60
	if w.width > 0 && w.width < 80 {
		boxWidth = w.width - 20
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Success).
		Width(boxWidth).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, details...))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		box,
		"",
		t.S().Success.Render(theme.Icons.ArrowRight+" Press Enter to continue to dashboard"),
	)
}

func (w *Wizard) renderHint() string {
	t := w.theme

	if w.prevFocused {
		return t.S().Muted.Render("[enter] go back  [esc] cancel")
	}
	if w.nextFocused {
		return t.S().Muted.Render("[enter] continue  [esc] cancel")
	}

	var hint string
	switch w.step {
	case StepSource:
		hint = w.sourceStep.Hint()
	case StepMachine:
		hint = "[tab] cycle  [↑/↓] select  [enter] next  [esc] back"
	case StepPersonal:
		hint = "[tab] cycle  [enter] next  [esc] back"
	case StepSettings:
		hint = "[tab] cycle  [↑/↓] select  [enter] next  [esc] back"
	case StepFeatures:
		hint = "[tab] cycle  [enter] next  [esc] back"
	case StepSummary:
		hint = "[enter] start setup  [esc] back"
	case StepSetup:
		hint = "Please wait..."
	case StepComplete:
		hint = "[enter] continue to dashboard"
	default:
		hint = "[enter] continue  [esc] back"
	}

	if w.stepUsesScrollable() && w.scrollable.NeedsScrolling() {
		hint += "  [pgup/pgdn] scroll"
	}

	return t.S().Muted.Render(hint)
}
