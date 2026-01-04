package wizard

import (
	"os"
	"path/filepath"

	"github.com/arthur404dev/dotts/pkg/vetru/components"
	"github.com/arthur404dev/dotts/pkg/vetru/components/input"
	"github.com/arthur404dev/dotts/pkg/vetru/keys"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type SourceMode string

const (
	SourceModeStarter  SourceMode = "starter"
	SourceModeBYOP     SourceMode = "byop"
	SourceModeForkUser SourceMode = "fork_user"
)

type BYOPType string

const (
	BYOPTypeRemote BYOPType = "remote"
	BYOPTypeLocal  BYOPType = "local"
)

type SourceSubStep int

const (
	SourceSubStepSelectMode SourceSubStep = iota
	SourceSubStepConfigure
	SourceSubStepSymlink
)

type SourceConfig struct {
	Mode SourceMode

	BYOPType  BYOPType
	RemoteURL string
	LocalPath string

	ForkFrom string

	CreateSymlink bool
	SymlinkPath   string
}

type SourceStep struct {
	theme  *theme.Theme
	width  int
	height int

	subStep    SourceSubStep
	focusIndex int
	completed  bool

	prevFocused bool
	nextFocused bool

	modeSelect *input.Select

	byopTypeControl *components.SegmentedControl
	remoteInput     *input.RepoInput
	localPathInput  *input.PathInput

	forkFromInput *input.RepoInput

	symlinkCheckbox *components.Checkbox
	symlinkPath     *input.PathInput

	config SourceConfig
}

const (
	dottsStarterRepo = "github.com/arthur404dev/dotts-starter"
	defaultSymlink   = "~/dotfiles"
)

func NewSourceStep(t *theme.Theme) *SourceStep {
	modeItems := []input.SelectItem{
		{ID: string(SourceModeStarter), Label: "Use dotts-starter (recommended)", Description: "Clone the official starter configs"},
		{ID: string(SourceModeBYOP), Label: "Bring your own", Description: "Use an existing repo or local folder"},
		{ID: string(SourceModeForkUser), Label: "Base on another user", Description: "Clone someone else's dotts-compatible config"},
	}
	modeSelect := input.NewSelect(t, "Where should dotts get your configurations?", modeItems)

	byopItems := []components.SegmentItem{
		{ID: string(BYOPTypeRemote), Label: "Remote repo"},
		{ID: string(BYOPTypeLocal), Label: "Local path"},
	}
	byopTypeControl := components.NewSegmentedControl(t, byopItems)

	remoteInput := input.NewRepoInput(t, "Repository URL", "github.com/username/dotfiles")
	remoteInput.SetHelp("No https:// needed")

	localPathInput := input.NewPathInput(t, "Path to existing dotfiles", "~/dotfiles")
	localPathInput.SetValue("~/")
	localPathInput.SetHelp("Will be copied to managed location")

	forkFromInput := input.NewRepoInput(t, "Source repository", "github.com/username/dotfiles")
	forkFromInput.SetHelp("A dotts-compatible repository to clone")

	symlinkCheckbox := components.NewCheckbox(t, "Create symlink for easy access").SetChecked(true)

	symlinkPath := input.NewPathInput(t, "Symlink path", defaultSymlink)
	symlinkPath.SetValue(defaultSymlink)
	symlinkPath.SetHelp("Where to create the symlink")

	return &SourceStep{
		theme:           t,
		subStep:         SourceSubStepSelectMode,
		modeSelect:      modeSelect,
		byopTypeControl: byopTypeControl,
		remoteInput:     remoteInput,
		localPathInput:  localPathInput,
		forkFromInput:   forkFromInput,
		symlinkCheckbox: symlinkCheckbox,
		symlinkPath:     symlinkPath,
		config: SourceConfig{
			Mode:          SourceModeStarter,
			BYOPType:      BYOPTypeRemote,
			CreateSymlink: true,
			SymlinkPath:   defaultSymlink,
		},
	}
}

func (s *SourceStep) SetSize(width, height int) {
	s.width = width
	s.height = height

	inputWidth := width - 10
	if inputWidth > 60 {
		inputWidth = 60
	}

	s.modeSelect.SetWidth(inputWidth)
	s.byopTypeControl.SetWidth(inputWidth)
	s.remoteInput.SetWidth(inputWidth)
	s.localPathInput.SetWidth(inputWidth)
	s.forkFromInput.SetWidth(inputWidth)
	s.symlinkPath.SetWidth(inputWidth)
}

func (s *SourceStep) Focus() tea.Cmd {
	return s.focusCurrent()
}

func (s *SourceStep) Blur() {
	s.blurAll()
}

func (s *SourceStep) Config() SourceConfig {
	s.config.Mode = SourceMode(s.modeSelect.SelectedID())
	s.config.BYOPType = BYOPType(s.byopTypeControl.SelectedID())
	s.config.RemoteURL = s.remoteInput.Value()
	s.config.LocalPath = s.localPathInput.Value()
	s.config.ForkFrom = s.forkFromInput.Value()
	s.config.CreateSymlink = s.symlinkCheckbox.Checked()
	s.config.SymlinkPath = s.symlinkPath.Value()
	return s.config
}

func (s *SourceStep) IsComplete() bool {
	return s.completed
}

func (s *SourceStep) CanGoBack() bool {
	return s.subStep > SourceSubStepSelectMode
}

func (s *SourceStep) CanGoNext() bool {
	return true
}

func (s *SourceStep) IsOnLastSubStep() bool {
	return s.subStep == SourceSubStepSymlink
}

func (s *SourceStep) GoNext() tea.Cmd {
	result, cmd := s.handleNext()
	*s = *result
	return cmd
}

func (s *SourceStep) GoBack() tea.Cmd {
	result, cmd := s.handleBack()
	*s = *result
	return cmd
}

func (s *SourceStep) inputCountForSubStep() int {
	switch s.subStep {
	case SourceSubStepSelectMode:
		return 1
	case SourceSubStepConfigure:
		switch s.config.Mode {
		case SourceModeStarter:
			return 0
		case SourceModeBYOP:
			return 2
		case SourceModeForkUser:
			return 1
		}
	case SourceSubStepSymlink:
		if s.symlinkCheckbox.Checked() {
			return 2
		}
		return 1
	}
	return 0
}

func (s *SourceStep) maxFieldsForSubStep() int {
	return s.inputCountForSubStep()
}

func (s *SourceStep) isLastField() bool {
	maxFields := s.maxFieldsForSubStep()
	if maxFields == 0 {
		return true
	}
	return s.focusIndex >= maxFields-1
}

func (s *SourceStep) IsPrevFocused() bool {
	return s.prevFocused
}

func (s *SourceStep) IsNextFocused() bool {
	return s.nextFocused
}

func (s *SourceStep) Update(msg tea.Msg) (*SourceStep, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.pathInputWantsKey(msg) {
			return s.updateCurrentInputs(msg)
		}

		switch {
		case keys.MatchesConfirm(msg):
			return s.handleEnter()

		case keys.MatchesNextField(msg):
			return s.nextField()

		case keys.MatchesPrevField(msg):
			return s.prevField()

		case keys.MatchesEscape(msg):
			return s.handleEscape()
		}
	}

	return s.updateCurrentInputs(msg)
}

