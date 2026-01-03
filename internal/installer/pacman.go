package installer

import (
	"context"
	"os"
	"strings"
)

type PacmanInstaller struct {
	BaseInstaller
}

func NewPacmanInstaller() *PacmanInstaller {
	return &PacmanInstaller{
		BaseInstaller: BaseInstaller{name: "pacman"},
	}
}

func (p *PacmanInstaller) Available() bool {
	return commandExists("pacman")
}

func (p *PacmanInstaller) Install(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"-S", "--noconfirm", "--needed"}
	args = append(args, packages...)

	cmd, cmdArgs := sudoWrap(p.NeedsSudo(), "pacman", args)
	_, err := runCommand(ctx, cmd, cmdArgs...)
	if err != nil {
		return &InstallError{Installer: p.name, Cause: err}
	}
	return nil
}

func (p *PacmanInstaller) Remove(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"-R", "--noconfirm"}
	args = append(args, packages...)

	cmd, cmdArgs := sudoWrap(p.NeedsSudo(), "pacman", args)
	_, err := runCommand(ctx, cmd, cmdArgs...)
	if err != nil {
		return &InstallError{Installer: p.name, Cause: err}
	}
	return nil
}

func (p *PacmanInstaller) IsInstalled(pkg string) bool {
	err := runCommandSilent(context.Background(), "pacman", "-Qi", pkg)
	return err == nil
}

func (p *PacmanInstaller) Update(ctx context.Context) error {
	cmd, args := sudoWrap(p.NeedsSudo(), "pacman", []string{"-Sy"})
	_, err := runCommand(ctx, cmd, args...)
	return err
}

func (p *PacmanInstaller) NeedsSudo() bool {
	return os.Geteuid() != 0
}

type YayInstaller struct {
	BaseInstaller
	helper string
}

func NewYayInstaller() *YayInstaller {
	helper := "yay"
	if commandExists("paru") {
		helper = "paru"
	}
	return &YayInstaller{
		BaseInstaller: BaseInstaller{name: "yay"},
		helper:        helper,
	}
}

func (y *YayInstaller) Available() bool {
	return commandExists(y.helper)
}

func (y *YayInstaller) Install(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"-S", "--noconfirm", "--needed"}
	args = append(args, packages...)

	_, err := runCommand(ctx, y.helper, args...)
	if err != nil {
		return &InstallError{Installer: y.name, Cause: err}
	}
	return nil
}

func (y *YayInstaller) Remove(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	args := []string{"-R", "--noconfirm"}
	args = append(args, packages...)

	_, err := runCommand(ctx, y.helper, args...)
	if err != nil {
		return &InstallError{Installer: y.name, Cause: err}
	}
	return nil
}

func (y *YayInstaller) IsInstalled(pkg string) bool {
	output, err := runCommand(context.Background(), y.helper, "-Qi", pkg)
	if err != nil {
		output, err = runCommand(context.Background(), y.helper, "-Ss", "^"+pkg+"$")
		if err != nil {
			return false
		}
		return strings.Contains(output, "[installed]")
	}
	return true
}

func (y *YayInstaller) Update(ctx context.Context) error {
	_, err := runCommand(ctx, y.helper, "-Sy")
	return err
}

func (y *YayInstaller) NeedsSudo() bool {
	return false
}
