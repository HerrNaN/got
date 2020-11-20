package filesystem

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (g *Got) AddPath(paths ...string) error {
	var matches []string
	for _, p := range paths {
		ms, err := filepath.Glob(p)
		if err != nil {
			return errors.Wrapf(err, "couldn't add path %s", p)
		}
		matches = append(matches, ms...)
	}
	for _, m := range matches {
		err := g.forAllFilesInRepo(m, func(path string, info os.FileInfo, err error) error {
			return g.addFile(path)
		})
		if err != nil {
			return errors.Wrapf(err, "couldn't add path %s", m)
		}
	}
	return nil
}

func (g *Got) addFile(filename string) error {
	rel, err := g.repoRel(filename)
	if err != nil {
		return errors.Wrapf(err, "couldn't add path %s", filename)
	}

	hash, err := g.HashFile(rel, true)
	if err != nil {
		return errors.Wrapf(err, "couldn't add path %s", filename)
	}

	err = g.Index.AddFile(rel, hash)
	if err != nil {
		return errors.Wrapf(err, "couldn't add path %s", filename)
	}
	return nil
}
