package disk

import (
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

func (o *Objects) Store(obj objects.Object) error {
	id := obj.ID()
	dir := string(id)[:2]
	file := filepath.Join(o.dir, ObjectsDir, dir, string(id)[2:])
	err := filesystem.MkDirIfIsNotExist(filepath.Join(o.dir, ObjectsDir, dir), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't store %s %s", obj.Type(), id)
	}
	buf, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrapf(err, "couldn't store %s %s", obj.Type(), id)
	}
	return ioutil.WriteFile(file, buf, os.ModePerm)
}

func (o *Objects) GetBlob(id objects.ID) (objects.Blob, error) {
	file := filepath.Join(o.dir, ObjectsDir, string(id)[:2], string(id)[2:])
	var obj objects.Blob
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Blob{}, errors.Wrapf(err, "couldn't get blob %s", id)
	}
	err = json.Unmarshal(bs, &obj)
	if err != nil {
		return objects.Blob{}, errors.Wrapf(err, "couldn't get blob %s", id)
	}
	return obj, nil
}

func (o *Objects) GetTree(id objects.ID) (objects.Tree, error) {
	file := filepath.Join(o.dir, ObjectsDir, string(id)[:2], string(id)[2:])
	var tree objects.Tree
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Tree{}, errors.Wrapf(err, "couldn't get tree %s", id)
	}
	err = json.Unmarshal(bs, &tree)
	if err != nil {
		return objects.Tree{}, errors.Wrapf(err, "couldn't get tree %s", id)
	}
	return tree, nil
}

func (o *Objects) GetCommit(id objects.ID) (objects.Commit, error) {
	file := filepath.Join(o.dir, ObjectsDir, string(id)[:2], string(id)[2:])
	var commit objects.Commit
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return objects.Commit{}, errors.Wrapf(err, "couldn't get commit %s", id)
	}
	err = json.Unmarshal(bs, &commit)
	if err != nil {
		return objects.Commit{}, errors.Wrapf(err, "couldn't get commit %s", id)
	}
	return commit, nil
}

func (o *Objects) TypeOf(id objects.ID) (objects.Type, error) {
	_, err := o.GetBlob(id)
	if err == nil {
		return objects.TypeBlob, nil
	}
	_, err = o.GetTree(id)
	if err == nil {
		return objects.TypeTree, nil
	}
	_, err = o.GetCommit(id)
	if err == nil {
		return objects.TypeCommit, nil
	}
	return "", fmt.Errorf("couldn't get type of %s", id)
}
