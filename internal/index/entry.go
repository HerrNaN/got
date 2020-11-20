package index

import (
	"fmt"
	"os"
	"strings"

	"got/internal/objects"
)

type Entry struct {
	Perm      os.FileMode
	EntryType objects.Type
	ID        objects.ID
	Name      string
}

func (e Entry) String() string {
	return fmt.Sprintf("%-10v %s %-46v %s", e.Perm, e.EntryType, e.ID, e.Name)
}

type Entries []Entry
type EntryMap map[string]Entry

func NewEntry(perm os.FileMode, entryType objects.Type, sum objects.ID, name string) Entry {
	return Entry{
		Perm:      perm,
		EntryType: entryType,
		ID:        sum,
		Name:      name,
	}
}

func (m EntryMap) Slice() Entries {
	var entries Entries
	for _, e := range m {
		entries = append(entries, e)
	}
	return entries
}

func (es Entries) Less(i, j int) bool {
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
