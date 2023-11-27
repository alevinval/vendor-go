package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/alevinval/vendor-go/internal/git"
	"github.com/alevinval/vendor-go/internal/lock"
	"github.com/alevinval/vendor-go/pkg/vending"
	"github.com/fatih/color"
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
		return fmt.Errorf("cannot create locks directory: %w", err)
	}

	reposDir := path.Join(cacheDir, REPOS_DIR)
	err = man.fs.MkdirAll(reposDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create repositories directory: %w", err)
	}

	return nil
}

func (man *Manager) Clean() error {
	err := man.fs.RemoveAll(man.preset.GetCacheDir())
	if err != nil {
		return fmt.Errorf("cannot remove cache directory: %w", err)
	}
	return man.Ensure()
}

func (man *Manager) LockCache() (*lock.Lock, error) {
	err := man.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot ensure paths: %w", err)
	}

	lockPath := man.getCacheLockPath()
	lock := lock.New(lockPath).WithWarn(
		color.RedString(
			"cannot acquire cache lock, are you running multiple instances in parallel?",
		),
	)
	err = lock.Acquire()
	if err != nil {
		return nil, fmt.Errorf("cannot acquire lock %q: %w", lockPath, err)
	}
	return lock, nil
}

func (man *Manager) GetRepository(dep *vending.Dependency) (*git.Repository, error) {
	lock, err := man.getRepositoryLock(dep)
	if err != nil {
		return nil, fmt.Errorf("cannot get repository lock: %w", err)
	}
	return git.NewRepository(
		man.getRepositoryPath(dep),
		lock,
		dep,
	), nil
}

func (man *Manager) getRepositoryLock(dep *vending.Dependency) (*lock.Lock, error) {
	err := man.Ensure()
	if err != nil {
		return nil, fmt.Errorf("cannot ensure paths: %w", err)
	}

	lockPath := man.getRepositoryLockPath(dep)
	lock := lock.New(lockPath)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire lock %q: %w", lockPath, err)
	}

	return lock, nil
}

func (man *Manager) getCacheLockPath() string {
	return path.Join(
		man.preset.GetCacheDir(),
		"LOCK",
	)
}

func (man *Manager) getRepositoryLockPath(dep *vending.Dependency) string {
	return path.Join(
		man.preset.GetCacheDir(),
		LOCKS_DIR,
		getDependencyID(dep),
	)
}

func (man *Manager) getRepositoryPath(dep *vending.Dependency) string {
	return path.Join(
		man.preset.GetCacheDir(),
		REPOS_DIR,
		getDependencyID(dep),
	)
}

func getDependencyID(dep *vending.Dependency) string {
	sha := sha1.New()
	sha.Write([]byte(dep.URL))
	data := sha.Sum(nil)
	return hex.EncodeToString(data)
}

type cacheFs struct{}

func (cacheFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (cacheFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
