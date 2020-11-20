package objects

import (
	"crypto/sha1"
	"fmt"
)

type Commit struct {
	TreeID   ID
	ParentID *ID
	Author   string
	Message  string
	Checksum ID
}

func NewCommit(treeID ID, parentID *ID, author string, message string) Commit {
	var checksum ID
	if parentID == nil {
		checksum = IdFromSum(sha1.Sum([]byte(string(treeID) + author + message)))
	} else {
		checksum = IdFromSum(sha1.Sum([]byte(string(treeID) + string(*parentID) + author + message)))
	}

	return Commit{
		TreeID:   treeID,
		ParentID: parentID,
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
	content += fmt.Sprintf("tree %s\n", c.TreeID)
	content += fmt.Sprintf("parent %s\n", c.ParentID)
	content += fmt.Sprintf("author %s\n", c.Author)
	content += fmt.Sprintf("message %s\n", c.Message)
	content += fmt.Sprintf("checksum %s\n", c.Checksum)
	return content
}

func (c Commit) ID() ID {
	return IdFromSum(sha1.Sum([]byte(c.Content())))
}
