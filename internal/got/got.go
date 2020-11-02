package got

import (
	"fmt"
	"io/ioutil"
	"os"

	"got/internal/index"
	"got/internal/objects"
)

const (
	RootDir = ".got/"
)

type Got struct {
	Objects objects.Objects
	Index   index.Index
}

func NewGot(objects objects.Objects, index index.Index) *Got {
	return &Got{Objects: objects, Index: index}
}

func (g *Got) HashFile(filename string, store bool) string {
	bs, _ := ioutil.ReadFile(filename)
	return g.Objects.HashObject(bs, store, objects.TypeBlob)
}

func (g *Got) AddToIndex(sum string, filename string) {
	g.Index.Update(sum, filename)
}

func (g *Got) WriteTree() string {
	fmt.Println("Writing tree...")
	var buf string
	for _, e := range g.Index.SortedEntries() {
		buf += fmt.Sprintf("%s\n", e.String())
	}
	buf = buf[:len(buf)-1] // Drop last new line
	sum := g.Objects.HashObject([]byte(buf), true, objects.TypeTree)
	fmt.Printf("%s\n\n", sum)
	return sum
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
	s, err := os.Stat(RootDir)
	if err != nil {
		return false
	}
	return s.IsDir()
}
