package wizard

import (
	"fmt"

	"github.com/arthur404dev/dotts/internal/config"
	"github.com/arthur404dev/dotts/internal/state"
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/internal/tui/styles"
	"github.com/arthur404dev/dotts/pkg/schema"
)

type WizardResult struct {
	Source   *SourceResult
	Machine  *MachineResult
	Settings *SettingsResult
	Features *FeaturesResult
	Auth     *AuthResult
}

type Wizard struct {
	sysInfo *system.SystemInfo
	state   *state.State
	paths   *state.Paths
}

func New() (*Wizard, error) {
	sysInfo, err := system.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %w", err)
	}

	st := state.New()
	paths := state.GetPaths()

	return &Wizard{
		sysInfo: sysInfo,
		state:   st,
		paths:   paths,
	}, nil
}

func (w *Wizard) Run() (*WizardResult, error) {
	result := &WizardResult{}

	sourceResult, err := RunSourceWizard()
	if err != nil {
		return nil, fmt.Errorf("source wizard failed: %w", err)
	}
	result.Source = sourceResult

	configPath, err := w.setupConfigSource(sourceResult)
	if err != nil {
		return nil, fmt.Errorf("failed to setup config source: %w", err)
	}

	loader := config.NewLoader(configPath)

	if err := validateConfig(loader); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	repoConfig, _ := loader.LoadRepoConfig()

	machineResult, err := RunMachineWizard(loader, w.sysInfo)
	if err != nil {
		return nil, fmt.Errorf("machine wizard failed: %w", err)
	}
	result.Machine = machineResult

	var machineType MachineType
	if machineResult.UseExisting {
		machineType = MachineTypeCustom
	} else {
		machineType = machineResult.Type
	}

	if !machineResult.UseExisting {
		settingsResult, err := RunSettingsWizard(machineType, w.sysInfo)
		if err != nil {
			return nil, fmt.Errorf("settings wizard failed: %w", err)
		}
		result.Settings = settingsResult
	}

	featuresResult, err := RunFeaturesWizard(repoConfig, loader, w.sysInfo)
	if err != nil {
		return nil, fmt.Errorf("features wizard failed: %w", err)
	}
	result.Features = featuresResult

	authResult, err := RunAuthWizard()
	if err != nil {
		return nil, fmt.Errorf("auth wizard failed: %w", err)
	}
	result.Auth = authResult

	w.updateState(result, configPath)

	return result, nil
}

func (w *Wizard) setupConfigSource(source *SourceResult) (string, error) {
	switch source.Type {
	case SourceTypeDefault:
		fmt.Println()
		fmt.Println(styles.Info("Cloning default configuration..."))

		src := config.DefaultSource()
		if err := src.Clone(w.paths.ConfigRepo); err != nil {
			return "", err
		}

		fmt.Println(styles.Success("Configuration cloned successfully!"))
		return w.paths.ConfigRepo, nil

	case SourceTypeFork:
		fmt.Println()
		fmt.Println(styles.Info("Creating new configuration from template..."))

		return source.LocalPath, nil

	case SourceTypeCustom:
		if source.IsLocal {
			return source.LocalPath, nil
		}

		fmt.Println()
		fmt.Println(styles.Info("Cloning custom configuration..."))

		src := config.NewGitSource(source.URL, "main")
		if err := src.Clone(w.paths.ConfigRepo); err != nil {
			return "", err
		}

		fmt.Println(styles.Success("Configuration cloned successfully!"))
		return w.paths.ConfigRepo, nil

	default:
		return "", fmt.Errorf("unknown source type: %s", source.Type)
	}
}

func (w *Wizard) updateState(result *WizardResult, configPath string) {
	if result.Source.IsLocal {
		w.state.SetConfigSource(state.SourceTypeLocal, "", result.Source.LocalPath, "")
	} else {
		url := result.Source.URL
		if url == "" {
			url = config.DefaultConfigRepo
		}
		w.state.SetConfigSource(state.SourceTypeGit, url, configPath, "main")
	}

	machineName := result.Machine.Name
	if result.Machine.UseExisting {
		machineName = result.Machine.ExistingName
	}

	w.state.SetMachine(
		machineName,
		result.Machine.Hostname,
		string(w.sysInfo.OS),
		string(w.sysInfo.Distro),
		string(result.Machine.Type),
	)

	if result.Settings != nil {
		w.state.SetSetting("monitors", result.Settings.Monitors)
		w.state.SetSetting("git_email", result.Settings.GitEmail)
		w.state.SetSetting("git_name", result.Settings.GitName)
	}

	for _, feature := range result.Features.Features {
		w.state.AddFeature(feature)
	}
}

func (w *Wizard) SaveState() error {
	return w.state.Save()
}

func (w *Wizard) GetState() *state.State {
	return w.state
}

func (w *Wizard) GetSystemInfo() *system.SystemInfo {
	return w.sysInfo
}

func validateConfig(loader *config.Loader) error {
	profiles, err := loader.ListProfiles()
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		return fmt.Errorf("no profiles found in config")
	}

	if !loader.ProfileExists("base") {
		return fmt.Errorf("missing required profile: base")
	}

	fmt.Println(styles.Success(fmt.Sprintf("Found %d profile(s)", len(profiles))))

	machines, _ := loader.ListMachines()
	if len(machines) > 0 {
		fmt.Println(styles.Success(fmt.Sprintf("Found %d machine config(s)", len(machines))))
	}

	return nil
}

func SaveRepoConfig(path string, cfg *schema.RepoConfig) error {
	return nil
}
