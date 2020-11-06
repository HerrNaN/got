package got

type Got interface {
	HashFile(filename string, store bool) string
	AddToIndex(filename string)
	WriteTree() string
	ReadTree(sum string) error
	CommitTree(msg string, tree string, parent string) string
	Status(wd string) ([]string, []string, error)
}
