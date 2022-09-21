package git

import (
	"io/fs"
	"path/filepath"

	"github.com/alevinval/vendor-go/internal/lock"
	"github.com/alevinval/vendor-go/pkg/vending"
)

type Repository struct {
	dep  *vending.Dependency
	git  *Git
	lock *lock.Lock
	path string
}

func NewRepository(path string, lock *lock.Lock, dep *vending.Dependency) *Repository {
	return &Repository{
		dep:  dep,
		git:  &Git{},
		lock: lock,
		path: path,
	}
}

func (r *Repository) Path() string {
	return r.path
}

func (r *Repository) OpenOrClone() error {
	return r.git.OpenOrClone(r.dep.URL, r.dep.Branch, r.Path())
}

func (r *Repository) Fetch() error {
	return r.git.Fetch(r.Path())
}

func (r *Repository) Reset(refname string) error {
	return r.git.Reset(r.Path(), refname)
}

func (r *Repository) GetCurrentCommit() (string, error) {
	return r.git.GetCurrentCommit(r.Path())
}

func (r *Repository) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(r.Path(), fn)
}

func (r *Repository) Lock() (*lock.Lock, error) {
	return r.lock, r.lock.Acquire()
}
