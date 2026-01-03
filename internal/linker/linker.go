package linker

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Linker interface {
	Link(source, target string) error
	Unlink(target string) error
	IsOurLink(target string) bool
	Status() (*LinkStatus, error)
}

type LinkEntry struct {
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"created_at"`
	IsDir     bool      `json:"is_dir"`
}

type LinkStatus struct {
	Links   []LinkEntry
	Broken  []string
	Foreign []string
}

type LinkResult struct {
	Linked  []LinkEntry
	Skipped []string
	Backed  []string
	Errors  []LinkError
}

func (r *LinkResult) Success() bool {
	return len(r.Errors) == 0
}

type LinkError struct {
	Source string
	Target string
	Err    error
}

func (e LinkError) Error() string {
	return fmt.Sprintf("failed to link %s -> %s: %v", e.Source, e.Target, e.Err)
}

type LinkProgress struct {
	Source  string
	Target  string
	Current int
	Total   int
	Status  LinkProgressStatus
	Message string
}

type LinkProgressStatus int

const (
	LinkPending LinkProgressStatus = iota
	LinkRunning
	LinkSuccess
	LinkSkipped
	LinkBackedUp
	LinkFailed
)

type ProgressCallback func(progress LinkProgress)

type LinkOptions struct {
	DryRun   bool
	Force    bool
	Backup   bool
	Progress ProgressCallback
}

func DefaultLinkOptions() LinkOptions {
	return LinkOptions{
		DryRun:   false,
		Force:    false,
		Backup:   true,
		Progress: nil,
	}
}

func expandPath(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func readLink(path string) (string, error) {
	return os.Readlink(path)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func ensureParentDir(path string) error {
	parent := filepath.Dir(path)
	return os.MkdirAll(parent, 0755)
}
