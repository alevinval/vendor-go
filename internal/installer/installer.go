package installer

import (
	"fmt"
	"os"
	"sync"

	"github.com/alevinval/vendor-go/internal/cache"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

type Installer struct {
	spec     *vending.Spec
	specLock *vending.SpecLock
	cache    *cache.Cache
}

func New(cache *cache.Cache, spec *vending.Spec, specLock *vending.SpecLock) *Installer {
	return &Installer{
		spec,
		specLock,
		cache,
	}
}

func (in *Installer) Install() error {
	return in.runInParallel(installFunc)
}

func (in *Installer) Update() error {
	return in.runInParallel(updateFunc)
}

func (in *Installer) runInParallel(action actionFunc) error {
	err := resetVendorDir(in.spec.VendorDir)
	if err != nil {
		return err
	}

	n := len(in.spec.Deps)
	out := make(chan *vending.DependencyLock, n)
	errors := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(n)

	for _, dep := range in.spec.Deps {
		go in.runInBackground(wg, action, dep, out, errors)
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

	for n > 0 {
		dependencyLock := <-out

		log.S().Infof("locking %s\n  🔒 %s",
			color.CyanString(dependencyLock.URL),
			color.YellowString(dependencyLock.Commit),
		)
		in.specLock.AddDependencyLock(dependencyLock)

		n--
	}

	return nil
}

func (in *Installer) runInBackground(
	wg *sync.WaitGroup,
	action actionFunc,
	dep *vending.Dependency,
	out chan *vending.DependencyLock,
	errors chan error,
) {
	defer wg.Done()

	repo, err := in.cache.GetRepository(dep)
	if err != nil {
		errors <- fmt.Errorf("cannot complete action: %w", err)
		return
	}

	lock, _ := in.specLock.FindByURL(dep.URL)
	dependencyInstaller := newDependencyInstaller(in.spec, dep, lock, repo)

	dependencyLock, err := action(dependencyInstaller)
	if err != nil {
		errors <- fmt.Errorf("cannot complete action: %w", err)
		return
	}

	out <- dependencyLock
	return
}

func resetVendorDir(vendorDir string) error {
	err := os.RemoveAll(vendorDir)
	if err != nil {
		return fmt.Errorf("cannot remove vendor dir: %w", err)
	}
	err = os.MkdirAll(vendorDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create vendor dir: %w", err)
	}
	return nil
}

type actionFunc = func(*dependencyInstaller) (*vending.DependencyLock, error)

func installFunc(installer *dependencyInstaller) (*vending.DependencyLock, error) {
	return installer.Install()
}

func updateFunc(installer *dependencyInstaller) (*vending.DependencyLock, error) {
	if installer.dep.Pinned {
		log.S().Infof("%s update for pinned dependency %s", color.RedString("skipping"), color.YellowString(installer.dep.URL))
		return installer.Install()
	} else {
		return installer.Update()
	}
}
