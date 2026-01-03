package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update configs and packages",
	Long: `Pull the latest configuration and apply changes.

This command will:
  1. Pull the latest changes from your config repo
  2. Show a diff of what will change
  3. Update system packages
  4. Apply new dotfile configurations
  5. Run any post-update scripts`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().Bool("dry-run", false, "Show what would change without applying")
	updateCmd.Flags().Bool("packages-only", false, "Only update packages")
	updateCmd.Flags().Bool("dotfiles-only", false, "Only update dotfiles")
	updateCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Println("Dry run mode - showing what would change...")
	} else {
		fmt.Println("Updating dotts configuration...")
	}
	return nil
}
