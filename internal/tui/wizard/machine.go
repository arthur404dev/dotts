package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/config"
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/pkg/vetru/styles"
)

type MachineType string

const (
	MachineTypeDesktop  MachineType = "desktop"
	MachineTypeNotebook MachineType = "notebook"
	MachineTypeServer   MachineType = "server"
	MachineTypeVM       MachineType = "vm"
	MachineTypeMacOS    MachineType = "macos"
	MachineTypeCustom   MachineType = "custom"
)

type MachineResult struct {
	Name         string
	Type         MachineType
	Hostname     string
	UseExisting  bool
	ExistingName string
}

func RunMachineWizard(loader *config.Loader, sysInfo *system.SystemInfo) (*MachineResult, error) {
	fmt.Println()
	fmt.Println(styles.Title("Machine Setup"))

	existingMachines, _ := loader.ListMachines()

	if len(existingMachines) > 0 {
		return runMachineSelectionWizard(existingMachines, sysInfo)
	}

	return runNewMachineWizard(sysInfo)
}

func runMachineSelectionWizard(existing []string, sysInfo *system.SystemInfo) (*MachineResult, error) {
	var choice string

	options := []huh.Option[string]{
		huh.NewOption("Create a new machine configuration", "new"),
	}

	for _, m := range existing {
		options = append(options, huh.NewOption(fmt.Sprintf("Use existing: %s", m), m))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Machine Configuration").
				Description(fmt.Sprintf("Found %d existing machine config(s)", len(existing))).
				Options(options...).
				Value(&choice),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	if choice == "new" {
		return runNewMachineWizard(sysInfo)
	}

	return &MachineResult{
		UseExisting:  true,
		ExistingName: choice,
	}, nil
}

func runNewMachineWizard(sysInfo *system.SystemInfo) (*MachineResult, error) {
	var (
		machineType string
		hostname    string
		machineName string
	)

	hostname = sysInfo.Hostname

	typeOptions := []huh.Option[string]{
		huh.NewOption("Desktop (multi-monitor, full WM setup)", "desktop"),
		huh.NewOption("Notebook (single monitor, portable)", "notebook"),
		huh.NewOption("Server (minimal, no GUI)", "server"),
		huh.NewOption("VM (lightweight)", "vm"),
	}

	if sysInfo.IsMacOS() {
		typeOptions = append([]huh.Option[string]{
			huh.NewOption("macOS (work machine)", "macos"),
		}, typeOptions...)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What type of machine is this?").
				Options(typeOptions...).
				Value(&machineType),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Hostname for this machine").
				Value(&hostname).
				Placeholder(sysInfo.Hostname),
			huh.NewInput().
				Title("Machine config name").
				Description("This will be used as the filename (e.g., 'my-laptop' -> machines/my-laptop.yaml)").
				Value(&machineName).
				Placeholder(hostname).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					for _, c := range s {
						if !isValidMachineNameChar(c) {
							return fmt.Errorf("machine name can only contain letters, numbers, and hyphens")
						}
					}
					return nil
				}),
		),
	).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	if machineName == "" {
		machineName = hostname
	}

	return &MachineResult{
		Name:     machineName,
		Type:     MachineType(machineType),
		Hostname: hostname,
	}, nil
}

func isValidMachineNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}

func (m MachineType) GetDefaultProfiles() []string {
	switch m {
	case MachineTypeDesktop:
		return []string{"base", "linux", "desktop"}
	case MachineTypeNotebook:
		return []string{"base", "linux", "notebook"}
	case MachineTypeServer:
		return []string{"base", "linux", "server"}
	case MachineTypeVM:
		return []string{"base", "linux", "vm"}
	case MachineTypeMacOS:
		return []string{"base", "darwin"}
	default:
		return []string{"base"}
	}
}
