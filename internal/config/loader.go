package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/arthur404dev/dotts/pkg/schema"
)

type Loader struct {
	basePath string
}

func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

func (l *Loader) LoadRepoConfig() (*schema.RepoConfig, error) {
	path := filepath.Join(l.basePath, "config.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &schema.RepoConfig{
				Name:    "dotfiles",
				Version: "1.0.0",
			}, nil
		}
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var config schema.RepoConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	return &config, nil
}

func (l *Loader) LoadProfile(name string) (*schema.Profile, error) {
	path := filepath.Join(l.basePath, "profiles", name+".yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile %s: %w", name, err)
	}

	var profile schema.Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile %s: %w", name, err)
	}

	profile.Name = name
	return &profile, nil
}

func (l *Loader) LoadMachine(name string) (*schema.Machine, error) {
	path := filepath.Join(l.basePath, "machines", name+".yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read machine %s: %w", name, err)
	}

	var machine schema.Machine
	if err := yaml.Unmarshal(data, &machine); err != nil {
		return nil, fmt.Errorf("failed to parse machine %s: %w", name, err)
	}

	return &machine, nil
}

func (l *Loader) LoadPackages(name string) (*schema.PackageManifest, error) {
	path := filepath.Join(l.basePath, "packages", name+".yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read packages %s: %w", name, err)
	}

	var packages schema.PackageManifest
	if err := yaml.Unmarshal(data, &packages); err != nil {
		return nil, fmt.Errorf("failed to parse packages %s: %w", name, err)
	}

	return &packages, nil
}

func (l *Loader) ListProfiles() ([]string, error) {
	return l.listYAMLFiles("profiles")
}

func (l *Loader) ListMachines() ([]string, error) {
	return l.listYAMLFiles("machines")
}

func (l *Loader) ListPackageManifests() ([]string, error) {
	return l.listYAMLFiles("packages")
}

func (l *Loader) ListConfigs() ([]string, error) {
	configsPath := filepath.Join(l.basePath, "configs")

	entries, err := os.ReadDir(configsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configs directory: %w", err)
	}

	var configs []string
	for _, entry := range entries {
		if entry.IsDir() {
			configs = append(configs, entry.Name())
		}
	}

	return configs, nil
}

func (l *Loader) GetConfigPath(name string) string {
	return filepath.Join(l.basePath, "configs", name)
}

func (l *Loader) GetScriptPath(name string) string {
	return filepath.Join(l.basePath, "scripts", name)
}

func (l *Loader) GetAssetsPath() string {
	return filepath.Join(l.basePath, "assets")
}

func (l *Loader) listYAMLFiles(subdir string) ([]string, error) {
	dirPath := filepath.Join(l.basePath, subdir)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s directory: %w", subdir, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			names = append(names, strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml"))
		}
	}

	return names, nil
}

func (l *Loader) ConfigExists(name string) bool {
	path := l.GetConfigPath(name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (l *Loader) ProfileExists(name string) bool {
	path := filepath.Join(l.basePath, "profiles", name+".yaml")
	_, err := os.Stat(path)
	return err == nil
}

func (l *Loader) MachineExists(name string) bool {
	path := filepath.Join(l.basePath, "machines", name+".yaml")
	_, err := os.Stat(path)
	return err == nil
}

func (l *Loader) PackagesExist(name string) bool {
	path := filepath.Join(l.basePath, "packages", name+".yaml")
	_, err := os.Stat(path)
	return err == nil
}
