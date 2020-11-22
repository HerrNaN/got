package filesystem

import (
	"crypto/sha1"
	"io/ioutil"
	"os"

	"got/internal/diff"
	"got/internal/diff/simple"
	"got/internal/index"
	"got/internal/objects"
	"got/internal/status"
)

type fileInfo struct {
	name string
	hash objects.ID
	perm os.FileMode
}

func (g *Got) Status() (*status.Status, error) {
	tree, err := g.statusTree()
	if err != nil {
		return nil, err
	}
	return tree.GetStatus(), nil
}

func (g *Got) statusTree() (*status.Tree, error) {
	headDiff, err := g.diffHead()
	if err != nil {
		return nil, err
	}
	workTreeDiff, untracked, err := g.diffFiles()
	if err != nil {
		return nil, err
	}
	tree := status.NewTree()

	for _, d := range headDiff {
		switch d.EditType {
		case diff.FileEditTypeInPlace:
			tree.AddFile(d.SrcPath, status.Changes{Head: status.Modified}, true)
		case diff.FileEditTypeDelete:
			tree.AddFile(d.SrcPath, status.Changes{Head: status.Deleted}, true)
		case diff.FileEditTypeCreate:
			tree.AddFile(d.DstPath, status.Changes{Head: status.Created}, true)
		}
	}

	for _, d := range workTreeDiff {
		switch d.EditType {
		case diff.FileEditTypeInPlace:
			tree.AddFile(d.SrcPath, status.Changes{Worktree: status.Modified}, true)
		case diff.FileEditTypeDelete:
			tree.AddFile(d.SrcPath, status.Changes{Worktree: status.Deleted}, true)
		case diff.FileEditTypeCreate:
			tree.AddFile(d.DstPath, status.Changes{Worktree: status.Created}, true)
		}
	}

	for _, d := range untracked {
		tree.AddFile(d, status.Changes{}, false)
	}

	return tree, nil
}

func (g *Got) diffHead() ([]*diff.FileDiff, error) {
	var diffs []*diff.FileDiff

	headTree, err := g.headTree()
	if err != nil {
		return nil, err
	}
	for _, ie := range g.Index.SortedEntries() {
		var d *diff.FileDiff
		if headTree != nil {
			d, err = g.diffEntryAgainstHead(ie, headTree)
			if err != nil {
				return nil, err
			}
		}
		if d == nil {
			d = diff.NewCreateFileDiff(ie.Perm, ie.ID, ie.Name)
		}
		diffs = append(diffs, d)
	}
	if headTree == nil {
		return diffs, nil
	}
	for _, te := range headTree.Entries {
		if !g.Index.HasEntryFor(te.Name) {
			diffs = append(diffs, diff.NewDeleteFileDiff(te.Mode, te.ID, te.Name))
		}
	}
	return diffs, nil
}

func (g *Got) diffFiles() ([]*diff.FileDiff, []string, error) {
	var untracked []string
	var diffs []*diff.FileDiff
	var files []*fileInfo
	err := g.forAllInRepo(g.dir, func(path string, info os.FileInfo, err error) error {
		path, err = g.repoRel(path)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			hash := objects.IdFromSum(sha1.Sum(bs))
			files = append(files, &fileInfo{
				name: path,
				hash: hash,
				perm: info.Mode(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	for _, ie := range g.Index.SortedEntries() {
		d, err := g.diffEntryAgainstFiles(ie, files)
		if err != nil {
			return nil, nil, err
		}
		if d == nil {
			d = diff.NewCreateFileDiff(ie.Perm, ie.ID, ie.Name)
		}
		diffs = append(diffs, d)
	}
	for _, f := range files {
		if !g.Index.HasEntryFor(f.name) {
			untracked = append(untracked, f.name)
		}
	}
	return diffs, untracked, nil
}

func (g *Got) diffEntryAgainstHead(ie index.Entry, headTree *objects.Tree) (*diff.FileDiff, error) {
	d := simple.Diff{}
	for _, te := range headTree.Entries {
		if ie.Name != te.Name {
			continue
		}
		if ie.ID == te.ID {
			return diff.NewUnmodifiedFileDiff(ie.Perm, ie.ID, ie.Name), nil
		}
		iBlob, err := g.Objects.GetBlob(ie.ID)
		if err != nil {
			return nil, err
		}
		tBlob, err := g.Objects.GetBlob(te.ID)
		if err != nil {
			return nil, err
		}
		_, err = d.DiffFiles([]byte(iBlob.Contents), []byte(tBlob.Contents))
		if err != nil {
			return nil, err
		}
		return diff.NewInPlaceFileDiff(te.Mode, ie.Perm, te.ID, ie.ID, ie.Name), nil
	}
	return nil, nil
}

func (g *Got) diffEntryAgainstFiles(ie index.Entry, files []*fileInfo) (*diff.FileDiff, error) {
	d := simple.Diff{}
	for _, f := range files {
		if ie.Name != f.name {
			continue
		}
		if ie.ID == f.hash {
			return diff.NewUnmodifiedFileDiff(f.perm, f.hash, f.name), nil
		}
		iBlob, err := g.Objects.GetBlob(ie.ID)
		if err != nil {
			return nil, err
		}
		contents, err := ioutil.ReadFile(f.name)
		if err != nil {
			return nil, err
		}
		_, err = d.DiffFiles([]byte(iBlob.Contents), contents)
		if err != nil {
			return nil, err
		}
		return diff.NewInPlaceFileDiff(f.perm, ie.Perm, f.hash, ie.ID, f.name), nil
	}
	return nil, nil
}
