package dashboard

import (
	"github.com/arthur404dev/dotts/internal/tui/components"
	"github.com/arthur404dev/dotts/internal/tui/keys"
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/internal/tui/messages"
	"github.com/arthur404dev/dotts/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	theme  *theme.Theme
	width  int
	height int

	tabs     *components.Tabs
	spinner  *components.Spinner
	progress float64

	checkbox1 *components.Checkbox
	checkbox2 *components.Checkbox
	toggle1   *components.Toggle
	toggle2   *components.Toggle
	confirm   *components.Confirm

	navMenu      *components.Menu
	navFocusItem int
}

func New(t *theme.Theme) *Dashboard {
	tabs := components.NewTabs(t, []components.TabItem{
		{ID: "feedback", Label: "Feedback"},
		{ID: "input", Label: "Input"},
		{ID: "navigation", Label: "Navigation"},
		{ID: "data", Label: "Data"},
		{ID: "layout", Label: "Layout"},
	})

	spinner := components.NewSpinner(t).SetLabel("Loading...")

	menu := components.NewMenu(t, []components.MenuItem{
		{ID: "dashboard", Label: "Dashboard", Icon: theme.Icons.Pending},
		{ID: "settings", Label: "Settings", Icon: theme.Icons.Pending},
		{ID: "status", Label: "Status", Icon: theme.Icons.Success},
		{ID: "disabled", Label: "Disabled", Icon: theme.Icons.Error, Enabled: false},
	})

	return &Dashboard{
		theme:     t,
		tabs:      tabs,
		spinner:   spinner,
		progress:  0.65,
		checkbox1: components.NewCheckbox(t, "Enable feature A").SetChecked(true),
		checkbox2: components.NewCheckbox(t, "Enable feature B"),
		toggle1:   components.NewToggle(t, "Dark mode").SetEnabled(true),
		toggle2:   components.NewToggle(t, "Notifications").SetLabels("Yes", "No"),
		confirm:   components.NewConfirm(t, "Proceed with setup?").SetDefaultYes(),
		navMenu:   menu,
	}
}

func (d *Dashboard) ID() messages.PageID {
	return app.PageDashboard
}

func (d *Dashboard) Title() string {
	return "Component Showcase"
}

func (d *Dashboard) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *Dashboard) Focus() tea.Cmd {
	d.tabs.Focus()
	return d.spinner.Init()
}

func (d *Dashboard) Blur() {
	d.tabs.Blur()
	d.blurAll()
}

func (d *Dashboard) Init() tea.Cmd {
	return d.spinner.Init()
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	d.spinner, cmd = d.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keys.MatchesNavLeft(msg):
			d.tabs.Prev()
			d.blurAll()
			return d, cmd

		case keys.MatchesNavRight(msg):
			d.tabs.Next()
			d.blurAll()
			return d, cmd

		case keys.MatchesListDown(msg):
			return d.handleListDown(cmd)

		case keys.MatchesListUp(msg):
			return d.handleListUp(cmd)

		case msg.String() == "+", msg.String() == "=":
			if d.progress < 1 {
				d.progress += 0.1
				if d.progress > 1 {
					d.progress = 1
				}
			}
			return d, cmd

		case msg.String() == "-", msg.String() == "_":
			if d.progress > 0 {
				d.progress -= 0.1
				if d.progress < 0 {
					d.progress = 0
				}
			}
			return d, cmd
		}
	}

	return d.updateCurrentTab(msg, cmd)
}

func (d *Dashboard) handleListDown(cmd tea.Cmd) (*Dashboard, tea.Cmd) {
	switch d.tabs.ActiveID() {
	case "input":
		d.blurAll()
		d.navFocusItem = (d.navFocusItem + 1) % 5
		d.focusInputItem()
	case "navigation":
		d.navMenu, _ = d.navMenu.Update(tea.KeyMsg{})
	}
	return d, cmd
}

func (d *Dashboard) handleListUp(cmd tea.Cmd) (*Dashboard, tea.Cmd) {
	switch d.tabs.ActiveID() {
	case "input":
		d.blurAll()
		d.navFocusItem--
		if d.navFocusItem < 0 {
			d.navFocusItem = 4
		}
		d.focusInputItem()
	case "navigation":
		d.navMenu, _ = d.navMenu.Update(tea.KeyMsg{})
	}
	return d, cmd
}

func (d *Dashboard) focusInputItem() {
	switch d.navFocusItem {
	case 0:
		d.checkbox1.Focus()
	case 1:
		d.checkbox2.Focus()
	case 2:
		d.toggle1.Focus()
	case 3:
		d.toggle2.Focus()
	case 4:
		d.confirm.Focus()
	}
}

func (d *Dashboard) blurAll() {
	d.checkbox1.Blur()
	d.checkbox2.Blur()
	d.toggle1.Blur()
	d.toggle2.Blur()
	d.confirm.Blur()
	d.navMenu.Blur()
}

