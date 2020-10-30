package objects

import (
	"os"
)

type Objects interface {
	HashObject(bs []byte, store bool, t Type) string
	Get(sum string) (Object, error)
	Store(sum string, bs []byte, t Type)
}

type Object struct {
	Type Type
	Bs   string
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
