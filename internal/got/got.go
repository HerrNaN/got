package got

import (
	"got/internal/objects"
	"got/internal/status"
)

type Got interface {

	// Calculates the checksum for the contents of a given file and optionally
	// stores the contents of the given file in a new file at
	// '.got/objects/hash[:2]/hash[2:]
	HashFile(filename string, store bool) (objects.ID, error)

	// Adds a given file to the index. NOTE: Does NOT create a blob object of
	// the file.
	UpdateIndex(filename string) error

	// Writes the contents of the index into a tree object and stores that
	// object.
	WriteTree() (objects.ID, error)

	// Reads a tree with a given checksum from the objects directory into the
	// index.
	ReadTree(id objects.ID) error

	// Creates a commit object from the given tree checksum, message and parent
	// commit checksum. NOTE: For the first commit the parent commit should be
	// empty.
	CommitTree(msg string, treeID objects.ID, parentID *objects.ID) (objects.ID, error)

	// Returns the checksum of the commit that the head is currently on.
	Head() (*objects.ID, error)

	// Creates a blob object(s) from the given files (or files if the given
	// path represents a directory) and then adds the file(s) into the index.
	AddPath(paths ...string) error

	// Unstages a file by removing it from the index and is replaced by an
	// earlier version of the file if one exists in the head tree.
	UnstagePath(paths ...string) error

	// Discards the changes to a file in the working tree.
	DiscardPath(paths ...string) error

	// Return the diff between index and HEAD for every path specified
	DiffIndexPath(paths ...string) (string, error)

	// Return the diff between working tree and index for every path described by
	// the given path specifiers
	DiffPathSpec(pathspecs ...string) (string, error)

	// Return the diff between working tree and index for every path specified
	DiffPath(paths ...string) (string, error)

	// Returns a list of untracked, staged and unstaged files from the working directory.
	// NOTE:
	//   tracked = Files that are tracked by the got repository.
	//   staged = Files which are up to date in the index.
	//   unstaged = Files that are not up to date in the index.
	Status() (*status.Status, error)

	// Commits the files in the index to the repository by:
	// 1. Writing the contents of the index into a tree object
	// 2. Creating a commit object from that tree object
	// 3. Updating the head
	Commit(message string) error
}
