package got

type Got interface {

	// Calculates the checksum for the contents of a given file and optionally
	// stores the contents of the given file in a new file at
	// '.got/objects/hash[:2]/hash[2:]
	HashFile(filename string, store bool) (string, error)

	// Adds a given file to the index. NOTE: Does NOT create a blob object of
	// the file.
	AddToIndex(filename string) error

	// Writes the contents of the index into a tree object and stores that
	// object.
	WriteTree() (string, error)

	// Reads a tree with a given checksum from the objects directory into the
	// index.
	ReadTree(sum string) error

	// Creates a commit object from the given tree checksum, message and parent
	// commit checksum. NOTE: For the first commit the parent commit should be
	// empty.
	CommitTree(msg string, tree string, parent string) (string, error)

	// Creates a blob object from the given file and then adds the file into
	// the index.
	Add(filename string) error

	// Returns a list of staged and unstaged files from the working directory.
	// NOTE:
	//   tracked = Files that are tracked by the got repository.
	//   staged = Files which are up to date in the index.
	//   unstaged = Files that are not up to date in the index.
	Status() ([]string, []string, error)

	// Returns the checksum of the commit that the head is currently on.
	Head() (string, error)

	// Commits the files in the index to the repository by:
	// 1. Writing the contents of the index into a tree object
	// 2. Creating a commit object from that tree object
	// 3. Updating the head
	Commit(message string) error
}