func (s *SourceStep) handleEnter() (*SourceStep, tea.Cmd) {
	if s.nextFocused {
		return s.handleNext()
	}
	if s.prevFocused {
		return s.handleBack()
	}

	s.updateCurrentInputs(tea.KeyMsg{Type: tea.KeySpace})

	if s.isLastField() {
		s.blurAll()
		s.nextFocused = true
		return s, nil
	}
	return s.nextField()
}

func (s *SourceStep) handleEscape() (*SourceStep, tea.Cmd) {
	if s.nextFocused {
		s.nextFocused = false
		maxFields := s.maxFieldsForSubStep()
		if maxFields > 0 {
			s.focusIndex = maxFields - 1
		} else {
			s.focusIndex = 0
		}
		return s, s.focusCurrent()
	}
	if s.prevFocused {
		s.prevFocused = false
		s.focusIndex = 0
		return s, s.focusCurrent()
	}

	if s.focusIndex == 0 && s.CanGoBack() {
		s.blurAll()
		s.prevFocused = true
		return s, nil
	}

	return s.prevField()
}

func (s *SourceStep) handleNext() (*SourceStep, tea.Cmd) {
	s.blurAll()
	switch s.subStep {
	case SourceSubStepSelectMode:
		s.config.Mode = SourceMode(s.modeSelect.SelectedID())
		s.subStep = SourceSubStepConfigure
		s.focusIndex = 0
		return s, s.focusCurrent()

	case SourceSubStepConfigure:
		s.subStep = SourceSubStepSymlink
		s.focusIndex = 0
		return s, s.focusCurrent()

	case SourceSubStepSymlink:
		s.completed = true
		return s, nil
	}
	return s, nil
}

