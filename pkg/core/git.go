package core

import (
	"github.com/alevinval/vendor-go/pkg/core/log"

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

func (g Git) CheckoutCommit(commit, path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	checkoutOpts := &git.CheckoutOptions{
		Hash:  plumbing.NewHash(commit),
		Force: true,
	}
	return wt.Checkout(checkoutOpts)
}

func (g Git) CheckoutBranch(branch, path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	checkoutOpts := &git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName("origin", branch),
		Force:  true,
	}
	return wt.Checkout(checkoutOpts)
}

func (g Git) Pull(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	pullOpts := &git.PullOptions{Force: true}
	return wt.Pull(pullOpts)
}

func (g Git) Fetch(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	fetchOpts := &git.FetchOptions{Force: true}
	return repo.Fetch(fetchOpts)
}
