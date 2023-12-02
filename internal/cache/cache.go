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
	locksDir = "locks"
	reposDir = "repos"
)

type (
	// Cache is used to mediate with the caching layer where all repositories
	// for the various dependencies are kept. It relies on a locking mechanism
	// to ensure safe access to the cache when things are executed concurrently
	// or different instances of the tool run at the same time.
	Cache struct {
		fs   fsOps
		path string
	}
)

// New cache manager pointing towards the path where it resides.
func New(path string) *Cache {
	return &Cache{
		fs:   &defaultFs{},
		path: path,
	}
}

// Init ensures proper directory structure of the cache exists.
func (c *Cache) Init() error {
	locksDir := path.Join(c.path, locksDir)
	if err := c.fs.MkdirAll(locksDir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create locks directory: %w", err)
	}

	reposDir := path.Join(c.path, reposDir)
	if err := c.fs.MkdirAll(reposDir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create repositories directory: %w", err)
	}

	return nil
}

// Reset removes all contents in the cache root directory and immediately
// initialises the cache again.
func (c *Cache) Reset() error {
	if err := c.fs.RemoveAll(c.path); err != nil {
		return fmt.Errorf("cannot remove cache directory: %w", err)
	}
	return c.Init()
}

// Lock acquires the global cache lock. This can be used to ensure only one agent
// is manipulating the cache and its contents.
func (c *Cache) Lock() (*lock.Lock, error) {
	if err := c.Init(); err != nil {
		return nil, fmt.Errorf("cannot initialize paths: %w", err)
	}

	lock := lock.New(c.lockPath()).WithWarn(
		color.RedString(
			"cannot acquire cache lock, are you running multiple instances in parallel?",
		),
	)
	if err := lock.Acquire(); err != nil {
		return nil, fmt.Errorf("cannot acquire lock %q: %w", c.lockPath(), err)
	}
	return lock, nil
}

// GetRepository returns git.Repository from the cache.
func (c *Cache) GetRepository(dep *vending.Dependency) (*git.Repository, error) {
	lock, err := c.repositoryLock(dep)
	if err != nil {
		return nil, fmt.Errorf("cannot get repository lock: %w", err)
	}
	return git.NewRepository(
		c.getRepositoryPath(dep),
		lock,
		dep,
	), nil
}

// repositoryLock returns lock.Lock for a given dependency.
func (c *Cache) repositoryLock(dep *vending.Dependency) (*lock.Lock, error) {
	if err := c.Init(); err != nil {
		return nil, fmt.Errorf("cannot initialize paths: %w", err)
	}
	repoPath := c.getRepositoryLockPath(dep)
	return lock.New(repoPath), nil
}

// lockPath returns path of the global cache lock.
func (c *Cache) lockPath() string {
	return path.Join(
		c.path,
		"LOCK",
	)
}

// getRepositoryLockPath returns path of the repository lock.
func (c *Cache) getRepositoryLockPath(dep *vending.Dependency) string {
	return path.Join(
		c.path,
		locksDir,
		getSha1(dep),
	)
}

// getRepositoryPath returns path of the git repository.
func (c *Cache) getRepositoryPath(dep *vending.Dependency) string {
	return path.Join(
		c.path,
		reposDir,
		getSha1(dep),
	)
}

func getSha1(dep *vending.Dependency) string {
	sha := sha1.New()
	sha.Write([]byte(dep.URL))
	data := sha.Sum(nil)
	return hex.EncodeToString(data)
}
