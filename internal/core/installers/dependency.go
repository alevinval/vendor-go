package installers

import (
	"fmt"
	"io/fs"
	"vendor-go/internal/core"
	"vendor-go/internal/core/log"
	"vendor-go/internal/core/paths"
)

var logger = log.GetLogger()

type DependencyInstaller struct {
	dep           *core.Dependency
	depLock       *core.DependencyLock
	repo          *core.Repository
	importFilesFn fs.WalkDirFunc
}

func NewDependencyInstaller(spec *core.Spec, dep *core.Dependency, depLock *core.DependencyLock, repo *core.Repository) *DependencyInstaller {
	selector := paths.NewPathSelector(spec, dep)
	importFilesFn := paths.ImportFileFunc(selector, repo.Path(), spec.VendorDir)

	return &DependencyInstaller{
		dep:           dep,
		depLock:       depLock,
		repo:          repo,
		importFilesFn: importFilesFn,
	}
}

func (d *DependencyInstaller) Install() (*core.DependencyLock, error) {
	err := d.repo.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %s", err)
	}

	if d.depLock != nil {
		logger.Infof("installing %s@%s", d.dep.URL, d.depLock.Commit)
		d.repo.CheckoutCommit(d.depLock.Commit)
	} else {
		logger.Infof("installing %s@%s", d.dep.URL, d.dep.Branch)
		d.repo.CheckoutBranch(d.dep.Branch)
		d.repo.Pull()
	}

	return d.importFiles()
}

func (d *DependencyInstaller) Update() (*core.DependencyLock, error) {
	err := d.repo.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %s", err)
	}

	logger.Infof("updating %s@%s", d.dep.URL, d.dep.Branch)

	d.repo.Fetch()
	d.repo.CheckoutBranch(d.dep.Branch)
	d.repo.Pull()

	return d.importFiles()
}

func (d *DependencyInstaller) importFiles() (*core.DependencyLock, error) {
	err := d.repo.WalkDir(d.importFilesFn)
	if err != nil {
		return nil, fmt.Errorf("cannot import files: %s", err)
	}

	commit, err := d.repo.GetCurrentCommit()
	if err != nil {
		return nil, fmt.Errorf("cannot get current commit: %s", err)
	}

	return core.NewDependencyLock(d.dep.URL, commit), nil
}
