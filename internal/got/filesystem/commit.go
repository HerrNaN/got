package filesystem

import (
	"fmt"

	"github.com/pkg/errors"

	"got/internal/objects"
)

func (g *Got) Commit(message string) error {
	headType, err := g.HeadType()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	if headType == HeadTypeRef {
		return g.commitAtRef(message)
	}
	return g.firstCommit(message)
}

func (g *Got) commitAtRef(message string) error {
	ref, err := g.HeadAsRef()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	currentCommitID, err := g.Refs.IDFromRef(ref)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	treeID, err := g.WriteTree()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	// Update branch head if it exists
	newCommitID, err := g.CommitTree(message, treeID, &currentCommitID)
	err = g.Refs.UpdateRef(ref, newCommitID)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	return nil
}

func (g *Got) firstCommit(message string) error {
	treeID, err := g.WriteTree()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	newCommitID, err := g.CommitTree(message, treeID, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	ref, err := g.Refs.CreateBranchAt("master", newCommitID)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	err = g.updateHeadWithRef(ref)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	return nil
}

func (g *Got) CommitTree(msg string, treeID objects.ID, parentID *objects.ID) (objects.ID, error) {
	commit := objects.NewCommit(treeID, parentID, "John Doe <john@doe.com> 0123456789 +0000", msg)
	fmt.Printf("Committing %s", treeID)
	if parentID != nil {
		fmt.Printf(" with parent %s", *parentID)
	}
	fmt.Println("...")
	var buf string
	buf += fmt.Sprintf("tree %s\n", treeID)
	if parentID != nil {
		buf += fmt.Sprintf("parent %s\n", *parentID)
	}
	buf += fmt.Sprintln("author John Doe <john@doe.com> 0123456789 +0000")
	buf += fmt.Sprintln("committer John Doe <john@doe.com> 0123456789 +0000")
	buf += fmt.Sprintf("\n%s", msg)
	return commit.ID(), g.Objects.Store(commit)
}
