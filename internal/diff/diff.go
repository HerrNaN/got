package diff

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type Diff interface {
	DiffBytes(a []byte, b []byte) BytesDiff
}

type FileDiff struct {
	EditType fileEditType
	SrcPerm  os.FileMode
	DstPerm  os.FileMode
	SrcHash  string
	DstHash  string
	SrcPath  string
	DstPath  string
}

func NewInPlaceFileDiff(oldPerm, newPerm os.FileMode, oldHash, newHash string, path string) FileDiff {
	return FileDiff{
		EditType: FileEditTypeInPlace,
		SrcPerm:  oldPerm,
		DstPerm:  newPerm,
		SrcHash:  oldHash,
		DstHash:  newHash,
		SrcPath:  path,
		DstPath:  "",
	}
}

func NewCreateFileDiff(perm os.FileMode, hash string, path string) FileDiff {
	return FileDiff{
		EditType: FileEditTypeCreate,
		SrcPerm:  0,
		DstPerm:  perm,
		SrcHash:  "",
		DstHash:  hash,
		SrcPath:  "",
		DstPath:  path,
	}
}

type fileEditType string

const (
	FileEditTypeInPlace fileEditType = "in-place edit"
	FileEditTypeCopy    fileEditType = "copy-edit"
	FileEditTypeRename  fileEditType = "rename-edit"
	FileEditTypeCreate  fileEditType = "create"
	FileEditTypeDelete  fileEditType = "delete"
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
	EditType editType
	Text     string
	ALine    int
	BLine    int
}

func NewLineEdit(editType editType, text string, ALine int, BLine int) LineEdit {
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

type editType string

const (
	INS editType = "+"
	DEL editType = "-"
	EQL editType = " "
)
