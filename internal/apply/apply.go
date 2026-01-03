package apply

import (
	"context"
	"fmt"

	"github.com/arthur404dev/dotts/internal/config"
	"github.com/arthur404dev/dotts/internal/installer"
	"github.com/arthur404dev/dotts/internal/linker"
	"github.com/arthur404dev/dotts/internal/state"
	"github.com/arthur404dev/dotts/internal/system"
	"github.com/arthur404dev/dotts/internal/tui/progress"
	"github.com/arthur404dev/dotts/internal/tui/styles"
)

type Applier struct {
	sysInfo    *system.SystemInfo
	paths      *state.Paths
	configPath string
	loader     *config.Loader
	resolver   *config.Resolver
	registry   *installer.Registry
	linker     *linker.SymlinkLinker
}

type ApplyOptions struct {
	DryRun       bool
	SkipPackages bool
	SkipDotfiles bool
	MachineName  string
}

type ApplyResult struct {
	PackageResults []installer.InstallResult
	LinkResult     *linker.LinkResult
	Errors         []error
}

func (r *ApplyResult) Success() bool {
	return len(r.Errors) == 0
}

func New(sysInfo *system.SystemInfo, configPath string) (*Applier, error) {
	paths := state.GetPaths()

	loader := config.NewLoader(configPath)

	lnk, err := linker.NewSymlinkLinker(paths.DataDir, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize linker: %w", err)
	}

	return &Applier{
		sysInfo:    sysInfo,
		paths:      paths,
		configPath: configPath,
		loader:     loader,
		resolver:   config.NewResolver(loader),
		registry:   installer.NewRegistry(sysInfo),
		linker:     lnk,
	}, nil
}

func (a *Applier) Apply(ctx context.Context, opts ApplyOptions) (*ApplyResult, error) {
	result := &ApplyResult{}

	resolved, err := a.resolver.ResolveMachine(opts.MachineName)
	if err != nil {
		resolved, err = a.resolver.ResolveProfile(opts.MachineName)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve configuration: %w", err)
		}
	}

	progress.PrintHeader("Applying Configuration")

	if !opts.SkipPackages && resolved.Packages != nil {
		fmt.Println()
		fmt.Println(styles.Info("Installing packages..."))

		plan := a.registry.CreatePlan(resolved.Packages)

		if !plan.IsEmpty() {
			if opts.DryRun {
				a.printPlan(plan)
			} else {
				results := a.installPackages(ctx, plan)
				result.PackageResults = results

				for _, r := range results {
					if !r.Success() {
						result.Errors = append(result.Errors, r.Error)
					}
				}
			}
		} else {
			fmt.Println(styles.Mute("  No packages to install"))
		}
	}

	if !opts.SkipDotfiles && len(resolved.Configs) > 0 {
		fmt.Println()
		fmt.Println(styles.Info("Linking dotfiles..."))

		linkOpts := linker.DefaultLinkOptions()
		linkOpts.DryRun = opts.DryRun

		for _, configName := range resolved.Configs {
			linkResult, err := a.linker.LinkConfig(configName, linkOpts)
			if err != nil {
				result.Errors = append(result.Errors, err)
				progress.PrintError(fmt.Sprintf("%s: %v", configName, err))
				continue
			}

			if result.LinkResult == nil {
				result.LinkResult = linkResult
			} else {
				result.LinkResult.Linked = append(result.LinkResult.Linked, linkResult.Linked...)
				result.LinkResult.Skipped = append(result.LinkResult.Skipped, linkResult.Skipped...)
				result.LinkResult.Backed = append(result.LinkResult.Backed, linkResult.Backed...)
				result.LinkResult.Errors = append(result.LinkResult.Errors, linkResult.Errors...)
			}

			if len(linkResult.Linked) > 0 {
				progress.PrintSuccess(fmt.Sprintf("%s: %d files linked", configName, len(linkResult.Linked)))
			} else if len(linkResult.Skipped) > 0 {
				fmt.Println(styles.Mute(fmt.Sprintf("  %s: already linked", configName)))
			}
		}

		if err := a.linker.Save(); err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	return result, nil
}

func (a *Applier) installPackages(ctx context.Context, plan *installer.InstallPlan) []installer.InstallResult {
	prog := progress.New()

	var stepIndices []int

	if len(plan.Nix) > 0 {
		idx := prog.AddStep(fmt.Sprintf("nix (%d packages)", len(plan.Nix)))
		stepIndices = append(stepIndices, idx)
	}
	if len(plan.Pacman) > 0 {
		idx := prog.AddStep(fmt.Sprintf("pacman (%d packages)", len(plan.Pacman)))
		stepIndices = append(stepIndices, idx)
	}
	if len(plan.AUR) > 0 {
		idx := prog.AddStep(fmt.Sprintf("aur (%d packages)", len(plan.AUR)))
		stepIndices = append(stepIndices, idx)
	}
	if len(plan.Apt) > 0 {
		idx := prog.AddStep(fmt.Sprintf("apt (%d packages)", len(plan.Apt)))
		stepIndices = append(stepIndices, idx)
	}
	if len(plan.Brew) > 0 {
		idx := prog.AddStep(fmt.Sprintf("brew (%d packages)", len(plan.Brew)))
		stepIndices = append(stepIndices, idx)
	}
	if len(plan.Cask) > 0 {
		idx := prog.AddStep(fmt.Sprintf("cask (%d packages)", len(plan.Cask)))
		stepIndices = append(stepIndices, idx)
	}

	fmt.Println(prog.Render())

	orchestrator := installer.NewOrchestrator(a.registry, nil)
	results := orchestrator.Execute(ctx, plan)

	for i, result := range results {
		if i < len(stepIndices) {
			if result.Success() {
				prog.SetStatus(stepIndices[i], progress.StepSuccess,
					fmt.Sprintf("%d installed, %d skipped", len(result.Installed), len(result.Skipped)))
			} else {
				prog.SetStatus(stepIndices[i], progress.StepFailed, result.Error.Error())
			}
		}
	}

	return results
}

func (a *Applier) printPlan(plan *installer.InstallPlan) {
	fmt.Println(styles.Mute("  [dry-run] Would install:"))

	if len(plan.Nix) > 0 {
		fmt.Printf("    nix: %v\n", plan.Nix)
	}
	if len(plan.Pacman) > 0 {
		fmt.Printf("    pacman: %v\n", plan.Pacman)
	}
	if len(plan.AUR) > 0 {
		fmt.Printf("    aur: %v\n", plan.AUR)
	}
	if len(plan.Apt) > 0 {
		fmt.Printf("    apt: %v\n", plan.Apt)
	}
	if len(plan.Brew) > 0 {
		fmt.Printf("    brew: %v\n", plan.Brew)
	}
	if len(plan.Cask) > 0 {
		fmt.Printf("    cask: %v\n", plan.Cask)
	}
}

func (a *Applier) GetLoader() *config.Loader {
	return a.loader
}

func (a *Applier) GetLinker() *linker.SymlinkLinker {
	return a.linker
}
