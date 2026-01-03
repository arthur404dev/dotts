# Architecture Decisions

> Record of key architectural decisions for dotts development

This document captures important technical decisions and their rationale. New decisions should be added at the bottom with a date and context.

---

## ADR-001: Commit Convention

**Date**: 2026-01-02  
**Status**: Accepted

### Context

Need a consistent commit message format for maintainability and changelog generation.

### Decision

Use [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

<body with bullet points>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, no logic change)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding or updating tests
- `build`: Build system or external dependencies
- `ci`: CI/CD configuration
- `chore`: Other changes (maintenance, tooling)

**Rules**:
- No emojis in commits
- Body should contain bullet points explaining changes
- Commits should be atomic and self-contained
- Scope is optional but recommended (e.g., `feat(installer):`)
- **Max 2 folders per commit** unless changes are strictly related and small
- Commits should be domain-specific and follow logical development order
- Someone navigating history should understand the project evolution

### Consequences

- Enables automated changelog generation
- Clear history for debugging and review
- Enforces thoughtful commits

---

## ADR-002: Symlinks First, Home Manager Later

**Date**: 2026-01-02  
**Status**: Accepted

### Context

Two approaches for dotfile management:
1. Simple symlinks (traditional, like Stow)
2. Nix Home Manager (declarative, atomic, rollback support)

### Decision

Implement **symlinks first** as the default backend, with Home Manager as an optional upgrade path.

**Rationale**:
- Lower barrier to entry (no Nix knowledge required)
- Works on all systems immediately
- Home Manager can be added as `--backend=home-manager` flag later
- Users can migrate when ready

### Implementation

```
internal/linker/
├── linker.go       # Interface: Linker { Link(), Unlink(), Status() }
├── symlink.go      # Default: creates ~/.config/X -> config-repo/configs/X
└── homemanager.go  # Future: generates home.nix, runs home-manager switch
```

**Symlink strategy**:
- Link entire directories when possible (e.g., `~/.config/nvim/`)
- Link individual files for root-level configs (e.g., `~/.zshrc`)
- Backup existing files to `~/.local/share/dotts/backups/` before linking
- Track linked files in state.json for clean unlink

### Consequences

- Faster initial development
- Users can start using dotts immediately
- Home Manager integration deferred but architecturally planned

---

## ADR-003: Package Installation Strategy

**Date**: 2026-01-02  
**Status**: Accepted

### Context

Need to support multiple package managers across platforms:
- Nix (cross-platform)
- Homebrew (macOS)
- pacman/yay/paru (Arch Linux)
- apt (Debian/Ubuntu)
- dnf (Fedora)

### Decision

Implement a **unified installer interface** with per-manager implementations.

**Priority order for package installation**:
1. **Nix packages first** (if Nix is available) - cross-platform, reproducible
2. **System packages second** - platform-specific essentials
3. **AUR/Cask/etc. last** - specialized packages

**Rationale**:
- Nix provides consistency across machines
- System packages handle things Nix can't (kernel modules, system services)
- Specialized repos (AUR, Cask) for GUI apps and platform-specific tools

### Implementation

```go
// internal/installer/installer.go
type Installer interface {
    Name() string
    Available() bool
    Install(packages []string) error
    Remove(packages []string) error
    IsInstalled(pkg string) bool
    Update() error
}

type InstallPlan struct {
    Nix    []string  // Install first
    System []string  // Install second
    AUR    []string  // Arch only
    Brew   []string  // macOS only
    Cask   []string  // macOS GUI apps
}
```

### Consequences

- Clear installation order
- Each manager is isolated and testable
- Easy to add new package managers

---

## ADR-004: Rollback Strategy Without Nix

**Date**: 2026-01-02  
**Status**: Accepted

### Context

Home Manager provides atomic rollback via Nix generations. For symlink-based setup, need alternative rollback mechanism.

### Decision

Implement **snapshot-based rollback** for symlink backend.

**Strategy**:
1. Before any destructive operation, create a snapshot
2. Snapshots stored in `~/.local/share/dotts/snapshots/`
3. Each snapshot contains:
   - Backup of replaced files
   - State at that point in time
   - Manifest of what was linked

**Commands**:
```bash
dotts rollback              # Rollback to previous snapshot
dotts rollback --list       # List available snapshots
dotts rollback <snapshot>   # Rollback to specific snapshot
```

### Implementation

```
~/.local/share/dotts/
├── state.json
├── backups/           # Individual file backups (before first link)
└── snapshots/
    ├── 2026-01-02T10-30-00/
    │   ├── manifest.json
    │   ├── state.json
    │   └── files/     # Backed up files
    └── 2026-01-02T15-45-00/
        └── ...
```

### Consequences

- Safe operations with undo capability
- Disk space usage for backups (mitigated by cleanup policy)
- Not as atomic as Nix, but practical for symlink approach

---

## ADR-005: Secret Management Approach

**Date**: 2026-01-02  
**Status**: Proposed (not yet implemented)

### Context

Users need to manage secrets (API keys, tokens) in their dotfiles without committing them to git.

### Decision

Support **age encryption** via sops-nix pattern, deferred to post-MVP.

**Approach**:
- Files ending in `.secret` or in `secrets/` directory are encrypted
- Use `age` for encryption (simpler than GPG)
- Decrypt on `dotts update`, re-encrypt on `dotts sync`
- Keys stored in `~/.config/dotts/age-key.txt`

**For MVP**: Document that secrets should use:
- Environment variables
- External secret managers (1Password CLI, Bitwarden CLI)
- `.gitignore`d local files

### Consequences

- MVP ships without secret management (acceptable for initial release)
- Clear path to add encryption later
- Users have workarounds in the meantime

---

## ADR-006: Error Handling Philosophy

**Date**: 2026-01-02  
**Status**: Accepted

### Context

CLI tools need clear, actionable error messages.

### Decision

Follow these error handling principles:

1. **Fail fast, fail clearly**: Stop on first error with clear message
2. **Suggest fixes**: Error messages should include remediation steps
3. **Dry-run by default for destructive operations**: Show what would happen
4. **Confirm before destructive changes**: Unless `--yes` flag provided

**Error format**:
```
Error: <what went wrong>

Cause: <why it happened>

Fix: <what to do about it>
```

### Implementation

```go
type DottsError struct {
    Message string
    Cause   string
    Fix     string
}
```

### Consequences

- Better user experience
- Reduced support burden
- Self-documenting errors

---

## Future Decisions (To Be Made)

These decisions are pending and will be documented when resolved:

- [ ] **CI/CD Pipeline**: GitHub Actions vs other
- [ ] **Release Strategy**: GoReleaser configuration, versioning scheme
- [ ] **Plugin System**: Allow custom installers/linkers?
- [ ] **Remote Config Auth**: How to handle private repos
- [ ] **Update Notification**: Check for new dotts versions?
