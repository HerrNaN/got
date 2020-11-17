package diff

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type Differ interface {
	DiffBytes(a []byte, b []byte) BytesDiff
	FilesDiff(a []byte, b []byte) bool
	DiffFiles(a []byte, b []byte) *FileDiff
}

type FileDiff struct {
	EditType FileEditType
	SrcPerm  os.FileMode
	DstPerm  os.FileMode
	SrcHash  string
	DstHash  string
	SrcPath  string
	DstPath  string
}

func NewInPlaceFileDiff(oldPerm, newPerm os.FileMode, oldHash, newHash string, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeInPlace,
		SrcPerm:  oldPerm,
		DstPerm:  newPerm,
		SrcHash:  oldHash,
		DstHash:  newHash,
		SrcPath:  path,
		DstPath:  "",
	}
}

func NewCreateFileDiff(perm os.FileMode, hash string, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeCreate,
		SrcPerm:  0,
		DstPerm:  perm,
		SrcHash:  "",
		DstHash:  hash,
		SrcPath:  "",
		DstPath:  path,
	}
}

func NewDeleteFileDiff(perm os.FileMode, hash string, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeDelete,
		SrcPerm:  perm,
		DstPerm:  0,
		SrcHash:  hash,
		DstHash:  "",
		SrcPath:  path,
		DstPath:  "",
	}
}

func NewUnmodifiedFileDiff(perm os.FileMode, hash string, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeUnmodified,
		SrcPerm:  perm,
		DstPerm:  perm,
		SrcHash:  hash,
		DstHash:  hash,
		SrcPath:  path,
		DstPath:  path,
	}
}

type FileEditType string

const (
	FileEditTypeInPlace    FileEditType = "modified"
	FileEditTypeCreate     FileEditType = "created"
	FileEditTypeDelete     FileEditType = "deleted"
	FileEditTypeUnmodified FileEditType = "unmodified"
)

type BytesDiff []LineEdit

func (bd BytesDiff) String() string {
	var buf string
	for _, le := range bd {
		switch le.EditType {
		case INS:
			buf += color.GreenString("%v\n", le)
		case DEL:
			buf += color.RedString("%v\n", le)
		case EQL:
			buf += fmt.Sprintf("%v\n", le)
		}
	}
	return buf
}

type LineEdit struct {
	EditType EditType
	Text     string
	ALine    int
	BLine    int
}

func NewLineEdit(editType EditType, text string, ALine int, BLine int) LineEdit {
	return LineEdit{EditType: editType, Text: text, ALine: ALine, BLine: BLine}
}

func (e LineEdit) String() string {
	var aline string
	if e.ALine >= 0 {
		aline = strconv.Itoa(e.ALine + 1)
	}
	var bline string
	if e.BLine >= 0 {
		bline = strconv.Itoa(e.BLine + 1)
	}
	return fmt.Sprintf("%v %-1s %-1s %-1s", e.EditType, aline, bline, e.Text)
}

type EditType string

const (
	INS EditType = "+"
	DEL EditType = "-"
	EQL EditType = " "
)
