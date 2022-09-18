package internal

import (
	"fmt"
	"os"
	"path"

	"github.com/alevinval/vendor-go/pkg/vendor"
)

const (
	LOCKS_DIR = "locks"
	REPOS_DIR = "repos"
)

func ensureCacheDir(preset vendor.Preset) error {
	cacheDir := preset.GetCacheDir()

	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return ensureCacheErr(err)
	}

	os.MkdirAll(path.Join(cacheDir, LOCKS_DIR), os.ModePerm)
	if err != nil {
		return ensureCacheErr(err)
	}

	os.MkdirAll(path.Join(cacheDir, REPOS_DIR), os.ModePerm)
	if err != nil {
		return ensureCacheErr(err)
	}

	return nil
}

func resetCacheDir(preset vendor.Preset) error {
	err := os.RemoveAll(preset.GetCacheDir())
	if err != nil {
		return fmt.Errorf("cannot remove cache: %w", err)
	}
	return ensureCacheDir(preset)
}

func getRepositoryPath(cacheDir string, dep *vendor.Dependency) string {
	return path.Join(cacheDir, REPOS_DIR, dep.ID())
}

func getRepositoryLockPath(cacheDir string, dep *vendor.Dependency) string {
	return path.Join(cacheDir, LOCKS_DIR, dep.ID())
}

func ensureCacheErr(err error) error {
	return fmt.Errorf("cannot bootstrap cache: %w", err)
}
