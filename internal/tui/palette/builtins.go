package palette

import (
	"github.com/arthur404dev/dotts/internal/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

// DefaultCommands returns the default command set
func DefaultCommands() []Command {
	return []Command{
		// Navigation
		{
			ID:       "nav-dashboard",
			Label:    "Dashboard",
			Category: CategoryNavigation,
			Keywords: []string{"home", "main"},
			Handler: func() tea.Cmd {
				return messages.Navigate(messages.PageDashboard)
			},
		},
		{
			ID:       "nav-status",
			Label:    "Status",
			Category: CategoryNavigation,
			Keywords: []string{"check", "state"},
			Handler: func() tea.Cmd {
				return messages.Navigate(messages.PageStatus)
			},
		},
		{
			ID:       "nav-update",
			Label:    "Update",
			Category: CategoryNavigation,
			Keywords: []string{"upgrade", "sync"},
			Handler: func() tea.Cmd {
				return messages.Navigate(messages.PageUpdate)
			},
		},
		{
			ID:       "nav-doctor",
			Label:    "Doctor",
			Category: CategoryNavigation,
			Keywords: []string{"health", "diagnose", "check"},
			Handler: func() tea.Cmd {
				return messages.Navigate(messages.PageDoctor)
			},
		},
		{
			ID:       "nav-settings",
			Label:    "Settings",
			Category: CategoryNavigation,
			Keywords: []string{"config", "preferences", "options"},
			Handler: func() tea.Cmd {
				return messages.Navigate(messages.PageSettings)
			},
		},

		// Actions
		{
			ID:       "action-init",
			Label:    "Run Init Wizard",
			Category: CategoryActions,
			Keywords: []string{"setup", "bootstrap", "configure", "wizard"},
			Handler: func() tea.Cmd {
				return messages.StartWizard()
			},
		},
		{
			ID:       "action-sync",
			Label:    "Sync Dotfiles",
			Category: CategoryActions,
			Keywords: []string{"push", "pull", "update"},
			Handler: func() tea.Cmd {
				return messages.Action(messages.ActionSync)
			},
		},

		// System
		{
			ID:       "sys-quit",
			Label:    "Quit",
			Category: CategorySystem,
			Shortcut: "ctrl+c",
			Keywords: []string{"exit", "close"},
			Handler: func() tea.Cmd {
				return tea.Quit
			},
		},
	}
}
