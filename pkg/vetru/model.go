package vetru

import (
	"github.com/arthur404dev/dotts/pkg/vetru/layout"
	"github.com/arthur404dev/dotts/pkg/vetru/messages"
	"github.com/arthur404dev/dotts/pkg/vetru/palette"
	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type AppMode int

const (
	ModeNormal AppMode = iota
	ModeWizard
)

type Config struct {
	Brand        string
	Version      string
	Theme        *theme.Theme
	Pages        map[messages.PageID]Page
	Commands     []palette.Command
	DefaultPage  messages.PageID
	WizardPageID messages.PageID
	OnInit       func() tea.Msg
}

type Model struct {
	width  int
	height int

	config Config
	theme  *theme.Theme
	keymap KeyMap

	header *layout.Header
	footer *layout.Footer

	activePageID messages.PageID
	pages        map[messages.PageID]Page

	palette     *palette.Palette
	paletteOpen bool

	metadata map[string]string

	mode AppMode

	ready    bool
	quitting bool
}

func NewModel(cfg Config) *Model {
	zone.NewGlobal()

	t := cfg.Theme
	if t == nil {
		t = theme.DefaultTheme()
	}
	theme.SetCurrent(t)

	km := DefaultKeyMap()

	header := layout.NewHeader(t, cfg.Brand, cfg.Version)
	footer := layout.NewFooter(t)
	footer.SetKeyMap(km)

	commands := cfg.Commands
	if len(commands) == 0 {
		commands = []palette.Command{}
	}
	pal := palette.New(t, commands)

	defaultPage := cfg.DefaultPage
	if defaultPage == "" && len(cfg.Pages) > 0 {
		for id := range cfg.Pages {
			defaultPage = id
			break
		}
	}

	return &Model{
		config:       cfg,
		theme:        t,
		keymap:       km,
		header:       header,
		footer:       footer,
		activePageID: defaultPage,
		pages:        cfg.Pages,
		palette:      pal,
		metadata:     make(map[string]string),
	}
}

func (m *Model) SetMetadata(key, value string) {
	m.metadata[key] = value
	m.updateHeaderInfo()
}

func (m *Model) updateHeaderInfo() {
	if len(m.metadata) == 0 {
		return
	}

	var parts []string
	for _, v := range m.metadata {
		if v != "" {
			parts = append(parts, v)
		}
	}

	if len(parts) > 0 {
		info := ""
		for i, p := range parts {
			if i > 0 {
				info += " " + theme.Icons.Dot + " "
			}
			info += p
		}
		m.header.SetSysInfo(theme.Icons.Dot + " " + info)
	}
}

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tea.EnterAltScreen)

	if m.config.OnInit != nil {
		cmds = append(cmds, func() tea.Msg { return m.config.OnInit() })
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true
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

	if m.paletteOpen {
		var cmd tea.Cmd
		m.palette, cmd = m.palette.Update(msg)
		return m, cmd
	}

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

	wizardID := m.config.WizardPageID
	if wizardID == "" {
		return m, nil
	}

	m.activePageID = wizardID

	if page, ok := m.pages[wizardID]; ok {
		return m, page.Focus()
	}

	return m, nil
}

func (m *Model) endWizard(success bool) (*Model, tea.Cmd) {
	m.mode = ModeNormal
	m.activePageID = m.config.DefaultPage

	if page, ok := m.pages[m.config.DefaultPage]; ok {
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
	m.pages[m.activePageID] = updated.(Page)

	return m, cmd
}

func (m *Model) updateLayout() {
	headerHeight := 1
	footerHeight := 1
	contentHeight := m.height - headerHeight - footerHeight - 2

	m.header.SetWidth(m.width)
	m.footer.SetWidth(m.width)
	m.palette.SetWidth(min(60, m.width-4))

	for _, page := range m.pages {
		page.SetSize(m.width-4, contentHeight)
	}
}

func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	header := m.header.View()

	page := m.pages[m.activePageID]
	content := page.View()

	footer := m.footer.View()

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)

	if m.paletteOpen {
		mainView = m.overlayPalette(mainView)
	}

	return zone.Scan(mainView)
}

func (m *Model) overlayPalette(base string) string {
	paletteView := m.palette.View()
	dimmed := layout.DimBackground(base)
	return layout.CenterOverlay(paletteView, dimmed, m.width, m.height)
}
