# Configuration Schema Reference

> Complete reference for all YAML configuration files in a dotts config repository

## Repository Structure

```
dotfiles/
├── config.yaml           # Repository metadata
├── profiles/             # Profile definitions
│   ├── base.yaml
│   ├── linux.yaml
│   ├── darwin.yaml
│   ├── desktop.yaml
│   └── notebook.yaml
├── packages/             # Package manifests
│   ├── common.yaml
│   ├── arch.yaml
│   └── darwin.yaml
├── machines/             # Machine-specific configs
│   └── mymachine.yaml
├── configs/              # Actual config files
│   ├── shell/
│   ├── editor/
│   ├── terminal/
│   ├── git/
│   └── tools/
├── scripts/              # Lifecycle scripts
│   └── setup-shell.sh
└── assets/               # Static assets
    └── wallpapers/
```

## config.yaml

Repository metadata and global settings.

```yaml
# Required
name: my-dotfiles

# Optional
author: "Your Name"
description: "Personal dotfiles managed by dotts"
version: "1.0.0"

# Default machine when hostname doesn't match any machine file
default_machine: example

# Available features (used in feature wizard)
features:
  - ssh        # SSH key generation
  - gpg        # GPG key setup
  - github     # GitHub CLI authentication
  - asdf       # asdf version manager
  - docker     # Docker/container tooling
  - gui        # GUI applications

# Minimum dotts version required
min_dotts_version: "0.1.0"
```

## Profile Schema

Profiles define what configs and packages to use.

```yaml
# profiles/base.yaml

# Profile name (should match filename without .yaml)
name: base

# Human-readable description
description: Universal base configuration

# Profiles to inherit from (optional)
# Inheritance is processed in order, later values override earlier
inherits:
  - some-other-profile

# Config directories to symlink (relative to configs/)
# Each entry is a directory that will be linked
configs:
  - shell       # Links configs/shell/
  - editor      # Links configs/editor/
  - git         # Links configs/git/

# Package groups to install (relative to packages/)
packages:
  - common      # Loads packages/common.yaml

# Key-value settings (accessible in templates and scripts)
settings:
  shell: fish
  editor: nvim
  theme: catppuccin
  monitors: 1

# Lifecycle scripts (relative to scripts/)
scripts:
  pre_install:
    - scripts/check-requirements.sh
  post_install:
    - scripts/setup-shell.sh
    - scripts/configure-editor.sh
  pre_update:
    - scripts/backup-current.sh
  post_update:
    - scripts/reload-shell.sh
```

### Profile Inheritance

```yaml
# profiles/desktop.yaml
name: desktop
inherits:
  - linux      # First, apply linux profile
               # (which inherits from base)

configs:
  - wm         # Add desktop-specific configs

settings:
  monitors: 3  # Override inherited setting
```

Resolution order (for desktop):
1. base (root)
2. linux (inherits base)
3. desktop (inherits linux)

Later profiles override settings from earlier ones.

## Machine Schema

Machine configs specify which profiles apply to a specific machine.

```yaml
# machines/mydesktop.yaml

# Machine identification
machine:
  hostname: mydesktop     # Must match system hostname
  description: Main workstation with 3 monitors

# Profiles to inherit (in order)
inherits:
  - desktop              # Will also get linux -> base

# Machine-specific setting overrides
settings:
  gpu: nvidia
  monitors: 3
  theme: catppuccin-mocha

# Features to enable on this machine
features:
  - ssh
  - github
  - asdf
  - docker
  - gui
```

## Package Manifest Schema

Package manifests define what to install per platform.

```yaml
# packages/common.yaml

# Nix packages (cross-platform, preferred)
nix:
  - ripgrep
  - fd
  - fzf
  - bat
  - eza
  - neovim
  - git
  - lazygit

# System packages (platform-specific)
system:
  arch:
    - base-devel
    - openssh
  debian:
    - build-essential
    - openssh-client
  ubuntu:
    - build-essential
  fedora:
    - @development-tools
  darwin: []    # macOS uses brew/cask instead

# AUR packages (Arch Linux only)
aur:
  - yay-bin
  - visual-studio-code-bin

# Homebrew formulae (macOS)
brew:
  - mas       # Mac App Store CLI
  - trash

# Homebrew casks (macOS GUI apps)
cask:
  - firefox
  - visual-studio-code
  - docker
  - raycast

# asdf plugins and versions
asdf:
  nodejs: "lts"      # Latest LTS
  python: "latest"   # Latest stable
  golang: "1.21.0"   # Specific version
```

