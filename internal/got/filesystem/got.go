package filesystem

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"got/internal/status"

	"got/internal/index/file"

	"got/internal/diff/simple"

	"got/internal/diff"

	"github.com/pkg/errors"

	"got/internal/index"
	"got/internal/objects"
	"got/internal/objects/disk"
	"got/internal/pkg/filesystem"
)

const (
	rootDir = ".got"
)

var objectsDir = filepath.Join(rootDir, disk.ObjectsDir)
var indexFile = filepath.Join(rootDir, file.IndexFile)
var headFile = filepath.Join(rootDir, "HEAD")

type Got struct {
	gotDir  string
	dir     string
	Objects objects.Objects
	Index   index.Index
}

func NewGot() (*Got, error) {
	dir, err := getRepositoryRoot()
	if err != nil {
		return nil, errors.New("repository not initialized")
	}
	gotDir := filepath.Join(dir, rootDir)
	if !IsInitialized(dir) {
		return nil, errors.New("repository not initialized")
	}
	i, err := file.ReadFromFile(gotDir)
	if err != nil {
		return nil, err
	}

	return &Got{
		gotDir:  gotDir,
		dir:     dir,
		Objects: disk.NewObjects(gotDir),
		Index:   i,
	}, nil
}

func (g *Got) HashFile(filename string, store bool) (string, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't hash file %s", filename)
	}
	blob := objects.NewBlob(bs)
	sum := blob.Hash()
	if store {
		err = g.Objects.StoreBlob(sum, bs)
		if err != nil {
			return "", errors.Wrapf(err, "couldn't hash file %s", filename)
		}
	}
	return sum, nil
}

func (g *Got) AddToIndex(filename string) error {
	filepath.Join(g.dir, filename)
	abs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(g.dir, abs)
	sum, err := g.HashFile(rel, false)
	if err != nil {
		return errors.Wrapf(err, "couldn't add file %s to index", filename)
	}
	return g.Index.AddFile(rel, sum)
}

func (g *Got) WriteTree() (string, error) {
	var entries []objects.TreeEntry
	for _, e := range g.Index.SortedEntries() {
		entries = append(entries, objects.TreeEntry{
			Mode:     e.Perm,
			Type:     e.EntryType,
			Name:     e.Name,
			Checksum: e.Sum,
		})
	}
	tree := objects.Tree{
		Entries: entries,
	}
	sum := tree.Hash()
	err := g.Objects.StoreTree(sum, tree.Entries)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't write tree")
	}
	return sum, nil
}

func (g *Got) ReadTree(sum string) error {
	tree, err := g.Objects.GetTree(sum)
	if err != nil {
		return errors.Wrapf(err, "couldn't read tree %s", sum)
	}
	g.Index.AddTreeContents(tree)
	return nil
}

func (g *Got) CommitTree(msg string, tree string, parent string) (string, error) {
	commit := objects.NewCommit(tree, parent, "John Doe <john@doe.com> 0123456789 +0000", msg)
	fmt.Printf("Committing %s", tree)
	if parent != "" {
		fmt.Printf(" with parent %s", parent)
	}
	fmt.Println("...")
	var buf string
	buf += fmt.Sprintf("tree %s\n", tree)
	if parent != "" {
		buf += fmt.Sprintf("parent %s\n", parent)
	}
	buf += fmt.Sprintln("author John Doe <john@doe.com> 0123456789 +0000")
	buf += fmt.Sprintln("committer John Doe <john@doe.com> 0123456789 +0000")
	buf += fmt.Sprintf("\n%s", msg)
	return g.Objects.StoreCommit(commit)
}

func (g *Got) Head() (string, error) {
	bs, err := ioutil.ReadFile(filepath.Join(g.dir, headFile))
	if err != nil {
		return "", errors.Wrap(err, "couldn't read HEAD file")
	}
	return string(bs), nil
}

func (g *Got) Add(filename string) error {
	rel, err := g.repoRel(filename)
	if err != nil {
		return errors.Wrap(err, "couldn't add file")
	}
	hash, err := g.HashFile(rel, true)
	if err != nil {
		return errors.Wrap(err, "couldn't add file")
	}
	return g.Index.AddFile(rel, hash)
}

func (g *Got) Status() (*status.Status, error) {
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

	return tree.GetStatus(), nil
}

func (g *Got) diffHead() ([]*diff.FileDiff, error) {
	var diffs []*diff.FileDiff

	headTree, err := g.headTree()
	if err != nil {
		return nil, err
	}
	for _, ie := range g.Index.SortedEntries() {
		d, err := g.diffEntryAgainstHead(ie, headTree)
		if err != nil {
			return nil, err
		}
		if d == nil {
			d = diff.NewCreateFileDiff(ie.Perm, ie.Sum, ie.Name)
		}
		diffs = append(diffs, d)
	}
	for _, te := range headTree.Entries {
		if !g.Index.HasEntryFor(te.Name) {
			diffs = append(diffs, diff.NewDeleteFileDiff(te.Mode, te.Checksum, te.Name))
		}
	}
	return diffs, nil
}

