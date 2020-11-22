package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"got/internal/refs"

	"github.com/pkg/errors"

	"got/internal/objects"
)

type HeadType string

const (
	HeadTypeID    = "ID"
	HeadTypeRef   = "ref"
	HeadTypeEmpty = "empty"
)

func (g *Got) HeadType() (HeadType, error) {
	bs, err := ioutil.ReadFile(filepath.Join(g.dir, headFile))
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get head type")
	}
	if len(bs) == 0 {
		return HeadTypeEmpty, nil
	}
	_, err = refs.RefFromString(string(bs))
	if err == nil {
		return HeadTypeRef, nil
	}
	_, err = objects.IdFromString(string(bs))
	if err == nil {
		return HeadTypeID, nil
	}
	return "", errors.New("couldn't parse head to any known type")
}

func (g *Got) HeadAsID() (objects.ID, error) {
	bs, err := ioutil.ReadFile(filepath.Join(g.dir, headFile))
	if err != nil {
		return "", errors.Wrap(err, "couldn't get head as ID")
	}
	id, err := objects.IdFromString(string(bs))
	if err != nil {
		return "", errors.Wrap(err, "couldn't get head as ID")
	}
	return id, nil
}

func (g *Got) HeadAsRef() (refs.Ref, error) {
	bs, err := ioutil.ReadFile(filepath.Join(g.dir, headFile))
	if err != nil {
		return "", errors.Wrap(err, "couldn't get head as Ref")
	}
	ref, err := refs.RefFromString(string(bs))
	if err != nil {
		return "", errors.Wrap(err, "couldn't get head as Ref")
	}
	return ref, nil
}

func (g *Got) idAtHead() (*objects.ID, error) {
	headType, err := g.HeadType()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get id at HEAD")
	}
	if headType == HeadTypeEmpty {
		return nil, nil
	}
	ref, err := g.HeadAsRef()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get id at HEAD")
	}
	id, err := g.Refs.IDFromRef(ref)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tree at HEAD")
	}
	return &id, nil
}

func (g *Got) headTree() (*objects.Tree, error) {
	headID, err := g.idAtHead()
	if err != nil {
		return nil, err
	}
	if headID == nil {
		return nil, nil
	}
	c, err := g.Objects.GetCommit(*headID)
	if err != nil {
		return nil, err
	}
	tree, err := g.Objects.GetTree(c.TreeID)
	return &tree, err

}

func (g *Got) updateHeadWithID(id objects.ID) error {
	err := ioutil.WriteFile(filepath.Join(g.dir, headFile), []byte(id), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't update HEAD with id %s", id)
	}
	return nil
}

func (g *Got) updateHeadWithRef(ref refs.Ref) error {
	err := ioutil.WriteFile(filepath.Join(g.dir, headFile), []byte(ref), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't update HEAD with ref %s", ref)
	}
	return nil
}
