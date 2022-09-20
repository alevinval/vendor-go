package internal

import (
	"io/fs"
	"path/filepath"

	"github.com/alevinval/vendor-go/internal/cache"
	"github.com/alevinval/vendor-go/internal/lock"
	"github.com/alevinval/vendor-go/pkg/vending"
)

type Repository struct {
	dep  *vending.Dependency
	git  *Git
	lock *lock.Lock
	path string
}

func NewRepository(manager *cache.Manager, dep *vending.Dependency) *Repository {
	return &Repository{
		dep:  dep,
		git:  &Git{},
		path: manager.GetRepositoryPath(dep),
		lock: lock.New(manager.GetRepositoryLockPath(dep)),
	}
}

func (r *Repository) Path() string {
	return r.path
}

func (r *Repository) Ensure() error {
	return r.git.OpenOrClone(r.dep.URL, r.dep.Branch, r.path)
}

func (r *Repository) Fetch() error {
	return r.git.Fetch(r.path)
}

func (r *Repository) Reset(refname string) error {
	return r.git.Reset(r.path, refname)
}

func (r *Repository) GetCurrentCommit() (string, error) {
	return r.git.GetCurrentCommit(r.path)
}

func (r *Repository) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(r.path, fn)
}

func (r *Repository) Acquire() error {
	return r.lock.Acquire()
}

func (r *Repository) Release() error {
	return r.lock.Release()
}
