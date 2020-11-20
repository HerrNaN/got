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
	id := tree.ID()
	err := g.Objects.StoreTree(id, tree.Entries)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't write tree")
	}
	return id, nil
}
