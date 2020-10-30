package memory

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"got/internal/objects"
)

type Objects map[string]objects.Object

func NewObjects() Objects {
	return make(map[string]objects.Object)
}

func (o Objects) HashObject(bs []byte, store bool, t objects.Type) string {
	sum := sha1.Sum(bs)
	stringSum := fmt.Sprintf("%x", sum)
	if store {
		o.Store(stringSum, bs, t)
	}
	return stringSum
}

func (o Objects) Get(sum string) (objects.Object, error) {
	obj, ok := o[sum]
	if ok {
		return obj, nil
	}
	return objects.Object{}, errors.New("object not found")
}

func (o Objects) Store(sum string, bs []byte, t objects.Type) {
	o[sum] = objects.Object{
		Type: t,
		Bs:   string(bs),
	}
}

func (o Objects) String() string {
	var buf string
	for sum, obj := range o {
		buf += fmt.Sprintf("# %-8v %v\n%s\n\n", sum[:8], obj.Type, obj.Bs)
	}
	return buf
}
