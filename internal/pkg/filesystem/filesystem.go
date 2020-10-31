package filesystem

import "os"

func MkDirIfIsNotExist(name string, perm os.FileMode) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		os.Mkdir(name, perm)
	}
}
