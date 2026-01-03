# Session History

> Complete context of what was built, decisions made, and the journey to this point

## Project Genesis

### The Problem

Managing dotfiles across multiple machines with different configurations:
- **Arch Linux desktop**: 3 monitors, NVIDIA GPU, powerful workstation
- **Arch Linux notebook**: Single display, integrated graphics, battery optimization
- **Linux VMs**: Various distributions, quick setup needed
- **macOS work machine**: Different package managers, different paths

Existing solutions had limitations:
- **Stow**: No conditionals, manual symlinking
- **Chezmoi**: Complex templating, steep learning curve
- **YADM**: Good alternates, but configs tied to tool
- **Nix Home Manager**: Powerful but requires learning Nix

### The Vision

Create a tool that:
1. Separates the CLI from the config repository
2. Uses readable YAML instead of complex templates
3. Supports profile inheritance for DRY configs
4. Has a beautiful TUI for bootstrap
5. Can leverage Nix Home Manager for power users

## Research Phase

### Tools Analyzed

| Tool | Pros | Cons |
|------|------|------|
| Chezmoi | Powerful templating, good docs | Complex, configs tied to tool structure |
| YADM | Alternate files, git-based | Limited inheritance, no TUI |
| Nix Home Manager | Atomic, rollback, reproducible | Steep learning curve |
| Dotbot | Simple YAML config | No inheritance, basic |
| rcm | Simple, tag-based | Limited conditionals |

### Key Inspirations

1. **YADM's alternate files** → Our `##suffix` pattern
2. **Chezmoi's templating** → Simplified to profile settings
3. **Home Manager's declarative approach** → Future backend option
4. **Modern CLI UX (Charm)** → Beautiful TUI wizard

## Architecture Decisions

### Decision 1: Separate Tool from Config

**Why**: Users should be able to fork/share configs without understanding tool internals.

**Implementation**:
- `dotts` (CLI) is installed system-wide
- `dotfiles` (configs) lives in user's repo
- Three source options: default, fork template, custom URL

### Decision 2: Profile-Based Inheritance

**Why**: Avoid repetition, make configs composable.

**Implementation**:
```
base → linux → desktop
     → darwin
     → notebook
```

Each level adds/overrides configs, packages, settings.

### Decision 3: YAML Over Templates

**Why**: Most users don't need Go templates. Simple key-value settings suffice.

**Implementation**:
- Profiles have `settings` map
- Settings are available as environment variables in scripts
- Complex logic goes in shell scripts, not templates

### Decision 4: Alternate Files Pattern

**Why**: YADM proved this works. Clean, explicit, no parsing required.

**Implementation**:
- `file##os.darwin` for macOS
- `file##profile.desktop` for desktop profile
- Scoring system for conflict resolution
- Multiple suffixes with AND logic

### Decision 5: Charm for TUI

**Why**: Beautiful, modern, actively maintained. Makes CLI feel premium.

**Implementation**:
- Lip Gloss for styling (Catppuccin theme)
- Huh for form inputs
- Bubble Tea for interactive components

## Implementation Timeline

### Phase 1: Project Structure (Completed)

1. Set up Go module with Cobra
2. Defined package structure
3. Created Makefile with common commands

### Phase 2: Core Logic (Completed)

1. **System detection** (`internal/system/`)
   - OS, distro, architecture detection
   - Package manager detection
   - Tool availability checks

2. **State management** (`internal/state/`)
   - XDG-compliant paths
   - JSON state persistence
   - Config source tracking

3. **Config loading** (`internal/config/`)
   - Git source cloning
   - YAML file parsing
   - Profile/machine/package loading
   - Inheritance resolution
   - Alternate file matching

### Phase 3: TUI Wizard (Completed)

1. **Styles** (`internal/tui/styles/`)
   - Catppuccin color palette
   - Consistent styling functions

2. **Wizard steps** (`internal/tui/wizard/`)
   - Source selection (default/fork/custom)
   - Machine configuration
   - Settings input
   - Feature selection
   - Authentication setup

