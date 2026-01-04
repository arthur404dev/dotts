package theme

// Icons defines unicode symbols used throughout the TUI
var Icons = struct {
	// Navigation
	Diagonal     string
	ArrowLeft    string
	ArrowRight   string
	ArrowDown    string
	Chevron      string
	ChevronLeft  string
	ChevronRight string
	ChevronDown  string

	// Status
	Success string
	Warning string
	Error   string
	Info    string
	Pending string
	Current string

	// Progress
	StepComplete string
	StepCurrent  string
	StepPending  string

	// UI
	Prompt   string
	Cursor   string
	Selected string
	Bullet   string
	Dot      string

	// Files
	Folder string
	File   string
}{
	// Navigation
	Diagonal:     "╱",
	ArrowLeft:    "←",
	ArrowRight:   "→",
	ArrowDown:    "↓",
	Chevron:      "›",
	ChevronLeft:  "‹",
	ChevronRight: "›",
	ChevronDown:  "⌄",

	// Status
	Success: "✓",
	Warning: "!",
	Error:   "✗",
	Info:    "i",
	Pending: "○",
	Current: "●",

	// Progress
	StepComplete: "✓",
	StepCurrent:  "●",
	StepPending:  "○",

	// UI
	Prompt:   ">",
	Cursor:   "▌",
	Selected: "▌",
	Bullet:   "•",
	Dot:      "·",

	// Files
	Folder: "📁",
	File:   "📄",
}
