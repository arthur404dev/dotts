package installer

import (
	"context"
	"sync"

	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/pkg/schema"
)

type Registry struct {
	installers map[string]Installer
	sysInfo    *system.SystemInfo
	mu         sync.RWMutex
}

func NewRegistry(sysInfo *system.SystemInfo) *Registry {
	r := &Registry{
		installers: make(map[string]Installer),
		sysInfo:    sysInfo,
	}
	r.registerDefaults()
	return r
}

func (r *Registry) registerDefaults() {
	r.Register(NewNixInstaller())
	r.Register(NewPacmanInstaller())
	r.Register(NewYayInstaller())
	r.Register(NewAptInstaller())
	r.Register(NewBrewInstaller())
}

func (r *Registry) Register(i Installer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.installers[i.Name()] = i
}

func (r *Registry) Get(name string) (Installer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.installers[name]
	return i, ok
}

func (r *Registry) Available() []Installer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []Installer
	for _, i := range r.installers {
		if i.Available() {
			available = append(available, i)
		}
	}
	return available
}

func (r *Registry) CreatePlan(manifest *schema.PackageManifest) *InstallPlan {
	plan := &InstallPlan{}

	plan.Nix = manifest.Nix

	switch r.sysInfo.Distro {
	case system.DistroArch:
		plan.Pacman = manifest.System.Arch
		plan.AUR = manifest.AUR
	case system.DistroDebian:
		plan.Apt = manifest.System.Debian
	case system.DistroUbuntu:
		if len(manifest.System.Ubuntu) > 0 {
			plan.Apt = manifest.System.Ubuntu
		} else {
			plan.Apt = manifest.System.Debian
		}
	case system.DistroFedora:
		plan.Dnf = manifest.System.Fedora
	}

	if r.sysInfo.OS == system.OSDarwin {
		plan.Brew = manifest.Brew
		plan.Cask = manifest.Cask
	}

	return plan
}

type Orchestrator struct {
	registry *Registry
	progress ProgressCallback
}

func NewOrchestrator(registry *Registry, progress ProgressCallback) *Orchestrator {
	return &Orchestrator{
		registry: registry,
		progress: progress,
	}
}

func (o *Orchestrator) Execute(ctx context.Context, plan *InstallPlan) []InstallResult {
	var results []InstallResult

	type installJob struct {
		name      string
		installer Installer
		packages  []string
	}

	var jobs []installJob

	if len(plan.Nix) > 0 {
		if i, ok := o.registry.Get("nix"); ok && i.Available() {
			jobs = append(jobs, installJob{"nix", i, plan.Nix})
		}
	}

	if len(plan.Pacman) > 0 {
		if i, ok := o.registry.Get("pacman"); ok && i.Available() {
			jobs = append(jobs, installJob{"pacman", i, plan.Pacman})
		}
	}

	if len(plan.AUR) > 0 {
		if i, ok := o.registry.Get("yay"); ok && i.Available() {
			jobs = append(jobs, installJob{"yay", i, plan.AUR})
		}
	}

	if len(plan.Apt) > 0 {
		if i, ok := o.registry.Get("apt"); ok && i.Available() {
			jobs = append(jobs, installJob{"apt", i, plan.Apt})
		}
	}

	if len(plan.Dnf) > 0 {
		if i, ok := o.registry.Get("dnf"); ok && i.Available() {
			jobs = append(jobs, installJob{"dnf", i, plan.Dnf})
		}
	}

	if len(plan.Brew) > 0 {
		if i, ok := o.registry.Get("brew"); ok && i.Available() {
			jobs = append(jobs, installJob{"brew", i, plan.Brew})
		}
	}

	if len(plan.Cask) > 0 {
		if i, ok := o.registry.Get("brew"); ok && i.Available() {
			jobs = append(jobs, installJob{"cask", i, plan.Cask})
		}
	}

	for _, job := range jobs {
		result := o.runInstall(ctx, job.name, job.installer, job.packages)
		results = append(results, result)
	}

	return results
}

func (o *Orchestrator) runInstall(ctx context.Context, name string, inst Installer, packages []string) InstallResult {
	result := InstallResult{
		Installer: name,
		Requested: packages,
	}

	toInstall, alreadyInstalled := filterInstalled(inst, packages)
	result.Skipped = alreadyInstalled

	if len(toInstall) == 0 {
		return result
	}

	for i, pkg := range toInstall {
		if o.progress != nil {
			o.progress(InstallProgress{
				Installer: name,
				Package:   pkg,
				Current:   i + 1,
				Total:     len(toInstall),
				Status:    StatusRunning,
			})
		}
	}

	err := inst.Install(ctx, toInstall)
	if err != nil {
		result.Error = err
		result.Failed = toInstall
	} else {
		result.Installed = toInstall
	}

	return result
}

func (o *Orchestrator) UpdateAll(ctx context.Context) error {
	for _, inst := range o.registry.Available() {
		if err := inst.Update(ctx); err != nil {
			return err
		}
	}
	return nil
}
