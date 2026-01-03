package state

import (
	"encoding/json"
	"os"
	"time"
)

type ConfigSourceType string

const (
	SourceTypeGit   ConfigSourceType = "git"
	SourceTypeLocal ConfigSourceType = "local"
)

type ConfigSource struct {
	Type       ConfigSourceType `json:"type"`
	URL        string           `json:"url"`
	Path       string           `json:"path"`
	Branch     string           `json:"branch"`
	LastPull   time.Time        `json:"last_pull"`
	LastCommit string           `json:"last_commit"`
}

type MachineInfo struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Distro   string `json:"distro"`
	Profile  string `json:"profile"`
}

type State struct {
	Version      string         `json:"version"`
	ConfigSource ConfigSource   `json:"config_source"`
	Machine      MachineInfo    `json:"machine"`
	Settings     map[string]any `json:"settings"`
	Features     []string       `json:"features"`
	LastApply    time.Time      `json:"last_apply"`
	paths        *Paths
}

func New() *State {
	return &State{
		Version:  "1.0.0",
		Settings: make(map[string]any),
		Features: []string{},
		paths:    GetPaths(),
	}
}

func Load() (*State, error) {
	paths := GetPaths()

	data, err := os.ReadFile(paths.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, err
	}

	state := &State{paths: paths}
	if err := json.Unmarshal(data, state); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *State) Save() error {
	if err := s.paths.EnsureDirectories(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.paths.StateFile, data, 0644)
}

func (s *State) IsInitialized() bool {
	return s.ConfigSource.URL != "" || s.ConfigSource.Path != ""
}

func (s *State) HasMachine() bool {
	return s.Machine.Name != ""
}

func (s *State) SetConfigSource(sourceType ConfigSourceType, url, path, branch string) {
	s.ConfigSource = ConfigSource{
		Type:   sourceType,
		URL:    url,
		Path:   path,
		Branch: branch,
	}
}

func (s *State) SetMachine(name, hostname, osType, distro, profile string) {
	s.Machine = MachineInfo{
		Name:     name,
		Hostname: hostname,
		OS:       osType,
		Distro:   distro,
		Profile:  profile,
	}
}

func (s *State) SetSetting(key string, value any) {
	if s.Settings == nil {
		s.Settings = make(map[string]any)
	}
	s.Settings[key] = value
}

func (s *State) GetSetting(key string) any {
	if s.Settings == nil {
		return nil
	}
	return s.Settings[key]
}

func (s *State) AddFeature(feature string) {
	for _, f := range s.Features {
		if f == feature {
			return
		}
	}
	s.Features = append(s.Features, feature)
}

func (s *State) HasFeature(feature string) bool {
	for _, f := range s.Features {
		if f == feature {
			return true
		}
	}
	return false
}

func (s *State) UpdateLastApply() {
	s.LastApply = time.Now()
}

func (s *State) UpdateLastPull(commit string) {
	s.ConfigSource.LastPull = time.Now()
	s.ConfigSource.LastCommit = commit
}

func (s *State) GetPaths() *Paths {
	return s.paths
}

func Exists() bool {
	paths := GetPaths()
	_, err := os.Stat(paths.StateFile)
	return err == nil
}
