package internal

import (
	"fmt"
	"os"

	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

type Installer struct {
	spec     *vending.Spec
	lock     *vending.SpecLock
	cacheDir string
}

func NewInstaller(cacheDir string, spec *vending.Spec, lock *vending.SpecLock) *Installer {
	return &Installer{
		cacheDir: cacheDir,
		spec:     spec,
		lock:     lock,
	}
}

func (in *Installer) Install() error {
	return in.run(installFunc)
}

func (in *Installer) Update() error {
	return in.run(updateFunc)
}

func (in *Installer) run(action actionFunc) error {
	resetVendorDir(in.spec.VendorDir)

	for _, dep := range in.spec.Deps {
		repo := NewRepository(in.cacheDir, dep)
		lock, _ := in.lock.Find(dep.URL)
		dependencyInstaller := newDependencyInstaller(in.spec, dep, lock, repo)

		dependencyLock, err := action(dependencyInstaller)
		if err != nil {
			return fmt.Errorf("cannot complete action: %w", err)
		}

		log.S().Infof("  🔒 %s", color.YellowString(dependencyLock.Commit))
		in.lock.AddDependencyLock(dependencyLock)
	}

	return nil
}

func resetVendorDir(vendorDir string) {
	os.RemoveAll(vendorDir)
	os.MkdirAll(vendorDir, os.ModePerm)
}

type actionFunc = func(*dependencyInstaller) (*vending.DependencyLock, error)

func installFunc(installer *dependencyInstaller) (*vending.DependencyLock, error) {
	return installer.Install()
}

func updateFunc(installer *dependencyInstaller) (*vending.DependencyLock, error) {
	return installer.Update()
}
