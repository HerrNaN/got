package status

import (
	"path/filepath"
	"strings"
)

type Tree []*Node

type Node struct {
	path    string
	cs      *Tree
	changes Changes
	tracked bool
}

func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) AddFile(path string, changes Changes, tracked bool) {
	dirs := strings.Split(path, string(filepath.Separator))
	t.addFile(dirs, changes, tracked)
}

func (t *Tree) addFile(splitPath []string, changes Changes, tracked bool) {
	// If the path points to a file in this tree
	if len(splitPath) == 1 {
		for _, n := range *t {
			if n.path == splitPath[0] {
				n.changes = changes.AddChanges(n.changes)
				return
			}
		}
		*t = append(*t, &Node{
			path:    splitPath[0],
			cs:      nil,
			changes: changes,
			tracked: tracked,
		})
		return
	}

	// Otherwise look for matching prefix
	var tree *Tree = nil
	for _, n := range *t {

		// If one found add the changes to that prefix
		if n.path == splitPath[0] {
			n.changes = changes.AddChanges(n.changes)
			if tracked {
				n.tracked = true
			}
			tree = n.cs
			break
		}
	}

	// If none was found create one
	if tree == nil {
		if len(splitPath) == 1 {

		}
		tree = NewTree()
		*t = append(*t, &Node{
			path:    splitPath[0],
			cs:      tree,
			changes: changes,
			tracked: tracked,
		})
	}

	tree.addFile(splitPath[1:], changes, tracked)
}

func (c Changes) AddChanges(changes Changes) Changes {
	if c.Head == "" {
		c.Head = changes.Head
	}
	if c.Worktree == "" {
		c.Worktree = changes.Worktree
	}
	return c
}

func (t *Tree) GetStatus() *Status {
	s := Status{}
	return t.getStatus("", &s)
}

func (t *Tree) getStatus(rel string, s *Status) *Status {
	for _, n := range *t {
		path := filepath.Join(rel, n.path)

		// If n is a file
		if n.cs == nil {
			if n.changes.Worktree != "" {
				s.unstaged = append(s.unstaged, Change{path, &n.changes.Worktree})
			}
			if n.changes.Head != "" {
				s.staged = append(s.staged, Change{path, &n.changes.Head})
			}
			continue
		}

		// If path isn't tracked just add the path into the
		if !n.tracked {
			s.untracked = append(s.untracked, Change{path + "/", nil})
			continue
		}

		// Other
		s = n.cs.getStatus(path, s)
	}
	return s
}
