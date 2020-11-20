package filesystem

import (
	"fmt"

	"github.com/pkg/errors"

	"got/internal/objects"
)

func (g *Got) Commit(message string) error {
	head, err := g.Head()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	tree, err := g.WriteTree()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	commitID, err := g.CommitTree(message, tree, head)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	err = g.moveHead(commitID)
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
