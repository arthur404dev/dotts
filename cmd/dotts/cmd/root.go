package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Set via ldflags at build time
var (
	Version   = "dev"
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "dotts",
	Short: "Universal Dotfiles Manager",
	Long: `dotts - Universal Dotfiles Manager

A powerful CLI tool for managing dotfiles across multiple machines
and platforms. dotts handles bootstrapping new systems, managing
configurations, and keeping everything in sync.

Features:
  • Interactive bootstrap wizard for new systems
  • Support for multiple config sources (default, fork, custom)
  • Cross-platform support (Linux, macOS)
  • Profile-based configuration inheritance
  • Automatic package management (Nix, Homebrew, pacman/yay)
  • Smart updates with diff preview

Get started:
  dotts init          Bootstrap a new system
  dotts update        Update configs and packages
  dotts status        Show current state

Documentation: https://dotts.sh/docs`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf(`dotts version %s
Built: %s
`, Version, BuildTime))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(machineCmd)
	rootCmd.AddCommand(syncCmd)
}

func isVerbose() bool {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	return verbose
}

func verboseLog(format string, args ...interface{}) {
	if isVerbose() {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}
