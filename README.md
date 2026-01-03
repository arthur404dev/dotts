# dotts

> Universal Dotfiles Manager - Beautiful CLI for managing dotfiles across multiple machines and platforms

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- **Interactive TUI** - Beautiful bootstrap wizard using [Charm](https://charm.sh) libraries
- **Profile Inheritance** - Compose configurations: `base → linux → desktop`
- **Alternate Files** - Machine-specific configs with `##` suffix pattern
- **Cross-Platform** - Linux (Arch, Debian, Ubuntu, Fedora) and macOS
- **Multiple Package Managers** - Nix, Homebrew, pacman/yay, apt, dnf
- **Separated Concerns** - CLI tool independent from your config repository

## Quick Start

```bash
# Install (once releases are available)
curl -fsSL https://dotts.sh/install.sh | bash

# Bootstrap your system
dotts init
```

## How It Works

### 1. Separate Tool from Configuration

```
┌─────────────────────────┐     ┌─────────────────────────┐
│       dotts (CLI)       │     │  dotts-config (Yours)   │
│                         │     │                         │
│  • System detection     │ ──► │  • profiles/            │
│  • Package installation │     │  • packages/            │
│  • Dotfile linking      │     │  • machines/            │
│  • TUI wizard           │     │  • configs/             │
└─────────────────────────┘     └─────────────────────────┘
```

### 2. Profile-Based Inheritance

```yaml
# profiles/desktop.yaml
name: desktop
inherits:
  - linux       # Gets linux configs
                # Which inherits from base
configs:
  - terminal
  - wm
settings:
  monitors: 3
  gpu: nvidia
```

### 3. Alternate Files for Machine-Specific Configs

```
kitty.conf                       # Default
kitty.conf##os.darwin            # macOS
kitty.conf##profile.desktop      # Desktop profile
kitty.conf##hostname.workstation # Specific machine
```

## Commands

| Command | Description |
|---------|-------------|
| `dotts init` | Interactive bootstrap wizard |
| `dotts update` | Update configs and packages |
| `dotts status` | Show current configuration state |
| `dotts doctor` | Check system health |
| `dotts config` | Manage config source |
| `dotts machine` | Manage machine configurations |
| `dotts sync` | Sync local changes to config repo |

## Configuration

### Repository Structure

```
your-dotts-config/
├── config.yaml           # Repo metadata
├── profiles/             # Profile definitions
│   ├── base.yaml
│   ├── linux.yaml
│   ├── darwin.yaml
│   ├── desktop.yaml
│   └── notebook.yaml
├── packages/             # Package manifests
│   ├── common.yaml       # Cross-platform (Nix)
│   ├── arch.yaml         # Arch Linux
│   └── darwin.yaml       # macOS (Homebrew)
├── machines/             # Machine-specific configs
│   └── myworkstation.yaml
└── configs/              # Actual config files
    ├── shell/
    ├── editor/
    ├── terminal/
    └── git/
```

### Example Machine Config

```yaml
# machines/workstation.yaml
machine:
  hostname: workstation
  description: Main desktop with 3 monitors

inherits:
  - desktop

settings:
  gpu: nvidia
  monitors: 3

features:
  - ssh
  - github
  - docker
```

## Development

```bash
# Clone
git clone https://github.com/arthur404dev/dotts.git
cd dotts

# Build
make build

# Run
./dotts --help

# Install locally
make install-local
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development guide.

## Documentation

| Document | Description |
|----------|-------------|
| [docs/README.md](docs/README.md) | Documentation index |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Technical design decisions |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | Development guide and next steps |
| [docs/CONFIG-SCHEMA.md](docs/CONFIG-SCHEMA.md) | YAML configuration reference |
| [docs/SESSION-HISTORY.md](docs/SESSION-HISTORY.md) | Project history and context |

## Inspirations

- [Chezmoi](https://www.chezmoi.io/) - Template-based dotfile management
- [YADM](https://yadm.io/) - Alternate files pattern
- [Nix Home Manager](https://github.com/nix-community/home-manager) - Declarative configuration
- [Charm](https://charm.sh/) - Beautiful terminal UIs

## Tech Stack

- **Go** - Single binary, cross-platform
- **Cobra** - CLI framework
- **Charm** (huh, lipgloss, bubbletea) - TUI components
- **Catppuccin** - Color theme

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Status

**Current Phase**: Core CLI complete, ready for package installer implementation.

- [x] CLI skeleton with all commands
- [x] System detection (OS, distro, arch, package manager)
- [x] State management (XDG paths)
- [x] Config loading (YAML, inheritance, alternates)
- [x] TUI wizard (source, machine, settings, features, auth)
- [x] Config repo template
- [x] Install script
- [ ] Package installation
- [ ] Dotfile linking
- [ ] Home Manager integration
- [ ] Tests
- [ ] Releases
