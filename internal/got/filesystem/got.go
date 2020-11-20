package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"got/internal/diff/simple"

	"got/internal/diff"

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
	Ignores map[string]bool
	Differ  diff.Differ
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
	ignores, err := readIgnores(dir)
	if err != nil {
		return nil, err
	}

	return &Got{
		gotDir:  gotDir,
		dir:     dir,
		Objects: disk.NewObjects(gotDir),
		Index:   i,
		Ignores: ignores,
		Differ:  simple.Diff{},
	}, nil
}

func (g *Got) HashFile(filename string, store bool) (objects.ID, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't hash file %s", filename)
	}
	blob := objects.NewBlob(bs)
	sum := blob.ID()
	if store {
		err = g.Objects.StoreBlob(sum, bs)
		if err != nil {
			return "", errors.Wrapf(err, "couldn't hash file %s", filename)
		}
	}
	return sum, nil
}

func (g *Got) Head() (*objects.ID, error) {
	bs, err := ioutil.ReadFile(filepath.Join(g.dir, headFile))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read HEAD file")
	}
	if len(bs) == 0 {
		return nil, nil
	}
	id, err := objects.IdFromString(string(bs))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read HEAD file")
	}
	return &id, nil
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

func (g *Got) moveHead(id objects.ID) error {
	return ioutil.WriteFile(filepath.Join(g.dir, headFile), []byte(id), os.ModePerm)
}

func (g *Got) headTree() (*objects.Tree, error) {
	head, err := g.Head()
	if err != nil {
		return nil, err
	}
	if head == nil {
		return nil, nil
	}
	c, err := g.Objects.GetCommit(*head)
	if err != nil {
		return nil, err
	}
	tree, err := g.Objects.GetTree(c.TreeID)
	return &tree, err

}

// The path parameter in f is relative to the repository root
func (g *Got) forAllInRepo(dir string, f func(path string, info os.FileInfo, err error) error) error {
	if g.isIgnored(dir) {
		return nil
	}
	return filepath.Walk(dir, func(rel string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Convert path to repository relative path
		rel, err = g.repoRel(rel)
		if err != nil {
			return err
		}

		// Ignore the .got directory
		if rel == rootDir {
			return filepath.SkipDir
		}

		// Ignore files and directories specified in ignore file
		if g.isIgnored(rel) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		return f(rel, info, err)
	})
}

// The path parameter in f is relative to the repository root
func (g *Got) forAllFilesInRepo(dir string, f func(path string, info os.FileInfo, err error) error) error {
	return g.forAllInRepo(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return f(path, info, err)
	})
}

func (g *Got) isIgnored(path string) bool {
	if g.Ignores[path] {
		return true
	}
	dirs := strings.Split(path, string(filepath.Separator))
	for i := range dirs {
		dir := filepath.Join(dirs[:i]...)
		if g.Ignores[dir] || g.Ignores[dir+string(filepath.Separator)] {
			return true
		}
	}
	return false
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

func readIgnores(dir string) (map[string]bool, error) {
	ignores := make(map[string]bool)
	if !filesystem.FileExists(filepath.Join(dir, ".gitignore")) {
		return ignores, nil
	}
	ignoreFile, err := ioutil.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read ignorefile")
	}

	for _, line := range strings.Split(string(ignoreFile), "\n") {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		matches, err := filepath.Glob(line)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't read ignorefile")
		}
		for _, m := range matches {
			ignores[m] = true
		}
	}
	return ignores, nil
}
