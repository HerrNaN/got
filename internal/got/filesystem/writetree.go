package filesystem

import (
	"github.com/pkg/errors"

	"got/internal/objects"
)

func (g *Got) WriteTree() (objects.ID, error) {
	var entries []objects.TreeEntry
	for _, e := range g.Index.SortedEntries() {
		entries = append(entries, objects.TreeEntry{
			Mode: e.Perm,
			Type: e.EntryType,
			Name: e.Name,
			ID:   e.ID,
		})
	}
	tree := objects.Tree{
		Entries: entries,
	}
	err := g.Objects.Store(tree)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't write tree")
	}
	return tree.ID(), nil
}
