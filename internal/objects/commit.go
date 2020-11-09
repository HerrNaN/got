package objects

import (
	"crypto/sha1"
	"fmt"
)

type Commit struct {
	TreeHash string
	Parent   string
	Author   string
	Message  string
	Checksum string
}

func NewCommit(treeHash string, parent string, author string, message string) Commit {
	checksum := fmt.Sprintf("%x", sha1.Sum([]byte(treeHash+parent+author+message)))
	return Commit{
		TreeHash: treeHash,
		Parent:   parent,
		Author:   author,
		Message:  message,
		Checksum: checksum,
	}
}

func (c Commit) Type() Type {
	return TypeCommit
}

func (c Commit) Content() string {
	var content string
	content += fmt.Sprintf("tree %s\n", c.TreeHash)
	content += fmt.Sprintf("parent %s\n", c.Parent)
	content += fmt.Sprintf("author %s\n", c.Author)
	content += fmt.Sprintf("message %s\n", c.Message)
	content += fmt.Sprintf("checksum %s\n", c.Checksum)
	return content
}

func (c Commit) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(c.Content())))
}
