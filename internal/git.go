package internal

import (
	"github.com/alevinval/vendor-go/internal/log"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var logger = log.GetLogger()

type Git struct{}

func (g Git) GetCurrentCommit(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func (g Git) OpenOrClone(url, branch, path string) error {
	_, err := git.PlainOpen(path)
	if err == nil {
		return nil
	}

	return g.Clone(url, branch, path)
}

func (g Git) Clone(url, branch, path string) error {
	logger.Infof("cloning %s...", url)
	cloneOpts := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	}
	_, err := git.PlainClone(path, false, cloneOpts)
	return err
}

func (g Git) Fetch(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	fetchOpts := &git.FetchOptions{
		Force: true,
		Tags:  git.AllTags,
	}
	err = repo.Fetch(fetchOpts)
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (g Git) Reset(path, refname string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(refname))
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	cleanOpts := &git.CleanOptions{
		Dir: true,
	}
	err = wt.Clean(cleanOpts)
	if err != nil {
		return err
	}

	opts := &git.ResetOptions{
		Commit: *hash,
		Mode:   git.HardReset,
	}
	return wt.Reset(opts)
}
