package file

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/pkg/errors"

	"got/internal/got"
	"got/internal/index"
	"got/internal/objects"
)

const (
	indexFile = got.RootDir + "index"
)

type Index struct {
	Version  int
	Entries  index.EntryMap
	Checksum string
}

func NewIndex() *Index {
	return &Index{
		Version: 0,
		Entries: make(index.EntryMap),
	}
}

func ReadFromFile() (*Index, error) {
	bs, _ := ioutil.ReadFile(indexFile)
	var i Index
	err := json.Unmarshal(bs, &i)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read index file")
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

func (i *Index) Update(sum string, name string) {
	stat, _ := os.Stat(name)
	if stat.IsDir() {
		i.Entries[name] = index.NewEntry(stat.Mode(), objects.TypeTree, sum, name)
	} else {
		i.Entries[name] = index.NewEntry(stat.Mode(), objects.TypeBlob, sum, name)
	}
	i.writeToFile()
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

func (i *Index) writeToFile() {
	i.updateChecksum()
	bs, _ := json.Marshal(*i)
	ioutil.WriteFile(indexFile, bs, os.ModePerm)
}
