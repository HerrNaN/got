package internal

import (
	"fmt"
	"io/ioutil"
	"os"

	iMemory "got/internal/index/memory"
	"got/internal/objects"
	oMemory "got/internal/objects/memory"
)

const (
	GotDir     = ".got"
	ObjectsDir = GotDir + "/objects"
)

var Objects = oMemory.NewObjects()
var Index = iMemory.NewIndex()

func HashFile(filename string, store bool) string {
	bs, _ := ioutil.ReadFile(filename)
	return Objects.HashObject(bs, store, objects.TypeBlob)
}

func AddToIndex(sum string, filename string) {
	Index.Update(sum, filename)
}

func WriteTree() string {
	fmt.Println("Writing tree...")
	var buf string
	for _, e := range Index.SortedEntries() {
		buf += fmt.Sprintf("%s\n", e.String())
	}
	buf = buf[:len(buf)-1] // Drop last new line
	sum := Objects.HashObject([]byte(buf), true, objects.TypeTree)
	fmt.Printf("%s\n\n", sum)
	return sum
}

func CommitTree(msg string, tree string, parent string) string {
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
	sum := Objects.HashObject([]byte(buf), true, objects.TypeCommit)
	fmt.Printf("%s\n\n", sum)
	return sum
}

func MkDirIfIsNotExist(name string, perm os.FileMode) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		os.Mkdir(name, perm)
	}
}
