package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"got/internal/index"
	"got/internal/index/file"
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
		err := filepath.Walk(m, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
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

func (g *Got) UnstagePath(paths ...string) error {
	var matches []string
	for _, p := range paths {
		ms, err := filepath.Glob(p)
		if err != nil {
			return errors.Wrapf(err, "couldn't unstage path %s", p)
		}
		matches = append(matches, ms...)
	}
	for _, p := range paths {
		err := filepath.Walk(p, func(localPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return g.unstageFile(localPath)
		})
		if err != nil {
			return errors.Wrapf(err, "couldn't unstage path %s", p)
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
	for _, p := range paths {
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return g.discardFile(path)
		})
		if err != nil {
			return errors.Wrapf(err, "couldn't discard path %s", p)
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
