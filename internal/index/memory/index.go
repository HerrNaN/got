package memory

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"got/internal/index"
	"got/internal/objects"
)

/*const (
	indexFile = "index"
)*/

func NewIndexEntry(perm os.FileMode, entryType objects.Type, sum string, name string) index.Entry {
	return index.Entry{
		Perm:      perm,
		EntryType: entryType,
		Sum:       sum,
		Name:      name,
	}
}

type IndexEntries []index.Entry
type EntryMap map[string]index.Entry

func (m EntryMap) toSlice() IndexEntries {
	var entries IndexEntries
	for _, e := range m {
		entries = append(entries, e)
	}
	return entries
}

type Index struct {
	Entries EntryMap
}

func NewIndex() *Index {
	return &Index{
		Entries: make(EntryMap),
	}
}

func (es IndexEntries) less(i, j int) bool {
	switch strings.Compare(es[i].Name, es[j].Name) {
	case -1:
		return true
	case 1:
		return false
	// Should be based on stage field when we get to merging
	default:
		return false
	}
}

func (i *Index) SortedEntries() []index.Entry {
	entries := i.Entries.toSlice()
	sort.Slice(entries, entries.less)
	return entries
}

func (i *Index) Update(sum string, name string) {
	stat, _ := os.Stat(name)
	if stat.IsDir() {
		i.Entries[name] = NewIndexEntry(stat.Mode(), objects.TypeTree, sum, name)
	} else {
		i.Entries[name] = NewIndexEntry(objects.NORM, objects.TypeBlob, sum, name)
	}
}

func (i *Index) String() string {
	var buf string
	for _, e := range i.SortedEntries() {
		buf += fmt.Sprintln(e)
	}
	return buf
}

/*
func (i *Index) Save() {
	bs, err := memory.Marshal(i)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(indexFile, bs, objects.NORM)
}

func Load() *Index {
	bs, err := ioutil.ReadFile(indexFile)
	if err != nil {
		return &Index{
			Entries: make(map[string]index.Entry),
		}
	}
	var i Index
	memory.Unmarshal(bs, &i)
	return &i
}
*/
