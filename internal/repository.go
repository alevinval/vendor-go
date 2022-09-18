package internal

import (
	"io/fs"
	"path/filepath"

	"github.com/alevinval/vendor-go/pkg/govendor"
)

type Repository struct {
	dep  *govendor.Dependency
	git  *Git
	path string
}

func NewRepository(cacheDir string, dep *govendor.Dependency) *Repository {
	return &Repository{
		path: filepath.Join(cacheDir, dep.ID()),
		dep:  dep,
		git:  &Git{},
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
