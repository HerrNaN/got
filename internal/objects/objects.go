package objects

import (
	"os"
)

type Objects interface {
	// Retrieves a Blob from a given ID
	GetBlob(id ID) (Blob, error)

	// Retrieves a Tree from a given ID
	GetTree(id ID) (Tree, error)

	// Retrieves a Commit from a given ID
	GetCommit(id ID) (Commit, error)

	// Stores the given Object
	Store(object Object) error

	// Returns the Type of the objects with the given ID
	TypeOf(id ID) (Type, error)
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
	ID() ID
}
