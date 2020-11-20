package filesystem

import "github.com/pkg/errors"

func (g *Got) ReadTree(sum string) error {
	tree, err := g.Objects.GetTree(sum)
	if err != nil {
		return errors.Wrapf(err, "couldn't read tree %s", sum)
	}
	g.Index.AddTreeContents(tree)
	return nil
}
