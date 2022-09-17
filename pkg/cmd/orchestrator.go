package cmd

import (
	"fmt"
	"os"

	"github.com/alevinval/vendor-go/internal/installers"
	"github.com/alevinval/vendor-go/pkg/govendor"
	"github.com/fatih/color"
)

type CmdOrchestrator struct {
	preset govendor.Preset
}

func NewOrchestrator(preset govendor.Preset) *CmdOrchestrator {
	return &CmdOrchestrator{preset}
}

func (co *CmdOrchestrator) Init() error {
	_, err := os.ReadFile(co.preset.GetSpecFilename())
	if err == nil {
		return fmt.Errorf("%q already exists? %w", co.preset.GetSpecFilename(), err)
	}

	spec := govendor.NewSpec(co.preset)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("failed initializing: %w", err)
	}

	logger.Infof("%s has been created", co.preset.GetSpecFilename())
	return nil
}

func (co *CmdOrchestrator) Install() error {
	spec, err := govendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	specLock, err := govendor.LoadSpecLock(co.preset)
	if err != nil {
		return fmt.Errorf("cannot install: %w", err)
	}

	cache := co.preset.GetCachePath()
	logger.Infof("repository cache located at %s", cache)

	m := installers.NewInstaller(co.preset, spec, specLock)
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

	logger.Infof("install success ✅")
	return nil
}

func (co *CmdOrchestrator) Update() error {
	spec, err := govendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	specLock, err := govendor.LoadSpecLock(co.preset)
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	cache := co.preset.GetCachePath()
	logger.Infof("repository cache located at %s", cache)

	m := installers.NewInstaller(co.preset, spec, specLock)
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

	logger.Infof("update success ✅")
	return nil
}

func (co *CmdOrchestrator) AddDependency(url, branch string) error {
	spec, err := govendor.LoadSpec(co.preset)
	if err != nil {
		return fmt.Errorf("cannot add dependency: %w", err)
	}

	dep := govendor.NewDependency(url, branch)
	spec.Add(dep)

	err = spec.Save()
	if err != nil {
		return fmt.Errorf("failed adding dependency: %w", err)
	}

	logger.Infof("added dependency %s@%s",
		color.CyanString(url),
		color.YellowString(branch),
	)
	return nil
}
