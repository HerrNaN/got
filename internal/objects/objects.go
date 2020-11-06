package objects

import (
	"os"
)

type Objects interface {
	// Computes the ID of the object and optionally stores the object as well
	HashObject(bs []byte, store bool, t Type) string

	// Retrieves a Blob from a given ID
	GetBlob(sum string) (Blob, error)

	// Retrieves a Tree from a given ID
	GetTree(sum string) (Tree, error)

	// Stores content as Blob with a given ID
	StoreBlob(sum string, bs []byte)

	// Stores content as Tree with a given ID
	StoreTree(sum string, entries []TreeEntry)

	// Returns the Type of the objects with the given ID
	TypeOf(sum string) (Type, error)
}

const (
	NORM os.FileMode = 100644
	EXEC os.FileMode = 100755
	SYMB os.FileMode = 120000
	DIR  os.FileMode = os.ModeDir + 100644
)

type Type string

const (
	TypeBlob   Type = "blob"
	TypeTree   Type = "tree"
	TypeCommit Type = "commit"
)

type Object interface {
	Type() Type
	Content() string
}
