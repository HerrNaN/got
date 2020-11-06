package file

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	"got/internal/index"
	"got/internal/objects"
)

const (
	IndexFile = "index"
)

type Index struct {
	Dir      string
	Version  int
	Entries  index.EntryMap
	Checksum string
}

func NewIndex(dir string) *Index {
	return &Index{
		Dir:     dir,
		Version: 0,
		Entries: make(index.EntryMap),
	}
}

func ReadFromFile(dir string) (*Index, error) {
	bs, err := ioutil.ReadFile(filepath.Join(dir, IndexFile))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read index file")
	}
	if len(bs) == 0 {
		return NewIndex(dir), nil
	}
	var i Index
	err = json.Unmarshal(bs, &i)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal index file")
	}
	if i.calculateChecksum() != i.Checksum {
		return nil, errors.New("index file corrupted")
	}
	return &i, nil
}

func (i *Index) SortedEntries() []index.Entry {
	entries := i.Entries.Slice()
	sort.Slice(entries, entries.Less)
	return entries
}

func (i *Index) AddFile(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "couldn't read file %s", filename)
	}
	stat, _ := os.Stat(filename)
	sum := fmt.Sprintf("%x", sha1.Sum(bs))
	i.Entries[filename] = index.NewEntry(stat.Mode(), objects.TypeBlob, sum, filename)
	return i.writeToFile()
}

func (i *Index) AddTreeContents(tree objects.Tree) error {
	for _, e := range tree.Entries {
		i.Entries[e.Name] = index.NewEntry(e.Mode, e.Type, e.Checksum, e.Name)
	}
	return i.writeToFile()
}

func (i *Index) AddTree(sum string, prefix string) error {
	i.Entries[prefix] = index.NewEntry(os.ModePerm, objects.TypeTree, sum, prefix)
	return i.writeToFile()
}

func (i *Index) GetEntryFor(name string) (index.Entry, error) {
	e, ok := i.Entries[name]
	if !ok {
		return index.Entry{}, fmt.Errorf("couldn't find entry for file %s", name)
	}
	return e, nil
}

func (i *Index) GetEntry(sum string, name string) (index.Entry, error) {
	e, err := i.GetEntryFor(name)
	if err != nil {
		return index.Entry{}, err
	}
	if e.Sum != sum {
		return index.Entry{}, fmt.Errorf("found entry does not match sum %s", sum)
	}
	return e, nil
}

func (i *Index) HasEntryFor(name string) bool {
	_, ok := i.Entries[name]
	return ok
}

func (i *Index) GetEntrySum(filename string) (string, error) {
	e, ok := i.Entries[filename]
	if !ok {
		return "", errors.New("entry not found")
	}
	return e.Sum, nil
}

func (i *Index) updateChecksum() {
	i.Checksum = i.calculateChecksum()
}

func (i *Index) calculateChecksum() string {
	var buf []byte
	for _, e := range i.SortedEntries() {
		buf = append(buf, e.String()...)
	}
	return fmt.Sprintf("%x", sha1.Sum(buf))
}

func (i *Index) writeToFile() error {
	i.updateChecksum()
	bs, _ := json.Marshal(*i)
	return ioutil.WriteFile(filepath.Join(i.Dir, IndexFile), bs, os.ModePerm)
}
