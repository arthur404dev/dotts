package tui

import (
	"fmt"

	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/internal/tui/layout"
	"github.com/arthur404dev/dotts/internal/tui/messages"
	"github.com/arthur404dev/dotts/internal/tui/pages"
	"github.com/arthur404dev/dotts/internal/tui/pages/dashboard"
	"github.com/arthur404dev/dotts/internal/tui/pages/doctor"
	"github.com/arthur404dev/dotts/internal/tui/pages/settings"
	"github.com/arthur404dev/dotts/internal/tui/pages/status"
	"github.com/arthur404dev/dotts/internal/tui/pages/update"
	"github.com/arthur404dev/dotts/internal/tui/pages/wizard"
	"github.com/arthur404dev/dotts/internal/tui/palette"
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const Version = "0.1.0"

type AppMode int

const (
	ModeNormal AppMode = iota
	ModeWizard
)

// Model is the root TUI model
type Model struct {
	// Dimensions
	width  int
	height int

	// Theme and styling
	theme *theme.Theme

	// Key bindings
	keymap KeyMap

	// Layout components
	header *layout.Header
	footer *layout.Footer

	// Navigation
	activePageID messages.PageID
	pages        map[messages.PageID]pages.Page

	// Command palette
	palette     *palette.Palette
	paletteOpen bool

	// System info
	sysInfo *system.SystemInfo

	mode AppMode

	// State
	ready    bool
	quitting bool
}

// NewModel creates a new root model
func NewModel() *Model {
	zone.NewGlobal()

	t := theme.DefaultTheme()
	theme.SetCurrent(t)

	km := DefaultKeyMap()

	// Create layout components
	header := layout.NewHeader(t, Version)
	footer := layout.NewFooter(t)
	footer.SetKeyMap(km)

	pageMap := map[messages.PageID]pages.Page{
		messages.PageDashboard: dashboard.New(t),
		messages.PageStatus:    status.New(t),
		messages.PageUpdate:    update.New(t),
		messages.PageDoctor:    doctor.New(t),
		messages.PageSettings:  settings.New(t),
		messages.PageWizard:    wizard.New(t),
	}

	// Create palette
	pal := palette.New(t, palette.DefaultCommands())

	return &Model{
		theme:        t,
		keymap:       km,
		header:       header,
		footer:       footer,
		activePageID: messages.PageDashboard,
		pages:        pageMap,
		palette:      pal,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.detectSystem,
		tea.EnterAltScreen,
	)
}

func (m *Model) detectSystem() tea.Msg {
	info, err := system.Detect()
	if err != nil {
		return systemDetectedMsg{err: err}
	}
	return systemDetectedMsg{info: info}
}

type systemDetectedMsg struct {
	info *system.SystemInfo
	err  error
}

// Update handles events
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true
		return m, nil

	case systemDetectedMsg:
		if msg.err == nil {
			m.sysInfo = msg.info
			m.updateSysInfo()
		}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.ForceQuit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.OpenPalette):
			if !m.paletteOpen {
				m.paletteOpen = true
				return m, m.palette.Focus()
			}

		case key.Matches(msg, m.keymap.Escape):
			if m.paletteOpen {
				m.paletteOpen = false
				return m, nil
			}

		case key.Matches(msg, m.keymap.Quit):
			if m.mode == ModeWizard {
				return m, nil
			}
			if !m.paletteOpen {
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, m.keymap.Help):
			if m.mode != ModeWizard {
				m.footer.ToggleFullHelp()
			}
			return m, nil
		}

	case messages.StartWizardMsg:
		return m.startWizard()

	case messages.WizardCompleteMsg:
		return m.endWizard(msg.Success)

	case messages.NavigateMsg:
		if m.mode == ModeWizard {
			return m, nil
		}
		return m.navigate(msg.Page)

	case messages.ClosePaletteMsg:
		m.paletteOpen = false
		return m, nil

	case messages.OpenPaletteMsg:
		if m.mode == ModeWizard {
			return m, nil
		}
		m.paletteOpen = true
		return m, m.palette.Focus()
	}

	// Route to palette if open
	if m.paletteOpen {
		var cmd tea.Cmd
		m.palette, cmd = m.palette.Update(msg)
		return m, cmd
	}

	// Route to current page
	return m.updateCurrentPage(msg)
}

func (m *Model) navigate(pageID messages.PageID) (*Model, tea.Cmd) {
	if page, ok := m.pages[m.activePageID]; ok {
		page.Blur()
	}

	m.activePageID = pageID
	m.paletteOpen = false

	if page, ok := m.pages[pageID]; ok {
		return m, page.Focus()
	}

	return m, nil
}

func (m *Model) startWizard() (*Model, tea.Cmd) {
	m.mode = ModeWizard
	m.paletteOpen = false
	m.activePageID = messages.PageWizard

	if page, ok := m.pages[messages.PageWizard]; ok {
		return m, page.Focus()
	}

	return m, nil
}

func (m *Model) endWizard(success bool) (*Model, tea.Cmd) {
	m.mode = ModeNormal
	m.activePageID = messages.PageDashboard

	if page, ok := m.pages[messages.PageDashboard]; ok {
		return m, page.Focus()
	}

	return m, nil
}

func (m *Model) updateCurrentPage(msg tea.Msg) (*Model, tea.Cmd) {
	page, ok := m.pages[m.activePageID]
	if !ok {
		return m, nil
	}

	updated, cmd := page.Update(msg)
	m.pages[m.activePageID] = updated.(pages.Page)

	return m, cmd
}

func (m *Model) updateLayout() {
	headerHeight := 1
	footerHeight := 1
	contentHeight := m.height - headerHeight - footerHeight - 2 // padding

	m.header.SetWidth(m.width)
	m.footer.SetWidth(m.width)
	m.palette.SetWidth(min(60, m.width-4))

	// Update all pages
	for _, page := range m.pages {
		page.SetSize(m.width-4, contentHeight)
	}
}

func (m *Model) updateSysInfo() {
	if m.sysInfo == nil {
		return
	}

	info := fmt.Sprintf("%s %s %s %s",
		theme.Icons.Dot,
		m.sysInfo.OS,
		theme.Icons.Dot,
		m.sysInfo.Arch,
	)
	if m.sysInfo.Distro != "" && m.sysInfo.Distro != "unknown" {
		info = fmt.Sprintf("%s %s %s %s %s %s",
			theme.Icons.Dot,
			m.sysInfo.OS,
			theme.Icons.Dot,
			m.sysInfo.Distro,
			theme.Icons.Dot,
			m.sysInfo.Arch,
		)
	}
	m.header.SetSysInfo(info)
}

// View renders the model
func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Header
	header := m.header.View()

	// Content (current page)
	page := m.pages[m.activePageID]
	content := page.View()

	// Footer
	footer := m.footer.View()

	// Compose main view
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)

	// Overlay palette if open
	if m.paletteOpen {
		mainView = m.overlayPalette(mainView)
	}

	return zone.Scan(mainView)
}

func (m *Model) overlayPalette(base string) string {
	paletteView := m.palette.View()

	// Dim the background and overlay the palette
	dimmed := layout.DimBackground(base)
	return layout.CenterOverlay(paletteView, dimmed, m.width, m.height)
}
