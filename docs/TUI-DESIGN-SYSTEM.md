# dotts TUI Design System Plan

This document outlines the design system plan for the `dotts` TUI. Inspired by modern component libraries like shadcn/ui, this system aims to provide a robust, theme-aware, and highly composable set of building blocks for creating a polished and consistent terminal user interface.

## 0. Framework Architecture (Vetru)

The TUI is designed as a reusable framework called **Vetru** ("Next.js for TUIs"). The architecture separates framework code from application-specific code.

### Package Structure

```
internal/tui/
├── model.go              # Framework: Generic shell with Config-based DI
├── keymap.go             # Framework: Default key bindings
├── app.go                # App: dotts entry point, pages, commands, system detection
├── app/
│   └── pages.go          # App: PageID/ActionID constants (PageDashboard, PageWizard, etc.)
│
├── components/           # Framework: 35+ reusable components
├── theme/                # Framework: Theme system (colors, spacing, icons)
├── keys/                 # Framework: Centralized key binding definitions
├── layout/               # Framework: Header, footer, overlay compositing
├── messages/             # Framework: Generic message types
├── palette/              # Framework: Command palette component
├── progress/             # Framework: Progress tracking
├── styles/               # Framework: Style utilities
│
├── pages/                # App: Page implementations
│   ├── page.go           # Framework: Page interface definition
│   ├── dashboard/        # App: Dashboard page
│   ├── status/           # App: Status page
│   └── ...
└── wizard/               # App: Wizard step implementations
```

### Config-Based Dependency Injection

The framework shell (`model.go`) accepts configuration via a `Config` struct:

```go
type Config struct {
    Brand        string                         // App name (e.g., "dotts")
    Version      string                         // App version
    Theme        *theme.Theme                   // Optional custom theme
    Pages        map[messages.PageID]pages.Page // Page implementations
    Commands     []palette.Command              // Command palette commands
    DefaultPage  messages.PageID                // Starting page
    WizardPageID messages.PageID                // Wizard page (optional)
    OnInit       func() tea.Msg                 // Initialization hook
}
```

### Separation Rules

| Package Type | Can Import | Cannot Import |
|--------------|------------|---------------|
| Framework (`components/`, `theme/`, `layout/`, etc.) | Other framework packages, external libs | `internal/system`, `internal/config`, app-specific code |
| App (`app.go`, `app/`, `pages/*/`, `wizard/`) | Framework packages, dotts business logic | N/A |

### Future: Standalone Vetru Package

The framework code can be extracted to a standalone repository:

```bash
go get github.com/yourusername/vetru
```

With CLI tooling for scaffolding:

```bash
vetru init my-tui-app
vetru add component spinner
vetru add page settings
```

## 1. Executive Summary

**Goal**: Create a complete TUI design system inspired by shadcn/ui.

**Two Main Objectives**:
1.  **Fix Overlay Behavior**: Implement true compositing to allow overlays (like the command palette) to render over background content without blanking it.
2.  **Create Comprehensive Component Library**: Develop a complete set of "building block" components that are theme-aware, interactive, and highly reusable.

## 2. Overlay Fix (Task 1)

The command palette currently blanks the background using `lipgloss.Place()`. We will implement true compositing using a `PlaceOverlay` function in `internal/tui/layout/compose.go`.

### Compositing Algorithm

The `PlaceOverlay` function renders `fg` on top of `bg` at position `(x, y)`:

1.  **Split** both `bg` and `fg` into lines.
2.  **For each line** in the overlay region:
    a.  **Truncate** background to get the left portion.
    b.  **Append** foreground line.
    c.  **TruncateLeft** background to get right portion.
    d.  **Append** right portion.
3.  **Lines outside** overlay region remain unchanged.

```go
// PlaceOverlay renders fg on top of bg at position (x, y)
// 1. Split both into lines
// 2. For each line in the overlay region:
//    a. Truncate background to get the left portion
//    b. Append foreground line
//    c. TruncateLeft background to get right portion
//    d. Append right portion
// 3. Lines outside overlay region remain unchanged
```

The system will also support an optional dimming of the background using `lipgloss.Faint()` to enhance focus on the overlay.

## 3. Theme Enhancements

We will introduce standardized spacing to ensure visual consistency across all pages and components.

### Spacing Constants (`internal/tui/theme/spacing.go`)

