package importer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/pkg/vending"
)

// Importer knows how to copy files from a source path to a destination path.
type Importer struct {
	repo *git.Repository
	spec *vending.Spec
	dep  *vending.Dependency
}

// New allocates a new Importer instance.
func New(repo *git.Repository, spec *vending.Spec, dep *vending.Dependency) *Importer {
	return &Importer{
		repo,
		spec,
		dep,
	}
}

// Import executes the import operation by copying files from the source to the
// destination.
func (imp *Importer) Import() error {
	collector, err := imp.collect()
	if err != nil {
		return fmt.Errorf("cannot import: %w", err)
	}
	err = collector.copyAll()
	if err != nil {
		return fmt.Errorf("cannot import: %w", err)
	}
	return nil
}

func (imp *Importer) collect() (*targetCollector, error) {
	selector := newSelector(imp.spec, imp.dep)
	targetCollector := &targetCollector{targets: []target{}}

	err := imp.repo.WalkDir(
		collectPathsFunc(
			imp.repo.Path(),
			imp.spec.VendorDir,
			selector,
			targetCollector,
		),
	)

	return targetCollector, err
}

func collectPathsFunc(srcRoot, dstRoot string, selector *Selector, collector *targetCollector) fs.WalkDirFunc {
	return func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("path collection interrupted: %w", err)
		}

		relativePath, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return fmt.Errorf("cannot get relative path: %w", err)
		}

		if isSelected, isTarget, isIgnored, _ := selector.Select(relativePath); isSelected {
			collector.add(
				target{
					srcRelative: relativePath,
					src:         path,
					dst:         filepath.Join(dstRoot, relativePath),
				},
			)
		} else if entry.IsDir() && (!isTarget || isIgnored) {
			return fs.SkipDir
		}

		return nil
	}
}
