package objects

import (
	"crypto/sha1"
	"encoding/json"
	"os"
)

type Tree struct {
	Entries []TreeEntry
}

type TreeEntry struct {
	Mode os.FileMode
	Type Type
	Name string
	ID   ID
}

func (t Tree) Type() Type {
	return TypeTree
}

func (t Tree) Content() string {
	bs, _ := json.Marshal(t)
	return string(bs)
}

func (t Tree) ID() ID {
	return IdFromSum(sha1.Sum([]byte(t.Content())))
}
