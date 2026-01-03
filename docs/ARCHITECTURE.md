# Architecture

> Technical design decisions and system architecture for dotts

## Design Philosophy

### 1. Separation of Tool and Configuration

Unlike chezmoi or YADM where configs live inside the tool's structure, dotts keeps them separate:

- **dotts** (this repo): The CLI tool, installable via single binary
- **dotts-config** (user repo): The configuration repository, forkable and customizable

This allows:
- Sharing configs without requiring others to understand the tool internals
- Updating the tool independently of configurations
- Multiple config sources (official, community, personal)

### 2. Profile-Based Inheritance

Instead of complex conditionals in each config file, dotts uses a profile system:

```yaml
# profiles/desktop.yaml
name: desktop
inherits:
  - linux      # Gets all linux configs
configs:
  - terminal   # Adds desktop-specific terminal config
settings:
  monitors: 3
```

This makes configurations:
- Composable (mix and match profiles)
- Readable (no inline conditionals)
- Debuggable (clear inheritance chain)

### 3. Alternate Files Pattern

Inspired by YADM's alternate files, but with cleaner syntax:

```
config.fish                    # Default for all systems
config.fish##os.darwin         # macOS override
config.fish##profile.desktop   # Desktop profile override
config.fish##hostname.mypc     # Specific machine override
```

**Scoring system** for conflict resolution:
| Suffix Type | Score |
|-------------|-------|
| `hostname.X` | 1000 |
| `profile.X` | 100 |
| `distro.X` | 50 |
| `os.X` | 10 |
| `default` | 1 |

Higher score wins. Non-matching suffixes disqualify the file entirely.

## System Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              dotts CLI                                   │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  cmd/dotts/                      internal/                               │
│  ├── main.go                     ├── system/                             │
│  └── cmd/                        │   └── detect.go      # OS detection   │
│      ├── root.go                 ├── state/                              │
│      ├── init.go    ─────────►   │   ├── paths.go       # XDG paths      │
│      ├── update.go               │   └── state.go       # JSON state     │
│      ├── status.go               ├── config/                             │
│      ├── doctor.go               │   ├── source.go      # Git clone      │
│      ├── config.go               │   ├── loader.go      # YAML parsing   │
│      ├── machine.go              │   ├── resolver.go    # Inheritance    │
│      └── sync.go                 │   └── alternates.go  # File matching  │
│                                  └── tui/                                │
│                                      ├── styles/        # Lip Gloss      │
│                                      └── wizard/        # Bubble Tea     │
│                                                                          │
│  pkg/schema/                     templates/                              │
│  ├── profile.go                  └── config-repo/       # Fork template  │
│  ├── machine.go                      ├── profiles/                       │
│  ├── packages.go                     ├── packages/                       │
│  └── config.go                       ├── machines/                       │
│                                      └── configs/                        │
└──────────────────────────────────────────────────────────────────────────┘
```

## Key Components

### System Detection (`internal/system/detect.go`)

Detects runtime environment:
- **OS**: Linux, Darwin (macOS)
- **Distro**: Arch (+ derivatives), Debian, Ubuntu, Fedora, NixOS, macOS
- **Architecture**: amd64, arm64
- **Package Manager**: pacman/yay/paru, apt, dnf, brew, nix
- **Available tools**: git, curl, nix, homebrew

### State Management (`internal/state/`)

Persists session state to `~/.local/share/dotts/state.json`:
- Config source (git URL or local path)
- Current machine name
- Enabled features
- Custom settings

Uses XDG Base Directory specification:
- Config: `~/.config/dotts/`
- Data: `~/.local/share/dotts/`
- Cache: `~/.cache/dotts/`

### Config Loading (`internal/config/`)

#### Source Management
- Clone git repositories
- Validate repo structure
- Support local directories for development

#### Profile Resolution
Resolves inheritance chain and merges:
- Configs (accumulated)
- Packages (merged, deduplicated)
- Settings (overridden by children)
- Scripts (accumulated in order)

#### Alternate Files
Scores and selects the best matching file variant based on current system context.

### TUI Wizard (`internal/tui/wizard/`)

Interactive bootstrap using Charm libraries:
- **Source wizard**: Default / Fork New / Custom URL
- **Machine wizard**: Select existing or create new
- **Settings wizard**: Monitors, git identity, etc.
- **Features wizard**: SSH, GPG, GitHub, asdf, docker, gui
- **Auth wizard**: SSH key generation, GitHub CLI auth

### Schema Definitions (`pkg/schema/`)

Go structs mapping to YAML configuration files:

```go
type Profile struct {
    Name        string
    Inherits    []string      // Profile inheritance
    Configs     []string      // Config directories to link
    Packages    []string      // Package groups to install
    Settings    map[string]any
    Scripts     ProfileScripts
}

type Machine struct {
    Machine  MachineInfo       // hostname, description
    Inherits []string          // Profiles to inherit
    Settings map[string]any    // Machine-specific overrides
    Features []string          // Enabled features
}

type PackageManifest struct {
    Nix    []string            // Nix packages (cross-platform)
    System SystemPackages      // OS-specific packages
    AUR    []string            // Arch User Repository
    Brew   []string            // Homebrew formulae
    Cask   []string            // Homebrew casks (GUI apps)
    Asdf   map[string]string   // asdf plugins and versions
}
```

## Technology Choices

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Language | Go | Single binary, cross-platform, good CLI ecosystem |
| CLI Framework | Cobra | Industry standard, good UX patterns |
| TUI | Charm (huh, lipgloss, bubbletea) | Beautiful, modern, well-maintained |
| Config Format | YAML | Human-readable, widely understood |
| Backend | Nix Home Manager (planned) | Declarative, reproducible, rollback support |
| Theme | Catppuccin | Consistent, modern, multi-platform |

## Inspirations

| Tool | What We Took |
|------|--------------|
| **Chezmoi** | Template-based config generation concept |
| **YADM** | Alternate files pattern (`##suffix`) |
| **Nix Home Manager** | Declarative package management, atomic updates |
| **Stow** | Symlink-based dotfile management |
| **Dotbot** | YAML-based configuration |

## Future Architecture Considerations

### Package Installation Pipeline

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Resolve   │ -> │   Install   │ -> │    Link     │
│  Packages   │    │  Packages   │    │  Dotfiles   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                  │                  │
       ▼                  ▼                  ▼
   Merge from        Nix/Homebrew/      Home Manager
   all profiles      pacman/apt         or symlinks
```

### Home Manager Integration (Planned)

Generate `home.nix` from resolved configuration:

```nix
{ config, pkgs, ... }:
{
  home.packages = with pkgs; [
    # Generated from packages/*.yaml
  ];
  
  home.file = {
    # Generated from configs/*
  };
}
```

Benefits:
- Atomic updates (all or nothing)
- Instant rollback to previous generation
- Reproducible across machines
- No conflicts or partial states

## Security Considerations

- Never store secrets in config repo
- SSH keys generated locally, not synced
- GitHub tokens use `gh auth login` flow
- Support for age/sops encryption (future)