func (d *Dashboard) updateCurrentTab(msg tea.Msg, spinnerCmd tea.Cmd) (*Dashboard, tea.Cmd) {
	var cmd tea.Cmd

	switch d.tabs.ActiveID() {
	case "input":
		switch d.navFocusItem {
		case 0:
			d.checkbox1, cmd = d.checkbox1.Update(msg)
		case 1:
			d.checkbox2, cmd = d.checkbox2.Update(msg)
		case 2:
			d.toggle1, cmd = d.toggle1.Update(msg)
		case 3:
			d.toggle2, cmd = d.toggle2.Update(msg)
		case 4:
			d.confirm, cmd = d.confirm.Update(msg)
		}
	case "navigation":
		d.navMenu.Focus()
		d.navMenu, cmd = d.navMenu.Update(msg)
	}

	if spinnerCmd != nil {
		return d, tea.Batch(cmd, spinnerCmd)
	}
	return d, cmd
}

func (d *Dashboard) View() string {
	t := d.theme

	title := t.S().Title.Render("Component Showcase")
	subtitle := t.S().Subtle.Render("dotts TUI Design System")

	tabBar := d.tabs.View()
	sep := components.NewSeparator(t).SetLength(d.width - 4).View()

	var content string
	switch d.tabs.ActiveID() {
	case "feedback":
		content = d.renderFeedbackTab()
	case "input":
		content = d.renderInputTab()
	case "navigation":
		content = d.renderNavigationTab()
	case "data":
		content = d.renderDataTab()
	case "layout":
		content = d.renderLayoutTab()
	}

	hints := d.renderHints()

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		tabBar,
		sep,
		"",
		content,
		"",
		hints,
	)

	container := lipgloss.NewStyle().
		Width(d.width).
		Height(d.height).
		Padding(1, 2)

	return container.Render(body)
}

func (d *Dashboard) renderHints() string {
	t := d.theme
	hints := "[←/→] tabs"

	switch d.tabs.ActiveID() {
	case "input":
		hints += " • [↑/↓] navigate • [space/enter] toggle"
	case "navigation":
		hints += " • [↑/↓] menu"
	case "feedback":
		hints += " • [+/-] progress"
	}

	hints += " • [ctrl+p] commands"
	return t.S().Muted.Render(hints)
}

func (d *Dashboard) renderFeedbackTab() string {
	t := d.theme
	colWidth := (d.width - 10) / 2
	if colWidth > 50 {
		colWidth = 50
	}

	col1 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Badges"),
		"",
		d.renderBadges(),
		"",
		t.S().Subtitle.Render("Progress Bar"),
		t.S().Muted.Render("Use +/- to adjust"),
		"",
		d.renderProgress(colWidth),
		"",
		t.S().Subtitle.Render("Spinner"),
		"",
		d.spinner.View(),
	)

	col2 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Alerts"),
		"",
		d.renderAlerts(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
}

func (d *Dashboard) renderBadges() string {
	t := d.theme
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Center,
		components.SuccessBadge(t, "Success").View(),
		" ",
		components.WarningBadge(t, "Warning").View(),
		" ",
		components.ErrorBadge(t, "Error").View(),
		" ",
		components.InfoBadge(t, "Info").View(),
	)

	row2 := lipgloss.JoinHorizontal(
		lipgloss.Center,
		components.NewBadge(t, "Default").View(),
		" ",
		components.NewBadge(t, "Primary").SetVariant(components.BadgePrimary).View(),
		" ",
		components.NewBadge(t, "Secondary").SetVariant(components.BadgeSecondary).View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, row1, "", row2)
}

func (d *Dashboard) renderAlerts() string {
	t := d.theme
	return lipgloss.JoinVertical(
		lipgloss.Left,
		components.InfoAlert(t, "This is an informational message").View(),
		"",
		components.SuccessAlert(t, "Operation completed successfully").View(),
		"",
		components.WarningAlert(t, "Some packages need updates").View(),
		"",
		components.ErrorAlert(t, "Failed to connect to repository").View(),
	)
}

func (d *Dashboard) renderProgress(width int) string {
	t := d.theme
	barWidth := width - 10
	if barWidth < 20 {
		barWidth = 20
	}

	bar1 := components.NewProgressBar(t).
		SetProgress(d.progress).
		SetWidth(barWidth).
		SetLabel("Sync")

	bar2 := components.NewProgressBar(t).
		SetProgress(0.33).
		SetWidth(barWidth).
		SetGradient(false).
		SetLabel("Download")

	bar3 := components.NewProgressBar(t).
		SetProgress(1.0).
		SetWidth(barWidth).
		SetLabel("Complete")

	return lipgloss.JoinVertical(lipgloss.Left, bar1.View(), "", bar2.View(), "", bar3.View())
}

func (d *Dashboard) renderInputTab() string {
	t := d.theme
	colWidth := (d.width - 10) / 2
	if colWidth > 50 {
		colWidth = 50
	}

	col1 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Checkboxes"),
		t.S().Muted.Render("Use ↑/↓ to navigate, space to toggle"),
		"",
		d.checkbox1.View(),
		d.checkbox2.View(),
		"",
		t.S().Subtitle.Render("Toggles"),
		"",
		d.toggle1.View(),
		d.toggle2.View(),
	)

	col2 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Confirm"),
		t.S().Muted.Render("Use ←/→ or y/n to select"),
		"",
		d.confirm.View(),
		"",
		t.S().Subtitle.Render("State Summary"),
		"",
		d.renderInputState(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
}

