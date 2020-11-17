package status

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	tree := NewTree()
	tree.AddFile("cmd/got/got.go", Changes{Head: Created}, true)
	tree.AddFile("cmd/got/got.go", Changes{Worktree: Modified}, true)
	tree.AddFile("internal/got/got.go", Changes{Worktree: Modified}, true)
	tree.AddFile("internal/cmd/add.go", Changes{}, false)
	tree.AddFile("internal/cmd/commit.go", Changes{}, false)
	tree.AddFile("internal/status/status.go", Changes{}, false)
	tree.AddFile("internal/status/tree.go", Changes{}, false)
	tree.AddFile("internal/status/tree_test.go", Changes{}, false)
	fmt.Printf("%v\n", tree.GetStatus())
}
