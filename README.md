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
curl -fsSL https://dotts.4o4.sh/install.sh | sh

dotts init
```

## How It Works

### 1. Separate Tool from Configuration

```
┌─────────────────────────┐     ┌─────────────────────────┐
│       dotts (CLI)       │     │    dotfiles (Yours)     │
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
your-dotfiles/
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

## Releasing a New Version

This guide walks through releasing a new version of dotts.

### Prerequisites

- All changes committed and pushed to `main`
- Tests passing (when implemented)
- Documentation updated in dotts-docs repo

### Step-by-Step Release Process

#### 1. Update Version References

Update any hardcoded version references in the codebase if needed.

#### 2. Create and Push a Git Tag

```bash
# Create annotated tag
git tag -a v0.2.0 -m "v0.2.0 - Brief description of release"

# Push the tag
git push origin v0.2.0
```

#### 3. GitHub Actions Handles the Rest

The `release.yaml` workflow automatically:
- Builds binaries for all platforms (linux/darwin × amd64/arm64)
- Creates a GitHub Release with the binaries
- Generates checksums

#### 4. Verify the Release

1. Check [GitHub Releases](https://github.com/arthur404dev/dotts/releases)
2. Verify all 4 binaries are attached:
   - `dotts_X.X.X_linux_amd64.tar.gz`
   - `dotts_X.X.X_linux_arm64.tar.gz`
   - `dotts_X.X.X_darwin_amd64.tar.gz`
   - `dotts_X.X.X_darwin_arm64.tar.gz`

#### 5. Test the Install Script

```bash
# Test fresh install
curl -fsSL https://dotts.4o4.sh/install.sh | sh

# Verify version
dotts --version
```

#### 6. Update Documentation (dotts-docs repo)

If this is a new minor/major version with breaking changes or new features:

```bash
cd ../dotts-docs

# Archive previous version docs
pnpm version:archive 0.1

# Update docs for new version
# ... edit content/docs/ ...

# Commit and deploy
git add .
git commit -m "chore: release v0.2, archive v0.1 docs"
git push
```

See the [dotts-docs README](https://github.com/arthur404dev/dotts-docs#releasing-a-new-version) for detailed documentation versioning steps.

### Release Checklist

- [ ] All changes committed to `main`
- [ ] Version tag created and pushed
- [ ] GitHub Release created automatically
- [ ] All platform binaries attached to release
- [ ] Install script tested with new version
- [ ] Documentation updated (if needed)
- [ ] Documentation archived (if new minor/major version)

### Versioning Scheme

dotts follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** (v0.1.0 → v0.2.0): New features, backward compatible
- **PATCH** (v0.1.0 → v0.1.1): Bug fixes, backward compatible

For documentation:
- **Minor/Major** releases: Archive previous docs, create new version
- **Patch** releases: Update existing docs in place

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
- [x] Releases (goreleaser + GitHub Actions)
- [ ] Package installation
- [ ] Dotfile linking
- [ ] Home Manager integration
- [ ] Tests