func (s *SourceStep) handleBack() (*SourceStep, tea.Cmd) {
	s.blurAll()
	switch s.subStep {
	case SourceSubStepSelectMode:
		return s, nil

	case SourceSubStepConfigure:
		s.subStep = SourceSubStepSelectMode
		s.focusIndex = 0
		return s, s.focusCurrent()

	case SourceSubStepSymlink:
		s.subStep = SourceSubStepConfigure
		s.focusIndex = 0
		return s, s.focusCurrent()
	}
	return s, nil
}

func (s *SourceStep) pathInputWantsKey(msg tea.KeyMsg) bool {
	return s.CapturesKey(msg)
}

func (s *SourceStep) CapturesKey(msg tea.KeyMsg) bool {
	key := msg.String()
	isNavKey := key == "tab" || key == " " || key == "up" || key == "down" ||
		key == "left" || key == "right" || key == "ctrl+n" || key == "ctrl+p"

	if !isNavKey {
		return false
	}

	switch s.subStep {
	case SourceSubStepConfigure:
		if s.config.Mode == SourceModeBYOP && s.focusIndex == 1 {
			if s.config.BYOPType == BYOPTypeLocal {
				return s.localPathInput.HasSuggestions()
			}
		}
	case SourceSubStepSymlink:
		if s.focusIndex == 1 && s.symlinkCheckbox.Checked() {
			return s.symlinkPath.HasSuggestions()
		}
	}

	return false
}

func (s *SourceStep) handleMouse(msg tea.MouseMsg) (*SourceStep, tea.Cmd) {
	switch s.subStep {
	case SourceSubStepSelectMode:
		if zone.Get(s.modeSelect.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
			s.blurAll()
			s.focusIndex = 0
			return s, s.focusCurrent()
		}

	case SourceSubStepConfigure:
		if s.config.Mode == SourceModeBYOP {
			if zone.Get(s.byopTypeControl.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
				s.blurAll()
				s.focusIndex = 0
				return s, s.focusCurrent()
			}
			if s.config.BYOPType == BYOPTypeLocal {
				if zone.Get(s.localPathInput.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
					s.blurAll()
					s.focusIndex = 1
					return s, s.focusCurrent()
				}
			} else {
				if zone.Get(s.remoteInput.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
					s.blurAll()
					s.focusIndex = 1
					return s, s.focusCurrent()
				}
			}
		} else if s.config.Mode == SourceModeForkUser {
			if zone.Get(s.forkFromInput.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
				s.blurAll()
				s.focusIndex = 0
				return s, s.focusCurrent()
			}
		}

	case SourceSubStepSymlink:
		if zone.Get(s.symlinkCheckbox.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
			s.blurAll()
			s.focusIndex = 0
			return s, s.focusCurrent()
		}
		if zone.Get(s.symlinkPath.ZoneID()).InBounds(msg) && msg.Action == tea.MouseActionRelease {
			s.blurAll()
			s.focusIndex = 1
			return s, s.focusCurrent()
		}
	}

	return s.forwardMouseToCurrentInput(msg)
}

func (s *SourceStep) forwardMouseToCurrentInput(msg tea.MouseMsg) (*SourceStep, tea.Cmd) {
	var cmd tea.Cmd

	switch s.subStep {
	case SourceSubStepSelectMode:
		s.modeSelect, cmd = s.modeSelect.Update(msg)

	case SourceSubStepConfigure:
		switch s.config.Mode {
		case SourceModeBYOP:
			if s.focusIndex == 0 {
				s.byopTypeControl, cmd = s.byopTypeControl.Update(msg)
			} else if s.focusIndex == 1 {
				if s.config.BYOPType == BYOPTypeRemote {
					s.remoteInput, cmd = s.remoteInput.Update(msg)
				} else {
					s.localPathInput, cmd = s.localPathInput.Update(msg)
				}
			}
		case SourceModeForkUser:
			if s.focusIndex == 0 {
				s.forkFromInput, cmd = s.forkFromInput.Update(msg)
			}
		}

	case SourceSubStepSymlink:
		if s.focusIndex == 0 {
			s.symlinkCheckbox, cmd = s.symlinkCheckbox.Update(msg)
		} else if s.focusIndex == 1 {
			s.symlinkPath, cmd = s.symlinkPath.Update(msg)
		}
	}

	return s, cmd
}

