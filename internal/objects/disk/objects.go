package disk

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"got/internal/objects"
	"got/internal/pkg/filesystem"
)

type Objects struct {
	dir string
}

func NewObjects(dir string) *Objects {
	return &Objects{
		dir: dir,
	}
}

const (
	ObjectsDir = "objects"
)

func (o *Objects) HashObject(bs []byte, store bool, t objects.Type) (string, error) {
	sum := sha1.Sum(bs)
	stringSum := fmt.Sprintf("%x", sum)
	if store {
		return stringSum, o.StoreBlob(stringSum, bs)
	}
	return stringSum, nil
}

func (o *Objects) GetBlob(sum string) (objects.Blob, error) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	var obj objects.Blob
	bs, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		return objects.Blob{}, err
	}
	return obj, nil
}

func (o *Objects) GetTree(sum string) (objects.Tree, error) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	var tree objects.Tree
	bs, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(bs, &tree)
	if err != nil {
		return objects.Tree{}, err
	}
	return tree, nil
}

func (o *Objects) StoreTree(sum string, entries []objects.TreeEntry) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	tree := objects.Tree{
		Entries: entries,
	}
	buf, _ := json.Marshal(tree)
	ioutil.WriteFile(file, buf, os.ModePerm)
}

func (o *Objects) StoreBlob(sum string, bs []byte) error {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	err := filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't store blob %s", sum)
	}
	blob := objects.NewBlob(bs)
	buf, _ := json.Marshal(blob)
	return ioutil.WriteFile(file, buf, os.ModePerm)
}

func (o *Objects) TypeOf(sum string) (objects.Type, error) {
	_, err := o.GetBlob(sum)
	if err == nil {
		return objects.TypeBlob, nil
	}
	_, err = o.GetTree(sum)
	if err == nil {
		return objects.TypeTree, nil
	}
	return "", fmt.Errorf("couldn't get type of %s", sum)
}
