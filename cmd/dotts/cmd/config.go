package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config source",
	Long: `Manage the dotts configuration source.

Without subcommands, shows the current config source.
Use subcommands to change or sync the config.`,
	RunE: runConfig,
}

var configSetCmd = &cobra.Command{
	Use:   "set <url>",
	Short: "Set config source URL",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigSet,
}

var configPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest from config repo",
	RunE:  runConfigPull,
}

var configPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local changes to config repo",
	RunE:  runConfigPush,
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPullCmd)
	configCmd.AddCommand(configPushCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	fmt.Println("Current config source:")
	fmt.Println("(Config info will be displayed here)")
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	url := args[0]
	fmt.Printf("Setting config source to: %s\n", url)
	return nil
}

func runConfigPull(cmd *cobra.Command, args []string) error {
	fmt.Println("Pulling latest config...")
	return nil
}

func runConfigPush(cmd *cobra.Command, args []string) error {
	fmt.Println("Pushing local changes...")
	return nil
}
