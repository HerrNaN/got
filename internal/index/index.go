package index

import "got/internal/objects"

type Index interface {
	// Retrieve a sorted list of the entries in the index.
	// NOTE: This list will be sorted with the Entries.Less function.
	SortedEntries() []Entry

	// Add the current contents of a file into the index
	AddFile(filename string) error

	// Add the contents of a tree object (but not the tree object itself)
	// into the index
	AddTreeContents(tree objects.Tree)

	// Add a tree object into the index
	AddTree(sum string, prefix string)

	// Returns true if the index contains an entry for a given file.
	HasEntryFor(filename string) bool
}
