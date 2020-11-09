package got

type Got interface {
	HashFile(filename string, store bool) (string, error)
	AddToIndex(filename string) error
	WriteTree() (string, error)
	ReadTree(sum string) error
	CommitTree(msg string, tree string, parent string) (string, error)
	Add(filename string) error
	Status() ([]string, []string, error)
}
