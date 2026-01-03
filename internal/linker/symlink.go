package linker

import (
	"os"
	"path/filepath"
	"time"

	"github.com/arthur404dev/dotts/internal/template"
)

type SymlinkLinker struct {
	manifest   *Manifest
	backup     *BackupManager
	configRoot string
}

func NewSymlinkLinker(dataDir, configRoot string) (*SymlinkLinker, error) {
	manifest, err := LoadManifest(dataDir)
	if err != nil {
		return nil, err
	}

	backup, err := NewBackupManager(dataDir)
	if err != nil {
		return nil, err
	}

	return &SymlinkLinker{
		manifest:   manifest,
		backup:     backup,
		configRoot: configRoot,
	}, nil
}

func (s *SymlinkLinker) LinkConfig(configName string, opts LinkOptions) (*LinkResult, error) {
	result := &LinkResult{}

	configPath := filepath.Join(s.configRoot, "configs", configName)
	if !pathExists(configPath) {
		return result, nil
	}

	return s.linkDirectory(configPath, opts)
}

func (s *SymlinkLinker) linkDirectory(sourceRoot string, opts LinkOptions) (*LinkResult, error) {
	result := &LinkResult{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(sourceRoot, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(homeDir, relPath)

		if info.IsDir() {
			if s.shouldLinkAsDirectory(sourcePath, sourceRoot) {
				if err := s.createLink(sourcePath, targetPath, true, opts, result); err != nil {
					result.Errors = append(result.Errors, LinkError{
						Source: sourcePath,
						Target: targetPath,
						Err:    err,
					})
				}
				return filepath.SkipDir
			}
			return nil
		}

		if err := s.createLink(sourcePath, targetPath, false, opts, result); err != nil {
			result.Errors = append(result.Errors, LinkError{
				Source: sourcePath,
				Target: targetPath,
				Err:    err,
			})
		}

		return nil
	})

	return result, err
}

func (s *SymlinkLinker) shouldLinkAsDirectory(dirPath, sourceRoot string) bool {
	relPath, err := filepath.Rel(sourceRoot, dirPath)
	if err != nil {
		return false
	}

	parts := filepath.SplitList(relPath)
	if len(parts) == 0 {
		return false
	}

	depth := len(filepath.SplitList(relPath))
	return depth >= 2
}

func (s *SymlinkLinker) createLink(source, target string, isDir bool, opts LinkOptions, result *LinkResult) error {
	isTemplate := false
	if !isDir && len(opts.TemplateValues) > 0 {
		isTemplate = s.hasTemplates(source)
	}

	if opts.DryRun {
		result.Linked = append(result.Linked, LinkEntry{
			Source:     source,
			Target:     target,
			CreatedAt:  time.Now(),
			IsDir:      isDir,
			IsTemplate: isTemplate,
		})
		return nil
	}

	if isSymlink(target) {
		existingSource, err := readLink(target)
		if err == nil && existingSource == source && !isTemplate {
			result.Skipped = append(result.Skipped, target)
			return nil
		}

		if s.manifest.HasEntry(target) || opts.Force {
			if err := os.Remove(target); err != nil {
				return err
			}
		} else if !isTemplate {
			result.Skipped = append(result.Skipped, target)
			return nil
		}
	}

	if pathExists(target) {
		if opts.Backup {
			backupPath, err := s.backup.Backup(target)
			if err != nil {
				return err
			}
			result.Backed = append(result.Backed, backupPath)
		}

		if err := os.RemoveAll(target); err != nil {
			return err
		}
	}

	if err := ensureParentDir(target); err != nil {
		return err
	}

	if isTemplate {
		if err := s.applyTemplate(source, target, opts.TemplateValues); err != nil {
			return err
		}
	} else {
		if err := os.Symlink(source, target); err != nil {
			return err
		}
	}

	entry := LinkEntry{
		Source:     source,
		Target:     target,
		CreatedAt:  time.Now(),
		IsDir:      isDir,
		IsTemplate: isTemplate,
	}
	result.Linked = append(result.Linked, entry)
	s.manifest.Add(entry)

	return nil
}

func (s *SymlinkLinker) hasTemplates(path string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return template.HasPlaceholdersBytes(content)
}

func (s *SymlinkLinker) applyTemplate(source, target string, values map[string]string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	processed := template.ApplyBytes(content, values)

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	return os.WriteFile(target, processed, info.Mode())
}

func (s *SymlinkLinker) Link(source, target string) error {
	source = expandPath(source)
	target = expandPath(target)

	result := &LinkResult{}
	isDir := isDir(source)

	return s.createLink(source, target, isDir, DefaultLinkOptions(), result)
}

func (s *SymlinkLinker) Unlink(target string) error {
	target = expandPath(target)

	if !isSymlink(target) {
		return nil
	}

	if !s.manifest.HasEntry(target) {
		return nil
	}

	if err := os.Remove(target); err != nil {
		return err
	}

	s.manifest.Remove(target)
	return nil
}

func (s *SymlinkLinker) UnlinkAll() error {
	for _, entry := range s.manifest.Entries() {
		if err := s.Unlink(entry.Target); err != nil {
			return err
		}
	}
	return nil
}

func (s *SymlinkLinker) IsOurLink(target string) bool {
	target = expandPath(target)
	return s.manifest.HasEntry(target)
}

func (s *SymlinkLinker) Status() (*LinkStatus, error) {
	status := &LinkStatus{}

	for _, entry := range s.manifest.Entries() {
		if !pathExists(entry.Target) {
			status.Broken = append(status.Broken, entry.Target)
			continue
		}

		if !isSymlink(entry.Target) {
			status.Foreign = append(status.Foreign, entry.Target)
			continue
		}

		actualSource, err := readLink(entry.Target)
		if err != nil || actualSource != entry.Source {
			status.Foreign = append(status.Foreign, entry.Target)
			continue
		}

		status.Links = append(status.Links, entry)
	}

	return status, nil
}

func (s *SymlinkLinker) Save() error {
	return s.manifest.Save()
}

func (s *SymlinkLinker) Restore(target string) error {
	return s.backup.Restore(target)
}
