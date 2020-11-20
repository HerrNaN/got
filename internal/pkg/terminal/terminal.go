package terminal

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func Height() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, errors.Wrap(err, "couldn't get terminal height")
	}
	heightText := strings.Split(string(out), " ")[0]
	height, err := strconv.Atoi(heightText)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't get terminal height")
	}
	return height, nil
}

func RunLess(s string) error {
	reader := strings.NewReader(s)
	cmd := exec.Command("less", "-R") // '-R' displays colors
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "couldn't run less")
	}
	return nil
}
