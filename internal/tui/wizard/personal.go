package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/personal"
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/internal/tui/styles"
	"github.com/arthur404dev/dotts/pkg/schema"
)

type PersonalResult struct {
	Name   string
	Email  string
	GitHub string
}

func RunPersonalWizard(sysInfo *system.SystemInfo) (*PersonalResult, error) {
	fmt.Println()
	fmt.Println(styles.Title("Personal Settings"))
	fmt.Println(styles.Subtitle("These settings personalize your dotfiles."))
	fmt.Println()

	existing, err := personal.Load()
	if err != nil {
		existing = &schema.PersonalConfig{}
	}

	var (
		name   = existing.User.Name
		email  = existing.User.Email
		github = existing.User.GitHub
	)

	if name == "" {
		name = sysInfo.Username
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Your full name").
				Description("Used for git commits and config files").
				Value(&name).
				Placeholder(sysInfo.Username).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Your email").
				Description("Used for git commits").
				Value(&email).
				Placeholder("you@example.com").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("email is required")
					}
					if !isValidEmail(s) {
						return fmt.Errorf("invalid email format")
					}
					return nil
				}),
			huh.NewInput().
				Title("GitHub username").
				Description("Optional, used for GitHub-related configs").
				Value(&github).
				Placeholder("username"),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	result := &PersonalResult{
		Name:   name,
		Email:  email,
		GitHub: github,
	}

	if err := savePersonalConfig(result); err != nil {
		fmt.Println(styles.Warn(fmt.Sprintf("Warning: could not save personal config: %v", err)))
	} else {
		fmt.Println(styles.Success(fmt.Sprintf("Personal settings saved to %s", personal.GetPath())))
	}

	return result, nil
}

func savePersonalConfig(result *PersonalResult) error {
	config := &schema.PersonalConfig{
		User: schema.UserInfo{
			Name:   result.Name,
			Email:  result.Email,
			GitHub: result.GitHub,
		},
	}

	return personal.Save(config)
}
