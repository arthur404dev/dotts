# Action Plan

> Prioritized implementation roadmap for dotts MVP

Last updated: 2026-01-02

---

## Overview

This document outlines the implementation path from current state to MVP release. Work is organized into phases with clear deliverables and success criteria.

---

## Current State

| Component | Status |
|-----------|--------|
| CLI Framework | Done |
| System Detection | Done |
| State Management | Done |
| Config Loading | Done |
| Profile Resolution | Done |
| Alternate Files | Done |
| TUI Wizard | Done |
| Package Installer | Not Started |
| Dotfile Linker | Not Started |
| Command Implementation | Stubs Only |
| Tests | Not Started |
| CI/CD | Not Started |
| Releases | Not Started |

---

## Phase 1: Core Functionality (HIGH PRIORITY)

**Goal**: Make `dotts init` actually do something beyond saving state.

### 1.1 Package Installer Infrastructure

**Files to create**:
```
internal/installer/
├── installer.go      # Interface and orchestrator
├── nix.go            # Nix package installer
├── brew.go           # Homebrew installer
├── pacman.go         # Arch Linux (pacman + AUR helpers)
└── apt.go            # Debian/Ubuntu
```

**Interface**:
```go
type Installer interface {
    Name() string
    Available() bool
    Install(packages []string) error
    Remove(packages []string) error  
    IsInstalled(pkg string) bool
    Update() error
}
```

**Success criteria**:
- [ ] Can install Nix packages via `nix-env` or `nix profile`
- [ ] Can install Homebrew packages via `brew`
- [ ] Can install Arch packages via `pacman` and AUR via `yay`/`paru`
- [ ] Graceful fallback when package manager unavailable
- [ ] Progress output during installation

### 1.2 Dotfile Linker

**Files to create**:
```
internal/linker/
├── linker.go         # Interface definition
├── symlink.go        # Symlink-based implementation
├── backup.go         # Backup before linking
└── manifest.go       # Track linked files
```

**Behavior**:
1. Scan config directories from resolved profile
2. For each config dir, determine link strategy:
   - Directory: `~/.config/X` -> `config-repo/configs/Y/.config/X`
   - File: `~/.zshrc` -> `config-repo/configs/shell/.zshrc`
3. If target exists and not our symlink:
   - Backup to `~/.local/share/dotts/backups/`
   - Record in manifest
4. Create symlink
5. Handle alternate files (select best match based on context)

**Success criteria**:
- [ ] Can link directories from config repo
- [ ] Can link individual files
- [ ] Backs up existing files before overwriting
- [ ] Tracks all links in state for clean removal
- [ ] Respects alternate file selection

### 1.3 Progress TUI

**Files to create**:
```
internal/tui/progress/
├── progress.go       # Progress display component
└── spinner.go        # Spinner for long operations
```

**Features**:
- Animated progress bar for package installation
- Step-by-step status during linking
- Clear success/failure indicators

---

## Phase 2: Command Implementation (MEDIUM PRIORITY)

**Goal**: Wire up all stub commands to real functionality.

### 2.1 Update Command (`cmd/dotts/cmd/update.go`)

**Flow**:
1. Pull latest from config repo (if git source)
2. Show diff of changes (new packages, changed configs)
3. Prompt for confirmation (unless `--yes`)
4. Run pre-update scripts
5. Install new packages
6. Update symlinks
7. Run post-update scripts
8. Show summary

**Flags**:
- `--dry-run`: Show what would change
- `--packages-only`: Only update packages
- `--dotfiles-only`: Only update dotfiles
- `--yes`: Skip confirmation

### 2.2 Status Command (`cmd/dotts/cmd/status.go`)

**Output**:
```
dotts status
============

Config Source: git@github.com:user/dotts-config.git (main)
Last Updated:  2026-01-02 10:30:00

Machine:       workstation
Profile Chain: base -> linux -> desktop
Features:      ssh, github, docker

Packages:
  Nix:     42 installed
  System:  15 installed
  AUR:     3 installed

Dotfiles:
  Linked:  12 directories, 5 files
  Pending: 0 changes

Health: OK
```

### 2.3 Doctor Command (`cmd/dotts/cmd/doctor.go`)

