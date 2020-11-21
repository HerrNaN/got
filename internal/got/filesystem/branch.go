package filesystem

import "github.com/pkg/errors"

func (g *Got) CreateBranch(newBranch string) error {
	id, err := g.idAtHead()
	if err != nil {
		return errors.Wrapf(err, "couldn't create branch %s", newBranch)
	}
	if id == nil {
		return errors.Errorf("cannot create branch before first commit")
	}
	_, err = g.Refs.CreateBranchAt(newBranch, *id)
	if err != nil {
		return errors.Wrapf(err, "couldn't create branch %s", newBranch)
	}
	return nil
}
