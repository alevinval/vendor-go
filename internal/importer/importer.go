package importer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/pkg/log"
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
		return fmt.Errorf("cannot collect: %w", err)
	}
	err = collector.copyAll()
	if err != nil {
		return fmt.Errorf("cannot copyAll: %w", err)
	}
	return nil
}

func (imp *Importer) collect() (*targetCollector, error) {
	selector := newSelector(imp.spec, imp.dep)
	targetCollector := &targetCollector{targets: []target{}}

	err := imp.repo.WalkDir(
		collectTargetsFunc(
			imp.repo.Path(),
			imp.spec.VendorDir,
			selector,
			targetCollector,
		),
	)

	return targetCollector, err
}

func collectTargetsFunc(
	srcRoot, dstRoot string,
	selector *Selector,
	collector *targetCollector,
) fs.WalkDirFunc {
	return func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			log.S().Warnf("skipping directory %s due to path error: %w", path, err)
			return nil
		}

		if strings.EqualFold(srcRoot, path) {
			return nil
		}

		pathRel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return fmt.Errorf("cannot get relative path: %w", err)
		}

		if entry.IsDir() && !selector.SelectDir(pathRel) {
			log.S().Debugf("  [skip] %s", pathRel)
			return fs.SkipDir
		} else if selector.SelectPath(pathRel) {
			collector.add(
				target{
					src:    path,
					srcRel: pathRel,
					dst:    filepath.Join(dstRoot, pathRel),
				},
			)
		}

		return nil
	}
}
