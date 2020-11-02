package objects

import (
	"os"
)

type Objects interface {
	HashObject(bs []byte, store bool, t Type) string
	GetBlob(sum string) (Blob, error)
	GetTree(sum string) (Tree, error)
	StoreBlob(sum string, bs []byte)
	StoreTree(sum string, entries []TreeEntry)
}

type Blob struct {
	Size    int
	Content string
}

type Tree struct {
	Entries []TreeEntry
}

type TreeEntry struct {
	Mode     os.FileMode
	Name     string
	Checksum string
}

const (
	NORM os.FileMode = 100644
	EXEC os.FileMode = 100755
	SYMB os.FileMode = 120000
)

type Type string

const (
	TypeBlob   Type = "blob"
	TypeTree   Type = "tree"
	TypeCommit Type = "commit"
)
