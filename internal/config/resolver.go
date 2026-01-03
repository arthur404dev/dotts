package config

import (
	"fmt"

	"github.com/arthur404dev/dotts/pkg/schema"
)

type ResolvedConfig struct {
	Configs  []string
	Packages *schema.PackageManifest
	Settings map[string]any
	Features []string
	Scripts  schema.ProfileScripts
}

type Resolver struct {
	loader   *Loader
	resolved map[string]*schema.Profile
}

func NewResolver(loader *Loader) *Resolver {
	return &Resolver{
		loader:   loader,
		resolved: make(map[string]*schema.Profile),
	}
}

func (r *Resolver) ResolveMachine(machineName string) (*ResolvedConfig, error) {
	machine, err := r.loader.LoadMachine(machineName)
	if err != nil {
		return nil, err
	}

	result := &ResolvedConfig{
		Configs:  []string{},
		Packages: &schema.PackageManifest{},
		Settings: make(map[string]any),
		Features: machine.Features,
	}

	for _, inherit := range machine.Inherits {
		if err := r.resolveProfile(inherit, result); err != nil {
			return nil, fmt.Errorf("failed to resolve profile %s: %w", inherit, err)
		}
	}

	for k, v := range machine.Settings {
		result.Settings[k] = v
	}

	return result, nil
}

func (r *Resolver) ResolveProfile(profileName string) (*ResolvedConfig, error) {
	result := &ResolvedConfig{
		Configs:  []string{},
		Packages: &schema.PackageManifest{},
		Settings: make(map[string]any),
		Features: []string{},
	}

	if err := r.resolveProfile(profileName, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resolver) resolveProfile(name string, result *ResolvedConfig) error {
	if _, already := r.resolved[name]; already {
		return nil
	}

	profile, err := r.loader.LoadProfile(name)
	if err != nil {
		return err
	}

	r.resolved[name] = profile

	for _, inherit := range profile.Inherits {
		if err := r.resolveProfile(inherit, result); err != nil {
			return fmt.Errorf("failed to resolve inherited profile %s: %w", inherit, err)
		}
	}

	for _, cfg := range profile.Configs {
		if !contains(result.Configs, cfg) {
			result.Configs = append(result.Configs, cfg)
		}
	}

	for _, pkgGroup := range profile.Packages {
		if r.loader.PackagesExist(pkgGroup) {
			pkgManifest, err := r.loader.LoadPackages(pkgGroup)
			if err != nil {
				return fmt.Errorf("failed to load package group %s: %w", pkgGroup, err)
			}
			result.Packages.Merge(pkgManifest)
		}
	}

	for k, v := range profile.Settings {
		if _, exists := result.Settings[k]; !exists {
			result.Settings[k] = v
		}
	}

	result.Scripts.PreInstall = append(result.Scripts.PreInstall, profile.Scripts.PreInstall...)
	result.Scripts.PostInstall = append(result.Scripts.PostInstall, profile.Scripts.PostInstall...)
	result.Scripts.PreUpdate = append(result.Scripts.PreUpdate, profile.Scripts.PreUpdate...)
	result.Scripts.PostUpdate = append(result.Scripts.PostUpdate, profile.Scripts.PostUpdate...)

	return nil
}

func (r *Resolver) GetInheritanceChain(machineName string) ([]string, error) {
	machine, err := r.loader.LoadMachine(machineName)
	if err != nil {
		return nil, err
	}

	var chain []string
	visited := make(map[string]bool)

	var walk func(profiles []string) error
	walk = func(profiles []string) error {
		for _, name := range profiles {
			if visited[name] {
				continue
			}
			visited[name] = true

			profile, err := r.loader.LoadProfile(name)
			if err != nil {
				return err
			}

			if err := walk(profile.Inherits); err != nil {
				return err
			}

			chain = append(chain, name)
		}
		return nil
	}

	if err := walk(machine.Inherits); err != nil {
		return nil, err
	}

	return chain, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
