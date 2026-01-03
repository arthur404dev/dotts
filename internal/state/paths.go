package state

import (
	"os"
	"path/filepath"
)

const AppName = "dotts"

type Paths struct {
	DataDir    string
	ConfigDir  string
	CacheDir   string
	StateFile  string
	ConfigRepo string
	LogsDir    string
	BackupsDir string
}

func GetPaths() *Paths {
	dataDir := getXDGDataHome()
	configDir := getXDGConfigHome()
	cacheDir := getXDGCacheHome()

	appDataDir := filepath.Join(dataDir, AppName)

	return &Paths{
		DataDir:    appDataDir,
		ConfigDir:  filepath.Join(configDir, AppName),
		CacheDir:   filepath.Join(cacheDir, AppName),
		StateFile:  filepath.Join(appDataDir, "state.json"),
		ConfigRepo: filepath.Join(appDataDir, "config"),
		LogsDir:    filepath.Join(appDataDir, "logs"),
		BackupsDir: filepath.Join(appDataDir, "backups"),
	}
}

func (p *Paths) EnsureDirectories() error {
	dirs := []string{
		p.DataDir,
		p.ConfigDir,
		p.CacheDir,
		p.LogsDir,
		p.BackupsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func getXDGDataHome() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share")
}

func getXDGConfigHome() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func getXDGCacheHome() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache")
}
