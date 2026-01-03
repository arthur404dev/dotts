package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/tui/styles"
)

type AuthResult struct {
	SetupSSH     bool
	SetupGitHub  bool
	SetupDoppler bool
}

func RunAuthWizard() (*AuthResult, error) {
	fmt.Println()
	fmt.Println(styles.Title("Authentication Setup (Optional)"))
	fmt.Println(styles.Mute("You can skip this and set up authentication later."))
	fmt.Println()

	var (
		setupSSH     bool
		setupGitHub  bool
		setupDoppler bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Set up SSH key?").
				Description("Generate an SSH key for GitHub and other services").
				Value(&setupSSH),
			huh.NewConfirm().
				Title("Authenticate with GitHub CLI?").
				Description("This will open a browser for authentication").
				Value(&setupGitHub),
			huh.NewConfirm().
				Title("Configure Doppler CLI?").
				Description("For secrets management").
				Affirmative("Yes").
				Negative("No").
				Value(&setupDoppler),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	return &AuthResult{
		SetupSSH:     setupSSH,
		SetupGitHub:  setupGitHub,
		SetupDoppler: setupDoppler,
	}, nil
}
