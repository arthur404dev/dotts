package wizard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/config"
	"github.com/arthur404dev/dotts/pkg/vetru/styles"
)

type SourceType string

const (
	SourceTypeDefault SourceType = "default"
	SourceTypeFork    SourceType = "fork"
	SourceTypeCustom  SourceType = "custom"
)

type SourceResult struct {
	Type      SourceType
	URL       string
	LocalPath string
	IsLocal   bool
}

func RunSourceWizard() (*SourceResult, error) {
	fmt.Println(styles.Banner())
	fmt.Println(styles.Title("Welcome to dotts!"))
	fmt.Println(styles.Subtitle("Let's set up your configuration source."))
	fmt.Println()

	var sourceType string

	sourceForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Where should dotts get your configurations from?").
				Options(
					huh.NewOption("Use default configs (arthur404dev/dotfiles)", "default").Selected(true),
					huh.NewOption("Start fresh with a template", "fork"),
					huh.NewOption("Use my existing config repo", "custom"),
				).
				Value(&sourceType),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := sourceForm.Run(); err != nil {
		return nil, err
	}

	switch SourceType(sourceType) {
	case SourceTypeDefault:
		return handleDefaultSource()
	case SourceTypeFork:
		return handleForkSource()
	case SourceTypeCustom:
		return handleCustomSource()
	default:
		return nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}

func handleDefaultSource() (*SourceResult, error) {
	fmt.Println()
	fmt.Println(styles.Info("Using default configuration from arthur404dev/dotfiles"))

	return &SourceResult{
		Type: SourceTypeDefault,
		URL:  config.DefaultConfigRepo,
	}, nil
}

func handleForkSource() (*SourceResult, error) {
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "dotfiles")

	var (
		localPath string
		name      string
		email     string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("New Config Setup").
				Description("We'll create a new config directory with starter templates you can customize."),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Where to create your config?").
				Value(&localPath).
				Placeholder(defaultPath).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					expanded := expandPath(s)
					parent := filepath.Dir(expanded)
					if _, err := os.Stat(parent); os.IsNotExist(err) {
						return fmt.Errorf("parent directory does not exist: %s", parent)
					}
					return nil
				}),
			huh.NewInput().
				Title("Your name (for git commits)").
				Value(&name).
				Placeholder("Your Name"),
			huh.NewInput().
				Title("Your email").
				Value(&email).
				Placeholder("you@example.com"),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	if localPath == "" {
		localPath = defaultPath
	}
	localPath = expandPath(localPath)

	fmt.Println()
	fmt.Println(styles.Info(fmt.Sprintf("Creating new config at %s", localPath)))
	fmt.Println(styles.Mute("After setup, you can push this to GitHub with:"))
	fmt.Println(styles.Mute(fmt.Sprintf("  cd %s && git remote add origin <your-repo-url>", localPath)))
	fmt.Println(styles.Mute("  git push -u origin main"))

	return &SourceResult{
		Type:      SourceTypeFork,
		LocalPath: localPath,
		IsLocal:   true,
	}, nil
}

func handleCustomSource() (*SourceResult, error) {
	var url string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your config repo URL").
				Description("Supports GitHub URLs, git@ SSH URLs, or local paths").
				Value(&url).
				Placeholder("https://github.com/username/dotfiles").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("URL is required")
					}
					if !config.IsValidGitURL(s) && !config.IsLocalPath(s) {
						return fmt.Errorf("invalid URL or path format")
					}
					return nil
				}),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	isLocal := config.IsLocalPath(url)

	if isLocal {
		url = expandPath(url)
		fmt.Println()
		fmt.Println(styles.Info(fmt.Sprintf("Using local config from %s", url)))

		return &SourceResult{
			Type:      SourceTypeCustom,
			LocalPath: url,
			IsLocal:   true,
		}, nil
	}

	fmt.Println()
	fmt.Println(styles.Info(fmt.Sprintf("Using custom config from %s", url)))

	return &SourceResult{
		Type: SourceTypeCustom,
		URL:  url,
	}, nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs
		}
	}
	return path
}
