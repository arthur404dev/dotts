package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health",
	Long: `Run diagnostic checks on your dotts installation.

Checks:
  • Required tools installed (git, nix, etc.)
  • Config repo accessible
  • Dotfiles properly linked
  • Package managers working
  • No broken symlinks`,
	RunE: runDoctor,
}

func init() {
	doctorCmd.Flags().Bool("fix", false, "Attempt to fix issues automatically")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Running dotts health checks...")
	fmt.Println("(Doctor checks will be implemented here)")
	return nil
}
