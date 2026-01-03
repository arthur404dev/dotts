package wizard

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/pkg/vetru/styles"
)

type SettingsResult struct {
	Monitors int
}

func RunSettingsWizard(machineType MachineType, sysInfo *system.SystemInfo) (*SettingsResult, error) {
	fmt.Println()
	fmt.Println(styles.Title("Configuration Settings"))

	var monitorsStr string

	defaultMonitors := "1"
	if machineType == MachineTypeDesktop {
		defaultMonitors = "3"
	}
	monitorsStr = defaultMonitors

	groups := []*huh.Group{}

	if machineType == MachineTypeDesktop || machineType == MachineTypeNotebook {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("How many monitors?").
				Value(&monitorsStr).
				Placeholder(defaultMonitors).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					n, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("must be a number")
					}
					if n < 1 || n > 10 {
						return fmt.Errorf("must be between 1 and 10")
					}
					return nil
				}),
		))
	}

	if len(groups) == 0 {
		return &SettingsResult{Monitors: 1}, nil
	}

	form := huh.NewForm(groups...).WithTheme(styles.GetHuhTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	monitors := 1
	if monitorsStr != "" {
		monitors, _ = strconv.Atoi(monitorsStr)
	}

	return &SettingsResult{
		Monitors: monitors,
	}, nil
}

func isValidEmail(email string) bool {
	hasAt := false
	hasDot := false
	atPos := -1

	for i, c := range email {
		if c == '@' {
			if hasAt {
				return false
			}
			hasAt = true
			atPos = i
		}
		if c == '.' && hasAt && i > atPos+1 {
			hasDot = true
		}
	}

	return hasAt && hasDot && atPos > 0 && atPos < len(email)-2
}
