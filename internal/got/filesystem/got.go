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

type Got struct {
	dir     string
	Objects objects.Objects
	Index   index.Index
}

func NewGot() (*Got, error) {
	dir, err := getRepositoryRoot()
	if err != nil {
		return nil, errors.New("repository not initialized")
	}

	if !IsInitialized(dir) {
		return nil, errors.New("repository not initialized")
	}
	i, err := file.ReadFromFile(filepath.Join(dir, rootDir))
	if err != nil {
		return nil, err
	}

	return &Got{
		dir:     dir,
		Objects: disk.NewObjects(dir),
		Index:   i,
	}, nil
}

func (g *Got) HashFile(filename string, store bool) string {
	bs, _ := ioutil.ReadFile(filename)
	return g.Objects.HashObject(bs, store, objects.TypeBlob)
}

func (g *Got) AddToIndex(filename string) error {
	filepath.Join(g.dir, filename)
	abs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(g.dir, abs)
	return g.Index.AddFile(rel)
}

func (g *Got) WriteTree() string {
	//fmt.Println("Writing tree...")
	var buf string
	for _, e := range g.Index.SortedEntries() {
		buf += fmt.Sprintf("%s\n", e.String())
	}
	buf = buf[:len(buf)-1] // Drop last new line
	sum := g.Objects.HashObject([]byte(buf), true, objects.TypeTree)
	//fmt.Printf("%s\n\n", sum)
	return sum
}

func (g *Got) ReadTree(sum string) error {
	tree, err := g.Objects.GetTree(sum)
	if err != nil {
		return errors.Wrapf(err, "couldn't read tree %s", sum)
	}
	g.Index.AddTreeContents(tree)
	return nil
}

func (g *Got) CommitTree(msg string, tree string, parent string) string {
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
	sum := g.Objects.HashObject([]byte(buf), true, objects.TypeCommit)
	fmt.Printf("%s\n\n", sum)
	return sum
}

func (g *Got) Status() ([]string, []string, error) {
	staged, err := g.staged(g.dir)
	if err != nil {
		return nil, nil, err
	}
	unstaged, err := g.unstaged(g.dir)
	if err != nil {
		return nil, nil, err
	}
	return staged, unstaged, nil
}

func (g *Got) unstaged(wd string) ([]string, error) {
	var unstaged []string
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == wd {
			return nil
		}
		path = path[len(wd)+1:]

		// Ignore .git and .got directories
		if info.Name() == ".git" || info.Name() == ".got" {
			return filepath.SkipDir
		}

		// Don't list the contents of a directory that doesn't have any staged files
		if info.IsDir() && !g.Index.HasDescendantsInIndex(path) {
			unstaged = append(unstaged, path+string(filepath.Separator))
			return filepath.SkipDir
		}

		// Only show the file paths
		if info.IsDir() {
			return nil
		}

		unstaged = append(unstaged, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return unstaged, nil
}

func (g *Got) staged(wd string) ([]string, error) {
	var staged []string
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == wd {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		path = path[len(wd)+1:]
		sum := g.HashFile(path, false)
		indexedSum, err := g.Index.GetEntrySum(path)
		if err != nil {
			return nil
		}
		if sum == indexedSum {
			staged = append(staged, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return staged, nil
}

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
		filesystem.DirExists(filepath.Join(dir, rootDir, disk.ObjectsDir)) &&
		filesystem.FileExists(filepath.Join(dir, rootDir, file.IndexFile))
}

func Initialize(dir string) error {
	if IsInitialized(dir) {
		return fmt.Errorf("Repository already exists for %s\n", dir)
	}
	err := filesystem.MkDirIfIsNotExist(filepath.Join(dir, rootDir), os.ModePerm)
	if err != nil {
		return err
	}
	err = filesystem.MkDirIfIsNotExist(filepath.Join(dir, rootDir, disk.ObjectsDir), os.ModePerm)
	if err != nil {
		return err
	}
	return filesystem.MkFileIfIsNotExist(filepath.Join(dir, rootDir, file.IndexFile))
}