type fileInfo struct {
	name string
	hash string
	perm os.FileMode
}

func (g *Got) diffFiles() ([]*diff.FileDiff, []string, error) {
	var untracked []string
	var diffs []*diff.FileDiff
	var files []*fileInfo
	err := filepath.Walk(g.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path, err = g.repoRel(path)
		if err != nil {
			return err
		}
		if path == ".git" || path == ".got" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			hash := fmt.Sprintf("%x", sha1.Sum(bs))
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
			d = diff.NewCreateFileDiff(ie.Perm, ie.Sum, ie.Name)
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
		if ie.Sum == te.Checksum {
			return diff.NewUnmodifiedFileDiff(ie.Perm, ie.Sum, ie.Name), nil
		}
		iBlob, err := g.Objects.GetBlob(ie.Sum)
		if err != nil {
			return nil, err
		}
		tBlob, err := g.Objects.GetBlob(te.Checksum)
		if err != nil {
			return nil, err
		}
		_, err = d.DiffFiles([]byte(iBlob.Contents), []byte(tBlob.Contents))
		if err != nil {
			return nil, err
		}
		return diff.NewInPlaceFileDiff(te.Mode, ie.Perm, te.Checksum, ie.Sum, ie.Name), nil
	}
	return nil, nil
}

func (g *Got) diffEntryAgainstFiles(ie index.Entry, files []*fileInfo) (*diff.FileDiff, error) {
	d := simple.Diff{}
	for _, f := range files {
		if ie.Name != f.name {
			continue
		}
		if ie.Sum == f.hash {
			return diff.NewUnmodifiedFileDiff(f.perm, f.hash, f.name), nil
		}
		iBlob, err := g.Objects.GetBlob(ie.Sum)
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
		return diff.NewInPlaceFileDiff(f.perm, ie.Perm, f.hash, ie.Sum, f.name), nil
	}
	return nil, nil
}

func (g *Got) repoRel(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	repoRel, err := filepath.Rel(g.dir, abs)
	if err != nil {
		return "", err
	}
	return repoRel, nil
}

func (g *Got) Commit(message string) error {
	head, err := g.Head()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	tree, err := g.WriteTree()
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	commitHash, err := g.CommitTree(message, tree, head)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	err = g.moveHead(commitHash)
	if err != nil {
		return errors.Wrap(err, "couldn't perform commit")
	}
	return nil
}

func (g *Got) moveHead(hash string) error {
	return ioutil.WriteFile(filepath.Join(g.dir, headFile), []byte(hash), os.ModePerm)
}

// Gets list of files that are currently tracked by the repository.
func (g *Got) trackedFiles() ([]string, error) {
	var trackedFiles []string
	head, err := g.Head()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tracked files")
	}
	if head == "" {
		return trackedFiles, nil
	}
	commit, err := g.Objects.GetCommit(head)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tracked files")
	}
	tree, err := g.Objects.GetTree(commit.TreeHash)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tracked files")
	}
	for _, te := range tree.Entries {
		trackedFiles = append(trackedFiles, te.Name)
	}
	return trackedFiles, nil
}

func (g *Got) headTree() (*objects.Tree, error) {
	head, err := g.Head()
	if err != nil {
		return nil, err
	}
	if head == "" {
		return nil, nil
	}
	c, err := g.Objects.GetCommit(head)
	if err != nil {
		return nil, err
	}
	tree, err := g.Objects.GetTree(c.TreeHash)
	return &tree, err

}

// Returns the closest ascendant that contains a '.got' directory
func getRepositoryRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "couldn't get working directory")
	}
	currentDir := wd
	for {
		if hasGotDir(currentDir) {
			return currentDir, nil
		}
		currentDir = filepath.Dir(currentDir)
		if currentDir == "/" {
			break
		}
	}
	return "", errors.New("no repository found")
}

func hasGotDir(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		if f.IsDir() && f.Name() == rootDir {
			return true
		}
	}
	return false
}

func IsInitialized(dir string) bool {
	return filesystem.DirExists(filepath.Join(dir, rootDir)) &&
		filesystem.DirExists(filepath.Join(dir, objectsDir)) &&
		filesystem.FileExists(filepath.Join(dir, indexFile)) &&
		filesystem.FileExists(filepath.Join(dir, headFile))
}

func Initialize(dir string) error {
	if IsInitialized(dir) {
		return fmt.Errorf("Repository already exists for %s\n", dir)
	}
	err := filesystem.MkDirIfIsNotExist(filepath.Join(dir, rootDir), os.ModePerm)
	if err != nil {
		return err
	}
	err = filesystem.MkDirIfIsNotExist(filepath.Join(dir, objectsDir), os.ModePerm)
	if err != nil {
		return err
	}
	err = filesystem.MkFileIfIsNotExist(filepath.Join(dir, indexFile))
	if err != nil {
		return err
	}
	err = filesystem.MkFileIfIsNotExist(filepath.Join(dir, headFile))
	if err != nil {
		return err
	}
	return nil
}
