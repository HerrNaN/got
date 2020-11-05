package memory

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/pkg/errors"

	"got/internal/index"
	"got/internal/objects"
)

type Index struct {
	Entries index.EntryMap
}

func NewIndex() *Index {
	return &Index{
		Entries: make(index.EntryMap),
	}
}

func (i *Index) SortedEntries() []index.Entry {
	entries := i.Entries.Slice()
	sort.Slice(entries, entries.Less)
	return entries
}

func (i *Index) HasEntryFor(filename string) bool {
	_, ok := i.Entries[filename]
	return ok
}

func (i *Index) AddFile(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "couldn't read file %s", filename)
	}
	stat, _ := os.Stat(filename)
	sum := fmt.Sprintf("%x", sha1.Sum(bs))
	i.Entries[filename] = index.NewEntry(stat.Mode(), objects.TypeBlob, sum, filename)
	return nil
}

func (i *Index) AddTreeContents(tree objects.Tree) {
	for _, e := range tree.Entries {
		i.Entries[e.Name] = index.NewEntry(e.Mode, e.Type, e.Checksum, e.Name)
	}
}

func (i *Index) AddTree(sum string, prefix string) {
	i.Entries[prefix] = index.NewEntry(os.ModePerm, objects.TypeTree, sum, prefix)
}

func (i *Index) String() string {
	var buf string
	for _, e := range i.SortedEntries() {
		buf += fmt.Sprintln(e)
	}
	return buf
}
