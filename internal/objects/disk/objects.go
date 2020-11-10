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
		err := o.StoreBlob(stringSum, bs)
		if err != nil {
			return "", errors.Wrapf(err, "couldn't hash object")
		}
	}
	return stringSum, nil
}

func (o *Objects) GetBlob(sum string) (objects.Blob, error) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	var obj objects.Blob
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Blob{}, errors.Wrapf(err, "couldn't get blob %s", sum)
	}
	err = json.Unmarshal(bs, &obj)
	if err != nil {
		return objects.Blob{}, errors.Wrapf(err, "couldn't get blob %s", sum)
	}
	return obj, nil
}

func (o *Objects) GetTree(sum string) (objects.Tree, error) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	var tree objects.Tree
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Tree{}, errors.Wrapf(err, "couldn't get tree %s", sum)
	}
	err = json.Unmarshal(bs, &tree)
	if err != nil {
		return objects.Tree{}, errors.Wrapf(err, "couldn't get tree %s", sum)
	}
	return tree, nil
}

func (o *Objects) GetCommit(sum string) (objects.Commit, error) {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	var commit objects.Commit
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Commit{}, errors.Wrapf(err, "couldn't get commit %s", sum)
	}
	err = json.Unmarshal(bs, &commit)
	if err != nil {
		return objects.Commit{}, errors.Wrapf(err, "couldn't get commit %s", sum)
	}
	return commit, nil
}

func (o *Objects) StoreBlob(sum string, bs []byte) error {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	err := filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't store blob %s", sum)
	}
	blob := objects.NewBlob(bs)
	buf, err := json.Marshal(blob)
	if err != nil {
		return errors.Wrap(err, "couldn't store blob")
	}
	return ioutil.WriteFile(file, buf, os.ModePerm)
}

func (o *Objects) StoreTree(sum string, entries []objects.TreeEntry) error {
	dir := sum[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, sum[2:])
	err := filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't store tree %s", sum)
	}
	tree := objects.Tree{
		Entries: entries,
	}
	buf, err := json.Marshal(tree)
	if err != nil {
		return errors.Wrapf(err, "couldn't store tree %s", sum)
	}
	err = ioutil.WriteFile(file, buf, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't store tree %s", sum)
	}
	return nil
}

func (o *Objects) StoreCommit(commit objects.Commit) (string, error) {
	hash := commit.Hash()
	dir := hash[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, hash[2:])
	err := filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't store commit %s", hash)
	}
	buf, err := json.Marshal(commit)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't store commit %s", hash)
	}
	return hash, ioutil.WriteFile(file, buf, os.ModePerm)
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
	_, err = o.GetCommit(sum)
	if err == nil {
		return objects.TypeCommit, nil
	}
	return "", fmt.Errorf("couldn't get type of %s", sum)
}
