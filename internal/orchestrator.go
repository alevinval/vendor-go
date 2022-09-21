// TODO: orchestrator will eventually be moved to pkg/vending, since it provides
// the basic instrument to orchestrate the tool operations, it makes sense
// that users have access to it.
package internal

import (
	"fmt"
	"os"

	"github.com/alevinval/vendor-go/internal/cache"
	"github.com/alevinval/vendor-go/internal/installer"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

type CmdOrchestrator struct {
	preset       vending.Preset
	cacheManager *cache.Manager
}

func NewOrchestrator(preset vending.Preset) *CmdOrchestrator {
	return &CmdOrchestrator{preset, cache.NewManager(preset)}
}

func (co *CmdOrchestrator) Init() error {
	_, err := os.ReadFile(co.preset.GetSpecFilename())
	if err == nil {
		return fmt.Errorf("%q already exists? %w", co.preset.GetSpecFilename(), err)
	}

	spec := vending.NewSpec(co.preset)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("failed initializing: %w", err)
	}

	log.S().Infof("%s has been created", co.preset.GetSpecFilename())
	return nil
}

func (co *CmdOrchestrator) Install() error {
	lock, err := co.cacheManager.LockCache()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}
	defer lock.Release()

	spec := vending.NewSpec(co.preset)
	err = spec.Load()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	specLock := vending.NewSpecLock(co.preset)
	err = specLock.Load()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	cacheDir := co.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s", cacheDir)

	ins := installer.New(co.cacheManager, spec, specLock)
	err = ins.Install()
	if err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	err = specLock.Save()
	if err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	log.S().Infof("install success ✅")
	return nil
}

func (co *CmdOrchestrator) Update() error {
	lock, err := co.cacheManager.LockCache()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}
	defer lock.Release()

	spec := vending.NewSpec(co.preset)
	err = spec.Load()
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	specLock := vending.NewSpecLock(co.preset)
	err = specLock.Load()
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	cacheDir := co.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s", cacheDir)

	ins := installer.New(co.cacheManager, spec, specLock)
	err = ins.Update()
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	err = specLock.Save()
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	log.S().Infof("update success ✅")
	return nil
}

func (co *CmdOrchestrator) AddDependency(url, branch string) error {
	spec := vending.NewSpec(co.preset)
	err := spec.Load()
	if err != nil {
		return fmt.Errorf("cannot add dependency: %w", err)
	}

	dep := vending.NewDependency(url, branch)
	spec.AddDependency(dep)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("failed adding dependency: %w", err)
	}

	log.S().Infof("added dependency %s@%s",
		color.CyanString(url),
		color.YellowString(branch),
	)
	return nil
}

func (co *CmdOrchestrator) CleanCache() error {
	err := co.cacheManager.Reset()
	if err != nil {
		return fmt.Errorf("cannot clean cache: %w", err)
	}
	return nil
}