| Constant | Value | Description |
| :--- | :--- | :--- |
| `None` | 0 | No spacing |
| `XS` | 1 | 1 cell / line |
| `SM` | 2 | 2 cells / lines |
| `MD` | 3 | 3 cells / lines |
| `LG` | 4 | 4 cells / lines |
| `XL` | 6 | 6 cells / lines |

## 4. Component Library (Flat Structure in `components/`)

All components will reside in `internal/tui/components/` and follow these principles:
- **Theme-aware**: Accept `*theme.Theme` in constructor.
- **Interactive**: Implement `tea.Model` where interactive.
- **Accessible**: Support mouse via `bubblezone` where appropriate.

### Primitives
- **Text**: Variants: title, subtitle, body, caption, code, muted.
- **Separator**: Horizontal/vertical dividers with optional label.
- **Spacer**: Flexible spacing (renders empty lines/spaces).
- **Badge**: Status pills (success, warning, error, info, default).

### Layout
- **Box**: Container with padding, border options (none, normal, rounded, double).
- **Card**: Elevated container with optional title/footer.
- **Stack**: `VStack`/`HStack` with gap control and alignment.
- **Center**: Centering utility for width/height.

### Feedback
- **Spinner**: Wrapper around `bubbles/spinner`, theme-aware.
- **ProgressBar**: Determinate progress with gradient, percentage, animation (uses `harmonica`).
- **Alert**: Inline alerts (info, success, warning, error) with icon and message.
- **Skeleton**: Loading placeholder with animation.
- **Toast**: Temporary notification (auto-dismiss).

### Input
- **TextInput**: Enhanced existing component.
- **Select**: Enhanced existing selection list.
- **Checkbox**: Single checkbox with label.
- **Toggle**: Toggle switch with on/off labels.
- **Confirm**: Yes/No confirmation component.

### Navigation
- **Stepper**: Horizontal wizard progress (Fixes the current cramped one).
  - Shows: `→ Source  ──  ○ Machine  ──  ○ Personal  ──  ○ Settings`
  - Proper gaps, connectors, clear visual states.
- **Tabs**: Tab bar with indicators.
- **Breadcrumb**: Navigation trail.
- **Menu**: Vertical menu/nav list.

### Data Display
- **List**: Wrapper around `bubbles/list`, theme-aware, filterable.
- **Table**: Simple data table with headers, optional borders/striping.
- **KeyValue**: Key-value pair display with aligned labels.
- **Tree**: Tree view for hierarchical data.

### Overlay
- **Modal**: Modal dialog using the new compositing.
- **Dialog**: Pre-built dialogs (confirm, alert, error).
- **Tooltip**: Hover tooltip (if feasible in TUI).

### Display
- **Logo**: Application branding (existing).
- **Banner**: Promotional/informational banner (existing).
- **Empty**: Empty state placeholder with icon, message, action hint.
- **Code**: Code block with optional line numbers.

## 5. Dependencies to Add
- `github.com/charmbracelet/harmonica` - For smooth, physics-based animations (especially in Progress Bar).

## 6. Implementation Tasks

1.  Implement overlay compositing (`layout/compose.go`)
2.  Fix command palette overlay in `model.go`
3.  Add spacing constants (`theme/spacing.go`)
4.  Create Stepper component
5.  Update wizard to use Stepper
6.  Create Box component
7.  Create Card component
8.  Create Stack component
9.  Create Center component
10. Create Text component
11. Create Separator component
12. Create Spacer component
13. Create Badge component
14. Create Spinner component
15. Create ProgressBar component
16. Create Alert component
17. Create Skeleton component
18. Create Toast component
19. Create Checkbox component
20. Create Toggle component
21. Create Confirm component
22. Create Tabs component
23. Create Breadcrumb component
24. Create Menu component
25. Create List component (wrapper)
26. Create Table component
27. Create KeyValue component
28. Create Tree component
29. Create Modal component
30. Create Dialog component
31. Create Empty component
32. Create Code component
33. Refactor wizard pages to use new components
34. Add harmonica dependency

## 7. Component Interface Patterns

### Static Component Pattern
```go
// Static component pattern
type Badge struct {
    theme   *theme.Theme
    variant BadgeVariant
    label   string
}

func NewBadge(t *theme.Theme, variant BadgeVariant, label string) *Badge
func (b *Badge) View() string
```

