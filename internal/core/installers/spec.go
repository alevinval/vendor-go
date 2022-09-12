package installers

import (
	"os"

	"github.com/alevinval/vendor-go/internal/core"
)

type SpecInstaller struct {
	spec      *core.Spec
	lock      *core.SpecLock
	cacheRoot string
}

var InstallFn = func(installer *DependencyInstaller) (*core.DependencyLock, error) {
	return installer.Install()
}

var UpdateFn = func(installer *DependencyInstaller) (*core.DependencyLock, error) {
	return installer.Update()
}

func NewSpecInstaller(cache string, spec *core.Spec, lock *core.SpecLock) *SpecInstaller {
	return &SpecInstaller{
		cacheRoot: cache,
		spec:      spec,
		lock:      lock,
	}
}

func (s *SpecInstaller) Install() error {
	s.spec.LoadPreset()
	return s.run(InstallFn)
}

func (s *SpecInstaller) Update() error {
	s.spec.LoadPreset()
	return s.run(UpdateFn)
}

func (s *SpecInstaller) run(action func(*DependencyInstaller) (*core.DependencyLock, error)) error {
	resetVendorPath(s.spec.VendorDir)

	for _, dep := range s.spec.Deps {
		repo := core.NewRepository(s.cacheRoot, dep)
		lock, _ := s.lock.Find(dep.URL)
		installer := NewDependencyInstaller(s.spec, dep, lock, repo)

		newLock, err := action(installer)
		if err != nil {
			return err
		}

		s.lock.Add(newLock)
	}

	return nil
}

func resetVendorPath(vendorPath string) {
	os.RemoveAll(vendorPath)
	os.MkdirAll(vendorPath, os.ModePerm)
}
