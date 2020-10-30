package index

import (
	"fmt"
	"os"

	"got/internal/objects"
)

type Index interface {
	SortedEntries() []Entry
	Update(sum string, name string)
}

type Entry struct {
	Perm      os.FileMode
	EntryType objects.Type
	Sum       string
	Name      string
}

func (e Entry) String() string {
	return fmt.Sprintf("%-10v %s %-46v %s", e.Perm, e.EntryType, e.Sum, e.Name)
}
