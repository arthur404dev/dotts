package installer

import (
	"context"
	"os"
	"strings"
)

type AptInstaller struct {
	BaseInstaller
}

func NewAptInstaller() *AptInstaller {
	return &AptInstaller{
		BaseInstaller: BaseInstaller{name: "apt"},
	}
}

func (a *AptInstaller) Available() bool {
	return commandExists("apt-get")
}

func (a *AptInstaller) Install(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"install", "-y"}
	args = append(args, packages...)

	cmd, cmdArgs := sudoWrap(a.NeedsSudo(), "apt-get", args)
	_, err := runCommand(ctx, cmd, cmdArgs...)
	if err != nil {
		return &InstallError{Installer: a.name, Cause: err}
	}
	return nil
}

func (a *AptInstaller) Remove(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"remove", "-y"}
	args = append(args, packages...)

	cmd, cmdArgs := sudoWrap(a.NeedsSudo(), "apt-get", args)
	_, err := runCommand(ctx, cmd, cmdArgs...)
	if err != nil {
		return &InstallError{Installer: a.name, Cause: err}
	}
	return nil
}

func (a *AptInstaller) IsInstalled(pkg string) bool {
	output, err := runCommand(context.Background(), "dpkg", "-s", pkg)
	if err != nil {
		return false
	}
	return strings.Contains(output, "Status: install ok installed")
}

func (a *AptInstaller) Update(ctx context.Context) error {
	cmd, args := sudoWrap(a.NeedsSudo(), "apt-get", []string{"update"})
	_, err := runCommand(ctx, cmd, args...)
	return err
}

func (a *AptInstaller) NeedsSudo() bool {
	return os.Geteuid() != 0
}
