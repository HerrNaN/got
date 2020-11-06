package got

type Got interface {
	HashFile(filename string, store bool) string
	AddToIndex(filename string) error
	WriteTree() string
	ReadTree(sum string) error
	CommitTree(msg string, tree string, parent string) string
	Status() ([]string, []string, error)
}
