package repository

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type Git struct {
	repo *git.Repository
}

// New initialize the repository and fetches the remote content.
func New(repositoryUrl string) (*Git, error) {
	r := &Git{}
	var err error
	r.repo, err = git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, err
	}

	if _, err := r.repo.CreateRemote(&config.RemoteConfig{Name: "upstream", URLs: []string{repositoryUrl}}); err != nil {
		return nil, err
	}
	if err = r.repo.Fetch(&git.FetchOptions{RemoteName: "upstream", Tags: git.AllTags}); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Git) ListUpstreamCommits(upstreamTagName, downstreamBranchName string) ([]*object.Commit, error) {
	upstreamReference, err := r.repo.Tag(upstreamTagName)
	if err != nil {
		return nil, fmt.Errorf("failed to checkout tag %q: %v", upstreamTagName, err)
	}

	downstreamReference, err := r.repo.Reference(plumbing.NewRemoteReferenceName("upstream", downstreamBranchName), true)
	if err != nil {
		return nil, err
	}
	if downstreamReference == nil {
		return nil, fmt.Errorf("downstream reference %s not found", downstreamBranchName)
	}

	upstreamCommitHash, _ := r.repo.ResolveRevision(plumbing.Revision(upstreamReference.Hash().String()))
	upstreamCommit, err := r.repo.CommitObject(plumbing.NewHash(upstreamCommitHash.String()))
	if err != nil {
		return nil, fmt.Errorf("unable to get upstream commit: %v", err)
	}

	downstreamCommit, err := r.repo.CommitObject(downstreamReference.Hash())
	if err != nil {
		return nil, fmt.Errorf("unable to get downstream commit: %v", err)
	}

	logIterator, err := r.repo.Log(&git.LogOptions{
		From: downstreamCommit.Hash,
	})
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit

	if err := logIterator.ForEach(func(commit *object.Commit) error {
		isAncestor, err := commit.IsAncestor(upstreamCommit)
		if err != nil {
			return err
		}
		// skip commits that are already present in upstream
		if isAncestor {
			return nil
		}
		// skip "merge" commits
		if commit.NumParents() > 1 {
			return nil
		}
		// skip "empty" commits
		if isEmptyCommit(commit) {
			return nil
		}
		commits = append(commits, commit)
		return nil
	}); err != nil {
		return nil, err
	}
	return commits, nil
}

func isEmptyCommit(commit *object.Commit) bool {
	stats, _ := commit.Stats()
	return len(stats) == 0
}
