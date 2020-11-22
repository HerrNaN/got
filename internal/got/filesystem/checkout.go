package filesystem

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// 1. Update HEAD
// 2. Update WT
func (g *Got) Checkout(branchName string, create bool) error {
	statusTree, err := g.statusTree()
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}
	if statusTree.HasChanges() {
		return errors.New("cannot checkout branch with uncommitted changes")
	}

	if create {
		id, err := g.idAtHead()
		if err != nil {
			return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
		}
		if id == nil {
			return errors.New("cannot create a new branch before first commit")
		}
		_, err = g.Refs.CreateBranchAt(branchName, *id)
		if err != nil {
			return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
		}
	} else {
		if !g.Refs.BranchExists(branchName) {
			return errors.Errorf("branch %s doesn't exist", branchName)
		}
	}

	id, err := g.Refs.IdAtBranch(branchName)
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}
	commit, err := g.Objects.GetCommit(id)
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}
	commitTree, err := g.Objects.GetTree(commit.TreeID)
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}

	// In the future there should probably be some rollback feature to handle if
	// this gets an error midway through. Until then a manual restore should
	// do the trick.
	for _, te := range commitTree.Entries {
		blob, err := g.Objects.GetBlob(te.ID)
		if err != nil {
			return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
		}
		err = ioutil.WriteFile(filepath.Join(g.dir, te.Name), []byte(blob.Contents), te.Mode)
		if err != nil {
			return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
		}
	}

	ref, err := g.Refs.BranchRef(branchName)
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}
	err = g.updateHeadWithRef(ref)
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout branch %s", branchName)
	}
	return nil
}
