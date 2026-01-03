package personal

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/arthur404dev/dotts/pkg/schema"
)

func GetPath() string {
	configDir := getXDGConfigHome()
	return filepath.Join(configDir, "dotts", "personal.yaml")
}

func Exists() bool {
	_, err := os.Stat(GetPath())
	return err == nil
}

func Load() (*schema.PersonalConfig, error) {
	path := GetPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &schema.PersonalConfig{}, nil
		}
		return nil, fmt.Errorf("failed to read personal config: %w", err)
	}

	var config schema.PersonalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse personal config: %w", err)
	}

	return &config, nil
}

func Save(config *schema.PersonalConfig) error {
	path := GetPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal personal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write personal config: %w", err)
	}

	return nil
}

func getXDGConfigHome() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}
