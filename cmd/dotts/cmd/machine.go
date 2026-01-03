package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Manage machine configurations",
	Long:  `List, switch, or create machine configurations.`,
	RunE:  runMachineList,
}

var machineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available machines",
	RunE:  runMachineList,
}

var machineSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a different machine config",
	Args:  cobra.ExactArgs(1),
	RunE:  runMachineSwitch,
}

var machineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new machine config",
	RunE:  runMachineCreate,
}

func init() {
	machineCmd.AddCommand(machineListCmd)
	machineCmd.AddCommand(machineSwitchCmd)
	machineCmd.AddCommand(machineCreateCmd)
}

func runMachineList(cmd *cobra.Command, args []string) error {
	fmt.Println("Available machines:")
	fmt.Println("(Machine list will be displayed here)")
	return nil
}

func runMachineSwitch(cmd *cobra.Command, args []string) error {
	name := args[0]
	fmt.Printf("Switching to machine: %s\n", name)
	return nil
}

func runMachineCreate(cmd *cobra.Command, args []string) error {
	fmt.Println("Creating new machine config...")
	return nil
}
