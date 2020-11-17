package simple

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"got/internal/diff"
)

type Diff struct{}

func (d Diff) DiffBytes(a []byte, b []byte) diff.BytesDiff {
	return diffBytes(a, b)
}

func (d Diff) FilesDiff(a []byte, b []byte) bool {
	oldHash := fmt.Sprintf("%x", sha1.Sum(a))
	newHash := fmt.Sprintf("%x", sha1.Sum(b))
	return oldHash != newHash
}

func (d Diff) DiffFiles(a, b []byte) (diff.FileEditType, error) {
	if a == nil && b == nil {
		return "", nil
	}
	return d.diffFiles(a, b)
}

func (d Diff) diffFiles(a []byte, b []byte) (diff.FileEditType, error) {
	if a == nil {
		return diff.FileEditTypeCreate, nil
	}
	if b == nil {
		return diff.FileEditTypeDelete, nil
	}
	byteDiff := diffBytes(a, b)
	for _, le := range byteDiff {
		if le.EditType != diff.EQL {
			return diff.FileEditTypeInPlace, nil
		}
	}
	return "", errors.New("couldn't determine diff type")
}

type point struct {
	A int
	B int
}

type path struct {
	From point
	To   point
}

func diffBytes(as, bs []byte) diff.BytesDiff {
	a := strings.Split(string(as), "\n")
	b := strings.Split(string(bs), "\n")
	n := len(a)
	m := len(b)
	var v = make([][]int, m+1)
	for i := range v {
		v[i] = make([]int, n+1)
	}

	// Create the edit matrix
	for i := 1; i < len(v); i++ {
		for j := 1; j < len(v[i]); j++ {
			if a[j-1] == b[i-1] {
				v[i][j] = v[i-1][j-1] + 1
				continue
			}
			v[i][j] = max(v[i-1][j], v[i][j-1])
		}
	}

	// Backtrack through the matrix to get the longest
	var points []point
	for {
		points = append(points, point{A: n, B: m})
		d := v[m][n]
		if n == 0 && m == 0 {
			break
		}
		if n == 0 {
			m--
		} else if m == 0 {
			n--
		} else if v[m-1][n] == d {
			m--
		} else if v[m][n-1] == d {
			n--
		} else if a[n-1] == b[m-1] && v[m-1][n-1] == d-1 {
			d--
			n--
			m--
		}
	}

	// Create the line diffs
	var edits []diff.LineEdit
	for _, p := range reverse(toPath(points)) {
		var aLine, bLine string
		if p.From.A != len(a) {
			aLine = a[p.From.A]
		}
		if p.From.B != len(b) {
			bLine = b[p.From.B]
		}

		if p.From.A == p.To.A {
			edits = append(edits, diff.NewLineEdit(diff.INS, bLine, -1, p.From.B))
		} else if p.From.B == p.To.B {
			edits = append(edits, diff.NewLineEdit(diff.DEL, aLine, p.From.A, -1))
		} else {
			edits = append(edits, diff.NewLineEdit(diff.EQL, aLine, p.From.A, p.From.B))
		}
	}
	return edits
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func toPath(points []point) []path {
	var paths []path
	for i := 0; i < len(points)-1; i++ {
		paths = append(paths, path{points[i+1], points[i]})
	}
	return paths
}

func reverse(ps []path) []path {
	var newPs []path
	for i := len(ps) - 1; i >= 0; i-- {
		newPs = append(newPs, ps[i])
	}
	return newPs
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.A, p.B)
}

func (p path) String() string {
	return fmt.Sprintf("%v -> %v", p.From, p.To)
}
