package tui

import (
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/internal/tui/app"
	"github.com/arthur404dev/dotts/internal/tui/pages/dashboard"
	"github.com/arthur404dev/dotts/internal/tui/pages/doctor"
	"github.com/arthur404dev/dotts/internal/tui/pages/settings"
	"github.com/arthur404dev/dotts/internal/tui/pages/status"
	"github.com/arthur404dev/dotts/internal/tui/pages/update"
	"github.com/arthur404dev/dotts/internal/tui/pages/wizard"
	"github.com/arthur404dev/dotts/pkg/vetru"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/palette"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	tea "github.com/charmbracelet/bubbletea"
)

func dottsCommands() []palette.Command {
	return []palette.Command{
		// Navigation
		{
			ID:       "nav-dashboard",
			Label:    "Dashboard",
			Category: palette.CategoryNavigation,
			Keywords: []string{"home", "main"},
			Handler: func() tea.Cmd {
				return messages.Navigate(app.PageDashboard)
			},
		},
		{
			ID:       "nav-status",
			Label:    "Status",
			Category: palette.CategoryNavigation,
			Keywords: []string{"check", "state"},
			Handler: func() tea.Cmd {
				return messages.Navigate(app.PageStatus)
			},
		},
		{
			ID:       "nav-update",
			Label:    "Update",
			Category: palette.CategoryNavigation,
			Keywords: []string{"upgrade", "sync"},
			Handler: func() tea.Cmd {
				return messages.Navigate(app.PageUpdate)
			},
		},
		{
			ID:       "nav-doctor",
			Label:    "Doctor",
			Category: palette.CategoryNavigation,
			Keywords: []string{"health", "diagnose", "check"},
			Handler: func() tea.Cmd {
				return messages.Navigate(app.PageDoctor)
			},
		},
		{
			ID:       "nav-settings",
			Label:    "Settings",
			Category: palette.CategoryNavigation,
			Keywords: []string{"config", "preferences", "options"},
			Handler: func() tea.Cmd {
				return messages.Navigate(app.PageSettings)
			},
		},

		// Actions
		{
			ID:       "action-init",
			Label:    "Run Init Wizard",
			Category: palette.CategoryActions,
			Keywords: []string{"setup", "bootstrap", "configure", "wizard"},
			Handler: func() tea.Cmd {
				return messages.StartWizard()
			},
		},
		{
			ID:       "action-sync",
			Label:    "Sync Dotfiles",
			Category: palette.CategoryActions,
			Keywords: []string{"push", "pull", "update"},
			Handler: func() tea.Cmd {
				return messages.Action(app.ActionSync)
			},
		},

		// System
		{
			ID:       "sys-quit",
			Label:    "Quit",
			Category: palette.CategorySystem,
			Shortcut: "ctrl+c",
			Keywords: []string{"exit", "close"},
			Handler: func() tea.Cmd {
				return tea.Quit
			},
		},
	}
}

const Version = "0.1.0"

type systemDetectedMsg struct {
	info *system.SystemInfo
	err  error
}

func Run() error {
	t := theme.DefaultTheme()

	pageMap := map[messages.PageID]vetru.Page{
		app.PageDashboard: dashboard.New(t),
		app.PageStatus:    status.New(t),
		app.PageUpdate:    update.New(t),
		app.PageDoctor:    doctor.New(t),
		app.PageSettings:  settings.New(t),
		app.PageWizard:    wizard.New(t),
	}

	model := vetru.NewModel(vetru.Config{
		Brand:        "dotts",
		Version:      Version,
		Theme:        t,
		Pages:        pageMap,
		Commands:     dottsCommands(),
		DefaultPage:  app.PageDashboard,
		WizardPageID: app.PageWizard,
		OnInit:       detectSystem,
	})

	p := tea.NewProgram(
		&appModel{Model: model},
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	_, err := p.Run()
	return err
}

type appModel struct {
	*vetru.Model
}

func (a *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case systemDetectedMsg:
		if msg.err == nil && msg.info != nil {
			a.SetMetadata("os", string(msg.info.OS))
			if msg.info.Distro != "" && msg.info.Distro != "unknown" {
				a.SetMetadata("distro", string(msg.info.Distro))
			}
			a.SetMetadata("arch", string(msg.info.Arch))
		}
		return a, nil
	}

	m, cmd := a.Model.Update(msg)
	a.Model = m.(*vetru.Model)
	return a, cmd
}

func detectSystem() tea.Msg {
	info, err := system.Detect()
	if err != nil {
		return systemDetectedMsg{err: err}
	}
	return systemDetectedMsg{info: info}
}
