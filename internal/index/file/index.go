package file

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"got/internal/got"
	"got/internal/index"
	"got/internal/objects"
)

const (
	indexFile = got.GotRootDir + "index"
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

func ReadFromFile() *Index {
	bs, _ := ioutil.ReadFile(indexFile)
	var i Index
	json.Unmarshal(bs, &i)
	return &i
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

func (i *Index) updateChecksum() {
	var buf []byte
	for _, e := range i.SortedEntries() {
		buf = append(buf, e.String()...)
	}
	i.Checksum = fmt.Sprintf("%x", sha1.Sum(buf))
}

func (i *Index) writeToFile() {
	i.updateChecksum()
	bs, _ := json.Marshal(*i)
	ioutil.WriteFile(indexFile, bs, os.ModePerm)
}
