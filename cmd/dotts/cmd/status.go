package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current state",
	Long: `Display the current dotts configuration state.

Shows:
  • Config source and last update
  • Current machine and profile
  • Enabled features
  • Package status
  • Any pending changes`,
	RunE: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("dotts status")
	fmt.Println("============")
	fmt.Println("(Status display will be implemented here)")
	return nil
}
