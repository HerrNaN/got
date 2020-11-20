package filesystem

import (
	"github.com/pkg/errors"

	"got/internal/objects"
)

func (g *Got) ReadTree(id objects.ID) error {
	tree, err := g.Objects.GetTree(id)
	if err != nil {
		return errors.Wrapf(err, "couldn't read tree %s", id)
	}
	g.Index.AddTreeContents(tree)
	return nil
}
