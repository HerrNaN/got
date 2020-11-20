package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (g *Got) UnstagePath(paths ...string) error {
	var matches []string
	for _, p := range paths {
		ms, err := filepath.Glob(p)
		if err != nil {
			return errors.Wrapf(err, "couldn't unstage path %s", p)
		}
		matches = append(matches, ms...)
	}
	for _, m := range matches {
		err := g.forAllFilesInRepo(m, func(localPath string, info os.FileInfo, err error) error {
			return g.unstageFile(localPath)
		})
		if err != nil {
			return errors.Wrapf(err, "couldn't unstage path %s", m)
		}
	}
	return nil
}

func (g *Got) unstageFile(filename string) error {
	rel, err := g.repoRel(filename)
	if err != nil {
		return errors.Wrapf(err, "couldn't unstage %s", filename)
	}
	headTree, err := g.headTree()
	if err != nil {
		return errors.Wrapf(err, "couldn't unstage %s", rel)
	}
	for _, te := range headTree.Entries {
		if te.Name == rel {
			return g.Index.AddFile(te.Name, te.Checksum)
		}
	}
	err = g.Index.RemoveFile(rel)
	if err != nil {
		return errors.Wrapf(err, "couldn't unstage %s", rel)
	}
	return nil
}

func (g *Got) DiscardPath(paths ...string) error {
	var matches []string
	for _, p := range paths {
		ms, err := filepath.Glob(p)
		if err != nil {
			return errors.Wrapf(err, "couldn't discard path %s", p)
		}
		matches = append(matches, ms...)
	}
	for _, m := range matches {
		err := g.forAllFilesInRepo(m, func(path string, info os.FileInfo, err error) error {
			return g.discardFile(path)
		})
		if err != nil {
			return errors.Wrapf(err, "couldn't discard path %s", m)
		}
	}
	return nil
}

func (g *Got) discardFile(filename string) error {
	rel, err := g.repoRel(filename)
	if err != nil {
		return errors.Wrapf(err, "couldn't discard changes in %s", filename)
	}
	headTree, err := g.headTree()
	if err != nil {
		return errors.Wrapf(err, "couldn't discard changes in %s", rel)
	}
	for _, te := range headTree.Entries {
		if te.Name == rel {
			blob, err := g.Objects.GetBlob(te.Checksum)
			if err != nil {
				return errors.Wrapf(err, "couldn't discard changes in %s", rel)
			}
			err = ioutil.WriteFile(filename, []byte(blob.Contents), te.Mode)
			if err != nil {
				return errors.Wrapf(err, "couldn't discard changes in %s", rel)
			}
		}
	}
	if !g.Index.HasEntryFor(rel) {
		return fmt.Errorf("%s did not match any file(s) know to got", filename)
	}
	return nil
}
