package objects

import (
	"crypto/sha1"
)

type Blob struct {
	Contents string `json:"contents"`
}

func NewBlob(bs []byte) Blob {
	return Blob{
		Contents: string(bs),
	}
}

func (b Blob) Type() Type {
	return TypeBlob
}

func (b Blob) Content() string {
	return b.Contents
}

func (b Blob) ID() ID {
	return IdFromSum(sha1.Sum([]byte(b.Content())))
}
