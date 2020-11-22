package status

import (
	"bytes"
	"fmt"

	"github.com/gookit/color"
)

type ChangeType string

const (
	UnModified ChangeType = "unmodified"
	Modified   ChangeType = "modified:"
	Created    ChangeType = "new file:"
	Deleted    ChangeType = "deleted: "
)

type Changes struct {
	Head     ChangeType
	Worktree ChangeType
}

type Change struct {
	path       string
	changeType *ChangeType
}

type Status struct {
	staged    []Change
	unstaged  []Change
	untracked []Change
}

func (s *Status) String() string {
	buf := bytes.NewBuffer(nil)

	if len(s.staged) > 0 {
		fmt.Fprintln(buf, "Changes to be committed:")
		fmt.Fprintln(buf, "  (use \"git restore --staged <file>...\" to unstage)")
		for _, staged := range s.staged {
			fmt.Fprint(buf, color.Green.Sprintf("        %s   %s\n", *staged.changeType, staged.path))
		}
		fmt.Fprintln(buf)
	}

	if len(s.unstaged) > 0 {
		fmt.Fprintln(buf, "Changes not staged for commit:")
		fmt.Fprintln(buf, "  (use \"git add <file>...\" to update what will be committed)")
		fmt.Fprintln(buf, "  (use \"git restore <file>...\" to discard changes in working directory)")
		for _, unstaged := range s.unstaged {
			fmt.Fprint(buf, color.Red.Sprintf("        %s   %s\n", *unstaged.changeType, unstaged.path))
		}
		fmt.Fprintln(buf)
	}

	if len(s.untracked) > 0 {
		fmt.Fprintln(buf, "Untracked files:")
		fmt.Fprintln(buf, "  (use \"git add <file>...\" to include in what will be committed)")
		for _, untracked := range s.untracked {
			fmt.Fprint(buf, color.Red.Sprintf("        %s\n", untracked.path))
		}
	}

	return buf.String()
}
