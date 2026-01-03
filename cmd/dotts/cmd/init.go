package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/arthur404dev/dotts/internal/state"
	"github.com/arthur404dev/dotts/internal/tui/styles"
	"github.com/arthur404dev/dotts/internal/tui/wizard"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap a new system",
	Long: `Initialize dotts on a new system.

This command runs an interactive wizard that guides you through:
  1. Selecting a config source (default, fork new, or custom)
  2. Choosing your machine type and profile
  3. Configuring settings and features
  4. Setting up authentication (SSH, GitHub)
  5. Installing packages and applying dotfiles`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().String("source", "", "Config source URL (skip source selection)")
	initCmd.Flags().String("machine", "", "Machine config name (skip machine selection)")
	initCmd.Flags().Bool("no-auth", false, "Skip authentication setup")
}

func runInit(cmd *cobra.Command, args []string) error {
	if state.Exists() {
		fmt.Println(styles.Warn("dotts is already initialized on this system."))
		fmt.Println(styles.Mute("Run 'dotts update' to update your configuration."))
		fmt.Println(styles.Mute("Or delete ~/.local/share/dotts/state.json to reinitialize."))
		return nil
	}

	wiz, err := wizard.New()
	if err != nil {
		return fmt.Errorf("failed to initialize wizard: %w", err)
	}

	result, err := wiz.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			fmt.Println()
			fmt.Println(styles.Warn("Setup cancelled."))
			return nil
		}
		return err
	}

	if err := wiz.SaveState(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Println()
	fmt.Println(styles.Title("Setup Summary"))
	fmt.Println(styles.StatusLine(styles.SuccessIcon, "Config Source", getSourceDescription(result.Source)))
	fmt.Println(styles.StatusLine(styles.SuccessIcon, "Machine", getMachineDescription(result.Machine)))

	if result.Settings != nil {
		fmt.Println(styles.StatusLine(styles.SuccessIcon, "Monitors", fmt.Sprintf("%d", result.Settings.Monitors)))
		if result.Settings.GitEmail != "" {
			fmt.Println(styles.StatusLine(styles.SuccessIcon, "Git Email", result.Settings.GitEmail))
		}
	}

	if len(result.Features.Features) > 0 {
		fmt.Println(styles.StatusLine(styles.SuccessIcon, "Features", fmt.Sprintf("%v", result.Features.Features)))
	}

	fmt.Println()
	fmt.Println(styles.Success("Configuration saved!"))
	fmt.Println()
	fmt.Println(styles.Info("Next steps:"))
	fmt.Println(styles.Mute("  1. Run 'dotts update' to apply your configuration"))
	fmt.Println(styles.Mute("  2. Run 'dotts doctor' to verify your setup"))

	if result.Auth.SetupSSH || result.Auth.SetupGitHub || result.Auth.SetupDoppler {
		fmt.Println()
		fmt.Println(styles.Warn("Authentication setup requested but not yet implemented."))
		fmt.Println(styles.Mute("This will be available in a future version."))
	}

	return nil
}

func getSourceDescription(source *wizard.SourceResult) string {
	switch source.Type {
	case wizard.SourceTypeDefault:
		return "Default (arthur404dev/dotts-config)"
	case wizard.SourceTypeFork:
		return fmt.Sprintf("New local config at %s", source.LocalPath)
	case wizard.SourceTypeCustom:
		if source.IsLocal {
			return fmt.Sprintf("Local: %s", source.LocalPath)
		}
		return fmt.Sprintf("Custom: %s", source.URL)
	default:
		return "Unknown"
	}
}

func getMachineDescription(machine *wizard.MachineResult) string {
	if machine.UseExisting {
		return fmt.Sprintf("Existing: %s", machine.ExistingName)
	}
	return fmt.Sprintf("%s (%s)", machine.Name, machine.Type)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
