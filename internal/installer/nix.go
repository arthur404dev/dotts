package installer

import (
	"context"
	"strings"
)

type NixInstaller struct {
	BaseInstaller
	useProfile bool
}

func NewNixInstaller() *NixInstaller {
	n := &NixInstaller{
		BaseInstaller: BaseInstaller{name: "nix"},
	}
	n.useProfile = n.hasNixProfile()
	return n
}

func (n *NixInstaller) Available() bool {
	return commandExists("nix")
}

func (n *NixInstaller) hasNixProfile() bool {
	return commandExists("nix") && n.Available()
}

func (n *NixInstaller) Install(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	if n.useProfile {
		return n.installWithProfile(ctx, packages)
	}
	return n.installWithEnv(ctx, packages)
}

func (n *NixInstaller) installWithProfile(ctx context.Context, packages []string) error {
	args := []string{"profile", "install", "--impure"}
	for _, pkg := range packages {
		args = append(args, "nixpkgs#"+pkg)
	}

	_, err := runCommand(ctx, "nix", args...)
	if err != nil {
		return &InstallError{Installer: n.name, Cause: err}
	}
	return nil
}

func (n *NixInstaller) installWithEnv(ctx context.Context, packages []string) error {
	args := []string{"-iA"}
	for _, pkg := range packages {
		args = append(args, "nixpkgs."+pkg)
	}

	_, err := runCommand(ctx, "nix-env", args...)
	if err != nil {
		return &InstallError{Installer: n.name, Cause: err}
	}
	return nil
}

func (n *NixInstaller) Remove(ctx context.Context, packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	if n.useProfile {
		for _, pkg := range packages {
			args := []string{"profile", "remove", "nixpkgs#" + pkg}
			if _, err := runCommand(ctx, "nix", args...); err != nil {
				return &InstallError{Installer: n.name, Package: pkg, Cause: err}
			}
		}
		return nil
	}

	args := append([]string{"-e"}, packages...)
	_, err := runCommand(ctx, "nix-env", args...)
	if err != nil {
		return &InstallError{Installer: n.name, Cause: err}
	}
	return nil
}

func (n *NixInstaller) IsInstalled(pkg string) bool {
	if n.useProfile {
		output, err := runCommand(context.Background(), "nix", "profile", "list")
		if err != nil {
			return false
		}
		return strings.Contains(output, pkg)
	}

	output, err := runCommand(context.Background(), "nix-env", "-q")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, pkg) || strings.Contains(line, pkg) {
			return true
		}
	}
	return false
}

func (n *NixInstaller) Update(ctx context.Context) error {
	_, err := runCommand(ctx, "nix-channel", "--update")
	return err
}

func (n *NixInstaller) NeedsSudo() bool {
	return false
}