func (s *SourceStep) nextField() (*SourceStep, tea.Cmd) {
	maxFields := s.maxFieldsForSubStep()
	if maxFields <= 1 {
		return s, nil
	}

	s.blurAll()
	s.focusIndex = (s.focusIndex + 1) % maxFields
	return s, s.focusCurrent()
}

func (s *SourceStep) prevField() (*SourceStep, tea.Cmd) {
	maxFields := s.maxFieldsForSubStep()
	if maxFields <= 1 {
		return s, nil
	}

	s.blurAll()
	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxFields - 1
	}
	return s, s.focusCurrent()
}

func (s *SourceStep) blurAll() {
	s.prevFocused = false
	s.nextFocused = false
	s.modeSelect.Blur()
	s.byopTypeControl.Blur()
	s.remoteInput.Blur()
	s.localPathInput.Blur()
	s.forkFromInput.Blur()
	s.symlinkCheckbox.Blur()
	s.symlinkPath.Blur()
}

func (s *SourceStep) focusCurrent() tea.Cmd {
	switch s.subStep {
	case SourceSubStepSelectMode:
		if s.focusIndex == 0 {
			return s.modeSelect.Focus()
		}

	case SourceSubStepConfigure:
		switch s.config.Mode {
		case SourceModeStarter:
			return nil

		case SourceModeBYOP:
			if s.focusIndex == 0 {
				return s.byopTypeControl.Focus()
			}
			if s.focusIndex == 1 {
				if s.config.BYOPType == BYOPTypeRemote {
					return s.remoteInput.Focus()
				}
				return s.localPathInput.Focus()
			}

		case SourceModeForkUser:
			if s.focusIndex == 0 {
				return s.forkFromInput.Focus()
			}
		}

	case SourceSubStepSymlink:
		if s.focusIndex == 0 {
			return s.symlinkCheckbox.Focus()
		}
		if s.focusIndex == 1 {
			return s.symlinkPath.Focus()
		}
	}
	return nil
}

func (s *SourceStep) updateCurrentInputs(msg tea.Msg) (*SourceStep, tea.Cmd) {
	if mouseMsg, ok := msg.(tea.MouseMsg); ok {
		return s.handleMouse(mouseMsg)
	}

	var cmd tea.Cmd

	switch s.subStep {
	case SourceSubStepSelectMode:
		if s.focusIndex == 0 {
			s.modeSelect, cmd = s.modeSelect.Update(msg)
			s.config.Mode = SourceMode(s.modeSelect.SelectedID())
		}

	case SourceSubStepConfigure:
		switch s.config.Mode {
		case SourceModeBYOP:
			if s.focusIndex == 0 {
				s.byopTypeControl, cmd = s.byopTypeControl.Update(msg)
				s.config.BYOPType = BYOPType(s.byopTypeControl.SelectedID())
			} else if s.focusIndex == 1 {
				if s.config.BYOPType == BYOPTypeRemote {
					s.remoteInput, cmd = s.remoteInput.Update(msg)
				} else {
					s.localPathInput, cmd = s.localPathInput.Update(msg)
				}
			}

		case SourceModeForkUser:
			if s.focusIndex == 0 {
				s.forkFromInput, cmd = s.forkFromInput.Update(msg)
			}
		}

	case SourceSubStepSymlink:
		if s.focusIndex == 0 {
			s.symlinkCheckbox, cmd = s.symlinkCheckbox.Update(msg)
			s.config.CreateSymlink = s.symlinkCheckbox.Checked()
		} else if s.focusIndex == 1 {
			s.symlinkPath, cmd = s.symlinkPath.Update(msg)
		}
	}

	return s, cmd
}

func (s *SourceStep) View() string {
	progress := s.renderSubProgress()

	var content string
	switch s.subStep {
	case SourceSubStepSelectMode:
		content = s.renderSelectMode()
	case SourceSubStepConfigure:
		content = s.renderConfigure()
	case SourceSubStepSymlink:
		content = s.renderSymlink()
	}

	return lipgloss.JoinVertical(lipgloss.Left, progress, "", content)
}

