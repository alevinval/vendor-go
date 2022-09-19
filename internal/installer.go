package internal

import (
	"fmt"
	"os"
	"sync"

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
	return in.runInParallel(installFunc)
}

func (in *Installer) Update() error {
	return in.runInParallel(updateFunc)
}

func (in *Installer) runInParallel(action actionFunc) error {
	resetVendorDir(in.spec.VendorDir)

	N := len(in.spec.Deps)
	out := make(chan *vending.DependencyLock, N)
	errors := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	for _, dep := range in.spec.Deps {
		go in.runInBackground(wg, action, dep, out)
	}

	completed := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		close(completed)
	}()

	var isCompleted bool
	for !isCompleted {
		select {
		case err := <-errors:
			return err
		case <-completed:
			isCompleted = true
		}
	}

	for N > 0 {
		dependencyLock := <-out

		log.S().Infof("locking %s\n  🔒 %s",
			color.CyanString(dependencyLock.URL),
			color.YellowString(dependencyLock.Commit),
		)
		in.lock.AddDependencyLock(dependencyLock)

		N--
	}

	return nil
}

func (in *Installer) runInBackground(wg *sync.WaitGroup, action actionFunc, dep *vending.Dependency, out chan *vending.DependencyLock) error {
	repo := NewRepository(in.cacheDir, dep)
	lock, _ := in.lock.FindByURL(dep.URL)
	dependencyInstaller := newDependencyInstaller(in.spec, dep, lock, repo)

	dependencyLock, err := action(dependencyInstaller)
	if err != nil {
		return fmt.Errorf("cannot complete action: %w", err)
	}

	out <- dependencyLock
	wg.Done()
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
