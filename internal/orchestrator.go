package internal

import (
	"fmt"
	"os"
	"path"

	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vendor"
	"github.com/fatih/color"
)

type CmdOrchestrator struct {
	preset vendor.Preset
}

func NewOrchestrator(preset vendor.Preset) *CmdOrchestrator {
	return &CmdOrchestrator{preset}
}

func (co *CmdOrchestrator) Init() error {
	_, err := os.ReadFile(co.preset.GetSpecFilename())
	if err == nil {
		return fmt.Errorf("%q already exists? %w", co.preset.GetSpecFilename(), err)
	}

	spec := vendor.NewSpec(co.preset)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("failed initializing: %w", err)
	}

	log.S().Infof("%s has been created", co.preset.GetSpecFilename())
	return nil
}

func (co *CmdOrchestrator) Install() error {
	lock, err := co.getCacheLock()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}
	defer lock.Release()

	spec, err := vendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	specLock, err := vendor.LoadSpecLock(co.preset)
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	cacheDir := co.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s", cacheDir)

	m := NewInstaller(cacheDir, spec, specLock)
	err = m.Install()
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
	lock, err := co.getCacheLock()
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}
	defer lock.Release()

	spec, err := vendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	specLock, err := vendor.LoadSpecLock(co.preset)
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	cacheDir := co.preset.GetCacheDir()
	log.S().Infof("repository cache located at %s", cacheDir)

	m := NewInstaller(cacheDir, spec, specLock)
	err = m.Update()
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
	spec, err := vendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot add dependency: %w", err)
	}

	dep := vendor.NewDependency(url, branch)
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
	err := resetCacheDir(co.preset)
	if err != nil {
		return fmt.Errorf("cannot clean cache: %w", err)
	}
	return nil
}

func (co *CmdOrchestrator) getCacheLock() (*Lock, error) {
	err := ensureCacheDir(co.preset)
	if err != nil {
		return nil, fmt.Errorf("cannot create cache: %w", err)
	}

	lockPath := path.Join(co.preset.GetCacheDir(), "LOCK")
	lock := NewLock(lockPath).
		WithWarn("cannot acquire cache lock, are you running multiple instances in parallel?")
	err = lock.Acquire()
	if err != nil {
		return nil, fmt.Errorf("cannot acquire cache lock %q: %w", lockPath, err)
	}
	return lock, nil
}