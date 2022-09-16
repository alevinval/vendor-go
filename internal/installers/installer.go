package installers

import (
	"os"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/pkg/govendor"
	"github.com/fatih/color"
)

type Installer struct {
	spec      *govendor.Spec
	lock      *govendor.SpecLock
	cacheRoot string
}

func NewInstaller(cache string, spec *govendor.Spec, lock *govendor.SpecLock) *Installer {
	return &Installer{
		cacheRoot: cache,
		spec:      spec,
		lock:      lock,
	}
}

func (in *Installer) Install() error {
	return in.run(installFunc)
}

func (in *Installer) Update() error {
	return in.run(updateFunc)
}

func (in *Installer) run(action actionFunc) error {
	resetVendorPath(in.spec.VendorDir)

	for _, dep := range in.spec.Deps {
		repo := internal.NewRepository(in.cacheRoot, dep)
		lock, _ := in.lock.Find(dep.URL)
		depInstaller := newDependencyInstaller(in.spec, dep, lock, repo)

		newLock, err := action(depInstaller)
		if err != nil {
			return err
		}

		logger.Infof("  🔒 %s", color.YellowString(newLock.Commit))
		in.lock.Add(newLock)
	}

	return nil
}

func resetVendorPath(vendorPath string) {
	os.RemoveAll(vendorPath)
	os.MkdirAll(vendorPath, os.ModePerm)
}

type actionFunc = func(*dependencyInstaller) (*govendor.DependencyLock, error)

func installFunc(installer *dependencyInstaller) (*govendor.DependencyLock, error) {
	return installer.Install()
}

func updateFunc(installer *dependencyInstaller) (*govendor.DependencyLock, error) {
	return installer.Update()
}