### Phase 4: Templates (Completed)

1. **Repository template** (`templates/config-repo/`)
   - `config.yaml` with metadata
   - Five profiles (base, linux, darwin, desktop, notebook)
   - Three package manifests (common, arch, darwin)
   - Example machine config
   - Starter config files for shell, git, editor, terminal, tools

2. **Install script** (`scripts/install.sh`)
   - System detection
   - Binary download
   - PATH configuration

### Phase 5: Documentation (Current)

Creating comprehensive docs for project continuation.

## What's Working

```bash
# Build and run
cd ~/software-development/dotts
make build
./dotts --help

# Commands available (init has full TUI, others are stubs)
dotts init      # Interactive bootstrap wizard
dotts update    # Stub
dotts status    # Stub
dotts doctor    # Stub
dotts config    # Stub
dotts machine   # Stub
dotts sync      # Stub
```

## What's Not Yet Working

1. **Package installation**: Commands defined but not wired to actual installers
2. **Dotfile linking**: No symlink creation or Home Manager integration
3. **Update/status commands**: Stubs only
4. **Tests**: No test coverage yet
5. **Releases**: No GoReleaser config

## File Locations

| Item | Location |
|------|----------|
| Source code | `~/software-development/dotts/` |
| Binary | `~/software-development/dotts/dotts` |
| Templates | `~/software-development/dotts/templates/config-repo/` |
| Docs | `~/software-development/dotts/docs/` |

## Key Files to Understand

| File | Purpose |
|------|---------|
| `cmd/dotts/cmd/init.go` | Bootstrap wizard entry point |
| `internal/tui/wizard/wizard.go` | Wizard orchestration |
| `internal/config/resolver.go` | Profile inheritance logic |
| `internal/config/alternates.go` | Alternate file scoring |
| `internal/system/detect.go` | System detection |
| `pkg/schema/*.go` | YAML type definitions |

## Continuation Prompt

Use this prompt to continue development in a new session:

```
Continue building the dotts CLI project at ~/software-development/dotts/

This is a Universal Dotfiles Manager with:
- Go CLI using Cobra + Charm (huh, bubbletea, lipgloss)
- Nix Home Manager as planned backend for dotfile management
- Separated CLI (dotts) and config repo (dotfiles)
- Alternate files pattern (##os.Darwin, ##profile.desktop) for machine-specific configs
- Interactive TUI wizard for bootstrap

COMPLETED:
- Full project structure and CLI skeleton
- System detection (OS, distro, arch, package manager)
- State management with XDG paths
- Config loading with YAML schemas
- Profile resolver with inheritance
- Alternate file resolution with scoring
- Complete TUI wizard (source, machine, settings, features, auth)
- Full config-repo template with starter configs
- Install script for curl|sh bootstrap
- Comprehensive documentation in docs/

NEXT PRIORITIES:
1. Implement internal/installer/ for package installation (Nix, Homebrew, pacman/yay)
2. Implement dotfile linking (symlinks or Home Manager)
3. Wire up update/status/doctor commands
4. Add tests
5. Create .goreleaser.yaml for releases

The project builds and runs: `cd ~/software-development/dotts && make build && ./dotts --help`

Read docs/DEVELOPMENT.md for detailed next steps.
```

## Lessons Learned

1. **Start with schemas**: Defining Go structs first made YAML loading straightforward
2. **Charm libraries are excellent**: Beautiful UX with minimal code
3. **Inheritance is tricky**: Order matters, need clear resolution rules
4. **Templates need real configs**: Users won't know what to put in empty files

## Open Questions for Future

1. **Home Manager vs symlinks**: Should symlinks be the default, with HM as upgrade path?
2. **Secret management**: How to handle encrypted secrets (age/sops)?
3. **Remote config sources**: Support for private repos, authentication?
4. **Rollback mechanism**: Without Nix, how to implement safe rollback?
