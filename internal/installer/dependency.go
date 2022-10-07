package installer

import (
	"fmt"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/internal/importer"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
)

type dependencyInstaller struct {
	dep     *vending.Dependency
	depLock *vending.DependencyLock
	repo    *git.Repository
	imp     *importer.Importer
}

func newDependencyInstaller(spec *vending.Spec, dep *vending.Dependency, depLock *vending.DependencyLock, repo *git.Repository) *dependencyInstaller {
	imp := importer.New(repo, spec, dep)

	return &dependencyInstaller{
		dep:     dep,
		depLock: depLock,
		repo:    repo,
		imp:     imp,
	}
}

func (d *dependencyInstaller) Install() (*vending.DependencyLock, error) {
	lock, err := d.repo.Lock()
	if err != nil {
		return nil, fmt.Errorf("cannot lock repository: %w", err)
	}
	defer lock.Release()

	err = d.repo.OpenOrClone()
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %w", err)
	}

	var refname, refnameLog string
	if d.depLock == nil {
		refname = d.dep.Branch
		refnameLog = refname
	} else {
		refname = d.depLock.Commit
		refnameLog = fmt.Sprintf("%.8s", d.depLock.Commit)
	}

	doReset := func(fetch bool) error {
		if fetch {
			err = d.repo.Fetch()
			if err != nil {
				return fmt.Errorf("cannot fetch repository: %w", err)
			}
		} else {
			log.S().Infof("installing %s@%s",
				color.CyanString(d.dep.URL),
				color.YellowString(
					refnameLog,
				),
			)
		}
		return d.repo.Reset(refname)
	}

	if err = doReset(false); err != nil {
		if err = doReset(true); err != nil {
			return nil, fmt.Errorf("cannot reset repository: %w", err)
		}
	}
	return d.importFiles()
}

func (d *dependencyInstaller) Update() (*vending.DependencyLock, error) {
	lock, err := d.repo.Lock()
	if err != nil {
		return nil, fmt.Errorf("cannot lock repository: %w", err)
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
		return nil, fmt.Errorf("cannot fetch repository: %w", err)
	}

	err = d.repo.Reset(d.dep.Branch)
	if err != nil {
		return nil, fmt.Errorf("cannot reset repository: %w", err)
	}

	return d.importFiles()
}

func (d *dependencyInstaller) importFiles() (*vending.DependencyLock, error) {
	err := d.imp.Import()
	if err != nil {
		return nil, fmt.Errorf("cannot import: %w", err)
	}

	commit, err := d.repo.GetCurrentCommit()
	if err != nil {
		return nil, fmt.Errorf("cannot get current commit: %w", err)
	}

	return vending.NewDependencyLock(d.dep.URL, commit), nil
}
