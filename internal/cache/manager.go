package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/alevinval/vendor-go/pkg/vending"
)

const (
	LOCKS_DIR = "locks"
	REPOS_DIR = "repos"
)

var _ CacheFS = (*cacheFs)(nil)

type CacheFS interface {
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
}

type Manager struct {
	fs     CacheFS
	preset vending.Preset
}

func NewManager(preset vending.Preset) *Manager {
	return &Manager{
		fs:     &cacheFs{},
		preset: preset,
	}
}

func (man *Manager) Ensure() error {
	cacheDir := man.preset.GetCacheDir()

	locksDir := path.Join(cacheDir, LOCKS_DIR)
	err := man.fs.MkdirAll(locksDir, os.ModePerm)
	if err != nil {
		return ensureCacheErr(err)
	}

	reposDir := path.Join(cacheDir, REPOS_DIR)
	err = man.fs.MkdirAll(reposDir, os.ModePerm)
	if err != nil {
		return ensureCacheErr(err)
	}

	return nil
}

func (man *Manager) Reset() error {
	err := man.fs.RemoveAll(man.preset.GetCacheDir())
	if err != nil {
		return fmt.Errorf("cannot remove cache: %w", err)
	}
	return man.Ensure()
}

func (man *Manager) GetRepositoryPath(dep *vending.Dependency) string {
	return path.Join(man.preset.GetCacheDir(), REPOS_DIR, getDependencyID(dep))
}

func (man *Manager) GetRepositoryLockPath(dep *vending.Dependency) string {
	return path.Join(man.preset.GetCacheDir(), LOCKS_DIR, getDependencyID(dep))
}

func ensureCacheErr(err error) error {
	return fmt.Errorf("cannot bootstrap cache: %w", err)
}

func getDependencyID(dep *vending.Dependency) string {
	sha := sha256.New()
	data := sha.Sum([]byte(dep.URL))
	return hex.EncodeToString(data)
}

type cacheFs struct{}

func (cacheFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (cacheFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
