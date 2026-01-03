package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/arthur404dev/dotts/internal/state"
)

const (
	DefaultConfigRepo = "https://github.com/arthur404dev/dotts-config"
	DefaultBranch     = "main"
)

type Source struct {
	URL     string
	Path    string
	Branch  string
	IsLocal bool
}

func NewGitSource(url, branch string) *Source {
	if branch == "" {
		branch = DefaultBranch
	}
	return &Source{
		URL:     url,
		Branch:  branch,
		IsLocal: false,
	}
}

func NewLocalSource(path string) *Source {
	return &Source{
		Path:    path,
		IsLocal: true,
	}
}

func DefaultSource() *Source {
	return NewGitSource(DefaultConfigRepo, DefaultBranch)
}

func (s *Source) Clone(destPath string) error {
	if s.IsLocal {
		return fmt.Errorf("cannot clone a local source")
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if _, err := os.Stat(destPath); err == nil {
		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", s.Branch, s.URL, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	s.Path = destPath
	return nil
}

func (s *Source) Pull() error {
	if s.IsLocal {
		return nil
	}

	if s.Path == "" {
		return fmt.Errorf("source path not set")
	}

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = s.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return nil
}

func (s *Source) GetCurrentCommit() (string, error) {
	if s.Path == "" {
		return "", fmt.Errorf("source path not set")
	}

	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = s.Path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (s *Source) Validate() error {
	path := s.Path
	if path == "" {
		return fmt.Errorf("source path not set")
	}

	requiredDirs := []string{"configs", "profiles", "packages"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(path, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("missing required directory: %s", dir)
		}
	}

	return nil
}

func (s *Source) GetConfigPath() string {
	return s.Path
}

func LoadFromState() (*Source, error) {
	st, err := state.Load()
	if err != nil {
		return nil, err
	}

	if !st.IsInitialized() {
		return nil, fmt.Errorf("dotts not initialized")
	}

	src := &Source{
		URL:     st.ConfigSource.URL,
		Path:    st.ConfigSource.Path,
		Branch:  st.ConfigSource.Branch,
		IsLocal: st.ConfigSource.Type == state.SourceTypeLocal,
	}

	return src, nil
}

func IsValidGitURL(url string) bool {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return true
	}
	if strings.HasPrefix(url, "git@") {
		return true
	}
	if strings.HasPrefix(url, "ssh://") {
		return true
	}
	return false
}

func IsLocalPath(path string) bool {
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~") || strings.HasPrefix(path, ".") {
		return true
	}
	return false
}
