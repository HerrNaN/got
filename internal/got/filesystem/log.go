package filesystem

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gookit/color"

	"github.com/pkg/errors"

	"got/internal/objects"
)

type Log []LogEntry
type LogEntry objects.Commit

func (l Log) String() string {
	buf := bytes.NewBuffer(nil)
	for _, le := range l {
		fmt.Fprintln(buf, le)
	}
	return buf.String()
}

func (le LogEntry) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, color.Yellow.Sprintf("commit %v\n", le.Checksum))
	fmt.Fprintf(buf, "Author: %v\n", le.Author)
	fmt.Fprintln(buf)
	for _, line := range strings.Split(le.Message, "\n") {
		fmt.Fprintf(buf, "    %s\n", line)
	}
	return buf.String()
}

func (g *Got) Log(n int) (Log, error) {
	headID, err := g.idAtHead()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't show log")
	}
	if headID == nil {
		return nil, nil
	}
	commit, err := g.Objects.GetCommit(*headID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't show log")
	}

	var log Log
	log = append(log, LogEntry(commit))
	n--
	for commit.ParentID != nil {
		if n == 0 {
			break
		}
		commit, err = g.Objects.GetCommit(*commit.ParentID)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't show log")
		}
		log = append(log, LogEntry(commit))
		n--
	}
	return log, nil
}