func (d *Dashboard) renderInputState() string {
	t := d.theme

	kv := components.NewKeyValue(t).
		Add("Feature A", boolToStr(d.checkbox1.Checked())).
		Add("Feature B", boolToStr(d.checkbox2.Checked())).
		Add("Dark mode", boolToStr(d.toggle1.Enabled())).
		Add("Notifications", boolToStr(d.toggle2.Enabled())).
		Add("Proceed", boolToStr(d.confirm.Value()))

	return kv.View()
}

func boolToStr(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (d *Dashboard) renderNavigationTab() string {
	t := d.theme
	colWidth := (d.width - 10) / 2
	if colWidth > 50 {
		colWidth = 50
	}

	stepper := components.NewStepper(t, []components.StepItem{
		{ID: "source", Label: "Source"},
		{ID: "machine", Label: "Machine"},
		{ID: "personal", Label: "Personal"},
		{ID: "settings", Label: "Settings"},
	})
	stepper.SetCurrent(1)

	breadcrumb := components.NewBreadcrumb(t, []components.BreadcrumbItem{
		{ID: "home", Label: "Home"},
		{ID: "settings", Label: "Settings"},
		{ID: "profile", Label: "Profile"},
	})

	col1 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Stepper"),
		"",
		stepper.View(),
		"",
		t.S().Subtitle.Render("Breadcrumb"),
		"",
		breadcrumb.View(),
		"",
		components.NewBreadcrumb(t, []components.BreadcrumbItem{
			{ID: "dotts", Label: "dotts"},
			{ID: "configs", Label: "configs"},
			{ID: "shell", Label: "shell"},
		}).SetSeparator(" / ").View(),
	)

	col2 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Menu"),
		t.S().Muted.Render("Use ↑/↓ to navigate"),
		"",
		d.navMenu.View(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
}

func (d *Dashboard) renderDataTab() string {
	t := d.theme
	colWidth := (d.width - 10) / 2
	if colWidth > 50 {
		colWidth = 50
	}

	kv := components.NewKeyValue(t).
		Add("OS", "Arch Linux").
		Add("Arch", "x86_64").
		Add("Host", "workstation").
		Add("Packages", "142")

	table := components.NewTable(t).
		SetHeaders("Package", "Version", "Status").
		AddRow("neovim", "0.9.5", "installed").
		AddRow("tmux", "3.4", "installed").
		AddRow("fish", "3.7.0", "pending").
		SetStriped(true)

	col1 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Key-Value"),
		"",
		kv.View(),
		"",
		t.S().Subtitle.Render("Table"),
		"",
		table.View(),
	)

	tree := components.NewTree(t).SetNodes([]components.TreeNode{
		{
			ID:       "configs",
			Label:    "configs",
			Expanded: true,
			Children: []components.TreeNode{
				{ID: "shell", Label: "shell", Icon: theme.Icons.Pending},
				{ID: "editor", Label: "editor", Icon: theme.Icons.Success},
				{
					ID:       "terminal",
					Label:    "terminal",
					Expanded: true,
					Children: []components.TreeNode{
						{ID: "kitty", Label: "kitty.conf"},
						{ID: "tmux", Label: "tmux.conf"},
					},
				},
			},
		},
		{ID: "profiles", Label: "profiles", Icon: theme.Icons.Info},
	})

	code := components.NewCode(t, `# config.yaml
source:
  type: git
  url: github.com/user/dotfiles
machine: workstation`).SetLineNumbers(true)

	col2 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Tree"),
		"",
		tree.View(),
		"",
		t.S().Subtitle.Render("Code"),
		"",
		code.View(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
}

func (d *Dashboard) renderLayoutTab() string {
	t := d.theme

	card := components.NewCard(t).
		SetTitle("System Info").
		SetContent("This is a card component with title, content, and footer.").
		SetFooter("Last updated: just now").
		SetWidth(40)

	focusedCard := components.NewCard(t).
		SetTitle("Focused Card").
		SetContent("Cards can have a focused state with highlighted border.").
		SetWidth(40).
		SetFocused(true)

	empty := components.NewEmpty(t).
		SetIcon("📦").
		SetTitle("No packages found").
		SetMessage("Install some packages to see them here").
		SetAction("Press 'i' to install")

	col1 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Cards"),
		"",
		card.View(),
		"",
		focusedCard.View(),
	)

	col2 := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Subtitle.Render("Empty State"),
		"",
		empty.View(),
		"",
		t.S().Subtitle.Render("Separators"),
		"",
		components.NewSeparator(t).SetLength(30).View(),
		"",
		components.NewSeparator(t).SetLength(30).SetLabel("Section").View(),
		"",
		components.NewSeparator(t).SetLength(30).SetChar("═").View(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
}
