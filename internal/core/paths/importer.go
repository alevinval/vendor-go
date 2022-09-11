package paths

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"vendor-go/internal/core/log"
)

var logger = log.GetLogger()

func ImportFileFunc(selector *PathSelector, srcRoot, dstRoot string) fs.WalkDirFunc {
	return func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			logger.Errorf("import interrupted: %s", err)
			return err
		}

		relativePath := strings.TrimPrefix(path, srcRoot)
		if !selector.Select(relativePath) {
			return nil
		}

		dst := filepath.Join(dstRoot, relativePath)
		logger.Debugf("  ..%s -> %s", relativePath, dst)

		dstDir := filepath.Dir(dst)
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create target path %s: %s", dstDir, err)
		}
		return copyFile(path, dst)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
