package internal

import (
	"io/fs"
	"path/filepath"

	"github.com/alevinval/vendor-go/pkg/core"
)

type Repository struct {
	dep  *core.Dependency
	git  *Git
	path string
}

func NewRepository(cache string, dep *core.Dependency) *Repository {
	return &Repository{
		path: filepath.Join(cache, dep.ID()),
		dep:  dep,
		git:  &Git{},
	}
}

func (r *Repository) Path() string {
	return r.path
}

func (r *Repository) CheckoutCommit(commit string) error {
	return r.git.CheckoutCommit(commit, r.path)
}

func (r *Repository) CheckoutBranch(branch string) error {
	return r.git.CheckoutBranch(branch, r.path)
}

func (r *Repository) Pull() error {
	return r.git.Pull(r.path)
}

func (r *Repository) Fetch() error {
	return r.git.Fetch(r.path)
}

func (r *Repository) GetCurrentCommit() (string, error) {
	return r.git.GetCurrentCommit(r.path)
}

func (r *Repository) Ensure() error {
	return r.git.OpenOrClone(r.dep.URL, r.dep.Branch, r.path)
}

func (r *Repository) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(r.path, fn)
}
