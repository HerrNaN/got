package filesystem

import (
	"os"
)

func MkDirIfIsNotExist(name string, perm os.FileMode) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return os.Mkdir(name, perm)
	}
	return nil
}

func MkFileIfIsNotExist(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		_, err := os.Create(name)
		return err
	}
	return nil
}

func DirExists(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExists(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !s.IsDir()
}
