package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"got/internal/diff"
)

func TestDiffFilesInPlace(t *testing.T) {
	var a = []byte("test\nhest")
	var b = []byte("hest\ntest")

	d := Diff{}
	fd, err := d.DiffFiles(a, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, fd, diff.FileEditTypeInPlace)
}

func TestDiffFilesCreate(t *testing.T) {
	var a []byte = nil
	var b = []byte("hest\ntest")

	d := Diff{}
	fd, err := d.DiffFiles(a, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, fd, diff.FileEditTypeCreate)
}

func TestDiffFilesDelete(t *testing.T) {
	var a = []byte("hest\ntest")
	var b []byte = nil

	d := Diff{}
	fd, err := d.DiffFiles(a, b)
	assert.Equal(t, err, nil)
	assert.Equal(t, fd, diff.FileEditTypeDelete)
}
