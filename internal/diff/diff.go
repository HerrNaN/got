package diff

import (
	"bytes"
	"fmt"
	"os"

	"got/internal/objects"

	"github.com/gookit/color"
)

type Differ interface {
	DiffBytes(a []byte, b []byte) BytesDiff
	FilesDiff(a []byte, b []byte) bool
	DiffFiles(a []byte, b []byte) (FileEditType, error)
}

type FileDiff struct {
	EditType FileEditType
	SrcPerm  os.FileMode
	DstPerm  os.FileMode
	SrcID    objects.ID
	DstID    objects.ID
	SrcPath  string
	DstPath  string
}

func NewInPlaceFileDiff(oldPerm, newPerm os.FileMode, oldHash, newHash objects.ID, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeInPlace,
		SrcPerm:  oldPerm,
		DstPerm:  newPerm,
		SrcID:    oldHash,
		DstID:    newHash,
		SrcPath:  path,
		DstPath:  "",
	}
}

func NewCreateFileDiff(perm os.FileMode, id objects.ID, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeCreate,
		SrcPerm:  0,
		DstPerm:  perm,
		SrcID:    "",
		DstID:    id,
		SrcPath:  "",
		DstPath:  path,
	}
}

func NewDeleteFileDiff(perm os.FileMode, id objects.ID, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeDelete,
		SrcPerm:  perm,
		DstPerm:  0,
		SrcID:    id,
		DstID:    "",
		SrcPath:  path,
		DstPath:  "",
	}
}

func NewUnmodifiedFileDiff(perm os.FileMode, id objects.ID, path string) *FileDiff {
	return &FileDiff{
		EditType: FileEditTypeUnmodified,
		SrcPerm:  perm,
		DstPerm:  perm,
		SrcID:    id,
		DstID:    id,
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

type Hunk struct {
	SrcStart int
	SrcEnd   int
	DstStart int
	DstEnd   int
	Edits    BytesDiff
}

type Hunks []Hunk

func (bd BytesDiff) Strip() Hunks {
	var hunks []Hunk
	var hunk *Hunk = nil
	for i := range bd {
		if bd.hasDiffWithInRangeAt(i, 3) {
			if hunk == nil {
				hunk = &Hunk{SrcStart: bd[i].ALine, DstStart: bd[i].BLine}
			}
			hunk.Edits = append(hunk.Edits, bd[i])
			hunk.SrcEnd = bd[i].ALine
			hunk.DstEnd = bd[i].BLine
		} else {
			if hunk != nil {
				hunks = append(hunks, *hunk)
				hunk = nil
			}
		}
	}
	if hunk != nil {
		hunks = append(hunks, *hunk)
	}
	return hunks
}

func (hs Hunks) String() string {
	buf := bytes.NewBuffer(nil)
	for _, h := range hs {
		fmt.Fprint(buf, h)
	}
	return buf.String()
}

func (bd BytesDiff) hasDiffWithInRangeAt(at int, within int) bool {
	if bd[at].EditType != EQL {
		return true
	}
	for i := within; i > 0; i-- {
		if (at-i >= 0 && bd[at-i].EditType != EQL) ||
			(at+i < len(bd) && bd[at+i].EditType != EQL) {
			return true
		}
	}
	return false
}

func (h Hunk) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, color.Cyan.Sprintf("@@ -%d,%d +%d,%d @@\n", h.SrcStart, h.SrcEnd, h.DstStart, h.DstEnd))
	fmt.Fprint(buf, h.Edits)
	return buf.String()
}

func (bd BytesDiff) String() string {
	buf := bytes.NewBuffer(nil)
	for _, le := range bd {
		fmt.Fprintln(buf, le)
	}
	return buf.String()
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
	col := color.Normal
	switch e.EditType {
	case INS:
		col = color.Green
	case DEL:
		col = color.Red
	}
	return col.Sprintf("%s%s", e.EditType, e.Text)
}

type EditType string

const (
	INS EditType = "+"
	DEL EditType = "-"
	EQL EditType = " "
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
