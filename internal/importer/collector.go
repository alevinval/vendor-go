package importer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alevinval/vendor-go/pkg/log"
)

type targetCollector struct {
	targets []target
}

type target struct {
	src         string
	dst         string
	srcRelative string
}

func (tc *targetCollector) add(t target) {
	tc.targets = append(tc.targets, t)
}

func (tc *targetCollector) copyAll() error {
	for _, target := range tc.targets {
		err := target.copy()
		if err != nil {
			return fmt.Errorf("cannot copy: %w", err)
		}
	}
	return nil
}

func (t *target) copy() error {
	log.S().Debugf("  [copy] ../%s -> %s", t.srcRelative, t.dst)

	dstDir := filepath.Dir(t.dst)
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create dstDir %q: %w", dstDir, err)
	}
	err = copyFile(t.src, t.dst)
	if err != nil {
		return fmt.Errorf("cannot copyFile: %w", err)
	}
	return nil
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