**Checks**:
- [ ] Config source accessible
- [ ] Required tools installed (git, package managers)
- [ ] All symlinks valid (not broken)
- [ ] No conflicting files
- [ ] Disk space adequate

**Output format**:
```
dotts doctor
============

[OK] Config source accessible
[OK] Git installed (2.43.0)
[OK] Nix installed (2.18.1)
[WARN] yay not found (AUR packages unavailable)
[OK] All 17 symlinks valid
[OK] No conflicts detected

Summary: 4 passed, 1 warning, 0 errors
```

### 2.4 Config Command (`cmd/dotts/cmd/config.go`)

**Subcommands**:
- `dotts config show`: Display current config source
- `dotts config set <url>`: Change config source
- `dotts config pull`: Pull latest without applying
- `dotts config path`: Print config repo path

### 2.5 Machine Command (`cmd/dotts/cmd/machine.go`)

**Subcommands**:
- `dotts machine list`: List available machines
- `dotts machine show [name]`: Show machine config
- `dotts machine set <name>`: Switch to different machine
- `dotts machine create <name>`: Create new machine config

### 2.6 Sync Command (`cmd/dotts/cmd/sync.go`)

**Flow**:
1. Detect local changes in config repo
2. Show diff
3. Prompt for commit message
4. Commit and push

---

## Phase 3: Testing (MEDIUM PRIORITY)

**Goal**: Establish test coverage for critical paths.

### 3.1 Unit Tests

**Priority test files**:
```
internal/config/resolver_test.go      # Profile inheritance
internal/config/alternates_test.go    # File matching
internal/installer/installer_test.go  # Mock-based tests
internal/linker/symlink_test.go       # Filesystem operations
pkg/schema/profile_test.go            # YAML parsing
```

**Test utilities needed**:
- Temporary directory fixtures
- Mock package manager
- Mock filesystem

### 3.2 Integration Tests

**Scenarios**:
- Full init wizard flow (headless mode)
- Package installation (with mock)
- Dotfile linking and unlinking
- Update with changes

---

## Phase 4: Release Pipeline (MEDIUM PRIORITY)

**Goal**: Automated multi-platform releases.

### 4.1 GoReleaser Configuration

**File**: `.goreleaser.yaml`

**Artifacts**:
- Linux amd64/arm64
- macOS amd64/arm64 (universal binary)
- Checksums
- Release notes from conventional commits

### 4.2 GitHub Actions

**Workflows**:
```
.github/workflows/
├── ci.yaml           # Test on push/PR
├── release.yaml      # Release on tag
└── lint.yaml         # Linting
```

### 4.3 Install Script Update

Update `scripts/install.sh` to work with actual releases.

---

## Phase 5: Polish (LOW PRIORITY)

### 5.1 Shell Completions

Generate completions for bash, zsh, fish via Cobra.

### 5.2 Man Pages

Generate man pages from command docs.

### 5.3 Homebrew Formula

Create formula for `brew install dotts`.

### 5.4 AUR Package

Create PKGBUILD for Arch users.

---

## Success Criteria for MVP

MVP is ready when:

- [ ] `dotts init` completes full setup (wizard -> packages -> dotfiles)
- [ ] `dotts update` pulls changes and applies them
- [ ] `dotts status` shows current state accurately
- [ ] `dotts doctor` validates setup
- [ ] Basic test coverage (>50% on critical paths)
- [ ] Works on: Arch Linux, Ubuntu, macOS
- [ ] Releases available on GitHub
- [ ] Install script works

---

## Timeline Estimate

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| Phase 1 (Core) | 3-5 days | None |
| Phase 2 (Commands) | 2-3 days | Phase 1 |
| Phase 3 (Testing) | 2-3 days | Phase 1-2 |
| Phase 4 (Release) | 1-2 days | Phase 1-3 |
| Phase 5 (Polish) | 1-2 days | Phase 4 |

**Total to MVP**: ~10-15 days of focused work

---

## Next Immediate Steps

1. ~~Initialize git repository~~
2. ~~Document architecture decisions~~
3. Implement `internal/installer/` package
4. Implement `internal/linker/` package
5. Wire up `dotts init` to use installer + linker
6. Test on local machine
7. Implement remaining commands
