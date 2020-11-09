package index

import "got/internal/objects"

type Index interface {
	// Retrieve a sorted list of the entries in the index.
	// NOTE: This list will be sorted with the Entries.Less function.
	SortedEntries() []Entry

	// Add the current contents of a file into the index. The file
	// path should be relative to the repository root.
	AddFile(filename string, sum string) error

	// Add the contents of a tree object (but not the tree object itself)
	// into the index
	AddTreeContents(tree objects.Tree) error

	// Add a tree object into the index
	AddTree(sum string, prefix string) error

	// Returns true if the index contains an entry for a given file.
	HasEntryFor(filename string) bool

	GetEntrySum(filename string) (string, error)
	HasDescendantsInIndex(dir string) bool
}
