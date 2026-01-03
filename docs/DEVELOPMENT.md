# Development Guide

> How to build, test, and continue development on dotts

## Project Structure

```
dotts/
├── cmd/dotts/              # CLI entry point
│   ├── main.go             # Main function
│   └── cmd/                # Cobra commands
│       ├── root.go         # Base command, flags
│       ├── init.go         # Bootstrap wizard
│       ├── update.go       # Update command (stub)
│       ├── status.go       # Status command (stub)
│       ├── doctor.go       # Health check (stub)
│       ├── config.go       # Config source management (stub)
│       ├── machine.go      # Machine management (stub)
│       └── sync.go         # Sync changes (stub)
│
├── internal/               # Private packages
│   ├── system/             # System detection
│   │   └── detect.go       # OS, distro, arch, pkg manager
│   ├── state/              # State management
│   │   ├── paths.go        # XDG paths
│   │   └── state.go        # JSON state persistence
│   ├── config/             # Configuration loading
│   │   ├── source.go       # Git source management
│   │   ├── loader.go       # YAML file loading
│   │   ├── resolver.go     # Profile inheritance resolution
│   │   └── alternates.go   # Alternate file matching
│   └── tui/                # Terminal UI
│       ├── styles/         # Lip Gloss theme
│       │   └── theme.go    # Catppuccin colors
│       └── wizard/         # Bootstrap wizard
│           ├── wizard.go   # Main wizard orchestration
│           ├── source.go   # Source selection
│           ├── machine.go  # Machine configuration
│           ├── settings.go # Settings input
│           ├── features.go # Feature selection
│           └── auth.go     # Authentication setup
│
├── pkg/schema/             # Public types
│   ├── profile.go          # Profile YAML schema
│   ├── machine.go          # Machine YAML schema
│   ├── packages.go         # Package manifest schema
│   └── config.go           # Repo config schema
│
├── templates/              # Embedded templates
│   ├── config-repo/        # Fork template
│   │   ├── config.yaml     # Repo metadata
│   │   ├── profiles/       # Profile definitions
│   │   ├── packages/       # Package manifests
│   │   ├── machines/       # Machine examples
│   │   └── configs/        # Actual config files
│   └── nix/                # Home Manager templates (empty)
│
├── scripts/                # Utility scripts
│   └── install.sh          # curl|sh installer
│
├── docs/                   # Documentation
│   ├── README.md           # Documentation index
│   ├── ARCHITECTURE.md     # Technical design
│   ├── DEVELOPMENT.md      # This file
│   ├── CONFIG-SCHEMA.md    # YAML reference
│   └── SESSION-HISTORY.md  # Development context
│
├── Makefile                # Build commands
├── go.mod                  # Go modules
├── go.sum                  # Dependency checksums
└── .gitignore              # Git ignore patterns
```

## Building

```bash
# Download dependencies
make deps

# Build binary
make build

# Run
./dotts --help

# Build for all platforms
make build-all
```

## Current Implementation Status

### Completed

| Component | Status | Notes |
|-----------|--------|-------|
| CLI skeleton | Done | All commands registered, flags defined |
| System detection | Done | OS, distro, arch, package manager |
| State management | Done | XDG paths, JSON persistence |
| Config loading | Done | YAML parsing, profile/machine/package loading |
| Profile resolution | Done | Inheritance chain, merging |
| Alternate files | Done | Scoring, matching, selection |
| TUI styles | Done | Catppuccin theme, Lip Gloss components |
| TUI wizard | Done | All 5 steps (source, machine, settings, features, auth) |
| Templates | Done | Full config-repo template with starter configs |
| Install script | Done | Curl-installable bootstrap |

### Not Yet Implemented

| Component | Priority | Description |
|-----------|----------|-------------|
| Package installer | High | Actually install packages via nix/brew/pacman |
| Dotfile linker | High | Create symlinks or generate Home Manager config |
| Home Manager setup | High | Install and configure Home Manager |
| Update command | Medium | Pull changes, show diff, apply |
| Status command | Medium | Show current state, pending changes |
| Doctor command | Medium | Validate setup, check dependencies |
| Sync command | Low | Push local changes back to config repo |
| GoReleaser config | Medium | Multi-platform release automation |
| Tests | Medium | Unit and integration tests |

## Next Steps (Priority Order)

### 1. Package Installer (`internal/installer/`)

Create installers for each package manager:

```go
// internal/installer/installer.go
type Installer interface {
    Install(packages []string) error
    IsInstalled(pkg string) bool
    Update() error
}

// Implementations:
// - NixInstaller
// - BrewInstaller  
// - PacmanInstaller (with AUR support via yay/paru)
// - AptInstaller
```

### 2. Dotfile Linker (`internal/linker/`)

Options:
- **Simple symlinks**: Direct symlinks from `~/.config/X` to config repo
- **Home Manager**: Generate `home.nix` for declarative management

Start with symlinks, add Home Manager as optional backend.

### 3. Progress Display (`internal/tui/progress/`)

Use Bubble Tea for animated progress:
- Package installation progress
- File linking status
- Overall completion

### 4. Command Implementations

Wire up the stub commands to real logic:

```go
// cmd/dotts/cmd/update.go
func runUpdate(cmd *cobra.Command, args []string) error {
    // 1. Pull config repo changes
    // 2. Show diff of what changed
    // 3. Prompt for confirmation
    // 4. Apply changes (reinstall packages, relink files)
}
```

### 5. Testing

Add tests for critical paths:
- Profile resolution
- Alternate file matching
- Package manifest merging

## Development Commands

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Install locally
make install-local  # to ~/.local/bin
make install        # to $GOPATH/bin
```

## Adding a New Command

1. Create file `cmd/dotts/cmd/mycommand.go`
2. Define command with Cobra
3. Register in `root.go`'s init function
4. Implement logic in `internal/`

```go
// cmd/dotts/cmd/mycommand.go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    RunE:  runMyCommand,
}

func init() {
    // Add flags here
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/charmbracelet/huh` | v0.8.0 | Form inputs |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | Styling |
| `github.com/charmbracelet/bubbletea` | v1.3.6 | TUI framework |
| `github.com/catppuccin/go` | v0.3.0 | Color palette |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing |

## Release Process (Future)

1. Create `.goreleaser.yaml` for multi-platform builds
2. Set up GitHub Actions for CI
3. Tag releases with semantic versioning
4. Publish to GitHub Releases
5. Update install.sh to point to releases

## Contributing Guidelines

1. Follow existing code patterns
2. Use the Charm libraries for any TUI work
3. Keep commands thin, logic in `internal/`
4. Match Catppuccin theme for any new UI elements
5. Document architecture decisions in this file
