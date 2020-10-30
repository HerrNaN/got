package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"got/internal"
	"got/internal/objects"
)

var sum string

func main() {
	fmt.Println(internal.Objects.HashObject([]byte("test content"), true, objects.TypeBlob))
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index.String())

	fmt.Println("\nCreating test.txt...")
	ioutil.WriteFile("test.txt", []byte("version 1"), os.ModePerm)
	sum = internal.HashFile("test.txt", true)
	fmt.Println(sum)
	internal.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree := internal.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit := internal.CommitTree("first commit", tree, "")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev := commit

	fmt.Println("\nUpdating test.txt...")
	ioutil.WriteFile("test.txt", []byte("version 2"), os.ModePerm)
	sum = internal.HashFile("test.txt", true)
	fmt.Println(sum)
	internal.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree = internal.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit = internal.CommitTree("second commit", tree, prev)
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev = commit

	fmt.Println("\nCreating new.txt...")
	ioutil.WriteFile("new.txt", []byte("new file"), os.ModePerm)
	sum = internal.HashFile("new.txt", true)
	fmt.Println(sum)
	internal.AddToIndex(sum, "test.txt")
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	//fmt.Printf("Index:\n%v\n", internal.Index)

	tree = internal.WriteTree()
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	commit = internal.CommitTree("third commit", tree, prev)
	//fmt.Printf("[Objects]:\n%v\n", internal.Objects)
	prev = commit

}
