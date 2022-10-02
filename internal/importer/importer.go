package importer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/alevinval/vendor-go/pkg/vending"
)

type Importer struct {
	repo *git.Repository
	spec *vending.Spec
	dep  *vending.Dependency
}

func New(repo *git.Repository, spec *vending.Spec, dep *vending.Dependency) *Importer {
	return &Importer{
		repo,
		spec,
		dep,
	}
}

func (imp *Importer) Import() error {
	selector := newSelector(imp.spec, imp.dep)
	importWalkDirFn := walkDirFunc(selector, imp.repo.Path(), imp.spec.VendorDir)
	return imp.repo.WalkDir(importWalkDirFn)
}

func walkDirFunc(selector *Selector, srcRoot, dstRoot string) fs.WalkDirFunc {
	return func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("import interrupted: %w", err)
		}

		relativePath := strings.TrimPrefix(path, srcRoot)
		if !selector.Select(relativePath) {
			return nil
		}

		dst := filepath.Join(dstRoot, relativePath)
		log.S().Debugf("  ..%s -> %s", relativePath, dst)

		dstDir := filepath.Dir(dst)
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create target path %q: %w", dstDir, err)
		}
		return copyFile(path, dst)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open %q: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("cannot create %q: %w", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("cannot copy %q => %q: %w", src, dst, err)
	}

	err = out.Close()
	if err != nil {
		return fmt.Errorf("cannot close %q: %w", dst, err)
	}

	return nil
}
