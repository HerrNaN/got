package filesystem

import (
	"github.com/pkg/errors"

	"got/internal/objects"
)

func (g *Got) WriteTree() (string, error) {
	var entries []objects.TreeEntry
	for _, e := range g.Index.SortedEntries() {
		entries = append(entries, objects.TreeEntry{
			Mode:     e.Perm,
			Type:     e.EntryType,
			Name:     e.Name,
			Checksum: e.Sum,
		})
	}
	tree := objects.Tree{
		Entries: entries,
	}
	sum := tree.Hash()
	err := g.Objects.StoreTree(sum, tree.Entries)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't write tree")
	}
	return sum, nil
}
