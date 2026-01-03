package installer

import (
	"context"
	"strings"
)

type BrewInstaller struct {
	BaseInstaller
	caskMode bool
}

func NewBrewInstaller() *BrewInstaller {
	return &BrewInstaller{
		BaseInstaller: BaseInstaller{name: "brew"},
		caskMode:      false,
	}
}

func NewCaskInstaller() *BrewInstaller {
	return &BrewInstaller{
		BaseInstaller: BaseInstaller{name: "cask"},
		caskMode:      true,
	}
}

func (b *BrewInstaller) Available() bool {
	return commandExists("brew")
}

func (b *BrewInstaller) Install(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"install"}
	if b.caskMode {
		args = append(args, "--cask")
	}
	args = append(args, packages...)

	_, err := runCommand(ctx, "brew", args...)
	if err != nil {
		return &InstallError{Installer: b.name, Cause: err}
	}
	return nil
}

func (b *BrewInstaller) Remove(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"uninstall"}
	if b.caskMode {
		args = append(args, "--cask")
	}
	args = append(args, packages...)

	_, err := runCommand(ctx, "brew", args...)
	if err != nil {
		return &InstallError{Installer: b.name, Cause: err}
	}
	return nil
}

func (b *BrewInstaller) IsInstalled(pkg string) bool {
	var args []string
	if b.caskMode {
		args = []string{"list", "--cask", pkg}
	} else {
		args = []string{"list", pkg}
	}

	err := runCommandSilent(context.Background(), "brew", args...)
	return err == nil
}

func (b *BrewInstaller) Update(ctx context.Context) error {
	_, err := runCommand(ctx, "brew", "update")
	return err
}

func (b *BrewInstaller) NeedsSudo() bool {
	return false
}

func (b *BrewInstaller) InstalledPackages() ([]string, error) {
	var args []string
	if b.caskMode {
		args = []string{"list", "--cask", "-1"}
	} else {
		args = []string{"list", "-1"}
	}

	output, err := runCommand(context.Background(), "brew", args...)
	if err != nil {
		return nil, err
	}

	var packages []string
	for _, line := range strings.Split(output, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			packages = append(packages, line)
		}
	}
	return packages, nil
}
