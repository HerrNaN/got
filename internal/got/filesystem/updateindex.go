package filesystem

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func (g *Got) UpdateIndex(filename string) error {
	filepath.Join(g.dir, filename)
	abs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(g.dir, abs)
	sum, err := g.HashFile(rel, false)
	if err != nil {
		return errors.Wrapf(err, "couldn't add file %s to index", filename)
	}
	return g.Index.AddFile(rel, sum)
}
