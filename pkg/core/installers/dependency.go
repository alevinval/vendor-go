package installers

import (
	"fmt"
	"io/fs"

	"github.com/alevinval/vendor-go/internal"
	"github.com/alevinval/vendor-go/internal/paths"
	"github.com/alevinval/vendor-go/pkg/core"
	"github.com/alevinval/vendor-go/pkg/core/log"
	"github.com/fatih/color"
)

var logger = log.GetLogger()

type dependencyInstaller struct {
	dep           *core.Dependency
	depLock       *core.DependencyLock
	repo          *internal.Repository
	importFilesFn fs.WalkDirFunc
}

func newDependencyInstaller(spec *core.Spec, dep *core.Dependency, depLock *core.DependencyLock, repo *internal.Repository) *dependencyInstaller {
	selector := paths.NewPathSelector(spec, dep)
	importFilesFn := paths.ImportFileFunc(selector, repo.Path(), spec.VendorDir)

	return &dependencyInstaller{
		dep:           dep,
		depLock:       depLock,
		repo:          repo,
		importFilesFn: importFilesFn,
	}
}

func (d *dependencyInstaller) Install() (*core.DependencyLock, error) {
	err := d.repo.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %s", err)
	}

	if d.depLock != nil {
		logger.Infof("installing %s", color.CyanString(d.dep.URL))
		d.repo.CheckoutCommit(d.depLock.Commit)
	} else {
		logger.Infof("installing %s@%s", color.CyanString(d.dep.URL), color.YellowString(d.dep.Branch))
		d.repo.CheckoutBranch(d.dep.Branch)
		d.repo.Pull()
	}

	return d.importFiles()
}

func (d *dependencyInstaller) Update() (*core.DependencyLock, error) {
	err := d.repo.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %s", err)
	}

	logger.Infof("updating %s@%s", color.CyanString("%s", d.dep.URL), color.YellowString(d.dep.Branch))

	d.repo.Fetch()
	d.repo.CheckoutBranch(d.dep.Branch)
	d.repo.Pull()

	return d.importFiles()
}

func (d *dependencyInstaller) importFiles() (*core.DependencyLock, error) {
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
