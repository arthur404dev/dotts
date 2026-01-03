package installer

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Installer defines the interface for package managers
type Installer interface {
	// Name returns the installer identifier (e.g., "nix", "pacman", "brew")
	Name() string

	// Available checks if this installer is usable on the current system
	Available() bool

	// Install installs the given packages
	Install(ctx context.Context, packages []string) error

	// Remove uninstalls the given packages
	Remove(ctx context.Context, packages []string) error

	// IsInstalled checks if a specific package is already installed
	IsInstalled(pkg string) bool

	// Update refreshes the package database/index
	Update(ctx context.Context) error

	// NeedsSudo returns true if this installer requires elevated privileges
	NeedsSudo() bool
}

// InstallPlan groups packages by their target installer
type InstallPlan struct {
	Nix    []string // Cross-platform Nix packages
	Pacman []string // Arch Linux official repos
	AUR    []string // Arch User Repository
	Apt    []string // Debian/Ubuntu
	Dnf    []string // Fedora
	Brew   []string // Homebrew formulae
	Cask   []string // Homebrew casks (GUI apps)
}

// IsEmpty returns true if there are no packages to install
func (p *InstallPlan) IsEmpty() bool {
	return len(p.Nix) == 0 &&
		len(p.Pacman) == 0 &&
		len(p.AUR) == 0 &&
		len(p.Apt) == 0 &&
		len(p.Dnf) == 0 &&
		len(p.Brew) == 0 &&
		len(p.Cask) == 0
}

// Total returns the total number of packages across all installers
func (p *InstallPlan) Total() int {
	return len(p.Nix) + len(p.Pacman) + len(p.AUR) +
		len(p.Apt) + len(p.Dnf) + len(p.Brew) + len(p.Cask)
}

// InstallResult tracks the outcome of an installation
type InstallResult struct {
	Installer string
	Requested []string
	Installed []string
	Skipped   []string // Already installed
	Failed    []string
	Error     error
}

// Success returns true if all requested packages were handled without error
func (r *InstallResult) Success() bool {
	return r.Error == nil && len(r.Failed) == 0
}

// InstallProgress reports installation progress
type InstallProgress struct {
	Installer string
	Package   string
	Current   int
	Total     int
	Status    ProgressStatus
	Message   string
}

type ProgressStatus int

const (
	StatusPending ProgressStatus = iota
	StatusRunning
	StatusSuccess
	StatusSkipped
	StatusFailed
)

func (s ProgressStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusSuccess:
		return "success"
	case StatusSkipped:
		return "skipped"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// ProgressCallback is called during installation to report progress
type ProgressCallback func(progress InstallProgress)

// BaseInstaller provides common functionality for installers
type BaseInstaller struct {
	name string
}

// Name returns the installer name
func (b *BaseInstaller) Name() string {
	return b.name
}

// commandExists checks if a command is available in PATH
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// runCommand executes a command and returns combined output
func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// runCommandSilent executes a command and discards output
func runCommandSilent(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Run()
}

// sudoWrap prepends sudo to args if needed
func sudoWrap(needsSudo bool, command string, args []string) (string, []string) {
	if needsSudo {
		return "sudo", append([]string{command}, args...)
	}
	return command, args
}

// filterInstalled returns only packages that are not already installed
func filterInstalled(installer Installer, packages []string) (toInstall, alreadyInstalled []string) {
	for _, pkg := range packages {
		if installer.IsInstalled(pkg) {
			alreadyInstalled = append(alreadyInstalled, pkg)
		} else {
			toInstall = append(toInstall, pkg)
		}
	}
	return
}

// InstallError wraps installation errors with context
type InstallError struct {
	Installer string
	Package   string
	Cause     error
	Output    string
}

func (e *InstallError) Error() string {
	if e.Package != "" {
		return fmt.Sprintf("%s: failed to install %s: %v", e.Installer, e.Package, e.Cause)
	}
	return fmt.Sprintf("%s: installation failed: %v", e.Installer, e.Cause)
}

func (e *InstallError) Unwrap() error {
	return e.Cause
}