func (s *SourceStep) renderSubProgress() string {
	t := s.theme

	steps := []string{"Choose", "Configure", "Symlink"}
	var parts []string

	for i, step := range steps {
		var style lipgloss.Style
		var icon string

		if i < int(s.subStep) {
			style = lipgloss.NewStyle().Foreground(t.Success)
			icon = theme.Icons.Success + " "
		} else if i == int(s.subStep) {
			style = lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
			icon = theme.Icons.Current + " "
		} else {
			style = lipgloss.NewStyle().Foreground(t.FgMuted)
			icon = theme.Icons.Pending + " "
		}

		parts = append(parts, style.Render(icon+step))
	}

	separator := lipgloss.NewStyle().Foreground(t.FgMuted).Render("  →  ")
	return lipgloss.JoinHorizontal(lipgloss.Center, parts[0], separator, parts[1], separator, parts[2])
}

func (s *SourceStep) renderSelectMode() string {
	t := s.theme

	title := t.S().Title.Render("Configuration Source")
	subtitle := t.S().Subtle.Render("Choose where to get your dotfiles configuration from.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		s.modeSelect.View(),
	)
}

func (s *SourceStep) renderConfigure() string {
	t := s.theme

	switch s.config.Mode {
	case SourceModeStarter:
		return s.renderStarterConfig()

	case SourceModeBYOP:
		title := t.S().Title.Render("Bring Your Own")
		subtitle := t.S().Subtle.Render("Use an existing repository or local folder.")

		var inputView string
		if s.config.BYOPType == BYOPTypeRemote {
			inputView = s.remoteInput.View()
		} else {
			inputView = s.localPathInput.View()

			warning := t.S().Warning.Render(
				theme.Icons.Info + " Your dotfiles will be copied to the managed location",
			)
			detail := t.S().Muted.PaddingLeft(2).Render("→ ~/.local/share/dotts/config (original stays untouched)")
			inputView = lipgloss.JoinVertical(lipgloss.Left, inputView, "", warning, detail)
		}

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			s.byopTypeControl.View(),
			"",
			inputView,
		)

	case SourceModeForkUser:
		title := t.S().Title.Render("Base on Another User")
		subtitle := t.S().Subtle.Render("Clone from another user's dotts-compatible repository.")

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			s.forkFromInput.View(),
		)
	}

	return ""
}

func (s *SourceStep) renderStarterConfig() string {
	t := s.theme

	title := t.S().Title.Render("dotts-starter")
	subtitle := t.S().Subtle.Render("Clone the official starter configuration.")

	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Success).
		Padding(1, 2).
		Width(50)

	info := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Success.Render(theme.Icons.Success+" Will clone:"),
		t.S().Text.PaddingLeft(2).Render(dottsStarterRepo),
		"",
		t.S().Success.Render(theme.Icons.Success+" To:"),
		t.S().Text.PaddingLeft(2).Render("~/.local/share/dotts/config"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		infoStyle.Render(info),
	)
}

func (s *SourceStep) renderSymlink() string {
	t := s.theme

	title := t.S().Title.Render("Symlink Configuration")
	subtitle := t.S().Subtle.Render("Create a symlink for easy access to your config.")

	var content string
	if s.symlinkCheckbox.Checked() {
		symlinkPathView := s.symlinkPath.View()

		expanded := expandSymlinkPath(s.symlinkPath.Value())
		if pathExists(expanded) {
			warning := t.S().Warning.Render(
				theme.Icons.Warning + " Path already exists! Will be backed up.",
			)
			symlinkPathView = lipgloss.JoinVertical(lipgloss.Left, symlinkPathView, warning)
		}

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			s.symlinkCheckbox.View(),
			"",
			symlinkPathView,
		)
	} else {
		content = s.symlinkCheckbox.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		content,
	)
}

func (s *SourceStep) Hint() string {
	if s.prevFocused {
		return "[enter] go back  [esc] cancel"
	}
	if s.nextFocused {
		if s.subStep == SourceSubStepSymlink {
			return "[enter] confirm & continue  [esc] cancel"
		}
		return "[enter] continue  [esc] cancel"
	}

	switch s.subStep {
	case SourceSubStepSelectMode:
		return "[tab] cycle  [↑/↓] select  [enter] next  [esc] back"
	case SourceSubStepConfigure:
		return "[tab] cycle  [enter] next  [esc] back"
	case SourceSubStepSymlink:
		return "[tab] cycle  [enter] next  [esc] back"
	}
	return "[tab] cycle  [enter] next  [esc] back"
}

func expandSymlinkPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
