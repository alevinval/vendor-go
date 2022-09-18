package internal

import (
	"fmt"

	"github.com/alevinval/vendor-go/pkg/log"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Git struct{}

func (g Git) GetCurrentCommit(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", gitOpenErr(err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("cannot get HEAD reference: %w", err)
	}

	return head.Hash().String(), nil
}

func (g Git) OpenOrClone(url, branch, path string) error {
	_, err := git.PlainOpen(path)
	if err == nil {
		return nil
	}

	err = g.Clone(url, branch, path)
	if err != nil {
		return fmt.Errorf("cannot clone: %w", err)
	}

	return nil
}

func (g Git) Clone(url, branch, path string) error {
	log.S().Infof("cloning %s...", url)
	cloneOpts := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	}
	_, err := git.PlainClone(path, false, cloneOpts)
	if err != nil {
		return gitCloneErr(err)
	}

	return nil
}

func (g Git) Fetch(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return gitOpenErr(err)
	}

	fetchOpts := &git.FetchOptions{
		Force: true,
		Tags:  git.AllTags,
	}
	err = repo.Fetch(fetchOpts)
	switch err {
	case git.NoErrAlreadyUpToDate:
		return nil
	case nil:
		return nil
	}

	return fmt.Errorf("cannot fetch: %w", err)
}

func (g Git) Reset(path, refname string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return gitOpenErr(err)
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(refname))
	if err != nil {
		return fmt.Errorf("cannot resolve %q: %w", refname, err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("cannot get worktree: %w", err)
	}

	cleanOpts := &git.CleanOptions{
		Dir: true,
	}
	err = wt.Clean(cleanOpts)
	if err != nil {
		return fmt.Errorf("cannot clean: %w", err)
	}

	opts := &git.ResetOptions{
		Commit: *hash,
		Mode:   git.HardReset,
	}
	err = wt.Reset(opts)
	if err != nil {
		return fmt.Errorf("cannot reset: %w", err)
	}

	return nil
}

func gitOpenErr(err error) error {
	return fmt.Errorf("cannot open repository: %w", err)
}

func gitCloneErr(err error) error {
	return fmt.Errorf("cannot clone: %w", err)
}
