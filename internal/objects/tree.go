package objects

import (
	"encoding/json"
	"os"
)

type Tree struct {
	Entries []TreeEntry
}

type TreeEntry struct {
	Mode     os.FileMode
	Type     Type
	Name     string
	Checksum string
}

func (t Tree) Type() Type {
	return TypeTree
}

func (t Tree) Content() string {
	bs, _ := json.Marshal(t)
	return string(bs)
}
