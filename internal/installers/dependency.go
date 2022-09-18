package installers

import (
	"fmt"
	"io/fs"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/internal/paths"
	"github.com/alevinval/vendor-go/pkg/govendor"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/fatih/color"
)

type dependencyInstaller struct {
	dep           *govendor.Dependency
	depLock       *govendor.DependencyLock
	repo          *internal.Repository
	importFilesFn fs.WalkDirFunc
}

func newDependencyInstaller(spec *govendor.Spec, dep *govendor.Dependency, depLock *govendor.DependencyLock, repo *internal.Repository) *dependencyInstaller {
	selector := paths.NewSelector(spec, dep)
	importFilesFn := paths.ImportFileFunc(selector, repo.Path(), spec.VendorDir)

	return &dependencyInstaller{
		dep:           dep,
		depLock:       depLock,
		repo:          repo,
		importFilesFn: importFilesFn,
	}
}

func (d *dependencyInstaller) Install() (*govendor.DependencyLock, error) {
	err := d.repo.Ensure()
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

func (d *dependencyInstaller) Update() (*govendor.DependencyLock, error) {
	err := d.repo.Ensure()
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

func (d *dependencyInstaller) importFiles() (*govendor.DependencyLock, error) {
	err := d.repo.WalkDir(d.importFilesFn)
	if err != nil {
		return nil, fmt.Errorf("cannot import files: %w", err)
	}

	commit, err := d.repo.GetCurrentCommit()
	if err != nil {
		return nil, fmt.Errorf("cannot get current commit: %w", err)
	}

	return govendor.NewDependencyLock(d.dep.URL, commit), nil
}
