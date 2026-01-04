package input

import (
	"regexp"
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	repoShorthandPattern = regexp.MustCompile(`^([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)$`)
	repoHostPattern      = regexp.MustCompile(`^(github\.com|gitlab\.com|bitbucket\.org)/([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)/?$`)
	repoFullURLPattern   = regexp.MustCompile(`^https?://`)
)

type RepoInput struct {
	model   textinput.Model
	theme   *theme.Theme
	zoneID  string
	label   string
	help    string
	focused bool
	hovered bool
	width   int

	defaultHost string
}

func NewRepoInput(t *theme.Theme, label, placeholder string) *RepoInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 40

	zoneID := "repo-input-" + strings.ReplaceAll(label, " ", "-")

	return &RepoInput{
		model:       ti,
		theme:       t,
		zoneID:      zoneID,
		label:       label,
		width:       50,
		defaultHost: "github.com",
	}
}

func (r *RepoInput) SetWidth(w int) *RepoInput {
	r.width = w
	r.model.Width = w - 4
	return r
}

func (r *RepoInput) SetLabel(label string) *RepoInput {
	r.label = label
	return r
}

func (r *RepoInput) SetHelp(help string) *RepoInput {
	r.help = help
	return r
}

func (r *RepoInput) SetPlaceholder(placeholder string) *RepoInput {
	r.model.Placeholder = placeholder
	return r
}

func (r *RepoInput) SetValue(value string) *RepoInput {
	r.model.SetValue(value)
	return r
}

func (r *RepoInput) SetDefaultHost(host string) *RepoInput {
	r.defaultHost = host
	return r
}

func (r *RepoInput) Value() string {
	return r.model.Value()
}

func (r *RepoInput) NormalizedURL() string {
	return normalizeRepoURL(r.model.Value(), r.defaultHost)
}

func (r *RepoInput) CloneURL() string {
	url := r.NormalizedURL()
	if url != "" && !strings.HasSuffix(url, ".git") {
		url += ".git"
	}
	return url
}

func (r *RepoInput) ZoneID() string {
	return r.zoneID
}

func (r *RepoInput) Focus() tea.Cmd {
	r.focused = true
	return r.model.Focus()
}

func (r *RepoInput) Blur() {
	r.focused = false
	r.model.Blur()
}

func (r *RepoInput) Focused() bool {
	return r.focused
}

func (r *RepoInput) Hovered() bool {
	return r.hovered
}

func (r *RepoInput) IsValid() bool {
	value := strings.TrimSpace(r.model.Value())
	if value == "" {
		return false
	}
	return repoShorthandPattern.MatchString(value) ||
		repoHostPattern.MatchString(value) ||
		repoFullURLPattern.MatchString(value)
}

func (r *RepoInput) Update(msg tea.Msg) (*RepoInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		inBounds := zone.Get(r.zoneID).InBounds(msg)

		switch msg.Action {
		case tea.MouseActionMotion:
			r.hovered = inBounds

		case tea.MouseActionRelease:
			if inBounds && !r.focused {
				r.focused = true
				return r, r.model.Focus()
			}
		}
	}

	if r.focused {
		var cmd tea.Cmd
		r.model, cmd = r.model.Update(msg)
		return r, cmd
	}

	return r, nil
}

func (r *RepoInput) View() string {
	th := r.theme
	var rows []string

	if r.label != "" {
		labelView := th.S().TextInput.Label.Render(r.label)
		rows = append(rows, labelView)
	}

	var borderStyle lipgloss.Style
	switch {
	case r.focused:
		borderStyle = th.S().TextInput.FocusedB.Width(r.width)
	case r.hovered:
		borderStyle = th.S().TextInput.Hovered.Width(r.width)
	default:
		borderStyle = th.S().TextInput.Normal.Width(r.width)
	}

	fieldView := borderStyle.Render(r.model.View())
	rows = append(rows, fieldView)

	value := r.model.Value()
	if value != "" {
		normalized := r.NormalizedURL()
		if normalized != "" && normalized != value {
			preview := th.S().Muted.Render(theme.Icons.ArrowRight + " " + normalized)
			rows = append(rows, preview)
		}

		if !r.IsValid() && len(value) > 3 {
			errMsg := th.S().Error.Render(theme.Icons.Warning + " Invalid repository format")
			rows = append(rows, errMsg)
		}
	}

	if r.help != "" && value == "" {
		helpView := th.S().TextInput.Help.Render(r.help)
		rows = append(rows, helpView)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return zone.Mark(r.zoneID, content)
}

func (r *RepoInput) Blink() tea.Msg {
	return textinput.Blink()
}

func normalizeRepoURL(input, defaultHost string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	input = strings.TrimSuffix(input, "/")
	input = strings.TrimSuffix(input, ".git")

	if repoFullURLPattern.MatchString(input) {
		return input
	}

	if repoHostPattern.MatchString(input) {
		return "https://" + input
	}

	if repoShorthandPattern.MatchString(input) {
		return "https://" + defaultHost + "/" + input
	}

	if strings.Contains(input, "/") {
		return "https://" + defaultHost + "/" + input
	}

	return ""
}
