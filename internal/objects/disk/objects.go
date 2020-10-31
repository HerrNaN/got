package disk

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"

	"got/internal/got"

	"got/internal/objects"
	"got/internal/pkg/filesystem"
)

type Objects struct{}

func NewObjects() *Objects {
	filesystem.MkDirIfIsNotExist(got.GotRootDir, os.ModePerm)
	filesystem.MkDirIfIsNotExist(objectsDir, os.ModePerm)
	return &Objects{}
}

const (
	objectsDir = got.GotRootDir + "objects/"
)

func (o *Objects) HashObject(bs []byte, store bool, t objects.Type) string {
	sum := sha1.Sum(bs)
	stringSum := fmt.Sprintf("%x", sum)
	if store {
		o.Store(stringSum, bs, t)
	}
	return stringSum
}

func (o *Objects) Get(sum string) (objects.Object, error) {
	dir := sum[:2]
	file := fmt.Sprintf("%s%s/%s", objectsDir, dir, sum[2:])
	var obj objects.Object
	bs, _ := ioutil.ReadFile(file)
	obj.Bs = string(bs)
	obj.Type = objects.TypeBlob
	return obj, nil
}

func (o *Objects) Store(sum string, bs []byte, t objects.Type) {
	dir := sum[:2]
	file := fmt.Sprintf("%s%s/%s", objectsDir, dir, sum[2:])
	filesystem.MkDirIfIsNotExist(objectsDir+dir, os.ModePerm)
	ioutil.WriteFile(file, bs, os.ModePerm)
}
