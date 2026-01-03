package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type OS string

const (
	OSLinux   OS = "linux"
	OSDarwin  OS = "darwin"
	OSUnknown OS = "unknown"
)

type Distro string

const (
	DistroArch    Distro = "arch"
	DistroDebian  Distro = "debian"
	DistroUbuntu  Distro = "ubuntu"
	DistroFedora  Distro = "fedora"
	DistroNixOS   Distro = "nixos"
	DistroMacOS   Distro = "macos"
	DistroUnknown Distro = "unknown"
)

type Arch string

const (
	ArchAMD64   Arch = "amd64"
	ArchARM64   Arch = "arm64"
	ArchUnknown Arch = "unknown"
)

type PackageManager string

const (
	PMPacman  PackageManager = "pacman"
	PMYay     PackageManager = "yay"
	PMParu    PackageManager = "paru"
	PMBrew    PackageManager = "brew"
	PMApt     PackageManager = "apt"
	PMDnf     PackageManager = "dnf"
	PMNix     PackageManager = "nix"
	PMUnknown PackageManager = "unknown"
)

type SystemInfo struct {
	OS             OS
	Distro         Distro
	Arch           Arch
	Hostname       string
	Username       string
	HomeDir        string
	PackageManager PackageManager
	HasNix         bool
	HasHomebrew    bool
	HasGit         bool
	HasCurl        bool
}

func Detect() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:      detectOS(),
		Arch:    detectArch(),
		HasNix:  commandExists("nix"),
		HasGit:  commandExists("git"),
		HasCurl: commandExists("curl"),
	}

	info.Distro = detectDistro(info.OS)
	info.PackageManager = detectPackageManager(info.OS, info.Distro)
	info.HasHomebrew = commandExists("brew")

	var err error
	info.Hostname, err = os.Hostname()
	if err != nil {
		info.Hostname = "unknown"
	}

	info.Username = os.Getenv("USER")
	if info.Username == "" {
		info.Username = "unknown"
	}

	info.HomeDir = os.Getenv("HOME")
	if info.HomeDir == "" {
		info.HomeDir, _ = os.UserHomeDir()
	}

	return info, nil
}

func detectOS() OS {
	switch runtime.GOOS {
	case "linux":
		return OSLinux
	case "darwin":
		return OSDarwin
	default:
		return OSUnknown
	}
}

func detectArch() Arch {
	switch runtime.GOARCH {
	case "amd64":
		return ArchAMD64
	case "arm64":
		return ArchARM64
	default:
		return ArchUnknown
	}
}

func detectDistro(osType OS) Distro {
	if osType == OSDarwin {
		return DistroMacOS
	}

	if osType != OSLinux {
		return DistroUnknown
	}

	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return DistroArch
	}

	if _, err := os.Stat("/etc/nixos"); err == nil {
		return DistroNixOS
	}

	osRelease, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return DistroUnknown
	}

	content := strings.ToLower(string(osRelease))

	if strings.Contains(content, "id=arch") || strings.Contains(content, "id_like=arch") ||
		strings.Contains(content, "id=cachyos") || strings.Contains(content, "id=endeavouros") ||
		strings.Contains(content, "id=manjaro") {
		return DistroArch
	}

	if strings.Contains(content, "id=ubuntu") {
		return DistroUbuntu
	}

	if strings.Contains(content, "id=debian") {
		return DistroDebian
	}

	if strings.Contains(content, "id=fedora") {
		return DistroFedora
	}

	return DistroUnknown
}

func detectPackageManager(osType OS, distro Distro) PackageManager {
	if osType == OSDarwin {
		if commandExists("brew") {
			return PMBrew
		}
		return PMUnknown
	}

	switch distro {
	case DistroArch:
		if commandExists("paru") {
			return PMParu
		}
		if commandExists("yay") {
			return PMYay
		}
		return PMPacman

	case DistroDebian, DistroUbuntu:
		return PMApt

	case DistroFedora:
		return PMDnf

	case DistroNixOS:
		return PMNix

	default:
		return PMUnknown
	}
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func (s *SystemInfo) IsLinux() bool {
	return s.OS == OSLinux
}

func (s *SystemInfo) IsMacOS() bool {
	return s.OS == OSDarwin
}

func (s *SystemInfo) IsArch() bool {
	return s.Distro == DistroArch
}

func (s *SystemInfo) SupportsAUR() bool {
	return s.Distro == DistroArch && (s.PackageManager == PMYay || s.PackageManager == PMParu)
}

func (s *SystemInfo) String() string {
	return strings.Join([]string{
		"OS: " + string(s.OS),
		"Distro: " + string(s.Distro),
		"Arch: " + string(s.Arch),
		"Hostname: " + s.Hostname,
		"User: " + s.Username,
		"Home: " + s.HomeDir,
		"Package Manager: " + string(s.PackageManager),
	}, "\n")
}
