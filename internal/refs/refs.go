package refs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"got/internal/pkg/filesystem"

	"github.com/pkg/errors"

	"got/internal/objects"
)

const Dir = "refs"
const HeadsDir = "heads"

var refRegex *regexp.Regexp

func init() {
	var err error
	refRegex, err = regexp.Compile(fmt.Sprintf("%s\\/%s\\/(\\/[a-zA-Z0-9]+|[a-zA-Z0-9])+", Dir, HeadsDir))
	if err != nil {
		panic("couldn't compile ref regex")
	}
}

// Example: 'refs/heads/master'
type Ref string

func RefFromString(s string) (Ref, error) {
	if refRegex.MatchString(s) {
		return Ref(s), nil
	}
	return "", errors.New("string is not a ref")
}

func (r Ref) Name() string {
	_, file := filepath.Split(string(r))
	return file
}

type Refs struct {
	gotDir string
}

func NewRefs(gotDir string) *Refs {
	return &Refs{gotDir}
}

func (r *Refs) IDFromRef(ref Ref) (objects.ID, error) {
	filename := filepath.Join(r.gotDir, string(ref))
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get head for branch %s", ref)
	}
	id, err := objects.IdFromString(string(bs))
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get head for branch %s", ref)
	}
	return id, nil
}

func (r *Refs) UpdateRef(ref Ref, id objects.ID) error {
	err := ioutil.WriteFile(filepath.Join(r.gotDir, string(ref)), []byte(id), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "couldn't update ref to %s", id)
	}
	return nil
}

func (r *Refs) CreateBranchAt(branchName string, id objects.ID) (Ref, error) {
	ref, err := RefFromString(filepath.Join(r.headsDir(), branchName))
	if err != nil {
		return "", errors.Wrapf(err, "couldn't create branch %s at %s", branchName, id)
	}
	if filesystem.FileExists(string(ref)) {
		return "", errors.Errorf("branch %s already exists", branchName)
	}
	err = ioutil.WriteFile(string(ref), []byte(id), os.ModePerm)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't create branch %s at %s", branchName, id)
	}
	return Ref(filepath.Join(Dir, HeadsDir, branchName)), nil
}

func (r *Refs) HeadForBranch(branchName string) (objects.ID, error) {
	filename := filepath.Join(r.headsDir(), branchName)
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get head for branch %s", branchName)
	}
	id, err := objects.IdFromString(string(bs))
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get head for branch %s", branchName)
	}
	return id, nil
}

func (r *Refs) Branches() ([]string, error) {
	var branches []string
	err := filepath.Walk(r.headsDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		_, branch := filepath.Split(path)
		branches = append(branches, branch)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get list of branches")
	}
	return branches, nil
}

func (r *Refs) headsDir() string {
	return filepath.Join(r.gotDir, Dir, HeadsDir)
}