### Interactive Component Pattern (implements `tea.Model`)
```go
// Interactive component pattern (implements tea.Model)
type Checkbox struct {
    theme   *theme.Theme
    label   string
    checked bool
    focused bool
    zoneID  string
}

func NewCheckbox(t *theme.Theme, label string) *Checkbox
func (c *Checkbox) Init() tea.Cmd
func (c *Checkbox) Update(msg tea.Msg) (*Checkbox, tea.Cmd)
func (c *Checkbox) View() string
func (c *Checkbox) Focus() tea.Cmd
func (c *Checkbox) Blur()
func (c *Checkbox) Checked() bool
func (c *Checkbox) SetChecked(checked bool)
```

## 8. Visual Examples

### Stepper (Before/After)

**Before (Cramped)**:
`→ Source  ○ Machine  ○ Personal  ○ Settings`

**After (Spacious)**:
`→ Source  ───────  ○ Machine  ───────  ○ Personal  ───────  ○ Settings`

### Card Layout
```
┌─ System Info ──────────────────────────────────────────────────┐
│                                                                │
│  OS:       Arch Linux                                          │
│  Arch:     x86_64                                              │
│  Host:     workstation                                         │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

### Alert Variants
```
[!] Error: Failed to clone repository
[i] Info: New version of dotts available
[v] Success: Configuration applied successfully
```

## 9. File Structure

### `internal/tui/components/` (Flat)
```
internal/tui/components/
├── alert.go
├── badge.go
├── banner.go          # existing
├── box.go
├── breadcrumb.go
├── card.go
├── center.go
├── checkbox.go
├── code.go
├── confirm.go
├── dialog.go
├── empty.go
├── help.go            # existing
├── input.go           # existing
├── keyvalue.go
├── list.go
├── logo.go            # existing (in logo/)
├── menu.go
├── modal.go
├── nav.go             # existing
├── progressbar.go     # existing (rename from progress.go)
├── selectlist.go      # existing
├── separator.go
├── skeleton.go
├── spacer.go
├── spinner.go
├── stack.go
├── stepper.go
├── sysinfo.go         # existing
├── table.go
├── tabs.go
├── text.go
├── toast.go
├── toggle.go
├── tooltip.go
└── tree.go
```

### `internal/tui/theme/`
```
internal/tui/theme/
├── gradient.go        # existing
├── icons.go           # existing
├── spacing.go         # NEW
├── styles.go          # existing
└── theme.go           # existing
```

### `internal/tui/layout/`
```
internal/tui/layout/
├── compose.go         # NEW - overlay compositing
├── footer.go          # existing
├── header.go          # existing
└── overlay.go         # existing
```

### `internal/tui/keys/`
```
internal/tui/keys/
└── keys.go            # NEW - centralized key bindings
```

## 10. Navigation Standards

All key bindings are defined in `internal/tui/keys/keys.go` as the single source of truth.

### Global Keys (Always Available)

| Key | Action |
|-----|--------|
| `ctrl+p` / `/` | Open command palette |
| `ctrl+c` | Force quit |
| `esc` | Close overlay / Go back |
| `ctrl+g` / `?` | Toggle full help |
| `q` | Quit (context-dependent) |

### Vertical Lists (Select, Menu, Palette)

| Key | Action |
|-----|--------|
| `↑` / `k` | Previous item (wraps) |
| `↓` / `j` | Next item (wraps) |
| `ctrl+p` | Previous item (in text inputs) |
| `ctrl+n` | Next item (in text inputs) |
| `home` | First item |
| `end` | Last item |
| `enter` | Select / Confirm |

### Horizontal Navigation (Tabs, Wizard Steps)

| Key | Action |
|-----|--------|
| `←` / `h` | Previous tab/step |
| `→` / `l` | Next tab/step |

### Form Navigation (Input Fields)

| Key | Action |
|-----|--------|
| `tab` | Next field |
| `shift+tab` | Previous field |
| `enter` | Next field / Submit |

### Design Principles

1. **Consistency**: Same keys do the same thing everywhere
2. **Vim-friendly**: `hjkl` navigation supported where appropriate
3. **Wrap-around**: Lists wrap from top to bottom and vice versa
4. **No conflicts**: `tab`/`shift+tab` reserved for forms, not tabs/lists
5. **Global access**: Command palette opens anywhere with `ctrl+p`
