package objects

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
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

func (t Tree) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(t.Content())))
}
