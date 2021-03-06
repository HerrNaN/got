package filesystem

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/pkg/errors"

	"got/internal/diff"
	"got/internal/objects"
	"got/internal/pkg/filesystem"
)

func (g *Got) DiffIndexPath(paths ...string) (string, error) {
	buf := bytes.NewBuffer(nil)
	for _, p := range paths {
		err := g.forAllFilesInRepo(p, func(path string, info os.FileInfo, err error) error {
			hs, err := g.diffIndexPath(path)
			if err != nil {
				return err
			}

			if hs != nil {
				fmt.Fprintf(buf, color.OpBold.Sprintf("--- a/%s\n", path))
				fmt.Fprintf(buf, color.OpBold.Sprintf("+++ b/%s\n", path))
				fmt.Fprint(buf, hs)
			}

			return nil
		})
		if err != nil {
			return "", errors.Wrapf(err, "couldn't diff index and HEAD %s", p)
		}
	}
	return buf.String(), nil
}

func (g *Got) diffIndexPath(path string) (diff.Hunks, error) {
	idx, err := g.getContentsFromIndex(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff index with HEAD %s", path)
	}
	headType, err := g.HeadType()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff index with HEAD %s", path)
	}
	if headType == HeadTypeEmpty {
		return nil, errors.New("no head to diff against")
	}
	ref, err := g.HeadAsRef()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff index with HEAD %s", path)
	}
	headID, err := g.Refs.IDFromRef(ref)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff index with HEAD %s", path)
	}
	hd, err := g.getContentsOfPathFromCommit(path, headID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff index with HEAD %s", path)
	}
	return g.Differ.DiffBytes(hd, idx).Strip(), nil
}

func (g *Got) DiffPathSpec(pathspecs ...string) (string, error) {
	var matches []string
	for _, ps := range pathspecs {
		ms, err := filepath.Glob(ps)
		if err != nil {
			return "", errors.Wrapf(err, "couldn't diff path %s", ps)
		}
		matches = append(matches, ms...)
	}
	return g.DiffPath(matches...)
}

func (g *Got) DiffPath(paths ...string) (string, error) {
	buf := bytes.NewBuffer(nil)
	for _, p := range paths {
		err := g.forAllFilesInRepo(p, func(path string, info os.FileInfo, err error) error {
			hs, err := g.diffPath(path)
			if err != nil {
				return err
			}
			if hs != nil {
				fmt.Fprintf(buf, color.OpBold.Sprintf("--- a/%s\n", path))
				fmt.Fprintf(buf, color.OpBold.Sprintf("+++ b/%s\n", path))
				fmt.Fprint(buf, hs)
			}
			return nil
		})
		if err != nil {
			return "", errors.Wrapf(err, "couldn't diff path %s", p)
		}
	}
	return buf.String(), nil
}

func (g *Got) diffPath(path string) (diff.Hunks, error) {
	abs := filepath.Join(g.dir, path)
	wt, err := g.getContentsFromWorkingTree(abs)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff path %s", abs)
	}

	idx, err := g.getContentsFromIndex(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff path %s", path)
	}

	return g.Differ.DiffBytes(idx, wt).Strip(), nil
}

func (g *Got) getContentsOfPathFromCommit(path string, commitID objects.ID) ([]byte, error) {
	commit, err := g.Objects.GetCommit(commitID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get contents for path %s in commit %s", path, commitID)
	}
	tree, err := g.Objects.GetTree(commit.TreeID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get contents for path %s in commit %s", path, commitID)
	}

	for _, te := range tree.Entries {
		if te.Name == path {
			blob, err := g.Objects.GetBlob(te.ID)
			if err != nil {
				return nil, errors.Wrapf(err, "couldn't get contents for path %s in commit %s", path, commitID)
			}
			return []byte(blob.Content()), nil
		}
	}
	return nil, nil
}

func (g *Got) getContentsFromWorkingTree(path string) ([]byte, error) {
	if !filesystem.FileExists(path) {
		return nil, nil
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff path %s", path)
	}

	return bs, nil
}

func (g *Got) getContentsFromIndex(path string) ([]byte, error) {
	if !g.Index.HasEntryFor(path) {
		return nil, nil
	}

	hash, err := g.Index.GetEntrySum(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff path %s", path)
	}

	blob, err := g.Objects.GetBlob(hash)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't diff path %s", path)
	}

	return []byte(blob.Contents), nil
}
