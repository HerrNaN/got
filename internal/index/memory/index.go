package memory

import (
	"fmt"
	"os"
	"sort"

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

func (i *Index) Update(sum string, name string) {
	stat, _ := os.Stat(name)
	if stat.IsDir() {
		i.Entries[name] = index.NewEntry(stat.Mode(), objects.TypeTree, sum, name)
	} else {
		i.Entries[name] = index.NewEntry(objects.NORM, objects.TypeBlob, sum, name)
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
