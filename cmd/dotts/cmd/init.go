package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/arthur404dev/dotts/internal/apply"
	"github.com/arthur404dev/dotts/internal/state"
	"github.com/arthur404dev/dotts/internal/tui/progress"
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

var (
	initSkipApply    bool
	initDryRun       bool
	initSkipPackages bool
	initSkipDotfiles bool
)

func init() {
	initCmd.Flags().String("source", "", "Config source URL (skip source selection)")
	initCmd.Flags().String("machine", "", "Machine config name (skip machine selection)")
	initCmd.Flags().Bool("no-auth", false, "Skip authentication setup")
	initCmd.Flags().BoolVar(&initSkipApply, "skip-apply", false, "Skip package installation and dotfile linking")
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Show what would be done without making changes")
	initCmd.Flags().BoolVar(&initSkipPackages, "skip-packages", false, "Skip package installation")
	initCmd.Flags().BoolVar(&initSkipDotfiles, "skip-dotfiles", false, "Skip dotfile linking")
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

	if initSkipApply {
		fmt.Println()
		fmt.Println(styles.Success("Configuration saved!"))
		fmt.Println()
		fmt.Println(styles.Info("Next steps:"))
		fmt.Println(styles.Mute("  1. Run 'dotts update' to apply your configuration"))
		fmt.Println(styles.Mute("  2. Run 'dotts doctor' to verify your setup"))
		return nil
	}

	configPath := getConfigPath(result.Source)
	machineName := getMachineName(result.Machine)

	applier, err := apply.New(wiz.GetSystemInfo(), configPath)
	if err != nil {
		return fmt.Errorf("failed to initialize applier: %w", err)
	}

	applyResult, err := applier.Apply(context.Background(), apply.ApplyOptions{
		DryRun:       initDryRun,
		SkipPackages: initSkipPackages,
		SkipDotfiles: initSkipDotfiles,
		MachineName:  machineName,
	})
	if err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	fmt.Println()
	if applyResult.Success() {
		progress.PrintSuccess("Setup complete!")
	} else {
		progress.PrintWarning(fmt.Sprintf("Setup completed with %d error(s)", len(applyResult.Errors)))
		for _, e := range applyResult.Errors {
			progress.PrintError(e.Error())
		}
	}

	if result.Auth.SetupSSH || result.Auth.SetupGitHub || result.Auth.SetupDoppler {
		fmt.Println()
		fmt.Println(styles.Warn("Authentication setup requested but not yet implemented."))
		fmt.Println(styles.Mute("This will be available in a future version."))
	}

	fmt.Println()
	fmt.Println(styles.Info("Useful commands:"))
	fmt.Println(styles.Mute("  dotts status  - Show current configuration state"))
	fmt.Println(styles.Mute("  dotts update  - Update packages and dotfiles"))
	fmt.Println(styles.Mute("  dotts doctor  - Verify your setup"))

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

func getConfigPath(source *wizard.SourceResult) string {
	if source.IsLocal || source.Type == wizard.SourceTypeFork {
		return source.LocalPath
	}
	paths := state.GetPaths()
	return paths.ConfigRepo
}

func getMachineName(machine *wizard.MachineResult) string {
	if machine.UseExisting {
		return machine.ExistingName
	}
	return machine.Name
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
