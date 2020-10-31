package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"got/internal/index/file"

	"got/internal/got"
	"got/internal/objects/disk"
)

var g = got.NewGot(disk.NewObjects(), file.ReadFromFile())
var sum string

func main() {

	//fmt.Println(g.HashObject([]byte("test content"), true, objects.TypeBlob))
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index.String())

	fmt.Println("\nCreating test.txt...")
	ioutil.WriteFile("test.txt", []byte("version 1"), os.ModePerm)
	sum = g.HashFile("test.txt", true)
	fmt.Println(sum)
	g.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree := g.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit := g.CommitTree("first commit", tree, "")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev := commit

	fmt.Println("\nUpdating test.txt...")
	ioutil.WriteFile("test.txt", []byte("version 2"), os.ModePerm)
	sum = g.HashFile("test.txt", true)
	fmt.Println(sum)
	g.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree = g.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit = g.CommitTree("second commit", tree, prev)
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev = commit

	fmt.Println("\nCreating new.txt...")
	ioutil.WriteFile("new.txt", []byte("new file"), os.ModePerm)
	sum = g.HashFile("new.txt", true)
	fmt.Println(sum)
	g.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree = g.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit = g.CommitTree("third commit", tree, prev)
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev = commit

}
