package disk

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"got/internal/got"

	"got/internal/objects"
	"got/internal/pkg/filesystem"
)

type Objects struct{}

func NewObjects() *Objects {
	filesystem.MkDirIfIsNotExist(objectsDir, os.ModePerm)
	return &Objects{}
}

const (
	objectsDir = got.RootDir + "objects/"
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
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		return objects.Object{}, err
	}
	return obj, nil
}

func (o *Objects) Store(sum string, bs []byte, t objects.Type) {
	dir := sum[:2]
	file := fmt.Sprintf("%s%s/%s", objectsDir, dir, sum[2:])
	filesystem.MkDirIfIsNotExist(objectsDir+dir, os.ModePerm)
	obj := objects.Object{
		Type: t,
		Size: len(bs),
		Bs:   string(bs),
	}
	buf, _ := json.Marshal(obj)
	ioutil.WriteFile(file, buf, os.ModePerm)
}
