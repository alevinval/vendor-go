package git

import (
	"fmt"

	"github.com/alevinval/vendor-go/pkg/log"
	"github.com/fatih/color"
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

	return g.Clone(url, branch, path)
}

func (g Git) Clone(url, branch, path string) error {
	log.S().Infof(
		"cloning %s...",
		color.CyanString(url),
	)
	cloneOpts := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	}
	_, err := git.PlainClone(path, false, cloneOpts)
	if err != nil {
		return fmt.Errorf("cannot clone: %w", err)
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

	return fmt.Errorf("cannot git fetch: %w", err)
}

func (g Git) Reset(path, refname string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return gitOpenErr(err)
	}

	revisions := []plumbing.Revision{
		plumbing.Revision(
			plumbing.NewRemoteReferenceName("origin", refname),
		),
		plumbing.Revision(refname),
	}

	var hash *plumbing.Hash
	for _, rev := range revisions {
		hash, err = repo.ResolveRevision(rev)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("cannot resolve revision %q: %w", refname, err)
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
		return fmt.Errorf("cannot git clean: %w", err)
	}

	opts := &git.CheckoutOptions{
		Hash:  *hash,
		Force: true,
	}
	err = wt.Checkout(opts)
	if err != nil {
		return fmt.Errorf("cannot git checkout: %w", err)
	}

	return nil
}

func gitOpenErr(err error) error {
	return fmt.Errorf("cannot open: %w", err)
}
