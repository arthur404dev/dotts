package schema

type PackageManifest struct {
	Nix    []string          `yaml:"nix,omitempty"`
	System SystemPackages    `yaml:"system,omitempty"`
	AUR    []string          `yaml:"aur,omitempty"`
	Brew   []string          `yaml:"brew,omitempty"`
	Cask   []string          `yaml:"cask,omitempty"`
	Asdf   map[string]string `yaml:"asdf,omitempty"`
}

type SystemPackages struct {
	Arch   []string `yaml:"arch,omitempty"`
	Debian []string `yaml:"debian,omitempty"`
	Ubuntu []string `yaml:"ubuntu,omitempty"`
	Fedora []string `yaml:"fedora,omitempty"`
	Darwin []string `yaml:"darwin,omitempty"`
}

func (p *PackageManifest) GetSystemPackagesFor(distro string) []string {
	switch distro {
	case "arch":
		return p.System.Arch
	case "debian":
		return p.System.Debian
	case "ubuntu":
		return p.System.Ubuntu
	case "fedora":
		return p.System.Fedora
	case "macos", "darwin":
		return p.System.Darwin
	default:
		return nil
	}
}

func (p *PackageManifest) HasNixPackages() bool {
	return len(p.Nix) > 0
}

func (p *PackageManifest) HasSystemPackages(distro string) bool {
	return len(p.GetSystemPackagesFor(distro)) > 0
}

func (p *PackageManifest) HasAURPackages() bool {
	return len(p.AUR) > 0
}

func (p *PackageManifest) HasBrewPackages() bool {
	return len(p.Brew) > 0
}

func (p *PackageManifest) HasCaskPackages() bool {
	return len(p.Cask) > 0
}

func (p *PackageManifest) HasAsdfTools() bool {
	return len(p.Asdf) > 0
}

func (p *PackageManifest) Merge(other *PackageManifest) {
	if other == nil {
		return
	}

	p.Nix = appendUnique(p.Nix, other.Nix)
	p.AUR = appendUnique(p.AUR, other.AUR)
	p.Brew = appendUnique(p.Brew, other.Brew)
	p.Cask = appendUnique(p.Cask, other.Cask)

	p.System.Arch = appendUnique(p.System.Arch, other.System.Arch)
	p.System.Debian = appendUnique(p.System.Debian, other.System.Debian)
	p.System.Ubuntu = appendUnique(p.System.Ubuntu, other.System.Ubuntu)
	p.System.Fedora = appendUnique(p.System.Fedora, other.System.Fedora)
	p.System.Darwin = appendUnique(p.System.Darwin, other.System.Darwin)

	if p.Asdf == nil {
		p.Asdf = make(map[string]string)
	}
	for k, v := range other.Asdf {
		if _, exists := p.Asdf[k]; !exists {
			p.Asdf[k] = v
		}
	}
}

func appendUnique(slice []string, items []string) []string {
	seen := make(map[string]bool)
	for _, s := range slice {
		seen[s] = true
	}

	for _, item := range items {
		if !seen[item] {
			slice = append(slice, item)
			seen[item] = true
		}
	}

	return slice
}
