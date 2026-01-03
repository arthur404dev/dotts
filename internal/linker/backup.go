package linker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupManager struct {
	backupDir string
	index     *BackupIndex
}

type BackupIndex struct {
	path    string
	Entries map[string]BackupEntry `json:"entries"`
}

type BackupEntry struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	BackedUpAt   time.Time `json:"backed_up_at"`
	IsDir        bool      `json:"is_dir"`
}

func NewBackupManager(dataDir string) (*BackupManager, error) {
	backupDir := filepath.Join(dataDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, err
	}

	index, err := loadBackupIndex(dataDir)
	if err != nil {
		return nil, err
	}

	return &BackupManager{
		backupDir: backupDir,
		index:     index,
	}, nil
}

func loadBackupIndex(dataDir string) (*BackupIndex, error) {
	path := filepath.Join(dataDir, "backup-index.json")
	index := &BackupIndex{
		path:    path,
		Entries: make(map[string]BackupEntry),
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return index, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, index); err != nil {
		return nil, err
	}

	return index, nil
}

func (b *BackupManager) Backup(path string) (string, error) {
	path = expandPath(path)

	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	baseName := filepath.Base(path)
	backupPath := filepath.Join(b.backupDir, timestamp, baseName)

	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return "", err
	}

	if info.IsDir() {
		if err := copyDir(path, backupPath); err != nil {
			return "", err
		}
	} else {
		if err := copyFile(path, backupPath); err != nil {
			return "", err
		}
	}

	b.index.Entries[path] = BackupEntry{
		OriginalPath: path,
		BackupPath:   backupPath,
		BackedUpAt:   time.Now(),
		IsDir:        info.IsDir(),
	}

	if err := b.saveIndex(); err != nil {
		return "", err
	}

	return backupPath, nil
}

func (b *BackupManager) Restore(originalPath string) error {
	originalPath = expandPath(originalPath)

	entry, ok := b.index.Entries[originalPath]
	if !ok {
		return fmt.Errorf("no backup found for %s", originalPath)
	}

	if err := os.RemoveAll(originalPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := ensureParentDir(originalPath); err != nil {
		return err
	}

	if entry.IsDir {
		if err := copyDir(entry.BackupPath, originalPath); err != nil {
			return err
		}
	} else {
		if err := copyFile(entry.BackupPath, originalPath); err != nil {
			return err
		}
	}

	delete(b.index.Entries, originalPath)
	return b.saveIndex()
}

func (b *BackupManager) HasBackup(originalPath string) bool {
	originalPath = expandPath(originalPath)
	_, ok := b.index.Entries[originalPath]
	return ok
}

func (b *BackupManager) List() []BackupEntry {
	entries := make([]BackupEntry, 0, len(b.index.Entries))
	for _, entry := range b.index.Entries {
		entries = append(entries, entry)
	}
	return entries
}

func (b *BackupManager) Clean(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)

	for path, entry := range b.index.Entries {
		if entry.BackedUpAt.Before(cutoff) {
			if err := os.RemoveAll(entry.BackupPath); err != nil {
				return err
			}
			delete(b.index.Entries, path)
		}
	}

	return b.saveIndex()
}

func (b *BackupManager) saveIndex() error {
	data, err := json.MarshalIndent(b.index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.index.path, data, 0644)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}
