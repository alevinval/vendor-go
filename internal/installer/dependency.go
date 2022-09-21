package installer

import (
	"fmt"
	"io/fs"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/internal/importer"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

type dependencyInstaller struct {
	dep             *vending.Dependency
	depLock         *vending.DependencyLock
	repo            *git.Repository
	importWalkDirFn fs.WalkDirFunc
}

func newDependencyInstaller(spec *vending.Spec, dep *vending.Dependency, depLock *vending.DependencyLock, repo *git.Repository) *dependencyInstaller {
	selector := importer.NewSelector(spec, dep)
	importWalkDirFn := importer.WalkDirFunc(selector, repo.Path(), spec.VendorDir)

	return &dependencyInstaller{
		dep:             dep,
		depLock:         depLock,
		repo:            repo,
		importWalkDirFn: importWalkDirFn,
	}
}

func (d *dependencyInstaller) Install() (*vending.DependencyLock, error) {
	lock, err := d.repo.Lock()
	if err != nil {
		return nil, fmt.Errorf("cannot acquire repository lock: %w", err)
	}
	defer lock.Release()

	err = d.repo.OpenOrClone()
	if err != nil {
		return nil, fmt.Errorf("cannot ensure repository: %w", err)
	}

	if d.depLock == nil {
		log.S().Infof("installing %s@%s",
			color.CyanString(d.dep.URL),
			color.YellowString(d.dep.Branch),
		)
		err = d.repo.Reset(d.dep.Branch)
	} else {
		log.S().Infof("installing %s@%s",
			color.CyanString(d.dep.URL),
			color.YellowString(
				fmt.Sprintf("%.8s", d.depLock.Commit),
			),
		)
		err = d.repo.Reset(d.depLock.Commit)
	}

	if err != nil {
		return nil, fmt.Errorf("reset failed: %w", err)
	}
	return d.importFiles()
}

func (d *dependencyInstaller) Update() (*vending.DependencyLock, error) {
	lock, err := d.repo.Lock()
	if err != nil {
		return nil, fmt.Errorf("cannot acquire repository lock: %w", err)
	}
	defer lock.Release()

	err = d.repo.OpenOrClone()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %s", err)
	}

	log.S().Infof("updating %s@%s",
		color.CyanString("%s", d.dep.URL),
		color.YellowString(d.dep.Branch),
	)

	err = d.repo.Fetch()
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	err = d.repo.Reset(d.dep.Branch)
	if err != nil {
		return nil, fmt.Errorf("reset failed: %w", err)
	}

	return d.importFiles()
}

func (d *dependencyInstaller) importFiles() (*vending.DependencyLock, error) {
	err := d.repo.WalkDir(d.importWalkDirFn)
	if err != nil {
		return nil, fmt.Errorf("cannot import files: %w", err)
	}

	commit, err := d.repo.GetCurrentCommit()
	if err != nil {
		return nil, fmt.Errorf("cannot get current commit: %w", err)
	}

	return vending.NewDependencyLock(d.dep.URL, commit), nil
}
