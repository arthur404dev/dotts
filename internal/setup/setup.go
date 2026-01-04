package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/arthur404dev/dotts/internal/state"
)

type SourceType string

const (
	SourceTypeStarter  SourceType = "starter"
	SourceTypeBYOP     SourceType = "byop"
	SourceTypeForkUser SourceType = "fork_user"
)

type BYOPSubType string

const (
	BYOPSubTypeRemote BYOPSubType = "remote"
	BYOPSubTypeLocal  BYOPSubType = "local"
)

type SourceConfig struct {
	Type SourceType

	// For BYOP
	BYOPSubType BYOPSubType
	RemoteURL   string
	LocalPath   string

	// For ForkUser
	ForkFrom string

	// Symlink options
	CreateSymlink bool
	SymlinkPath   string
}

type SetupResult struct {
	ConfigPath  string
	SymlinkPath string
	BackupPath  string
	ClonedFrom  string
	CopiedFrom  string
}

type Setup struct {
	paths *state.Paths
}

func New() *Setup {
	return &Setup{
		paths: state.GetPaths(),
	}
}

const (
	StarterRepoURL = "https://github.com/arthur404dev/dotts-starter.git"
)

func (s *Setup) Execute(config SourceConfig) (*SetupResult, error) {
	result := &SetupResult{
		ConfigPath: s.paths.ConfigRepo,
	}

	if err := s.paths.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	if pathExists(s.paths.ConfigRepo) {
		backupPath, err := s.backupExisting(s.paths.ConfigRepo)
		if err != nil {
			return nil, fmt.Errorf("failed to backup existing config: %w", err)
		}
		result.BackupPath = backupPath
	}

	switch config.Type {
	case SourceTypeStarter:
		if err := s.cloneRepo(StarterRepoURL, s.paths.ConfigRepo); err != nil {
			return nil, fmt.Errorf("failed to clone starter: %w", err)
		}
		result.ClonedFrom = StarterRepoURL

	case SourceTypeBYOP:
		switch config.BYOPSubType {
		case BYOPSubTypeRemote:
			url := normalizeGitURL(config.RemoteURL)
			if err := s.cloneRepo(url, s.paths.ConfigRepo); err != nil {
				return nil, fmt.Errorf("failed to clone repo: %w", err)
			}
			result.ClonedFrom = url

		case BYOPSubTypeLocal:
			sourcePath := expandPath(config.LocalPath)
			if err := s.copyDirectory(sourcePath, s.paths.ConfigRepo); err != nil {
				return nil, fmt.Errorf("failed to copy local path: %w", err)
			}
			result.CopiedFrom = sourcePath
		}

	case SourceTypeForkUser:
		url := normalizeGitURL(config.ForkFrom)
		if err := s.cloneRepo(url, s.paths.ConfigRepo); err != nil {
			return nil, fmt.Errorf("failed to clone fork: %w", err)
		}
		result.ClonedFrom = url
	}

	if config.CreateSymlink && config.SymlinkPath != "" {
		symlinkPath := expandPath(config.SymlinkPath)
		if err := s.createSymlink(s.paths.ConfigRepo, symlinkPath); err != nil {
			return nil, fmt.Errorf("failed to create symlink: %w", err)
		}
		result.SymlinkPath = symlinkPath
	}

	return result, nil
}

func (s *Setup) cloneRepo(url, dest string) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "--depth=1", url, dest)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func (s *Setup) copyDirectory(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return os.MkdirAll(dst, info.Mode())
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath, info.Mode())
	})
}

func (s *Setup) backupExisting(path string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(s.paths.BackupsDir, fmt.Sprintf("config-%s", timestamp))

	if err := os.MkdirAll(s.paths.BackupsDir, 0755); err != nil {
		return "", err
	}

	if err := os.Rename(path, backupPath); err != nil {
		return "", err
	}

	return backupPath, nil
}

func (s *Setup) createSymlink(target, link string) error {
	if pathExists(link) {
		linkInfo, err := os.Lstat(link)
		if err != nil {
			return err
		}

		if linkInfo.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(link); err != nil {
				return err
			}
		} else {
			timestamp := time.Now().Format("20060102-150405")
			backupPath := link + ".backup-" + timestamp
			if err := os.Rename(link, backupPath); err != nil {
				return err
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		return err
	}

	return os.Symlink(target, link)
}

func normalizeGitURL(input string) string {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "git@") {
		return input
	}

	if !strings.HasSuffix(input, ".git") {
		input += ".git"
	}

	return "https://" + input
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, mode)
}

func (s *Setup) ConfigRepoPath() string {
	return s.paths.ConfigRepo
}

func (s *Setup) ConfigRepoExists() bool {
	return pathExists(s.paths.ConfigRepo)
}

func (s *Setup) ValidateSource(config SourceConfig) error {
	switch config.Type {
	case SourceTypeStarter:
		return nil

	case SourceTypeBYOP:
		switch config.BYOPSubType {
		case BYOPSubTypeRemote:
			if config.RemoteURL == "" {
				return fmt.Errorf("remote URL is required")
			}
		case BYOPSubTypeLocal:
			if config.LocalPath == "" {
				return fmt.Errorf("local path is required")
			}
			expanded := expandPath(config.LocalPath)
			if !pathExists(expanded) {
				return fmt.Errorf("local path does not exist: %s", expanded)
			}
		}

	case SourceTypeForkUser:
		if config.ForkFrom == "" {
			return fmt.Errorf("fork source URL is required")
		}
	}

	return nil
}