### Package Resolution

When multiple package manifests are loaded, they are merged:
- Lists are concatenated and deduplicated
- asdf map entries don't override (first wins)

## Alternate Files

Files can have platform/profile-specific variants using `##` suffix:

```
configs/terminal/.config/kitty/
├── kitty.conf                      # Default
├── kitty.conf##os.darwin           # macOS override
├── kitty.conf##profile.desktop     # Desktop profile
├── kitty.conf##hostname.mypc       # Specific machine
└── kitty.conf##profile.notebook,os.linux  # Combined conditions
```

### Suffix Types

| Suffix | Example | Matches |
|--------|---------|---------|
| `os.X` | `##os.darwin` | Operating system (linux, darwin) |
| `distro.X` | `##distro.arch` | Linux distribution |
| `profile.X` | `##profile.desktop` | Active profile |
| `hostname.X` | `##hostname.mypc` | Machine hostname |
| `default` | `##default` | Fallback (low priority) |

### Scoring

Multiple matches are resolved by score (highest wins):

| Suffix Type | Score |
|-------------|-------|
| hostname | 1000 |
| profile | 100 |
| distro | 50 |
| os | 10 |
| default | 1 |

**Important**: If any suffix condition doesn't match, the entire file is disqualified.

### Combined Suffixes

Multiple conditions can be combined with commas:

```
config.fish##os.linux,profile.desktop
```

This file is only used when:
- OS is Linux AND
- Profile is desktop

All conditions must match for the file to be considered.

## Config Directory Structure

Each config directory should mirror the target structure:

```
configs/shell/
├── .config/
│   └── fish/
│       └── config.fish     # -> ~/.config/fish/config.fish
└── .zshrc                  # -> ~/.zshrc
```

The entire directory is linked/copied, maintaining structure.

## Settings Reference

Common settings used across profiles:

| Setting | Type | Description |
|---------|------|-------------|
| `shell` | string | Default shell (fish, zsh, bash) |
| `editor` | string | Default editor (nvim, vim, code) |
| `terminal` | string | Terminal emulator (kitty, wezterm, alacritty) |
| `theme` | string | Color scheme (catppuccin, dracula, etc.) |
| `monitors` | int | Number of displays |
| `gpu` | string | GPU vendor (nvidia, amd, intel) |
| `compositor` | string | Window manager (hyprland, sway, etc.) |
| `bar` | string | Status bar (waybar, polybar, etc.) |

Settings are accessible in:
- Lifecycle scripts (as environment variables)
- Templates (future feature)
- dotts status output

## Features Reference

Features toggle optional functionality:

| Feature | Description |
|---------|-------------|
| `ssh` | Generate SSH keys, configure agent |
| `gpg` | Set up GPG keys |
| `github` | Authenticate with GitHub CLI |
| `asdf` | Install asdf version manager |
| `docker` | Install Docker, add user to group |
| `gui` | Install GUI applications |

Features affect:
- Which packages are installed
- What setup scripts run
- Available wizard options

## Example: Complete Config Repo

```yaml
# config.yaml
name: my-dotfiles
author: arthur404dev
features: [ssh, github, asdf, docker, gui]

# profiles/base.yaml
name: base
configs: [shell, git, editor]
packages: [common]
settings:
  shell: fish
  editor: nvim

# profiles/linux.yaml
name: linux
inherits: [base]
configs: [tools]
packages: [arch]

# profiles/desktop.yaml
name: desktop
inherits: [linux]
configs: [terminal, wm]
settings:
  monitors: 3

# machines/workstation.yaml
machine:
  hostname: workstation
inherits: [desktop]
settings:
  gpu: nvidia
features: [ssh, github, asdf, docker, gui]
```
