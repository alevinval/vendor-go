package core

import (
	"io/fs"
	"path/filepath"
)

type Repository struct {
	path string
	dep  *Dependency
	Git  *Git
}

func NewRepository(cache string, dep *Dependency) *Repository {
	return &Repository{
		path: filepath.Join(cache, dep.ID()),
		dep:  dep,
		Git:  &Git{},
	}
}

func (r *Repository) Path() string {
	return r.path
}

func (r *Repository) CheckoutCommit(commit string) error {
	return r.Git.CheckoutCommit(commit, r.path)
}

func (r *Repository) CheckoutBranch(branch string) error {
	return r.Git.CheckoutBranch(branch, r.path)
}

func (r *Repository) Pull() error {
	return r.Git.Pull(r.path)
}

func (r *Repository) Fetch() error {
	return r.Git.Fetch(r.path)
}

func (r *Repository) GetCurrentCommit() (string, error) {
	return r.Git.GetCurrentCommit(r.path)
}

func (r *Repository) Ensure() error {
	return r.Git.OpenOrClone(r.dep.URL, r.dep.Branch, r.path)
}

func (r *Repository) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(r.path, fn)
}
