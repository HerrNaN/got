package got

import (
	"fmt"
	"io/ioutil"

	"got/internal/index"
	"got/internal/objects"
)

const (
	GotRootDir = ".got/"
)

type Got struct {
	objects objects.Objects
	index   index.Index
}

func NewGot(objects objects.Objects, index index.Index) *Got {
	return &Got{objects: objects, index: index}
}

func (g *Got) HashFile(filename string, store bool) string {
	bs, _ := ioutil.ReadFile(filename)
	return g.objects.HashObject(bs, store, objects.TypeBlob)
}

func (g *Got) AddToIndex(sum string, filename string) {
	g.index.Update(sum, filename)
}

func (g *Got) WriteTree() string {
	fmt.Println("Writing tree...")
	var buf string
	for _, e := range g.index.SortedEntries() {
		buf += fmt.Sprintf("%s\n", e.String())
	}
	buf = buf[:len(buf)-1] // Drop last new line
	sum := g.objects.HashObject([]byte(buf), true, objects.TypeTree)
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
	sum := g.objects.HashObject([]byte(buf), true, objects.TypeCommit)
	fmt.Printf("%s\n\n", sum)
	return sum
}
