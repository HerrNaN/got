package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"got/internal/got"

	"got/internal/index"
	"got/internal/index/file"
	"got/internal/objects"
	"got/internal/objects/disk"
)

type Got struct {
	Objects objects.Objects
	Index   index.Index
}

func NewGot() (*Got, error) {
	if !IsInitialized() {
		return nil, errors.New("repository not initialized")
	}
	i, err := file.ReadFromFile()
	if err != nil {
		return nil, err
	}
	return &Got{
		Objects: disk.NewObjects(),
		Index:   i,
	}, nil
}

func (g *Got) HashFile(filename string, store bool) string {
	bs, _ := ioutil.ReadFile(filename)
	return g.Objects.HashObject(bs, store, objects.TypeBlob)
}

func (g *Got) AddToIndex(filename string) {
	g.Index.AddFile(filename)
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

func IsInitialized() bool {
	s, err := os.Stat(got.RootDir)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Initialize(wd string) error {
	_, err := os.Stat(got.RootDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(got.RootDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "couldn't setup got directory")
		}
		return nil
	}
	return fmt.Errorf("Repository already exists for %s\n", wd)
}
