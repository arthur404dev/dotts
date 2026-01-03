package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/config"
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/pkg/vetru/styles"
	"github.com/arthur404dev/dotts/pkg/schema"
)

type FeaturesResult struct {
	Features []string
}

type FeatureOption struct {
	Name        string
	Description string
	Default     bool
	LinuxOnly   bool
	MacOSOnly   bool
}

var DefaultFeatures = []FeatureOption{
	{
		Name:        "development",
		Description: "Docker, asdf, development tools",
		Default:     true,
	},
	{
		Name:        "gaming",
		Description: "Steam, Lutris, MangoHud",
		LinuxOnly:   true,
	},
	{
		Name:        "heimdall",
		Description: "Personal shell command suite",
		LinuxOnly:   true,
	},
	{
		Name:        "ai-tools",
		Description: "OpenCode, Claude configurations",
	},
}

func RunFeaturesWizard(repoConfig *schema.RepoConfig, loader *config.Loader, sysInfo *system.SystemInfo) (*FeaturesResult, error) {
	fmt.Println()
	fmt.Println(styles.Title("Feature Selection"))

	features := getAvailableFeatures(repoConfig, sysInfo)

	if len(features) == 0 {
		fmt.Println(styles.Mute("No additional features available."))
		return &FeaturesResult{Features: []string{}}, nil
	}

	var selectedFeatures []string

	options := make([]huh.Option[string], 0, len(features))
	for _, f := range features {
		opt := huh.NewOption(fmt.Sprintf("%s - %s", f.Name, f.Description), f.Name)
		if f.Default {
			opt = opt.Selected(true)
		}
		options = append(options, opt)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select additional features to enable").
				Options(options...).
				Value(&selectedFeatures),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	return &FeaturesResult{
		Features: selectedFeatures,
	}, nil
}

func getAvailableFeatures(repoConfig *schema.RepoConfig, sysInfo *system.SystemInfo) []FeatureOption {
	var available []FeatureOption

	configFeatures := make(map[string]bool)
	if repoConfig != nil {
		for _, f := range repoConfig.Features {
			configFeatures[f] = true
		}
	}

	for _, f := range DefaultFeatures {
		if len(configFeatures) > 0 && !configFeatures[f.Name] {
			continue
		}

		if f.LinuxOnly && !sysInfo.IsLinux() {
			continue
		}

		if f.MacOSOnly && !sysInfo.IsMacOS() {
			continue
		}

		available = append(available, f)
	}

	return available
}
