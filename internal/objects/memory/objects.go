package memory

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"got/internal/objects"
)

type Objects struct {
	blobs map[string]objects.Blob
	trees map[string]objects.Tree
}

func NewObjects() *Objects {
	return &Objects{
		blobs: make(map[string]objects.Blob),
		trees: make(map[string]objects.Tree),
	}
}

func (o *Objects) HashObject(bs []byte, store bool, t objects.Type) string {
	sum := sha1.Sum(bs)
	stringSum := fmt.Sprintf("%x", sum)
	if store {
		o.StoreBlob(stringSum, bs)
	}
	return stringSum
}

func (o *Objects) GetBlob(sum string) (objects.Blob, error) {
	blob, ok := o.blobs[sum]
	if ok {
		return blob, nil
	}
	return objects.Blob{}, errors.New("object not found")
}

func (o *Objects) GetTree(sum string) (objects.Tree, error) {
	tree, ok := o.trees[sum]
	if ok {
		return tree, nil
	}
	return objects.Tree{}, errors.New("object not found")
}

func (o *Objects) StoreTree(sum string, entries []objects.TreeEntry) {
	o.trees[sum] = objects.Tree{
		Entries: entries,
	}
}

func (o *Objects) StoreBlob(sum string, bs []byte) {
	o.blobs[sum] = objects.Blob{
		Size:    len(bs),
		Content: string(bs),
	}
}

func (o *Objects) String() string {
	var buf string
	for sum, blob := range o.blobs {
		buf += fmt.Sprintf("# %-8v %v\n%s\n\n", sum[:8], objects.TypeBlob, blob.Content)
	}
	for sum, tree := range o.trees {
		buf += fmt.Sprintf("# %-8v %v\n%v\n\n", sum[:8], objects.TypeTree, tree.Entries)
	}
	return buf
}
