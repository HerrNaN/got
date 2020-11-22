package filesystem

import (
	"bytes"
	"fmt"

	"got/internal/refs"

	"github.com/gookit/color"

	"github.com/pkg/errors"
)

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

func (g *Got) DeleteBranch(branchName string) error {
	err := g.Refs.DeleteRef(branchName)
	if err != nil {
		return errors.Wrapf(err, "couldn't delete branch %s", branchName)
	}
	return nil
}

func (g *Got) ListBranches() (Branches, error) {
	branches, err := g.Refs.Branches()
	if err != nil {
		return Branches{}, errors.Wrapf(err, "couldn't list branches")
	}
	headRef, err := g.HeadAsRef()
	if err != nil {
		return Branches{}, errors.Wrapf(err, "couldn't list branches")
	}
	return Branches{branches, headRef}, nil
}

type Branches struct {
	list    []string
	current refs.Ref
}

func (bs Branches) String() string {
	buf := bytes.NewBuffer(nil)
	for _, b := range bs.list {
		if b == bs.current.Name() {
			fmt.Fprintln(buf, "* "+color.Green.Sprint(b))
		} else {
			fmt.Fprintln(buf, "  "+b)
		}
	}
	return buf.String()
}
