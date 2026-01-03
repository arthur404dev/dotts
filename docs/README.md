# dotts Documentation

> Universal Dotfiles Manager - Manage your dotfiles across multiple machines and platforms

## Overview

**dotts** is a CLI tool designed to solve the challenge of managing dotfiles across multiple machines with different configurations. It separates the tool (this CLI) from the configuration repository, allowing you to use pre-made configs or create your own.

## Quick Navigation

| Document | Description |
|----------|-------------|
| [Architecture](./ARCHITECTURE.md) | Technical design and system components |
| [Decisions](./DECISIONS.md) | Architecture Decision Records (ADRs) |
| [Action Plan](./ACTION-PLAN.md) | Prioritized implementation roadmap |
| [Development Guide](./DEVELOPMENT.md) | How to build, test, and contribute |
| [Config Schema](./CONFIG-SCHEMA.md) | YAML configuration reference |
| [Session History](./SESSION-HISTORY.md) | Project context and development history |

## Quick Start

### Installation

```bash
# One-liner install (once releases are available)
curl -fsSL https://dotts.sh/install.sh | bash

# Or build from source
git clone https://github.com/arthur404dev/dotts.git
cd dotts
make build
./dotts --help
```

### Bootstrap a New System

```bash
# Interactive wizard
dotts init

# This will:
# 1. Ask you to choose a config source (default, fork template, or custom)
# 2. Detect your system (OS, distro, architecture)
# 3. Let you select/create a machine configuration
# 4. Configure features (SSH, GPG, GitHub CLI, etc.)
# 5. Set up your dotfiles and install packages
```

## Key Concepts

### Separation of Concerns

```
┌─────────────────────────────────────────────────────────────────┐
│                         dotts (CLI Tool)                        │
│  - System detection                                             │
│  - Package installation                                         │
│  - Dotfile linking                                              │
│  - TUI wizard                                                   │
└─────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    dotts-config (Your Configs)                  │
│  - Profile definitions (base, linux, darwin, desktop, notebook) │
│  - Package manifests (nix, arch, darwin)                        │
│  - Machine-specific settings                                    │
│  - Actual config files (shell, editor, terminal, etc.)          │
└─────────────────────────────────────────────────────────────────┘
```

### Profile Inheritance

Profiles can inherit from other profiles, creating a hierarchy:

```
base
├── linux
│   ├── desktop (multi-monitor workstation)
│   └── notebook (laptop with battery optimization)
└── darwin (macOS)
```

### Alternate Files

Machine-specific config variations using `##` suffix:

```
monitors.conf              # Default
monitors.conf##profile.desktop    # Used on desktop profile
monitors.conf##profile.notebook   # Used on notebook profile
monitors.conf##os.darwin          # Used on macOS
monitors.conf##hostname.mypc      # Used on specific hostname
```

Priority (highest to lowest): hostname > profile > distro > os > default

## Commands

| Command | Description |
|---------|-------------|
| `dotts init` | Bootstrap wizard for new systems |
| `dotts update` | Update configs and packages |
| `dotts status` | Show current configuration state |
| `dotts doctor` | Check system health and dependencies |
| `dotts config` | Manage config source |
| `dotts machine` | Manage machine configurations |
| `dotts sync` | Sync local changes back to config repo |

## Project Status

**Current Phase**: Core CLI complete, templates created, ready for installer implementation.

See [Development Guide](./DEVELOPMENT.md) for what's next and how to continue.

## Target Use Cases

1. **Multi-device developer** with Arch desktop (3 monitors), Arch notebook, macOS work machine
2. **VM/container users** who need quick consistent setups
3. **Dotfiles sharers** who want others to easily use their configs
4. **Configuration experimenters** who want safe rollbacks

## License

MIT License - see LICENSE file for details.
