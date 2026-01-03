package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local changes to config repo",
	Long: `Sync local dotfile edits back to the config repository.

This command detects changes you've made to dotfiles and
copies them back to the config repo for committing.`,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().Bool("dry-run", false, "Show what would be synced")
	syncCmd.Flags().BoolP("all", "a", false, "Sync all changed files")
}

func runSync(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Println("Dry run - files that would be synced:")
	} else {
		fmt.Println("Syncing local changes to config repo...")
	}
	return nil
}
